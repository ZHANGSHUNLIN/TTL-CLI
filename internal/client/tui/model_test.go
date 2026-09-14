package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	clientapp "ttl-cli/internal/client/app"
	"ttl-cli/internal/core/resource"
	"ttl-cli/internal/i18n"

	tea "github.com/charmbracelet/bubbletea"
)

type fakeService struct {
	resources []clientapp.Resource
	findErr   error
	writeErr  error
	deleted   bool
}

func (s *fakeService) ListResources() ([]clientapp.Resource, error) { return s.resources, s.findErr }
func (s *fakeService) FindResources(query string) ([]clientapp.Resource, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	var result []clientapp.Resource
	for _, resource := range s.resources {
		if strings.Contains(resource.Key.Key, query) || strings.Contains(resource.Value.Val, query) {
			result = append(result, resource)
		}
	}
	if len(result) == 0 {
		return nil, &clientapp.ServiceError{Kind: clientapp.ErrorNotFound, Message: "not found"}
	}
	return result, nil
}
func (s *fakeService) CreateResource(key, value string, tags []string) (clientapp.Resource, error) {
	if s.writeErr != nil {
		return clientapp.Resource{}, s.writeErr
	}
	resource := clientapp.Resource{Key: resource.ValJsonKey{Key: key, Type: resource.ORIGIN}, Value: resource.ValJson{Val: value, Tag: tags}}
	s.resources = append(s.resources, resource)
	return resource, nil
}
func (s *fakeService) UpdateResourceValue(key, value string) (clientapp.Resource, error) {
	if s.writeErr != nil {
		return clientapp.Resource{}, s.writeErr
	}
	for index := range s.resources {
		if s.resources[index].Key.Key == key {
			s.resources[index].Value.Val = value
			return s.resources[index], nil
		}
	}
	return clientapp.Resource{}, errors.New("missing")
}
func (s *fakeService) AddResourceTags(key string, tags []string) (clientapp.Resource, error) {
	if s.writeErr != nil {
		return clientapp.Resource{}, s.writeErr
	}
	for index := range s.resources {
		if s.resources[index].Key.Key == key {
			s.resources[index].Value.Tag = append(s.resources[index].Value.Tag, tags...)
			return s.resources[index], nil
		}
	}
	return clientapp.Resource{}, errors.New("missing")
}
func (s *fakeService) DeleteResourceTag(key, tag string) (clientapp.Resource, error) {
	return clientapp.Resource{}, s.writeErr
}
func (s *fakeService) DeleteResourceWithCleanup(string) (clientapp.DeleteResult, error) {
	if s.writeErr != nil {
		return clientapp.DeleteResult{}, s.writeErr
	}
	s.deleted = true
	s.resources = nil
	return clientapp.DeleteResult{}, nil
}

func TestModel_EmptyAndSearchNoResultsAreObservable(t *testing.T) {
	service := &fakeService{}
	model := NewModel(service, 80, 24)
	model = updateModel(t, model, model.Init()())
	if view := model.View(); !strings.Contains(view, "No resources yet") {
		t.Fatalf("View() = %q", view)
	}

	model.query = "missing"
	model = updateModel(t, model, model.loadResourcesCmd(model.query)())
	if view := model.View(); !strings.Contains(view, "No results for \"missing\"") {
		t.Fatalf("View() = %q", view)
	}
}

func TestModel_UsesLocalizedStableCopy(t *testing.T) {
	i18n.Reset()
	if err := i18n.InitWithLanguage("zh-CN"); err != nil {
		t.Fatalf("i18n.InitWithLanguage() error = %v", err)
	}
	t.Cleanup(i18n.Reset)

	model := NewModel(&fakeService{}, 80, 24)
	model = updateModel(t, model, model.Init()())
	view := model.View()
	if !strings.Contains(view, "TTL 资源") || !strings.Contains(view, "暂无资源") {
		t.Fatalf("localized view = %q", view)
	}
}

func TestModel_SaveFailureKeepsEditorAndDraft(t *testing.T) {
	writeErr := errors.New("disk full")
	service := &fakeService{resources: testResources(), writeErr: writeErr}
	model := NewModel(service, 80, 24)
	model = updateModel(t, model, loadMsg{resources: service.resources})
	model.startEdit()
	model.value.SetValue("unsaved draft")
	model.dirty = true
	updated, cmd := model.updateEditor(tea.KeyMsg{Type: tea.KeyCtrlS})
	model = updated.(Model)
	model = updateModel(t, model, cmd())
	if model.screen != editScreen || model.value.Value() != "unsaved draft" || !errors.Is(model.err, writeErr) {
		t.Fatalf("save failure lost state: screen=%v value=%q err=%v", model.screen, model.value.Value(), model.err)
	}
}

