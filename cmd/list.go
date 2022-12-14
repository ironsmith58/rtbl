/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"rtbl/tables"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available tables",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		showAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
			return
		}
		env_root := os.Getenv("RTBL_ROOT")
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(root) == 0 {
			root = env_root
		}
		if len(root) == 0 { // even after checking getopt and the flag
			root = "."
		}

		var paths []string

		if len(args) > 0 {
			for _, name := range args {
				var path string
				if root != "" {
					path = root + "/Tables/" + name
				} else {
					path = name
				}
				tpaths, err := tables.FindTables(path)
				if err != nil {
					fmt.Println(err)
				}
				paths = append(paths, tpaths...)
			}
		} else {
			paths, err = tables.FindTables(root + "/Tables")
			if err != nil {
				fmt.Println(err)
			}
		}
		tabs := tables.NewTablesByCategory(paths)
		tables.PrintPaths(tabs, showAll)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolP("all", "a", false, "Show hidden catagores/tables also")
}
