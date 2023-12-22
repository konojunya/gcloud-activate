package cmd

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Selection of configurable accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if list flag is set
		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			return err
		}
		if list {
			return printAccountList()
		}

		// default action
		return selectAccount()
	},
}

func extractEmails(stdout []byte) []string {
	output := string(stdout)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(output, -1)

	return emails
}

func execGcloudAuthList() ([]byte, error) {
	list, err := exec.Command("gcloud", "auth", "list").Output()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func selectAccount() error {
	list, err := execGcloudAuthList()
	if err != nil {
		return err
	}

	emails := extractEmails(list)

	f, err := fzf.New()
	if err != nil {
		return err
	}

	index, err := f.Find(emails, func(i int) string {
		return emails[i]
	})
	if err != nil {
		return err
	}

	selectedAccount := emails[index[0]]
	err = exec.Command("gcloud", "config", "set", "account", selectedAccount).Run()
	if err != nil {
		return err
	}

	fmt.Printf("Activate account: %s\n", string(selectedAccount))

	return nil
}

func printAccountList() error {
	list, err := execGcloudAuthList()
	if err != nil {
		return err
	}

	emails := extractEmails(list)

	for _, email := range emails {
		fmt.Println(email)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(accountCmd)

	// print account list
	accountCmd.Flags().BoolP("list", "l", false, "List all accounts")
}
