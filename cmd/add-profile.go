package cmd

import (
	"bufio"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addProfileCmd represents the add-profile command
var addProfileCmd = &cobra.Command{
	Use:   "add-profile",
	Short: "Add a new Git profile",
	Long: `Add a new Git profile with details like username, email, and SSH private key.

Example:
	teagit add-profile
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := addProfile(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Git profile added successfully!")
	},
}

func init() {
	rootCmd.AddCommand(addProfileCmd)
}

func validateSSHKey(keyPath string) error {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH private key: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "OPENSSH PRIVATE KEY" {
		return errors.New("invalid SSH private key format")
	}

	return nil
}

func generateSSHKey(email string) (string, string, error) {
	tinyID := fmt.Sprintf("%d", time.Now().UnixNano())[:8]
	privateKeyPath := filepath.Join(os.Getenv("HOME"), ".ssh", "teagit_"+tinyID+"_ed25519")
	publicKeyPath := privateKeyPath + ".pub"

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", privateKeyPath, "-N", "")
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to generate SSH key: %v", err)
	}

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read generated SSH private key: %v", err)
	}
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read generated SSH public key: %v", err)
	}

	return string(privateKey), string(publicKey), nil
}

func sanitizeEmail(email string) string {
	return strings.ReplaceAll(email, ".", "_")
}

func addProfile() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Git username: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter Git email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	sanitizedEmail := sanitizeEmail(email)

	fmt.Print("Do you have an SSH private key? (y/n): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var sshKey, publicKey string
	var err error

	if choice == "y" {
		fmt.Print("Enter SSH private key path: ")
		sshKey, _ = reader.ReadString('\n')
		sshKey = strings.TrimSpace(sshKey)

		if err := validateSSHKey(sshKey); err != nil {
			return fmt.Errorf("invalid SSH private key: %v", err)
		}
	} else {
		fmt.Println("Generating a new SSH key...")
		sshKey, publicKey, err = generateSSHKey(email)
		if err != nil {
			return err
		}
		fmt.Println("New SSH key generated. Add the following public key to GitHub:")
		fmt.Println(publicKey)
	}

	fmt.Print("Enter GPG signing key (optional, press Enter to skip): ")
	signingKey, _ := reader.ReadString('\n')
	signingKey = strings.TrimSpace(signingKey)

	configDir := filepath.Join(os.Getenv("HOME"), ".teagit")
	configFile := filepath.Join(configDir, "config.yaml")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("teagit is not initialized. Run 'teagit init' first")
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	profiles := viper.GetStringMap("profiles")
	if _, exists := profiles[sanitizedEmail]; exists {
		return fmt.Errorf("profile with email '%s' already exists", email)
	}

	profile := map[string]string{
		"name":        name,
		"email":       email,
		"ssh_key":     sshKey,
		"ssh_pub_key": publicKey,
		"signing_key": signingKey,
	}
	viper.Set(fmt.Sprintf("profiles.%s", sanitizedEmail), profile)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}

