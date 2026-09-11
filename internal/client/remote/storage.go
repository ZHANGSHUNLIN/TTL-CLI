package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	corestorage "ttl-cli/internal/core/storage"
	"ttl-cli/models"
)

var _ corestorage.Storage = (*Storage)(nil)

// Storage implements the shared storage contract over the ttl HTTP API.
type Storage struct {
	apiURL     string
	apiKey     string
	timeout    int
	httpClient *http.Client
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type resourceDTO struct {
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Tags      []string `json:"tags"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
}

func NewStorage(apiURL, apiKey string, timeout int) *Storage {
	return &Storage{apiURL: strings.TrimRight(apiURL, "/"), apiKey: apiKey, timeout: timeout}
}

func (s *Storage) Init() error {
	s.httpClient = &http.Client{Timeout: time.Duration(s.timeout) * time.Second}
	return nil
}

func (s *Storage) Close() error {
	if s.httpClient != nil {
		s.httpClient.CloseIdleConnections()
	}
	return nil
}

func (s *Storage) doRequest(method, path string, body any) (*apiResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, s.apiURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(respBody))
	}
	if apiResp.Code != 0 {
		return &apiResp, fmt.Errorf("API 错误: %s", apiResp.Message)
	}
	return &apiResp, nil
}

func (s *Storage) GetAllResources() (map[models.ValJsonKey]models.ValJson, error) {
	apiResp, err := s.doRequest(http.MethodGet, "/api/v1/resources", nil)
	if err != nil {
		return nil, err
	}
	var dtos []resourceDTO
	if err := json.Unmarshal(apiResp.Data, &dtos); err != nil {
		return nil, fmt.Errorf("解析资源列表失败: %w", err)
	}
	sort.Slice(dtos, func(i, j int) bool { return dtos[i].CreatedAt > dtos[j].CreatedAt })

	resources := make(map[models.ValJsonKey]models.ValJson, len(dtos))
	for _, dto := range dtos {
		tags := dto.Tags
		if tags == nil {
			tags = []string{}
		}
		resources[models.ValJsonKey{Key: dto.Key, Type: models.ORIGIN}] = models.ValJson{
			Val: dto.Value, Tag: tags, CreatedAt: dto.CreatedAt, UpdatedAt: dto.UpdatedAt,
		}
	}
	return resources, nil
}

func (s *Storage) SaveResource(key models.ValJsonKey, value models.ValJson) error {
	body := map[string]any{
		"key": key.Key, "value": value.Val, "tags": value.Tag,
		"createdAt": value.CreatedAt, "updatedAt": time.Now().Unix(),
	}
	_, err := s.doRequest(http.MethodPost, "/api/v1/resources", body)
	return err
}

func (s *Storage) DeleteResource(key models.ValJsonKey) error {
	_, err := s.doRequest(http.MethodDelete, "/api/v1/resources/"+key.Key, nil)
	return err
}

func (s *Storage) UpdateResource(key models.ValJsonKey, value models.ValJson) error {
	_, err := s.doRequest(http.MethodPut, "/api/v1/resources/"+key.Key, map[string]string{"value": value.Val})
	return err
}

func (s *Storage) GetTagStats() ([]models.TagStat, error) {
	resources, err := s.GetAllResources()
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]models.TagStat)
	for key, value := range resources {
		if key.Type != models.ORIGIN {
			continue
		}
		for _, tag := range value.Tag {
			stat := tagMap[tag]
			stat.Tag = tag
			stat.Count++
			stat.ResourceKeys = append(stat.ResourceKeys, key.Key)
			tagMap[tag] = stat
		}
	}
	stats := make([]models.TagStat, 0, len(tagMap))
	for _, stat := range tagMap {
		stats = append(stats, stat)
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].Tag < stats[j].Tag })
	return stats, nil
}

func (s *Storage) SaveAuditRecord(models.AuditRecord) error { return nil }

func (s *Storage) GetAuditStats() (models.AuditStats, error) {
	apiResp, err := s.doRequest(http.MethodGet, "/api/v1/audit/stats", nil)
	if err != nil {
		return models.AuditStats{ByOperation: map[string]int{}, ByResource: map[string]int{}}, err
	}
	var stats models.AuditStats
	if err := json.Unmarshal(apiResp.Data, &stats); err != nil {
		return models.AuditStats{ByOperation: map[string]int{}, ByResource: map[string]int{}}, err
	}
	return stats, nil
}

func (s *Storage) DeleteAuditRecords(string) error { return nil }
func (s *Storage) GetAllAuditRecords() ([]models.AuditRecord, error) {
	return []models.AuditRecord{}, nil
}
func (s *Storage) SaveHistoryRecord(models.HistoryRecord) error { return nil }

func (s *Storage) GetAllHistoryRecords() ([]models.HistoryRecord, error) {
	apiResp, err := s.doRequest(http.MethodGet, "/api/v1/history", nil)
	if err != nil {
		return []models.HistoryRecord{}, err
	}
	var records []models.HistoryRecord
	if err := json.Unmarshal(apiResp.Data, &records); err != nil {
		return []models.HistoryRecord{}, err
	}
	return records, nil
}

func (s *Storage) GetHistoryRecord(int, models.SortOrder) (models.HistoryRecord, error) {
	return models.HistoryRecord{}, fmt.Errorf("云端存储不支持按索引获取历史记录")
}
func (s *Storage) GetHistoryStats() (models.HistoryStats, error) {
	return models.HistoryStats{ByOperation: map[string]int{}, ByResource: map[string]int{}}, nil
}
func (s *Storage) DeleteHistoryRecords(string) error    { return nil }
func (s *Storage) SaveLogRecord(models.LogRecord) error { return nil }
func (s *Storage) GetLogRecords(string, string) ([]models.LogRecord, error) {
	return []models.LogRecord{}, nil
}
func (s *Storage) DeleteLogRecord(int64) error { return nil }
