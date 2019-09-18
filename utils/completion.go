package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash/powershell/zsh]",
	Short: "Installs bash/powershell/zsh completion commands",
	Long: `Installs bash/powershell/zsh completion commands.
Please note that the completions may be incomplete or broken, see https://github.com/spf13/cobra`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"bash", "powershell", "zsh"},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "bash" {
			if viper.GetBool("utils.completion.source") {
				cmd.Root().GenBashCompletion(os.Stdout)
			} else {
				file, err := getConfigPath("completion.sh")
				if err != nil {
					fmt.Printf("Error saving completion file: %s\n", err)
					os.Exit(1)
				}
				err = cmd.Root().GenBashCompletionFile(file)
				if err != nil {
					fmt.Printf("Error saving completion file: %s\n", err)
					os.Exit(1)
				}

				// Get the value of $HOME
				home, err := os.UserHomeDir()
				if err != nil {
					fmt.Printf("Failed to get $HOME location: %s\n", err)
					os.Exit(1)
				}
				bashrc := filepath.Join(home, ".bashrc")

				absFile, err := filepath.Abs(file)
				if err != nil {
					fmt.Printf("Failed to resolve path: %s\n", err)
					os.Exit(1)
				}
				command := ". " + absFile
				commandFound := false
				// Check for existing text in bashrc
				data, err := ioutil.ReadFile(bashrc)
				if err == nil {
					commandFound = strings.Contains(string(data), command)
				}
				if !commandFound {
					// Append to bashrc
					os.MkdirAll(filepath.Dir(bashrc), os.ModePerm)
					f, err := os.OpenFile(bashrc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Printf("Failed to open bashrc: %s\n", err)
						os.Exit(1)
					}
					_, err = f.WriteString("\r\n" + command + "\r\n")
					if err != nil {
						fmt.Printf("Failed to write to bashrc: %s\n", err)
						f.Close()
						os.Exit(1)
					}
					f.Close()
				}
				fmt.Println("Completions installed! Restart your shell to load them.")
			}
		} else if args[0] == "powershell" {
			if viper.GetBool("utils.completion.source") {
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			} else {
				file, err := getConfigPath("completion.ps1")
				if err != nil {
					fmt.Printf("Error saving completion file: %s\n", err)
					os.Exit(1)
				}
				err = cmd.Root().GenPowerShellCompletionFile(file)
				if err != nil {
					fmt.Printf("Error saving completion file: %s\n", err)
					os.Exit(1)
				}

				// Get the value of $PROFILE (not an environment variable!!)
				profile, err := getPowershellProfile()
				if err != nil {
					fmt.Printf("Failed to get profile location: %s\n", err)
					os.Exit(1)
				}

				absFile, err := filepath.Abs(file)
				if err != nil {
					fmt.Printf("Failed to resolve path: %s\n", err)
					os.Exit(1)
				}
				command := ". " + absFile
				commandFound := false
				// Check for existing text in profile
				data, err := ioutil.ReadFile(profile)
				if err == nil {
					commandFound = strings.Contains(string(data), command)
				}
				if !commandFound {
					// Append to profile
					os.MkdirAll(filepath.Dir(profile), os.ModePerm)
					f, err := os.OpenFile(profile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Printf("Failed to open profile: %s\n", err)
						os.Exit(1)
					}
					_, err = f.WriteString("\r\n" + command + "\r\n")
					if err != nil {
						fmt.Printf("Failed to write to profile: %s\n", err)
						f.Close()
						os.Exit(1)
					}
					f.Close()
				}
				fmt.Println("Completions installed! Restart your shell to load them.")
			}
		} else if args[0] == "zsh" {
			if viper.GetBool("utils.completion.source") {
				cmd.Root().GenZshCompletion(os.Stdout)
			} else {
				file, err := getConfigPath("completion.zsh")
				if err != nil {
					fmt.Printf("Error saving completion file: %s\n", err)
					os.Exit(1)
				}
				err = cmd.Root().GenZshCompletionFile(file)
				if err != nil {
					fmt.Printf("Error saving completion file: %s\n", err)
					os.Exit(1)
				}

				// If someone knows how to do this automagically, please PR it!!
				fmt.Println("Completions saved to " + file)
				fmt.Println("You need to put this file in your $fpath manually!")
			}
		}
	},
}

func getConfigPath(fileName string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "packwiz")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fileName), nil
}

func getPowershellProfile() (string, error) {
	ps, err := exec.LookPath("powershell.exe")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(ps, "-NoProfile", "-NonInteractive", "-Command", "$PROFILE")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func init() {
	utilsCmd.AddCommand(completionCmd)

	completionCmd.Flags().Bool("source", false, "Output the source of the commands to be installed, rather than installing them")
	viper.BindPFlag("utils.completion.source", completionCmd.Flags().Lookup("source"))
}
