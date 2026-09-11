package bbolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
	"ttl-cli/conf"
	"ttl-cli/crypto"
	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"

	"go.etcd.io/bbolt"
)

type Storage = corestorage.Storage

type LocalStorage struct {
	db            *bbolt.DB
	dbPath        string
	confFile      string
	encryptionKey []byte
	encrypted     bool
	timeout       int
}

var _ corestorage.Storage = (*LocalStorage)(nil)

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (ls *LocalStorage) SetDBPath(path string) {
	ls.dbPath = path
}

func (ls *LocalStorage) SetConfigFile(path string) {
	ls.confFile = path
}

func (ls *LocalStorage) SetTimeout(timeout int) {
	ls.timeout = timeout
}

func (ls *LocalStorage) Init() error {
	if ls.dbPath == "" {
		dbPath, err := GetDBPath(ls.confFile, "local")
		if err != nil {
			return err
		}
		ls.dbPath = dbPath
	}

	if err := os.MkdirAll(filepath.Dir(ls.dbPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	timeoutSec := ls.timeout
	if timeoutSec <= 0 {
		timeoutSec = 5
	}
	db, err := bbolt.Open(ls.dbPath, 0600, &bbolt.Options{Timeout: time.Duration(timeoutSec) * time.Second})
	if err != nil {
		if err == bbolt.ErrTimeout {
			return fmt.Errorf("数据库文件被锁定，请检查是否有其他 ttl 进程正在运行")
		}
		return fmt.Errorf("打开数据库失败: %w", err)
	}
	ls.db = db

	key, err := crypto.LoadKey()
	if err != nil {
		return fmt.Errorf("加载密钥失败: %w", err)
	}
	if key != nil {
		ls.encryptionKey = key
		ls.encrypted = true
	}

	return ls.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("resources"))
		return err
	})
}

func (ls *LocalStorage) Close() error {
	if ls.db != nil {
		return ls.db.Close()
	}
	return nil
}

func (ls *LocalStorage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	resources := make(map[models.ValJsonKey]models.ValJson)
	var resourceList []struct {
		key models.ValJsonKey
		val models.ValJson
	}

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("resources"))
		if bucket == nil {
			return errors.New("资源桶不存在")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var key models.ValJsonKey
			if err := json.Unmarshal(k, &key); err != nil {
				return fmt.Errorf("解析key失败: %w", err)
			}

			var val models.ValJson
			if err := json.Unmarshal(v, &val); err != nil {
				return fmt.Errorf("解析value失败: %w", err)
			}

			if ls.encrypted && crypto.IsEncrypted(val.Val) {
				decrypted, err := crypto.Decrypt(ls.encryptionKey, val.Val)
				if err != nil {
					return fmt.Errorf("解密val失败: %w", err)
				}
				val.Val = decrypted
			}

			resourceList = append(resourceList, struct {
				key models.ValJsonKey
				val models.ValJson
			}{key, val})
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(resourceList, func(i, j int) bool {
		return resourceList[i].val.CreatedAt > resourceList[j].val.CreatedAt
	})

	for _, item := range resourceList {
		resources[item.key] = item.val
	}

	return resources, nil
}

func (ls *LocalStorage) SaveResource(key models.ValJsonKey, value models.ValJson) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("resources"))
		if bucket == nil {
			return errors.New("资源桶不存在")
		}

		keyBytes, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("序列化key失败: %w", err)
		}

		existingVal := bucket.Get(keyBytes)
		now := time.Now().Unix()

		var saveValue models.ValJson
		if existingVal != nil {
			var existing models.ValJson
			if err := json.Unmarshal(existingVal, &existing); err == nil {
				saveValue = value
				saveValue.CreatedAt = existing.CreatedAt
				saveValue.UpdatedAt = now
			} else {
				saveValue = value
				saveValue.CreatedAt = now
				saveValue.UpdatedAt = now
			}
		} else {
			saveValue = value
			saveValue.CreatedAt = now
			saveValue.UpdatedAt = now
		}

		if ls.encrypted {
			encryptedVal, err := crypto.Encrypt(ls.encryptionKey, saveValue.Val)
			if err != nil {
				return fmt.Errorf("加密val失败: %w", err)
			}
			saveValue = models.ValJson{Val: encryptedVal, Tag: saveValue.Tag, CreatedAt: saveValue.CreatedAt, UpdatedAt: saveValue.UpdatedAt}
		}

		valBytes, err := json.Marshal(saveValue)
		if err != nil {
			return fmt.Errorf("序列化value失败: %w", err)
		}

		return bucket.Put(keyBytes, valBytes)
	})
}