func TestModel_DeleteRequiresConfirmationAndFailureKeepsResource(t *testing.T) {
	service := &fakeService{resources: testResources(), writeErr: errors.New("locked")}
	model := NewModel(service, 80, 24)
	model = updateModel(t, model, loadMsg{resources: service.resources})
	model = updateModel(t, model, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	if model.screen != deleteScreen || service.deleted {
		t.Fatalf("delete was not gated: screen=%v deleted=%v", model.screen, service.deleted)
	}
	updated, cmd := model.updateDelete(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	model = updated.(Model)
	model = updateModel(t, model, cmd())
	if model.screen != deleteScreen || len(model.resources) != 1 || model.err == nil {
		t.Fatalf("delete failure state = %+v", model)
	}
}

func TestModel_DirtyEditorRequiresDiscardConfirmation(t *testing.T) {
	model := NewModel(&fakeService{resources: testResources()}, 80, 24)
	model.resources = testResources()
	model.loading = false
	model.startEdit()
	model.dirty = true
	updated, _ := model.updateEditor(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(Model)
	if model.screen != discardScreen {
		t.Fatalf("screen = %v, want discard", model.screen)
	}
	updated, _ = model.updateDiscard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if updated.(Model).screen != editScreen {
		t.Fatalf("discard cancel screen = %v", updated.(Model).screen)
	}
}

func TestModel_NarrowAndLongDetailDoNotPanic(t *testing.T) {
	resources := testResources()
	resources[0].Value.Val = strings.Repeat("长内容", 100)
	model := NewModel(&fakeService{resources: resources}, 50, 12)
	model.resources = resources
	model.loading = false
	if view := model.View(); !strings.Contains(view, "Enter details") {
		t.Fatalf("narrow view = %q", view)
	}
	model.screen = detailScreen
	model.detailOffset = 10
	if view := model.View(); !strings.Contains(view, "scroll") {
		t.Fatalf("detail view = %q", view)
	}
}

func TestModel_OpenFromDetailQuitsAfterSuccessfulOpen(t *testing.T) {
	service := &fakeService{resources: testResources()}
	model := NewModel(service, 80, 24)
	model.resources = service.resources
	model.loading = false
	model.screen = detailScreen
	opened := ""
	model.openResource = func(value string) error {
		opened = value
		return nil
	}

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	model = updated.(Model)
	if !model.busy || cmd == nil {
		t.Fatalf("open state = busy:%v cmd:%v, want async open", model.busy, cmd != nil)
	}
	updated, quitCmd := model.Update(cmd())
	model = updated.(Model)
	if opened != "value" {
		t.Fatalf("opened value = %q, want %q", opened, "value")
	}
	if quitCmd == nil {
		t.Fatal("successful open did not request TUI exit")
	}
	if _, ok := quitCmd().(tea.QuitMsg); !ok {
		t.Fatalf("successful open command returned %T, want tea.QuitMsg", quitCmd())
	}
}

func TestModel_OpenFailureKeepsDetailScreen(t *testing.T) {
	service := &fakeService{resources: testResources()}
	model := NewModel(service, 80, 24)
	model.resources = service.resources
	model.loading = false
	model.screen = detailScreen
	openErr := errors.New("open failed")
	model.openResource = func(string) error { return openErr }

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'o'}})
	model = updated.(Model)
	updated, quitCmd := model.Update(cmd())
	model = updated.(Model)
	if model.screen != detailScreen || model.busy || !errors.Is(model.err, openErr) || quitCmd != nil {
		t.Fatalf("open failure state = screen:%v busy:%v err:%v quit:%v", model.screen, model.busy, model.err, quitCmd != nil)
	}
}

func TestOpenExternalResource_ExtractsMarkdownTarget(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS opener behavior")
	}
	tmp := t.TempDir()
	argsFile := filepath.Join(tmp, "args")
	openScript := filepath.Join(tmp, "open")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s' \"$1\" > %q\n", argsFile)
	if err := os.WriteFile(openScript, []byte(script), 0755); err != nil {
		t.Fatalf("write fake open: %v", err)
	}
	t.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	if err := openExternalResource("[example](https://example.com/path)"); err != nil {
		t.Fatalf("openExternalResource() error = %v", err)
	}
	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("read fake open args: %v", err)
	}
	if string(got) != "https://example.com/path" {
		t.Fatalf("open target = %q, want %q", got, "https://example.com/path")
	}
}

