package app

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"ttl-cli/models"
	"ttl-cli/util"
)

// ErrorKind identifies client-service failures without coupling callers to
// localized error messages.
type ErrorKind string
type ErrorOperation string

const (
	ErrorNotFound  ErrorKind = "not_found"
	ErrorConflict  ErrorKind = "conflict"
	ErrorAmbiguous ErrorKind = "ambiguous"
	ErrorSystem    ErrorKind = "system_error"
)

const (
	ErrorRead   ErrorOperation = "read"
	ErrorSave   ErrorOperation = "save"
	ErrorUpdate ErrorOperation = "update"
	ErrorDelete ErrorOperation = "delete"
)

// ServiceError is a classifiable client-service failure.
type ServiceError struct {
	Kind       ErrorKind
	Operation  ErrorOperation
	Message    string
	Candidates []string
	Err        error
}

func (e *ServiceError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Kind)
}

func (e *ServiceError) Unwrap() error { return e.Err }

// ErrorKindOf returns the stable category carried by a service error.
func ErrorKindOf(err error) (ErrorKind, bool) {
	var serviceErr *ServiceError
	if !errors.As(err, &serviceErr) {
		return "", false
	}
	return serviceErr.Kind, true
}

// Resource pairs a canonical resource key with its value.
type Resource struct {
	Key   models.ValJsonKey
	Value models.ValJson
}

// ListResources returns origin resources in a deterministic order.
func (s *Service) ListResources() ([]Resource, error) {
	resources, err := s.GetAllResources()
	if err != nil {
		return nil, systemError(ErrorRead, "failed to read resources", err)
	}

	result := make([]Resource, 0, len(resources))
	for key, value := range resources {
		if key.Type == models.ORIGIN {
			result = append(result, Resource{Key: key, Value: value})
		}
	}
	sortResources(result)
	return result, nil
}

// SearchOptions controls which resource fields participate in substring matching.
type SearchOptions struct {
	IncludeValue bool
	IncludeTags  bool
}

// FindResources applies the existing key, value, and tag substring matching rules.
func (s *Service) FindResources(query string) ([]Resource, error) {
	return s.FindResourcesWithOptions(query, SearchOptions{IncludeValue: true, IncludeTags: true})
}

// FindResourcesWithOptions searches origin resources by key and optional fields.
// The key is always searched; callers opt into value and tag matching explicitly.
func (s *Service) FindResourcesWithOptions(query string, options SearchOptions) ([]Resource, error) {
	if query == "" {
		return s.ListResources()
	}
	resources, err := s.GetAllResources()
	if err != nil {
		return nil, systemError(ErrorRead, "failed to read resources", err)
	}

	result := make([]Resource, 0)
	for key, value := range resources {
		if key.Type != models.ORIGIN {
			continue
		}
		matched := util.ContainsIgnoreCase(key.Key, query)
		if !matched && options.IncludeValue {
			matched = util.ContainsIgnoreCase(value.Val, query)
		}
		if !matched && options.IncludeTags {
			for _, tag := range value.Tag {
				if util.ContainsIgnoreCase(tag, query) {
					matched = true
					break
				}
			}
		}
		if matched {
			result = append(result, Resource{Key: key, Value: value})
		}
	}

	if len(result) == 0 {
		return nil, &ServiceError{Kind: ErrorNotFound, Message: fmt.Sprintf("resource not found: %s", query)}
	}
	sortResources(result)
	return result, nil
}

// CreateResource creates one origin resource and rejects duplicate keys.
func (s *Service) CreateResource(key, value string, tags []string) (Resource, error) {
	resourceKey := models.ValJsonKey{Key: key, Type: models.ORIGIN}
	resources, err := s.GetAllResources()
	if err != nil {
		return Resource{}, systemError(ErrorRead, "failed to read resources", err)
	}
	if _, exists := resources[resourceKey]; exists {
		return Resource{}, &ServiceError{Kind: ErrorConflict, Message: fmt.Sprintf("resource already exists: %s", key)}
	}

	if err := s.SaveResource(resourceKey, models.ValJson{Val: value, Tag: util.RemoveDuplicates(tags)}); err != nil {
		return Resource{}, systemError(ErrorSave, "failed to save resource", err)
	}
	return s.resourceByKey(resourceKey)
}