func (ls *LocalStorage) DeleteResource(key models.ValJsonKey) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("resources"))
		if bucket == nil {
			return errors.New("资源桶不存在")
		}

		keyBytes, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("序列化key失败: %w", err)
		}

		return bucket.Delete(keyBytes)
	})
}

func (ls *LocalStorage) UpdateResource(key models.ValJsonKey, newValue models.ValJson) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("resources"))
		if bucket == nil {
			return errors.New("资源桶不存在")
		}

		keyBytes, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("序列化key失败: %w", err)
		}

		saveValue := newValue
		if ls.encrypted {
			encryptedVal, err := crypto.Encrypt(ls.encryptionKey, newValue.Val)
			if err != nil {
				return fmt.Errorf("加密val失败: %w", err)
			}
			saveValue = models.ValJson{Val: encryptedVal, Tag: newValue.Tag}
		}

		valBytes, err := json.Marshal(saveValue)
		if err != nil {
			return fmt.Errorf("序列化value失败: %w", err)
		}

		return bucket.Put(keyBytes, valBytes)
	})
}

func (ls *LocalStorage) GetTagStats() ([]models.TagStat, error) {
	resources, err := ls.GetAllResources()
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]models.TagStat)

	for key, val := range resources {
		if key.Type != models.ORIGIN {
			continue
		}

		for _, tag := range val.Tag {
			if stat, exists := tagMap[tag]; exists {
				stat.Count++
				stat.ResourceKeys = append(stat.ResourceKeys, key.Key)
				tagMap[tag] = stat
			} else {
				tagMap[tag] = models.TagStat{
					Tag:          tag,
					Count:        1,
					ResourceKeys: []string{key.Key},
				}
			}
		}
	}

	stats := make([]models.TagStat, 0, len(tagMap))
	for _, stat := range tagMap {
		stats = append(stats, stat)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Tag < stats[j].Tag
	})

	return stats, nil
}

func (ls *LocalStorage) IsEncryptionEnabled() bool {
	return ls.encrypted
}

func (ls *LocalStorage) EnableEncryption() error {
	if ls.encrypted {
		return errors.New("加密已启用")
	}

	key, err := crypto.LoadKey()
	if err != nil {
		return fmt.Errorf("加载密钥失败: %w", err)
	}
	if key == nil {
		key, err = crypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("生成密钥失败: %w", err)
		}
		if err := crypto.SaveKey(key); err != nil {
			return fmt.Errorf("保存密钥失败: %w", err)
		}
	}

	ls.encryptionKey = key
	ls.encrypted = true

	return ls.migrateToEncrypted()
}

func (ls *LocalStorage) DisableEncryption() error {
	if !ls.encrypted {
		return errors.New("加密未启用")
	}

	err := ls.migrateToPlain()
	if err != nil {
		return fmt.Errorf("解密数据失败: %w", err)
	}

	if err := crypto.DeleteKey(); err != nil {
		return fmt.Errorf("删除密钥失败: %w", err)
	}

	ls.encryptionKey = nil
	ls.encrypted = false
	return nil
}

func (ls *LocalStorage) migrateToEncrypted() error {
	resources := make(map[models.ValJsonKey]models.ValJson)

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("resources"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var key models.ValJsonKey
			if err := json.Unmarshal(k, &key); err != nil {
				return err
			}

			var val models.ValJson
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}

			if !crypto.IsEncrypted(val.Val) {
				resources[key] = val
			}
			return nil
		})
	})

	if err != nil {
		return err
	}

	if len(resources) == 0 {
		return nil
	}

	for key, value := range resources {
		if err := ls.SaveResource(key, value); err != nil {
			return fmt.Errorf("加密资源 [%s] 失败: %w", key.Key, err)
		}
	}

	return nil
}

