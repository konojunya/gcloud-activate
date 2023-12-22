package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gcloud-activate",
	Short: "gcloud-activate is gcloud cli alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		stdout, err := exec.Command("gcloud", "info").Output()
		if err != nil {
			return err
		}

		output := string(stdout)
		line := strings.Split(output, "\n")
		var account string
		var projectId string

		for _, v := range line {
			if strings.Contains(v, "Account:") {
				account = strings.TrimSpace(strings.Split(v, ":")[1])
				account = strings.Trim(account, "[]")
			}
			if strings.Contains(v, "Project:") {
				projectId = strings.TrimSpace(strings.Split(v, ":")[1])
				projectId = strings.Trim(projectId, "[]")
			}
		}

		if account == "" || projectId == "" {
			fmt.Println(output)
		}

		fmt.Printf("Account: %s\n", account)
		fmt.Printf("Project: %s\n", projectId)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
