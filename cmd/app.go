/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Repo type is: app")

		os.Setenv("REPO_TYPE", "app")

		repoPath, _ := cmd.Flags().GetString("repo-path")
		branch, _ := cmd.Flags().GetString("branch")
		baseBranch, _ := cmd.Flags().GetString("base-branch")

		if len(repoPath) > 0 {
			os.Setenv("LOCAL_REPO_PATH", repoPath)
			fmt.Println("Set repoPath to: ", repoPath)
		}
		if len(branch) > 0 {
			os.Setenv("BRANCH", branch)
			fmt.Println("Set branch env var to: ", branch)
		}
		fmt.Println(baseBranch)
		if len(baseBranch) > 0 {
			fmt.Println("Base Branch is: ", baseBranch)
			os.Setenv("BASE_BRANCH", baseBranch)
			os.Setenv("GIT_EVENT", "pr_main")
		} else {
			os.Setenv("GIT_EVENT", "commit_main")
		}
	},
}

func init() {
	rootCmd.AddCommand(appCmd)

	appCmd.Flags().String("repo-path", "", "The path to the repo")
	appCmd.Flags().String("branch", "", "The branch to check")
	appCmd.Flags().String("base-branch", "", "The base-branch to check if needed")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
