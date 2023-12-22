package cmd

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/koki-develop/go-fzf"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Selection of configurable projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if list flag is set
		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			return err
		}
		if list {
			return printProjectList()
		}

		// default action
		return selectProject()
	},
}

func extractProjectIds(stdout []byte) []string {
	lines := strings.Split(string(stdout), "\n")
	regex := regexp.MustCompile(`^\s*([^\s]+)\s+`)

	var projectIDs []string
	for i, line := range lines {
		if i == 0 {
			// skip header
			continue
		}

		matches := regex.FindStringSubmatch(line)
		if len(matches) >= 2 {
			projectID := matches[1]
			projectIDs = append(projectIDs, projectID)
		}
	}

	return projectIDs
}

func execGcloudProjectsList() ([]byte, error) {
	list, err := exec.Command("gcloud", "projects", "list").Output()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func selectProject() error {
	list, err := execGcloudProjectsList()
	if err != nil {
		return err
	}

	projectIDs := extractProjectIds(list)

	f, err := fzf.New()
	if err != nil {
		return err
	}

	index, err := f.Find(projectIDs, func(i int) string {
		return projectIDs[i]
	})
	if err != nil {
		return err
	}

	selectedProject := projectIDs[index[0]]
	err = exec.Command("gcloud", "config", "set", "project", selectedProject).Run()
	if err != nil {
		return err
	}

	fmt.Printf("Activate project: %s\n", string(selectedProject))

	return nil
}

func printProjectList() error {
	list, err := execGcloudProjectsList()
	if err != nil {
		return err
	}

	projectIDs := extractProjectIds(list)

	for _, projectID := range projectIDs {
		fmt.Println(projectID)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// print projects list
	projectCmd.Flags().BoolP("list", "l", false, "List all projects")
}
