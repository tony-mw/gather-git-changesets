/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Repo type is: terraform")

		os.Setenv("REPO_TYPE", "terraform")

		repoPath, _ := cmd.Flags().GetString("repo-path")
		branch, _ := cmd.Flags().GetString("branch")
		baseBranch, _ := cmd.Flags().GetString("base-branch")

		if len(repoPath) > 0 {
			os.Setenv("LOCAL_REPO_PATH", repoPath)
		}
		if len(branch) > 0 {
			os.Setenv("BRANCH", branch)
		}
		if len(baseBranch) > 0 {
			os.Setenv("BASE_BRANCH", baseBranch)
			os.Setenv("GIT_EVENT", "pr_main")
		} else {
			os.Setenv("GIT_EVENT", "commit_main")
		}
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	terraformCmd.Flags().String("repo-path", "repo-path", "The path to the repo")
	terraformCmd.Flags().String("branch", "branch", "The branch to check")
	terraformCmd.Flags().String("base-branch", "base-branch", "The base-branch to check if needed")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// terraformCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// terraformCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