func TestModel_BrowseListPagesAndKeepsControlsVisible(t *testing.T) {
	resources := make([]clientapp.Resource, 0, 30)
	for index := 0; index < 30; index++ {
		resources = append(resources, clientapp.Resource{
			Key:   resource.ValJsonKey{Key: fmt.Sprintf("resource-%02d", index), Type: resource.ORIGIN},
			Value: resource.ValJson{Val: "value"},
		})
	}
	model := NewModel(&fakeService{resources: resources}, 100, 14)
	model = updateModel(t, model, loadMsg{resources: resources})
	view := model.View()
	if !strings.Contains(view, "resource-00") || strings.Contains(view, "resource-29") {
		t.Fatalf("first page view = %q", view)
	}
	if !strings.Contains(view, "PgUp/PgDn page") || !strings.Contains(view, "Page 1/") {
		t.Fatalf("pagination controls missing: %q", view)
	}

	updated, _ := model.updateBrowse(tea.KeyMsg{Type: tea.KeyPgDown})
	model = updated.(Model)
	view = model.View()
	if !strings.Contains(view, "resource-06") || !strings.Contains(view, "Page 2/") {
		t.Fatalf("second page view = %q", view)
	}
	if strings.Contains(view, "resource-00") {
		t.Fatalf("second page still shows first item: %q", view)
	}
}

func TestModel_LastListPageReportsLastPage(t *testing.T) {
	resources := make([]clientapp.Resource, 0, 52)
	for index := 0; index < 52; index++ {
		resources = append(resources, clientapp.Resource{
			Key:   resource.ValJsonKey{Key: fmt.Sprintf("resource-%02d", index), Type: resource.ORIGIN},
			Value: resource.ValJson{Val: "value"},
		})
	}
	model := NewModel(&fakeService{resources: resources}, 180, 40)
	model = updateModel(t, model, loadMsg{resources: resources})
	updated, _ := model.updateBrowse(tea.KeyMsg{Type: tea.KeyPgDown})
	model = updated.(Model)
	if !strings.Contains(model.View(), "Page 2/2") || model.selected != 32 {
		t.Fatalf("last page state: selected=%d offset=%d view=%q", model.selected, model.listOffset, model.View())
	}
	updated, _ = model.updateBrowse(tea.KeyMsg{Type: tea.KeyEnd})
	model = updated.(Model)
	if !strings.Contains(model.View(), "Page 2/2") || model.selected != len(resources)-1 {
		t.Fatalf("end page state: selected=%d offset=%d view=%q", model.selected, model.listOffset, model.View())
	}
}

func TestModel_MouseWheelScrollsWithinTUI(t *testing.T) {
	resources := make([]clientapp.Resource, 0, 4)
	for index := 0; index < 4; index++ {
		resources = append(resources, clientapp.Resource{
			Key:   resource.ValJsonKey{Key: fmt.Sprintf("resource-%d", index), Type: resource.ORIGIN},
			Value: resource.ValJson{Val: fmt.Sprintf("value-%d", index)},
		})
	}
	model := NewModel(&fakeService{resources: resources}, 80, 12)
	model = updateModel(t, model, loadMsg{resources: resources})

	updated, _ := model.Update(tea.MouseMsg{Button: tea.MouseButtonWheelDown})
	model = updated.(Model)
	if model.selected != 1 {
		t.Fatalf("wheel down selected = %d, want 1", model.selected)
	}
	updated, _ = model.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})
	model = updated.(Model)
	if model.selected != 0 {
		t.Fatalf("wheel up selected = %d, want 0", model.selected)
	}

	model.screen = detailScreen
	model.detailOffset = 1
	updated, _ = model.Update(tea.MouseMsg{Button: tea.MouseButtonWheelUp})
	model = updated.(Model)
	if model.detailOffset != 0 {
		t.Fatalf("detail wheel up offset = %d, want 0", model.detailOffset)
	}
	updated, _ = model.Update(tea.MouseMsg{Button: tea.MouseButtonWheelDown})
	model = updated.(Model)
	if model.detailOffset != 1 {
		t.Fatalf("detail wheel down offset = %d, want 1", model.detailOffset)
	}
}

func TestModel_SaveShortcutMatchesOperatingSystem(t *testing.T) {
	model := NewModel(&fakeService{resources: testResources()}, 80, 24)
	model.resources = testResources()
	model.loading = false
	model.startEdit()
	view := model.View()
	want := "Ctrl+S save"
	if runtime.GOOS == "darwin" {
		want = "Command+S save"
	}
	if !strings.Contains(view, want) {
		t.Fatalf("editor shortcut = %q, want %q", view, want)
	}
}

func TestSplitTagsTrimsAndDeduplicates(t *testing.T) {
	got := splitTags(" work, ,home,work ")
	if strings.Join(got, ",") != "work,home" {
		t.Fatalf("splitTags() = %v", got)
	}
}

func updateModel(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()
	updated, _ := model.Update(msg)
	return updated.(Model)
}

func testResources() []clientapp.Resource {
	return []clientapp.Resource{{Key: resource.ValJsonKey{Key: "note", Type: resource.ORIGIN}, Value: resource.ValJson{Val: "value", Tag: []string{"work"}, CreatedAt: 1, UpdatedAt: 2}}}
}
