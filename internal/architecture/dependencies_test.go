package architecture

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

type listedPackage struct {
	ImportPath string
	Imports    []string
	Deps       []string
}

func TestClientAndServerDoNotImportEachOther(t *testing.T) {
	packages := listPackages(t, "../client/...", "../server/...")
	for _, pkg := range packages {
		for _, imported := range pkg.Imports {
			switch {
			case strings.HasPrefix(pkg.ImportPath, "ttl-cli/internal/client/") && strings.HasPrefix(imported, "ttl-cli/internal/server/"):
				t.Errorf("client package %s imports server package %s", pkg.ImportPath, imported)
			case strings.HasPrefix(pkg.ImportPath, "ttl-cli/internal/client/") && imported == "ttl-cli/db":
				t.Errorf("client package %s imports legacy db facade", pkg.ImportPath)
			case strings.HasPrefix(pkg.ImportPath, "ttl-cli/internal/server/") && strings.HasPrefix(imported, "ttl-cli/internal/client/"):
				t.Errorf("server package %s imports client package %s", pkg.ImportPath, imported)
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

func TestServerBinary_DependencyBoundary(t *testing.T) {
	packages := listPackages(t, "../../cmd/ttl-server")
	for _, pkg := range packages {
		for _, dependency := range append(pkg.Imports, pkg.Deps...) {
			if strings.HasPrefix(dependency, "ttl-cli/internal/client/") ||
				dependency == "ttl-cli/internal/client/cli/commands" ||
				dependency == "ttl-cli/db" ||
				dependency == "ttl-cli/internal/client/sync" {
				t.Errorf("server binary %s depends on forbidden package %s", pkg.ImportPath, dependency)
			}
		}
	}
}

func TestLegacyPackagesRemoved(t *testing.T) {
	packages := listPackages(t, "../...", "../../cmd/ttl", "../../cmd/ttl-server")
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
