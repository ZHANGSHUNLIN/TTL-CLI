package sync

import (
	"testing"
	"ttl-cli/internal/core/resource"
)

func makeResources(items map[string]string) map[resource.ValJsonKey]resource.ValJson {
	resources := make(map[resource.ValJsonKey]resource.ValJson, len(items))
	for k, v := range items {
		resources[resource.ValJsonKey{Key: k, Type: resource.ORIGIN}] = resource.ValJson{Val: v, Tag: []string{}}
	}
	return resources
}

func makeResourcesWithTags(items map[string]struct {
	Val  string
	Tags []string
}) map[resource.ValJsonKey]resource.ValJson {
	resources := make(map[resource.ValJsonKey]resource.ValJson, len(items))
	for k, v := range items {
		resources[resource.ValJsonKey{Key: k, Type: resource.ORIGIN}] = resource.ValJson{Val: v.Val, Tag: v.Tags}
	}
	return resources
}

func TestDiff_BothEmpty(t *testing.T) {
	local := make(map[resource.ValJsonKey]resource.ValJson)
	remote := make(map[resource.ValJsonKey]resource.ValJson)

	diff := ComputeDiff(local, remote)
	if !diff.InSync {
		t.Error("两端都为空，应该是 InSync")
	}
	if len(diff.LocalOnly) != 0 || len(diff.RemoteOnly) != 0 || len(diff.Conflicts) != 0 {
		t.Error("不应有任何差异项")
	}
}

func TestDiff_LocalOnly(t *testing.T) {
	local := makeResources(map[string]string{"a": "v1", "b": "v2"})
	remote := make(map[resource.ValJsonKey]resource.ValJson)

	diff := ComputeDiff(local, remote)
	if diff.InSync {
		t.Error("不应 InSync")
	}
	if len(diff.LocalOnly) != 2 {
		t.Errorf("期望 2 个 local_only，实际: %d", len(diff.LocalOnly))
	}
	if len(diff.RemoteOnly) != 0 {
		t.Errorf("不应有 remote_only，实际: %d", len(diff.RemoteOnly))
	}
}

func TestDiff_RemoteOnly(t *testing.T) {
	local := make(map[resource.ValJsonKey]resource.ValJson)
	remote := makeResources(map[string]string{"x": "v1"})

	diff := ComputeDiff(local, remote)
	if diff.InSync {
		t.Error("不应 InSync")
	}
	if len(diff.RemoteOnly) != 1 {
		t.Errorf("期望 1 个 remote_only，实际: %d", len(diff.RemoteOnly))
	}
	if diff.RemoteOnly[0].Key != "x" {
		t.Errorf("期望 key=x，实际: %s", diff.RemoteOnly[0].Key)
	}
}

func TestDiff_Conflict(t *testing.T) {
	local := makeResources(map[string]string{"shared": "local-val"})
	remote := makeResources(map[string]string{"shared": "remote-val"})

	diff := ComputeDiff(local, remote)
	if diff.InSync {
		t.Error("不应 InSync")
	}
	if len(diff.Conflicts) != 1 {
		t.Fatalf("期望 1 个 conflict，实际: %d", len(diff.Conflicts))
	}
	c := diff.Conflicts[0]
	if c.Key != "shared" {
		t.Errorf("期望 key=shared，实际: %s", c.Key)
	}
	if c.LocalVal.Val != "local-val" {
		t.Errorf("期望 local=local-val，实际: %s", c.LocalVal.Val)
	}
	if c.RemoteVal.Val != "remote-val" {
		t.Errorf("期望 remote=remote-val，实际: %s", c.RemoteVal.Val)
	}
}

