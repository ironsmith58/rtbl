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

		env_root := os.Getenv("RTBL_TABLE_ROOT")
		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(root) == 0 {
			root = env_root
		}
		if len(root) > 0 && root[0] != '/' {
			root = env_root + "/" + root
		}
		if len(args) > 0 {
			for _, name := range args {
				reg, err := tables.List(root + "/" + name)
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, tab := range reg {
					fmt.Println(tab.Name)
				}
			}
		} else {
			reg, err := tables.List(root)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, tab := range reg {
				fmt.Println(tab.Name)
			}
		}
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
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
