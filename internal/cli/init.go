package cli

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed configs/dkmcp.example.yaml
var exampleConfig []byte

var (
	initWorkspace string
	initForce     bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate dkmcp.yaml config from built-in template",
	Long: `Generate a dkmcp.yaml configuration file in {workspace}/.sandbox/config/
from the built-in template.

Example:
  dkmcp init --workspace ~/projects/my-app`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initWorkspace, "workspace", "", "Target workspace directory (required)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing config file")
	_ = initCmd.MarkFlagRequired("workspace")
}

func runInit(cmd *cobra.Command, args []string) error {
	absWorkspace, err := filepath.Abs(initWorkspace)
	if err != nil {
		return fmt.Errorf("invalid workspace path: %w", err)
	}

	configDir := filepath.Join(absWorkspace, ".sandbox", "config")
	configPath := filepath.Join(configDir, "dkmcp.yaml")

	if _, statErr := os.Stat(configPath); statErr == nil {
		if !initForce {
			return fmt.Errorf("config already exists: %s\nUse --force to overwrite", configPath)
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return fmt.Errorf("failed to check config file: %w", statErr)
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, exampleConfig, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Created: %s\n\n", configPath)
	fmt.Println("Edit the file to configure containers and permissions.")
	fmt.Printf("Then run:\n  dkmcp serve --workspace %s\n", absWorkspace)
	return nil
}
