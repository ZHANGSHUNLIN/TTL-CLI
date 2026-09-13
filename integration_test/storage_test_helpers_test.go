package integration_test

import (
	"ttl-cli/internal/client/app"
	"ttl-cli/internal/core/resource"
	corestorage "ttl-cli/internal/core/storage"
)

// testDatabase keeps integration setup explicit without exposing a production
// global storage facade.
type testDatabase struct {
	Stor    corestorage.Storage
	service *app.Service
}

var testDB testDatabase

func (d *testDatabase) InitDB(storageType, cloudAPIURL, cloudAPIKey string, cloudTimeout int, confFile string) error {
	storage, err := app.OpenStorage(storageType, cloudAPIURL, cloudAPIKey, cloudTimeout, confFile)
	if err != nil {
		return err
	}
	d.Stor = storage
	d.service = app.NewService(storage)
	return nil
}

func (d *testDatabase) CloseDB() error {
	if d.Stor == nil {
		return nil
	}
	err := d.Stor.Close()
	d.Stor = nil
	d.service = nil
	return err
}

func (d *testDatabase) svc() *app.Service {
	if d.service == nil {
		d.service = app.NewService(d.Stor)
	}
	return d.service
}

func (d *testDatabase) GetAllResources() (map[resource.ValJsonKey]resource.ValJson, error) {
	return d.svc().GetAllResources()
}
func (d *testDatabase) SaveResource(k resource.ValJsonKey, v resource.ValJson) error {
	return d.svc().SaveResource(k, v)
}
func (d *testDatabase) DeleteResource(k resource.ValJsonKey) error {
	return d.svc().DeleteResource(k)
}
func (d *testDatabase) UpdateResource(k resource.ValJsonKey, v resource.ValJson) error {
	return d.svc().UpdateResource(k, v)
}
func (d *testDatabase) GetTagStats() ([]resource.TagStat, error)    { return d.svc().GetTagStats() }
func (d *testDatabase) RecordAudit(k, op string) error              { return d.svc().RecordAudit(k, op) }
func (d *testDatabase) GetAuditStats() (resource.AuditStats, error) { return d.svc().GetAuditStats() }
func (d *testDatabase) GetAllAuditRecords() ([]resource.AuditRecord, error) {
	return d.svc().GetAllAuditRecords()
}
func (d *testDatabase) DeleteAuditRecords(k string) error { return d.svc().DeleteAuditRecords(k) }
func (d *testDatabase) GetAllHistoryRecords() ([]resource.HistoryRecord, error) {
	return d.svc().GetAllHistoryRecords()
}
func (d *testDatabase) DeleteHistoryRecords(k string) error { return d.svc().DeleteHistoryRecords(k) }
func (d *testDatabase) GetHistoryRecords(i int) (resource.HistoryRecord, error) {
	return d.svc().GetHistoryRecord(i, resource.Descending)
}
func (d *testDatabase) RecordCommandHistory(op, key string, debug bool) error {
	return d.svc().RecordCommandHistory(op, key, debug)
}
func (d *testDatabase) SaveLogRecord(v resource.LogRecord) error { return d.svc().SaveLogRecord(v) }
func (d *testDatabase) GetLogRecords(a, b string) ([]resource.LogRecord, error) {
	return d.svc().GetLogRecords(a, b)
}
func (d *testDatabase) DeleteLogRecord(id int64) error { return d.svc().DeleteLogRecord(id) }
