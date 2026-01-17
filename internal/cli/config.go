/*
Copyright © 2026 Orkflow Authors
*/
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Config represents the orka configuration
type Config struct {
	APIKey   string `yaml:"api_key,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Model    string `yaml:"model,omitempty"`
	Provider string `yaml:"provider,omitempty"`
}

var (
	apiKey   string
	name     string
	model    string
	provider string
	show     bool
	global   bool
	local    bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure orka settings",
	Long: `Configure orka settings like API key, project name, model, and provider.

Configuration can be stored globally or locally:
  --global    Save to ~/.orka.yaml (user-wide, default)
  --local     Save to ./.orka.yaml (project-specific)

Local config takes precedence over global config when running orka.

Examples:
  orka config --api "sk-xxx" --global     Set global API key
  orka config --name "my-project" --local Set project-specific name
  orka config --model "gpt-4"             Set default model (global)
  orka config --show                      Show current configuration
  orka config --show --local              Show local project config`,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine which config to use
		configPath := getConfigPathWithScope()

		// If --show flag, display current config
		if show {
			showConfigWithScope()
			return
		}

		// Check if at least one config flag is provided
		if apiKey == "" && name == "" && model == "" && provider == "" {
			fmt.Println("Error: No configuration option provided.")
			fmt.Println()
			fmt.Println("Available options:")
			fmt.Println("  --api <key>        Set API key")
			fmt.Println("  --name <name>      Set project name")
			fmt.Println("  --model <model>    Set default model")
			fmt.Println("  --provider <name>  Set default provider")
			fmt.Println("  --show             Show current configuration")
			fmt.Println()
			fmt.Println("Scope options:")
			fmt.Println("  --global           Save to ~/.orka.yaml (default)")
			fmt.Println("  --local            Save to ./.orka.yaml (project)")
			fmt.Println()
			fmt.Println("Example:")
			fmt.Println("  orka config --api \"sk-xxx\" --provider \"openai\" --global")
			fmt.Println("  orka config --name \"my-project\" --local")
			os.Exit(1)
		}

		// Load existing config or create new
		config := loadConfig(configPath)

		// Update config with provided values
		if apiKey != "" {
			config.APIKey = apiKey
			fmt.Printf("✓ API key set\n")
		}
		if name != "" {
			config.Name = name
			fmt.Printf("✓ Name set to: %s\n", name)
		}
		if model != "" {
			config.Model = model
			fmt.Printf("✓ Model set to: %s\n", model)
		}
		if provider != "" {
			config.Provider = provider
			fmt.Printf("✓ Provider set to: %s\n", provider)
		}

		// Save config
		if err := saveConfig(configPath, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		scope := "global"
		if local {
			scope = "local"
		}
		fmt.Printf("\nConfiguration saved to: %s (%s)\n", configPath, scope)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVar(&apiKey, "api", "", "API key for the AI provider")
	configCmd.Flags().StringVar(&name, "name", "", "Project name")
	configCmd.Flags().StringVar(&model, "model", "", "Default AI model to use")
	configCmd.Flags().StringVar(&provider, "provider", "", "AI provider (openai, anthropic, etc.)")
	configCmd.Flags().BoolVar(&show, "show", false, "Show current configuration")
	configCmd.Flags().BoolVar(&global, "global", false, "Use global config (~/.orka.yaml)")
	configCmd.Flags().BoolVar(&local, "local", false, "Use local config (./.orka.yaml)")
}

func getGlobalConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".orka.yaml")
}

func getLocalConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(cwd, ".orka.yaml")
}

func getConfigPathWithScope() string {
	if local {
		return getLocalConfigPath()
	}
	// Default to global
	return getGlobalConfigPath()
}

func loadConfig(path string) *Config {
	config := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		// Config doesn't exist yet, return empty config
		return config
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not parse existing config: %v\n", err)
		return &Config{}
	}

	return config
}

// LoadEffectiveConfig loads config with local taking precedence over global
func LoadEffectiveConfig() *Config {
	globalConfig := loadConfig(getGlobalConfigPath())
	localConfig := loadConfig(getLocalConfigPath())

	// Merge: local overrides global
	effective := &Config{
		APIKey:   globalConfig.APIKey,
		Name:     globalConfig.Name,
		Model:    globalConfig.Model,
		Provider: globalConfig.Provider,
	}

	if localConfig.APIKey != "" {
		effective.APIKey = localConfig.APIKey
	}
	if localConfig.Name != "" {
		effective.Name = localConfig.Name
	}
	if localConfig.Model != "" {
		effective.Model = localConfig.Model
	}
	if localConfig.Provider != "" {
		effective.Provider = localConfig.Provider
	}

	return effective
}

func saveConfig(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func showConfigWithScope() {
	if local {
		// Show only local config
		path := getLocalConfigPath()
		fmt.Println("=== Local Configuration ===")
		fmt.Printf("Config file: %s\n\n", path)
		printConfig(loadConfig(path))
	} else if global {
		// Show only global config
		path := getGlobalConfigPath()
		fmt.Println("=== Global Configuration ===")
		fmt.Printf("Config file: %s\n\n", path)
		printConfig(loadConfig(path))
	} else {
		// Show effective config (merged)
		fmt.Println("=== Effective Configuration ===")
		fmt.Printf("Global: %s\n", getGlobalConfigPath())
		fmt.Printf("Local:  %s\n\n", getLocalConfigPath())
		printConfig(LoadEffectiveConfig())
	}
}

func printConfig(config *Config) {
	if config.APIKey != "" {
		// Mask the API key for security
		masked := config.APIKey
		if len(masked) > 8 {
			masked = masked[:4] + "..." + masked[len(masked)-4:]
		}
		fmt.Printf("API Key:  %s\n", masked)
	} else {
		fmt.Println("API Key:  (not set)")
	}

	if config.Name != "" {
		fmt.Printf("Name:     %s\n", config.Name)
	} else {
		fmt.Println("Name:     (not set)")
	}

	if config.Model != "" {
		fmt.Printf("Model:    %s\n", config.Model)
	} else {
		fmt.Println("Model:    (not set)")
	}

	if config.Provider != "" {
		fmt.Printf("Provider: %s\n", config.Provider)
	} else {
		fmt.Println("Provider: (not set)")
	}
}
