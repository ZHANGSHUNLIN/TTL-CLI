package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"
	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
	storagebbolt "ttl-cli/internal/storage/bbolt"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db          *sql.DB
	dbPath      string
	confFile    string
	journalMode string
	cacheSize   int
	busyTimeout int
	synchronous string
}

var _ corestorage.Storage = (*SQLiteStorage)(nil)

func NewSQLiteStorage() *SQLiteStorage {
	journalMode := "WAL"
	if runtime.GOOS == "windows" {
		journalMode = "DELETE"
	}
	return &SQLiteStorage{
		journalMode: journalMode,
		cacheSize:   -64000,
		busyTimeout: 5000,
		synchronous: "NORMAL",
	}
}

func (s *SQLiteStorage) SetDBPath(path string) {
	s.dbPath = path
}

func (s *SQLiteStorage) SetConfigFile(path string) {
	s.confFile = path
}

func (s *SQLiteStorage) Init() error {
	if s.dbPath == "" {
		dbPath, err := storagebbolt.GetDBPath(s.confFile, "sqlite")
		if err != nil {
			return err
		}
		s.dbPath = dbPath
	}

	if err := os.MkdirAll(filepath.Dir(s.dbPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	var err error
	s.db, err = sql.Open("sqlite", s.dbPath+"?_pragma=journal_mode("+s.journalMode+")&_pragma=cache_size("+fmt.Sprintf("%d", s.cacheSize)+")&_pragma=busy_timeout("+fmt.Sprintf("%d", s.busyTimeout)+")&_pragma=synchronous("+s.synchronous+")")
	if err != nil {
		return fmt.Errorf("打开数据库失败: %w", err)
	}

	if err := s.db.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	return s.createTables()
}

func (s *SQLiteStorage) createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS resources (
			key TEXT NOT NULL,
			type TEXT NOT NULL,
			origin_key TEXT,
			value TEXT NOT NULL,
			tags TEXT,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			PRIMARY KEY (key, type)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_resources_created ON resources(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_resources_origin ON resources(type) WHERE type = 'ORIGIN'`,
		`CREATE INDEX IF NOT EXISTS idx_resources_tag ON resources(type) WHERE type = 'TAG'`,
		`CREATE TABLE IF NOT EXISTS audit (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			resource_key TEXT NOT NULL,
			operation TEXT NOT NULL,
			timestamp INTEGER NOT NULL,
			count INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_resource ON audit(resource_key)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit(timestamp)`,
		`CREATE TABLE IF NOT EXISTS history (
			id INTEGER PRIMARY KEY,
			resource_key TEXT,
			operation TEXT NOT NULL,
			timestamp INTEGER NOT NULL,
			time_str TEXT NOT NULL,
			command TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_history_timestamp ON history(timestamp)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY,
			content TEXT NOT NULL,
			tags TEXT,
			created_at TEXT NOT NULL,
			date TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_date ON logs(date)`,
	}

	for _, table := range tables {
		if _, err := s.db.Exec(table); err != nil {
			return fmt.Errorf("创建表失败: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *SQLiteStorage) GetAllResources() (map[resource.ValJsonKey]resource.ValJson, error) {
	resources := make(map[resource.ValJsonKey]resource.ValJson)

	rows, err := s.db.Query("SELECT key, type, origin_key, value, tags, created_at, updated_at FROM resources ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("查询资源失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key, typ, originKey, value, tagsStr string
		var createdAt, updatedAt int64
		if err := rows.Scan(&key, &typ, &originKey, &value, &tagsStr, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("扫描资源行失败: %w", err)
		}

		var tags []string
		if tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
				return nil, fmt.Errorf("解析标签失败: %w", err)
			}
		}

		keyType := resource.ORIGIN
		if typ == "TAG" {
			keyType = resource.TAG
		}

		vjk := resource.ValJsonKey{
			Key:       key,
			Type:      keyType,
			OriginKey: originKey,
		}

		resources[vjk] = resource.ValJson{
			Val:       value,
			Tag:       tags,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return resources, nil
}

func (s *SQLiteStorage) SaveResource(key resource.ValJsonKey, value resource.ValJson) error {
	typ := "ORIGIN"
	if key.Type == resource.TAG {
		typ = "TAG"
	}

	tagsJSON, err := json.Marshal(value.Tag)
	if err != nil {
		return fmt.Errorf("序列化标签失败: %w", err)
	}

	now := time.Now().Unix()

	var existingCreatedAt int64
	err = s.db.QueryRow(
		"SELECT created_at FROM resources WHERE key = ? AND type = ?",
		key.Key, typ,
	).Scan(&existingCreatedAt)

	if err == sql.ErrNoRows {
		_, err = s.db.Exec(
			`INSERT INTO resources (key, type, origin_key, value, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			key.Key, typ, key.OriginKey, value.Val, string(tagsJSON), now, now,
		)
	} else if err != nil {
		return fmt.Errorf("查询资源失败: %w", err)
	} else {
		_, err = s.db.Exec(
			`UPDATE resources SET origin_key = ?, value = ?, tags = ?, updated_at = ? WHERE key = ? AND type = ?`,
			key.OriginKey, value.Val, string(tagsJSON), now, key.Key, typ,
		)
	}

	if err != nil {
		return fmt.Errorf("保存资源失败: %w", err)
	}

	return nil
}

func (s *SQLiteStorage) DeleteResource(key resource.ValJsonKey) error {
	typ := "ORIGIN"
	if key.Type == resource.TAG {
		typ = "TAG"
	}

	_, err := s.db.Exec(
		`DELETE FROM resources WHERE key = ? AND type = ?`,
		key.Key, typ,
	)
	if err != nil {
		return fmt.Errorf("删除资源失败: %w", err)
	}

	return nil
}

func (s *SQLiteStorage) UpdateResource(key resource.ValJsonKey, newValue resource.ValJson) error {
	return s.SaveResource(key, newValue)
}

func (s *SQLiteStorage) GetTagStats() ([]resource.TagStat, error) {
	resources, err := s.GetAllResources()
	if err != nil {
		return nil, err
	}

	tagMap := make(map[string]resource.TagStat)

	for key, val := range resources {
		if key.Type != resource.ORIGIN {
			continue
		}

		for _, tag := range val.Tag {
			if stat, exists := tagMap[tag]; exists {
				stat.Count++
				stat.ResourceKeys = append(stat.ResourceKeys, key.Key)
				tagMap[tag] = stat
			} else {
				tagMap[tag] = resource.TagStat{
					Tag:          tag,
					Count:        1,
					ResourceKeys: []string{key.Key},
				}
			}
		}
	}

	stats := make([]resource.TagStat, 0, len(tagMap))
	for _, stat := range tagMap {
		stats = append(stats, stat)
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Tag < stats[j].Tag
	})

	return stats, nil
}

func (s *SQLiteStorage) SaveAuditRecord(record resource.AuditRecord) error {
	_, err := s.db.Exec(
		`INSERT INTO audit (resource_key, operation, timestamp, count) VALUES (?, ?, ?, ?)`,
		record.ResourceKey, record.Operation, record.Timestamp, record.Count,
	)
	if err != nil {
		return fmt.Errorf("保存审计记录失败: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) GetAuditStats() (resource.AuditStats, error) {
	stats := resource.AuditStats{
		ByOperation: make(map[string]int),
		ByResource:  make(map[string]int),
	}

	rows, err := s.db.Query("SELECT operation, resource_key, count FROM audit")
	if err != nil {
		return stats, fmt.Errorf("查询审计统计失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var operation, resourceKey string
		var count int
		if err := rows.Scan(&operation, &resourceKey, &count); err != nil {
			return stats, fmt.Errorf("扫描审计行失败: %w", err)
		}

		stats.TotalOperations += count
		stats.ByOperation[operation] += count
		stats.ByResource[resourceKey] += count
	}

	return stats, nil
}

func (s *SQLiteStorage) GetAllAuditRecords() ([]resource.AuditRecord, error) {
	var records []resource.AuditRecord

	rows, err := s.db.Query("SELECT resource_key, operation, timestamp, count FROM audit ORDER BY timestamp DESC")
	if err != nil {
		return nil, fmt.Errorf("查询审计记录失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record resource.AuditRecord
		if err := rows.Scan(&record.ResourceKey, &record.Operation, &record.Timestamp, &record.Count); err != nil {
			return nil, fmt.Errorf("扫描审计行失败: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *SQLiteStorage) DeleteAuditRecords(resourceKey string) error {
	_, err := s.db.Exec(`DELETE FROM audit WHERE resource_key = ?`, resourceKey)
	if err != nil {
		return fmt.Errorf("删除审计记录失败: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) SaveHistoryRecord(record resource.HistoryRecord) error {
	_, err := s.db.Exec(
		`INSERT INTO history (id, resource_key, operation, timestamp, time_str, command) VALUES (?, ?, ?, ?, ?, ?)`,
		record.ID, record.ResourceKey, record.Operation, record.Timestamp, record.TimeStr, record.Command,
	)
	if err != nil {
		return fmt.Errorf("保存历史记录失败: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) GetAllHistoryRecords() ([]resource.HistoryRecord, error) {
	var records []resource.HistoryRecord

	rows, err := s.db.Query(`SELECT id, resource_key, operation, timestamp, time_str, command FROM history ORDER BY timestamp DESC`)
	if err != nil {
		return nil, fmt.Errorf("查询历史记录失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record resource.HistoryRecord
		if err := rows.Scan(&record.ID, &record.ResourceKey, &record.Operation, &record.Timestamp, &record.TimeStr, &record.Command); err != nil {
			return nil, fmt.Errorf("扫描历史行失败: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *SQLiteStorage) GetHistoryRecord(index int, order resource.SortOrder) (resource.HistoryRecord, error) {
	var orderBy string
	if order == resource.Descending {
		orderBy = "timestamp DESC"
	} else {
		orderBy = "timestamp ASC"
	}

	var record resource.HistoryRecord
	err := s.db.QueryRow(`
		SELECT id, resource_key, operation, timestamp, time_str, command
		FROM history
		ORDER BY `+orderBy+`
		LIMIT 1 OFFSET ?
	`, index).Scan(&record.ID, &record.ResourceKey, &record.Operation, &record.Timestamp, &record.TimeStr, &record.Command)

	if err != nil {
		if err == sql.ErrNoRows {
			return resource.HistoryRecord{}, fmt.Errorf("index %d out of bounds", index)
		}
		return resource.HistoryRecord{}, fmt.Errorf("查询历史记录失败: %w", err)
	}

	return record, nil
}

func (s *SQLiteStorage) GetHistoryStats() (resource.HistoryStats, error) {
	stats := resource.HistoryStats{
		ByOperation: make(map[string]int),
		ByResource:  make(map[string]int),
	}

	records, err := s.GetAllHistoryRecords()
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

func (s *SQLiteStorage) DeleteHistoryRecords(resourceKey string) error {
	_, err := s.db.Exec(`DELETE FROM history WHERE resource_key = ?`, resourceKey)
	if err != nil {
		return fmt.Errorf("删除历史记录失败: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) SaveLogRecord(record resource.LogRecord) error {
	tagsJSON, err := json.Marshal(record.Tags)
	if err != nil {
		return fmt.Errorf("序列化标签失败: %w", err)
	}

	_, err = s.db.Exec(
		`INSERT INTO logs (id, content, tags, created_at, date) VALUES (?, ?, ?, ?, ?)`,
		record.ID, record.Content, string(tagsJSON), record.CreatedAt, record.Date,
	)
	if err != nil {
		return fmt.Errorf("保存日志记录失败: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) GetLogRecords(startDate, endDate string) ([]resource.LogRecord, error) {
	var records []resource.LogRecord

	query := `SELECT id, content, tags, created_at, date FROM logs WHERE 1=1`
	args := []interface{}{}

	if startDate != "" {
		query += ` AND date >= ?`
		args = append(args, startDate)
	}
	if endDate != "" {
		query += ` AND date <= ?`
		args = append(args, endDate)
	}

	query += ` ORDER BY id DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("查询日志记录失败: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var record resource.LogRecord
		var tagsStr string
		if err := rows.Scan(&record.ID, &record.Content, &tagsStr, &record.CreatedAt, &record.Date); err != nil {
			return nil, fmt.Errorf("扫描日志行失败: %w", err)
		}

		if tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &record.Tags); err != nil {
				return nil, fmt.Errorf("解析标签失败: %w", err)
			}
		}

		records = append(records, record)
	}

	return records, nil
}

func (s *SQLiteStorage) DeleteLogRecord(id int64) error {
	_, err := s.db.Exec(`DELETE FROM logs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除日志记录失败: %w", err)
	}
	return nil
}