func TestDiff_ConflictByTag(t *testing.T) {
	type item struct {
		Val  string
		Tags []string
	}
	local := makeResourcesWithTags(map[string]struct {
		Val  string
		Tags []string
	}{"k": {Val: "same", Tags: []string{"a"}}})

	remote := makeResourcesWithTags(map[string]struct {
		Val  string
		Tags []string
	}{"k": {Val: "same", Tags: []string{"a", "b"}}})

	diff := ComputeDiff(local, remote)
	if len(diff.Conflicts) != 1 {
		t.Errorf("tag 不同也应为 conflict，实际 conflicts: %d", len(diff.Conflicts))
	}
}

func TestDiff_Mixed(t *testing.T) {
	local := makeResources(map[string]string{
		"only-local":  "v1",
		"shared-same": "same",
		"shared-diff": "local-v",
	})
	remote := makeResources(map[string]string{
		"only-remote": "v2",
		"shared-same": "same",
		"shared-diff": "remote-v",
	})

	diff := ComputeDiff(local, remote)
	if diff.InSync {
		t.Error("不应 InSync")
	}
	if len(diff.LocalOnly) != 1 {
		t.Errorf("期望 1 个 local_only，实际: %d", len(diff.LocalOnly))
	}
	if len(diff.RemoteOnly) != 1 {
		t.Errorf("期望 1 个 remote_only，实际: %d", len(diff.RemoteOnly))
	}
	if len(diff.Conflicts) != 1 {
		t.Errorf("期望 1 个 conflict，实际: %d", len(diff.Conflicts))
	}
}

func TestDiff_SameData_InSync(t *testing.T) {
	local := makeResources(map[string]string{"a": "v1", "b": "v2"})
	remote := makeResources(map[string]string{"a": "v1", "b": "v2"})

	diff := ComputeDiff(local, remote)
	if !diff.InSync {
		t.Error("数据相同应该 InSync")
	}
}

type mockStorage struct {
	resources map[resource.ValJsonKey]resource.ValJson
}

func newMockStorage(data map[string]string) *mockStorage {
	ms := &mockStorage{resources: make(map[resource.ValJsonKey]resource.ValJson)}
	for k, v := range data {
		ms.resources[resource.ValJsonKey{Key: k, Type: resource.ORIGIN}] = resource.ValJson{Val: v, Tag: []string{}}
	}
	return ms
}

func (ms *mockStorage) Init() error  { return nil }
func (ms *mockStorage) Close() error { return nil }
func (ms *mockStorage) GetAllResources() (map[resource.ValJsonKey]resource.ValJson, error) {
	return ms.resources, nil
}
func (ms *mockStorage) SaveResource(key resource.ValJsonKey, value resource.ValJson) error {
	ms.resources[key] = value
	return nil
}
func (ms *mockStorage) DeleteResource(key resource.ValJsonKey) error {
	delete(ms.resources, key)
	return nil
}
func (ms *mockStorage) UpdateResource(key resource.ValJsonKey, newValue resource.ValJson) error {
	ms.resources[key] = newValue
	return nil
}
func (ms *mockStorage) GetTagStats() ([]resource.TagStat, error) {
	return nil, nil
}
func (ms *mockStorage) SaveAuditRecord(_ resource.AuditRecord) error { return nil }
func (ms *mockStorage) GetAuditStats() (resource.AuditStats, error) {
	return resource.AuditStats{}, nil
}
func (ms *mockStorage) GetAllAuditRecords() ([]resource.AuditRecord, error) {
	return nil, nil
}
func (ms *mockStorage) DeleteAuditRecords(_ string) error                { return nil }
func (ms *mockStorage) SaveHistoryRecord(_ resource.HistoryRecord) error { return nil }
func (ms *mockStorage) GetAllHistoryRecords() ([]resource.HistoryRecord, error) {
	return nil, nil
}
func (ms *mockStorage) GetHistoryRecord(_ int, _ resource.SortOrder) (resource.HistoryRecord, error) {
	return resource.HistoryRecord{}, nil
}
func (ms *mockStorage) GetHistoryStats() (resource.HistoryStats, error) {
	return resource.HistoryStats{}, nil
}
func (ms *mockStorage) DeleteHistoryRecords(_ string) error      { return nil }
func (ms *mockStorage) SaveLogRecord(_ resource.LogRecord) error { return nil }
func (ms *mockStorage) GetLogRecords(_, _ string) ([]resource.LogRecord, error) {
	return nil, nil
}
func (ms *mockStorage) DeleteLogRecord(_ int64) error { return nil }
func TestPull_Success(t *testing.T) {
	localStorage := newMockStorage(map[string]string{"local-only": "v1", "shared": "old"})
	remoteStorage := newMockStorage(map[string]string{"remote-only": "v2", "shared": "new"})

	localRes, _ := localStorage.GetAllResources()
	remoteRes, _ := remoteStorage.GetAllResources()
	diff := ComputeDiff(localRes, remoteRes)

	err := ExecutePull(diff, localStorage, remoteStorage, false)
	if err != nil {
		t.Fatalf("Pull 失败: %v", err)
	}

	finalLocal, _ := localStorage.GetAllResources()
	if len(finalLocal) != 2 {
		t.Errorf("pull 后期望 2 个资源，实际: %d", len(finalLocal))
	}

	sharedKey := resource.ValJsonKey{Key: "shared", Type: resource.ORIGIN}
	if finalLocal[sharedKey].Val != "new" {
		t.Errorf("shared 应被远程覆盖为 new，实际: %s", finalLocal[sharedKey].Val)
	}

	localOnlyKey := resource.ValJsonKey{Key: "local-only", Type: resource.ORIGIN}
	if _, exists := finalLocal[localOnlyKey]; exists {
		t.Error("local-only 应被删除")
	}
}

