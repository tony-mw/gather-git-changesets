/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// commitMainCmd represents the commitMain command
var commitMainCmd = &cobra.Command{
	Use:   "commitMain",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("GIT_EVENT", "commit_main")
		repo_path, _ := cmd.Flags().GetString("repo-path")
		branch, _ := cmd.Flags().GetString("branch")
		root_dir, _ := cmd.Flags().GetString("root-dir")
		base_branch, _ := cmd.Flags().GetString("base-branch")
		fmt.Println(repo_path, branch, root_dir)
		if len(repo_path) > 0 {
			os.Setenv("LOCAL_REPO_PATH", repo_path)
		}
		if len(branch) > 0 {
			os.Setenv("BRANCH", branch)
		}
		if len(root_dir) > 0 {
			os.Setenv("ROOT_DIR", root_dir)
		}
		if len(base_branch) > 0 {
			os.Setenv("BASE_BRANCH", base_branch)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitMainCmd)
	commitMainCmd.Flags().String("repo-path", "repo-path", "The path to the repo")
	commitMainCmd.Flags().String("branch", "branch", "The branch to check")
	commitMainCmd.Flags().String("base-branch", "base-branch", "The base-branch to check if needed")
	commitMainCmd.Flags().String("root-dir", "root-dir", "Root directory to look for changes")


	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitMainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitMainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