// UpdateResourceValue updates a resource value while preserving its tags.
func (s *Service) UpdateResourceValue(key, value string) (Resource, error) {
	resourceKey := models.ValJsonKey{Key: key, Type: models.ORIGIN}
	existing, err := s.resourceByKey(resourceKey)
	if err != nil {
		return Resource{}, err
	}
	if err := s.UpdateResource(resourceKey, models.ValJson{
		Val:       value,
		Tag:       existing.Value.Tag,
		CreatedAt: existing.Value.CreatedAt,
	}); err != nil {
		return Resource{}, systemError(ErrorUpdate, "failed to update resource", err)
	}
	return s.resourceByKey(resourceKey)
}

// AddResourceTags adds tags without duplicates.
func (s *Service) AddResourceTags(key string, tags []string) (Resource, error) {
	resourceKey := models.ValJsonKey{Key: key, Type: models.ORIGIN}
	existing, err := s.resourceByKey(resourceKey)
	if err != nil {
		return Resource{}, err
	}
	existing.Value.Tag = util.RemoveDuplicates(append(existing.Value.Tag, tags...))
	if err := s.SaveResource(resourceKey, existing.Value); err != nil {
		return Resource{}, systemError(ErrorSave, "failed to save resource tags", err)
	}
	return s.resourceByKey(resourceKey)
}

// DeleteResourceTag removes all occurrences of one tag.
func (s *Service) DeleteResourceTag(key, tag string) (Resource, error) {
	resourceKey := models.ValJsonKey{Key: key, Type: models.ORIGIN}
	existing, err := s.resourceByKey(resourceKey)
	if err != nil {
		return Resource{}, err
	}
	tags := make([]string, 0, len(existing.Value.Tag))
	for _, existingTag := range existing.Value.Tag {
		if existingTag != tag {
			tags = append(tags, existingTag)
		}
	}
	existing.Value.Tag = tags
	if err := s.SaveResource(resourceKey, existing.Value); err != nil {
		return Resource{}, systemError(ErrorSave, "failed to save resource tags", err)
	}
	return s.resourceByKey(resourceKey)
}

// DeleteResult reports non-fatal cleanup failures after a resource deletion.
type DeleteResult struct {
	HistoryCleanupError error
	AuditCleanupError   error
}

// DeleteResourceWithCleanup deletes one origin resource after best-effort auxiliary cleanup.
func (s *Service) DeleteResourceWithCleanup(key string) (DeleteResult, error) {
	resourceKey := models.ValJsonKey{Key: key, Type: models.ORIGIN}
	if _, err := s.resourceByKey(resourceKey); err != nil {
		return DeleteResult{}, err
	}
	historyErr, auditErr := s.CleanupResourceHistory(key)
	if err := s.DeleteResource(resourceKey); err != nil {
		return DeleteResult{}, systemError(ErrorDelete, "failed to delete resource", err)
	}
	return DeleteResult{HistoryCleanupError: historyErr, AuditCleanupError: auditErr}, nil
}

// GetResource returns one origin resource by exact key.
func (s *Service) GetResource(key string) (Resource, error) {
	return s.resourceByKey(models.ValJsonKey{Key: key, Type: models.ORIGIN})
}

func (s *Service) resourceByKey(key models.ValJsonKey) (Resource, error) {
	resources, err := s.GetAllResources()
	if err != nil {
		return Resource{}, systemError(ErrorRead, "failed to read resources", err)
	}
	value, exists := resources[key]
	if !exists {
		return Resource{}, &ServiceError{Kind: ErrorNotFound, Message: fmt.Sprintf("resource not found: %s", key.Key)}
	}
	return Resource{Key: key, Value: value}, nil
}

func sortResources(resources []Resource) {
	sort.Slice(resources, func(i, j int) bool {
		if resources[i].Value.CreatedAt != resources[j].Value.CreatedAt {
			return resources[i].Value.CreatedAt > resources[j].Value.CreatedAt
		}
		left, right := displayKey(resources[i].Key), displayKey(resources[j].Key)
		if left != right {
			return strings.Compare(left, right) < 0
		}
		return resources[i].Key.Key < resources[j].Key.Key
	})
}

func displayKey(key models.ValJsonKey) string {
	if key.Type == models.TAG && key.OriginKey != "" {
		return key.OriginKey
	}
	return key.Key
}

func systemError(operation ErrorOperation, message string, err error) error {
	return &ServiceError{Kind: ErrorSystem, Operation: operation, Message: message, Err: err}
}
