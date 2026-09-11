package tenant

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var validIDRegex = regexp.MustCompile(`^[a-z0-9_-]{1,32}$`)

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	CreatedAt int64  `json:"created_at"`
	Active    bool   `json:"active"`
}

type UserStore struct {
	filePath string
	users    []User
	mu       sync.RWMutex
}

func NewUserStore(filePath string) *UserStore {
	return &UserStore{filePath: filePath}
}

func (s *UserStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.users = []User{}
			return nil
		}
		return fmt.Errorf("failed to read users.json: %w", err)
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to parse users.json: %w", err)
	}
	s.users = users
	return nil
}

func (s *UserStore) save() error {
	data, err := json.MarshalIndent(s.users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize users.json: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return fmt.Errorf("failed to create user data directory: %w", err)
	}
	return os.WriteFile(s.filePath, data, 0600)
}

func (s *UserStore) AddUser(id, name string) (User, error) {
	if !validIDRegex.MatchString(id) {
		return User{}, fmt.Errorf("ID is invalid, only [a-z0-9_-] allowed, max 32 chars")
	}
	if name == "" {
		return User{}, fmt.Errorf("name cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.users {
		if user.ID == id {
			return User{}, fmt.Errorf("user %s already exists", id)
		}
	}

	key, err := generateAPIKey()
	if err != nil {
		return User{}, fmt.Errorf("failed to generate API Key: %w", err)
	}
	user := User{ID: id, Name: name, APIKey: key, CreatedAt: time.Now().Unix(), Active: true}
	s.users = append(s.users, user)
	if err := s.save(); err != nil {
		s.users = s.users[:len(s.users)-1]
		return User{}, err
	}
	return user, nil
}

func (s *UserStore) FindByAPIKey(apiKey string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.users {
		if s.users[i].APIKey == apiKey {
			user := s.users[i]
			return &user
		}
	}
	return nil
}

func (s *UserStore) FindByID(id string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.users {
		if s.users[i].ID == id {
			user := s.users[i]
			return &user
		}
	}
	return nil
}

func (s *UserStore) ListUsers() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, len(s.users))
	copy(result, s.users)
	return result
}

func (s *UserStore) SetActive(id string, active bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.users {
		if s.users[i].ID == id {
			s.users[i].Active = active
			return s.save()
		}
	}
	return fmt.Errorf("user %s does not exist", id)
}

func (s *UserStore) ResetKey(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.users {
		if s.users[i].ID == id {
			newKey, err := generateAPIKey()
			if err != nil {
				return "", fmt.Errorf("failed to generate API Key: %w", err)
			}
			s.users[i].APIKey = newKey
			if err := s.save(); err != nil {
				return "", err
			}
			return newKey, nil
		}
	}
	return "", fmt.Errorf("user %s does not exist", id)
}

func (s *UserStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.users {
		if s.users[i].ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("user %s does not exist", id)
}

func generateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
