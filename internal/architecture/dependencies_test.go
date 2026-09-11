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
}

func TestClientAndServerDoNotImportEachOther(t *testing.T) {
	packages := listPackages(t, "../client/...", "../server/...")
	for _, pkg := range packages {
		for _, imported := range pkg.Imports {
			switch {
			case pkg.ImportPath != "ttl-cli/internal/client/cli" && strings.HasPrefix(pkg.ImportPath, "ttl-cli/internal/client/") && strings.HasPrefix(imported, "ttl-cli/internal/server/"):
				t.Errorf("client package %s imports server package %s", pkg.ImportPath, imported)
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

func TestServerCommandDoesNotDependOnClientOrLegacyDB(t *testing.T) {
	packages := listPackages(t, "../../cmd/ttl-server", "../server/...")
	for _, pkg := range packages {
		for _, imported := range pkg.Imports {
			if strings.HasPrefix(imported, "ttl-cli/internal/client/") {
				t.Errorf("server package %s imports client package %s", pkg.ImportPath, imported)
			}
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