func (ls *LocalStorage) migrateToPlain() error {
	resources := make(map[models.ValJsonKey]models.ValJson)

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("resources"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var key models.ValJsonKey
			if err := json.Unmarshal(k, &key); err != nil {
				return err
			}

			var val models.ValJson
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}

			if crypto.IsEncrypted(val.Val) {
				resources[key] = val
			}
			return nil
		})
	})

	if err != nil {
		return err
	}

	if len(resources) == 0 {
		return nil
	}

	oldKey := ls.encryptionKey
	ls.encryptionKey = nil
	ls.encrypted = false

	for key, value := range resources {
		decryptedVal, err := crypto.Decrypt(oldKey, value.Val)
		if err != nil {
			return fmt.Errorf("解密资源 [%s] 失败: %w", key.Key, err)
		}
		value.Val = decryptedVal
		if err := ls.SaveResource(key, value); err != nil {
			return fmt.Errorf("保存资源 [%s] 失败: %w", key.Key, err)
		}
	}

	return nil
}

func GetDBPath(confFile string, storageType string) (string, error) {
	var (
		ttlConf models.TtlIni
		err     error
	)
	if confFile != "" {
		ttlConf, err = conf.GetTtlConfFromFile(confFile)
	} else {
		ttlConf, err = conf.GetTtlConf()
	}
	if err != nil {
		return "", err
	}

	workspaceName := ttlConf.Workspace
	if workspaceName == "" {
		workspaceName = "default"
	}

	var baseDir string
	var dbPath string

	if ws, ok := ttlConf.Workspaces[workspaceName]; ok && ws.DbPath != "" {
		return ws.DbPath, nil
	} else if ttlConf.DbPath != "" {
		dbPath = ttlConf.DbPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("获取用户目录失败: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".ttl")

		switch storageType {
		case "sqlite":
			return filepath.Join(baseDir, "data.db"), nil
		case "local", "bbolt":
			return filepath.Join(baseDir, "data.bbolt"), nil
		default:
			return filepath.Join(baseDir, "data.db"), nil
		}
	}

	if dbPath != "" {
		ext := filepath.Ext(dbPath)
		basePath := dbPath[:len(dbPath)-len(ext)]

		switch storageType {
		case "sqlite":
			return basePath + ".db", nil
		case "local", "bbolt":
			return basePath + ".bbolt", nil
		default:
			return basePath + ".db", nil
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %w", err)
	}
	baseDir = filepath.Join(homeDir, ".ttl")

	switch storageType {
	case "sqlite":
		return filepath.Join(baseDir, "data.db"), nil
	case "local", "bbolt":
		return filepath.Join(baseDir, "data.bbolt"), nil
	default:
		return filepath.Join(baseDir, "data.db"), nil
	}
}

func (ls *LocalStorage) SaveAuditRecord(record models.AuditRecord) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("audit"))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte("audit"))
		}

		key := fmt.Sprintf("%s_%s_%d", record.ResourceKey, record.Operation, record.Timestamp)
		valBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("序列化审计记录失败: %w", err)
		}

		return bucket.Put([]byte(key), valBytes)
	})
}

func (ls *LocalStorage) GetAuditStats() (models.AuditStats, error) {
	stats := models.AuditStats{
		ByOperation: make(map[string]int),
		ByResource:  make(map[string]int),
	}

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("audit"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var record models.AuditRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return fmt.Errorf("解析审计记录失败: %w", err)
			}

			stats.TotalOperations += record.Count
			stats.ByOperation[record.Operation] += record.Count
			stats.ByResource[record.ResourceKey] += record.Count

			return nil
		})
	})

	return stats, err
}

func (ls *LocalStorage) GetAllAuditRecords() ([]models.AuditRecord, error) {
	var records []models.AuditRecord

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("audit"))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(_, v []byte) error {
			var record models.AuditRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return fmt.Errorf("解析审计记录失败: %w", err)
			}
			records = append(records, record)
			return nil
		})
	})
	return records, err
}

func (ls *LocalStorage) DeleteAuditRecords(resourceKey string) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("audit"))
		if bucket == nil {
			return nil
		}

		var toDelete [][]byte
		err := bucket.ForEach(func(k, v []byte) error {
			var record models.AuditRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}

			if record.ResourceKey == resourceKey {
				toDelete = append(toDelete, k)
			}
			return nil
		})

		if err != nil {
			return err
		}

		for _, key := range toDelete {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}

		return nil
	})
}

