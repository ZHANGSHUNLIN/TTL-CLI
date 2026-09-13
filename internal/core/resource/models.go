package resource

const (
	Origin = iota
	Tag
)

// Legacy spellings are kept as names within the canonical resource package so
// callers can migrate without changing the persisted representation.
const (
	ORIGIN = Origin
	TAG    = Tag
)

const Version = "0.0.3"

type Key struct {
	Key       string `json:"key"`
	Type      int    `json:"type"`
	OriginKey string `json:"originKey"`
}

type Value struct {
	Val       string   `json:"val"`
	Tag       []string `json:"tag"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
}

type AuditRecord struct {
	ResourceKey string `json:"resourceKey"`
	Operation   string `json:"operation"`
	Timestamp   int64  `json:"timestamp"`
	Count       int    `json:"count"`
}

type AuditStats struct {
	TotalOperations int            `json:"totalOperations"`
	ByOperation     map[string]int `json:"byOperation"`
	ByResource      map[string]int `json:"byResource"`
}

type HistoryRecord struct {
	ID          int64  `json:"id"`
	ResourceKey string `json:"resourceKey"`
	Operation   string `json:"operation"`
	Timestamp   int64  `json:"timestamp"`
	TimeStr     string `json:"timeStr"`
	Command     string `json:"command"`
	Args        string `json:"args"`
}

type HistoryStats struct {
	TotalRecords int             `json:"totalRecords"`
	Records      []HistoryRecord `json:"records"`
	ByOperation  map[string]int  `json:"byOperation"`
	ByResource   map[string]int  `json:"byResource"`
}

type LogRecord struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	Date      string   `json:"date"`
}

type TagStat struct {
	Tag          string   `json:"tag"`
	Count        int      `json:"count"`
	ResourceKeys []string `json:"resourceKeys"`
}

type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

type ValJsonKey = Key
type ValJson = Value
