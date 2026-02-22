// init_test.go contains unit tests for the init command (dkmcp init).
//
// init_test.goはinitコマンド（dkmcp init）のユニットテストを含みます。
package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunInit_CreatesConfigFile tests that runInit creates the config file from the embedded template.
//
// TestRunInit_CreatesConfigFileは、runInitが埋め込みテンプレートから設定ファイルを作成することをテストします。
func TestRunInit_CreatesConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	origWorkspace := initWorkspace
	origForce := initForce
	defer func() {
		initWorkspace = origWorkspace
		initForce = origForce
	}()

	initWorkspace = tmpDir
	initForce = false

	if err := runInit(initCmd, nil); err != nil {
		t.Fatalf("runInit() unexpected error: %v", err)
	}

	configPath := filepath.Join(tmpDir, ".sandbox", "config", "dkmcp.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("expected config file at %s, got error: %v", configPath, err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read created config: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected config file to have content (embedded template should not be empty)")
	}
}

// TestRunInit_CreatesDirectory tests that runInit creates .sandbox/config/ if it does not exist.
//
// TestRunInit_CreatesDirectoryは、.sandbox/config/が存在しない場合に runInit がディレクトリを作成することをテストします。
func TestRunInit_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	// .sandbox/config/ does NOT exist yet

	origWorkspace := initWorkspace
	origForce := initForce
	defer func() {
		initWorkspace = origWorkspace
		initForce = origForce
	}()

	initWorkspace = tmpDir
	initForce = false

	if err := runInit(initCmd, nil); err != nil {
		t.Fatalf("runInit() unexpected error: %v", err)
	}

	configDir := filepath.Join(tmpDir, ".sandbox", "config")
	if _, err := os.Stat(configDir); err != nil {
		t.Errorf("expected config directory to be created at %s: %v", configDir, err)
	}
}

// TestRunInit_RefusesOverwriteWithoutForce tests that runInit returns an error
// when the config file already exists and --force is not set.
//
// TestRunInit_RefusesOverwriteWithoutForceは、設定ファイルが既に存在し--forceが指定されていない場合に
// runInitがエラーを返すことをテストします。
func TestRunInit_RefusesOverwriteWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".sandbox", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "dkmcp.yaml")
	if err := os.WriteFile(configPath, []byte("existing content"), 0644); err != nil {
		t.Fatal(err)
	}

	origWorkspace := initWorkspace
	origForce := initForce
	defer func() {
		initWorkspace = origWorkspace
		initForce = origForce
	}()

	initWorkspace = tmpDir
	initForce = false

	err := runInit(initCmd, nil)
	if err == nil {
		t.Fatal("expected error when config exists and --force not set")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Errorf("expected error to mention --force, got: %v", err)
	}

	// Verify the existing file was not modified.
	// 既存ファイルが変更されていないことを確認します。
	data, _ := os.ReadFile(configPath)
	if string(data) != "existing content" {
		t.Error("existing config file was unexpectedly overwritten")
	}
}

// TestRunInit_ForceOverwritesExisting tests that runInit overwrites the existing config
// when --force is set.
//
// TestRunInit_ForceOverwritesExistingは、--forceが指定された場合に runInit が既存の設定を上書きすることをテストします。
func TestRunInit_ForceOverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".sandbox", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "dkmcp.yaml")
	if err := os.WriteFile(configPath, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}

	origWorkspace := initWorkspace
	origForce := initForce
	defer func() {
		initWorkspace = origWorkspace
		initForce = origForce
	}()

	initWorkspace = tmpDir
	initForce = true

	if err := runInit(initCmd, nil); err != nil {
		t.Fatalf("runInit() unexpected error: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	if string(data) == "old content" {
		t.Error("expected config file to be overwritten, but it was not")
	}
	if len(data) == 0 {
		t.Error("expected overwritten config to have content")
	}
}