func (ls *LocalStorage) SaveHistoryRecord(record models.HistoryRecord) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("history"))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte("history"))
		}

		key := fmt.Sprintf("%d_%d", record.Timestamp, record.ID)
		valBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("序列化历史记录失败: %w", err)
		}

		return bucket.Put([]byte(key), valBytes)
	})
}

func (ls *LocalStorage) GetAllHistoryRecords() ([]models.HistoryRecord, error) {
	var records []models.HistoryRecord

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("history"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var record models.HistoryRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return fmt.Errorf("解析历史记录失败: %w", err)
			}
			records = append(records, record)
			return nil
		})
	})

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp > records[j].Timestamp
	})

	return records, err
}

func (ls *LocalStorage) GetHistoryRecord(index int, order models.SortOrder) (models.HistoryRecord, error) {
	var record models.HistoryRecord
	var found bool

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("history"))
		if bucket := tx.Bucket([]byte("history")); bucket == nil {
			return fmt.Errorf("history bucket does not exist")
		}

		cursor := bucket.Cursor()

		total := 0
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			total++
		}

		if index < 0 || index >= total {
			return nil
		}

		var targetIndex int
		if order == models.Descending {
			targetIndex = total - 1 - index
		} else if order == models.Ascending {
			targetIndex = index
		}

		idx := 0
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if idx == targetIndex {
				if err := json.Unmarshal(v, &record); err != nil {
					return fmt.Errorf("failed to unmarshal record: %w", err)
				}
				found = true
				return nil
			}
			idx++
		}

		return nil
	})

	if err != nil {
		return models.HistoryRecord{}, err
	}

	if !found {
		return models.HistoryRecord{}, fmt.Errorf("index %d out of bounds", index)
	}

	return record, nil
}

func (ls *LocalStorage) GetHistoryStats() (models.HistoryStats, error) {
	stats := models.HistoryStats{
		ByOperation: make(map[string]int),
		ByResource:  make(map[string]int),
	}

	records, err := ls.GetAllHistoryRecords()
	if err != nil {
		return stats, err
	}

	stats.TotalRecords = len(records)
	stats.Records = records

	for _, record := range records {
		stats.ByOperation[record.Operation]++
		stats.ByResource[record.ResourceKey]++
	}

	return stats, nil
}

func (ls *LocalStorage) DeleteHistoryRecords(resourceKey string) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("history"))
		if bucket == nil {
			return nil
		}

		var toDelete [][]byte
		err := bucket.ForEach(func(k, v []byte) error {
			var record models.HistoryRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}

			if record.ResourceKey == resourceKey {
				toDelete = append(toDelete, k)
			}
			return nil
		})

		if err != nil {
			return err
		}

		for _, key := range toDelete {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}

		return nil
	})
}

func (ls *LocalStorage) SaveLogRecord(record models.LogRecord) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte("logs"))
		}

		key := fmt.Sprintf("%s_%d", record.Date, record.ID)
		valBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("序列化日志记录失败: %w", err)
		}

		return bucket.Put([]byte(key), valBytes)
	})
}

func (ls *LocalStorage) GetLogRecords(startDate, endDate string) ([]models.LogRecord, error) {
	var records []models.LogRecord

	err := ls.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var record models.LogRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return fmt.Errorf("解析日志记录失败: %w", err)
			}

			if startDate != "" && record.Date < startDate {
				return nil
			}
			if endDate != "" && record.Date > endDate {
				return nil
			}

			records = append(records, record)
			return nil
		})
	})

	sort.Slice(records, func(i, j int) bool {
		return records[i].ID > records[j].ID
	})

	return records, err
}

func (ls *LocalStorage) DeleteLogRecord(id int64) error {
	return ls.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("logs"))
		if bucket == nil {
			return fmt.Errorf("未找到该日志记录")
		}

		var targetKey []byte
		err := bucket.ForEach(func(k, v []byte) error {
			var record models.LogRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}
			if record.ID == id {
				targetKey = make([]byte, len(k))
				copy(targetKey, k)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if targetKey == nil {
			return fmt.Errorf("未找到该日志记录")
		}
		return bucket.Delete(targetKey)
	})
}
