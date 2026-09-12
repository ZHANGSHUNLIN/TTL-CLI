package app

import (
	"fmt"

	"ttl-cli/internal/client/remote"
	corestorage "ttl-cli/internal/core/storage"
	storagebbolt "ttl-cli/internal/storage/bbolt"
)

func MigrateData(sourceType, targetType, sourceAPIURL, sourceAPIKey string, sourceTimeout int, cloudAPIURL, cloudAPIKey string, cloudTimeout int, debug bool, srcConfFile, dstConfFile string) error {
	fmt.Printf("Starting data migration: %s -> %s\n", sourceType, targetType)

	sourceStorage, err := migrationStorage(sourceType, sourceAPIURL, sourceAPIKey, sourceTimeout, srcConfFile)
	if err != nil {
		return err
	}
	if err := sourceStorage.Init(); err != nil {
		return fmt.Errorf("failed to initialize source storage: %w", err)
	}
	defer sourceStorage.Close()

	var targetStorage corestorage.Storage
	switch targetType {
	case "local":
		storage := storagebbolt.NewLocalStorage()
		storage.SetConfigFile(dstConfFile)
		targetStorage = storage
	case "cloud":
		if cloudAPIURL == "" || cloudAPIKey == "" {
			return fmt.Errorf("target cloud storage requires API URL and key")
		}
		targetStorage = remote.NewStorage(cloudAPIURL, cloudAPIKey, cloudTimeout)
	default:
		return fmt.Errorf("unsupported target storage type: %s", targetType)
	}
	if err := targetStorage.Init(); err != nil {
		return fmt.Errorf("failed to initialize target storage: %w", err)
	}
	defer targetStorage.Close()

	fmt.Println("Reading data from source storage...")
	resources, err := sourceStorage.GetAllResources()
	if err != nil {
		return fmt.Errorf("failed to read source data: %w", err)
	}
	fmt.Printf("Found %d resources to migrate\n", len(resources))

	successCount, failCount := 0, 0
	for key, value := range resources {
		if err := targetStorage.SaveResource(key, value); err != nil {
			fmt.Printf("Failed to migrate resource [%s]: %v\n", key.Key, err)
			failCount++
			continue
		}
		successCount++
		if debug {
			fmt.Printf("Successfully migrated resource: %s\n", key.Key)
		}
	}
	fmt.Printf("Migration completed! Success: %d, Failed: %d\n", successCount, failCount)
	if failCount > 0 {
		return fmt.Errorf("some resources migration failed")
	}
	return nil
}

func migrationStorage(storageType, apiURL, apiKey string, timeout int, confFile string) (corestorage.Storage, error) {
	switch storageType {
	case "local":
		storage := storagebbolt.NewLocalStorage()
		storage.SetConfigFile(confFile)
		return storage, nil
	case "cloud":
		if apiURL == "" || apiKey == "" {
			return nil, fmt.Errorf("source cloud storage requires API URL and key")
		}
		return remote.NewStorage(apiURL, apiKey, timeout), nil
	default:
		return nil, fmt.Errorf("unsupported source storage type: %s", storageType)
	}
}
