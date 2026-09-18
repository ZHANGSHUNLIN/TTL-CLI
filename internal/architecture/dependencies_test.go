package architecture

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type listedPackage struct {
	ImportPath string
	Imports    []string
	Deps       []string
}

func TestClientDoesNotImportRemovedPackages(t *testing.T) {
	packages := listPackages(t, "../client/...")
	for _, pkg := range packages {
		for _, imported := range pkg.Imports {
			switch {
			case strings.HasPrefix(pkg.ImportPath, "ttl-cli/internal/client/") && strings.HasPrefix(imported, "ttl-cli/internal/server/"):
				t.Errorf("client package %s imports removed server package %s", pkg.ImportPath, imported)
			case strings.HasPrefix(pkg.ImportPath, "ttl-cli/internal/client/") && imported == "ttl-cli/db":
				t.Errorf("client package %s imports legacy db facade", pkg.ImportPath)
			}
		}
	}
}

func TestCoreDoesNotImportAdapters(t *testing.T) {
	packages := listPackages(t, "../core/...")
	for _, pkg := range packages {
		for _, imported := range pkg.Imports {
			if strings.HasPrefix(imported, "ttl-cli/internal/client/") ||
				strings.HasPrefix(imported, "ttl-cli/internal/server/") ||
				strings.HasPrefix(imported, "ttl-cli/internal/storage/") {
				t.Errorf("core package %s imports adapter package %s", pkg.ImportPath, imported)
			}
		}
	}
}

func TestBackendImplementationRemoved(t *testing.T) {
	for _, path := range []string{"../../cmd/ttl-server", "../server"} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("backend path must not exist: %s", path)
		}
	}
}

func TestLegacyPackagesRemoved(t *testing.T) {
	packages := listPackages(t, "../...", "../../cmd/ttl")
	legacy := map[string]bool{
		"ttl-cli/command": true,
		"ttl-cli/models":  true,
		"ttl-cli/db":      true,
		"ttl-cli/sync":    true,
		"ttl-cli/conf":    true,
		"ttl-cli/crypto":  true,
		"ttl-cli/i18n":    true,
		"ttl-cli/util":    true,
	}
	for _, pkg := range packages {
		if legacy[pkg.ImportPath] {
			t.Errorf("legacy package still exists: %s", pkg.ImportPath)
		}
	}
}

func listPackages(t *testing.T, patterns ...string) []listedPackage {
	t.Helper()
	args := append([]string{"list", "-json"}, patterns...)
	cmd := exec.Command("go", args...)
	cmd.Dir = "."
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list %v: %v", patterns, err)
	}

	decoder := json.NewDecoder(strings.NewReader(string(output)))
	var packages []listedPackage
	for decoder.More() {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); err != nil {
			t.Fatalf("decode go list output: %v", err)
		}
		packages = append(packages, pkg)
	}
	return packages
}