func TestPush_Success(t *testing.T) {
	localStorage := newMockStorage(map[string]string{"local-only": "v1", "shared": "local-v"})
	remoteStorage := newMockStorage(map[string]string{"remote-only": "v2", "shared": "remote-v"})

	localRes, _ := localStorage.GetAllResources()
	remoteRes, _ := remoteStorage.GetAllResources()
	diff := ComputeDiff(localRes, remoteRes)

	err := ExecutePush(diff, localStorage, remoteStorage, false)
	if err != nil {
		t.Fatalf("Push 失败: %v", err)
	}

	finalRemote, _ := remoteStorage.GetAllResources()
	if len(finalRemote) != 2 {
		t.Errorf("push 后期望 2 个资源，实际: %d", len(finalRemote))
	}

	sharedKey := resource.ValJsonKey{Key: "shared", Type: resource.ORIGIN}
	if finalRemote[sharedKey].Val != "local-v" {
		t.Errorf("shared 应被本地覆盖为 local-v，实际: %s", finalRemote[sharedKey].Val)
	}

	remoteOnlyKey := resource.ValJsonKey{Key: "remote-only", Type: resource.ORIGIN}
	if _, exists := finalRemote[remoteOnlyKey]; exists {
		t.Error("remote-only 应被删除")
	}
}

func TestDryRun(t *testing.T) {
	localStorage := newMockStorage(map[string]string{"a": "v1"})
	remoteStorage := newMockStorage(map[string]string{"b": "v2"})

	localRes, _ := localStorage.GetAllResources()
	remoteRes, _ := remoteStorage.GetAllResources()
	diff := ComputeDiff(localRes, remoteRes)

	_ = ExecutePull(diff, localStorage, remoteStorage, true)

	finalLocal, _ := localStorage.GetAllResources()
	if len(finalLocal) != 1 {
		t.Errorf("dry-run 不应修改数据，期望 1 个资源，实际: %d", len(finalLocal))
	}
	aKey := resource.ValJsonKey{Key: "a", Type: resource.ORIGIN}
	if _, exists := finalLocal[aKey]; !exists {
		t.Error("dry-run 不应删除资源 a")
	}
}
