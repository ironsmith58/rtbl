/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"rtbl/tables"

	"github.com/spf13/cobra"
)

// rollCmd represents the roll command
var rollCmd = &cobra.Command{
	Use:   "roll",
	Short: "Create random number from DND dice",
	Long: `Examples:
	$ rtbl roll 3D6 3d6+1
	$ rtbl roll 2d4
	$ rtbl roll 4d20*1000
	
	To roll multiple of the same dice put 'number dash' before the dice spec, without spaces.
	
	$ rtbl roll 4-1d4
	
	Titles or other text can be mixed with dice
	
	$ ./rtbl roll init 1d8 dmg 12-1d8
	init 
	5 
	dmg 
	8 8 4 8 4 6 1 2 5 5 1 4 `,
	Run: func(cmd *cobra.Command, args []string) {
		t := tables.NewTable("dummy")
		for j := 0; j < len(args); j++ {
			// is this a multi roll
			idx := strings.Index(args[j], "-")
			var end = 1
			if idx != -1 {
				rep, err := strconv.ParseInt(args[j][:idx], 10, 32)
				if err != nil {
					fmt.Printf("non-numeric repeat: %s\n", args[j])
				} else {
					end = int(rep)
					args[j] = args[j][idx+1:]
				}
			}
			for r := 0; r < end; r++ {
				res, err := tables.BuiltinCall(t, "Dice", args[j])
				if err != nil {
					fmt.Printf("\n%s ", args[j])
				} else {
					fmt.Printf("%s ", res)
				}
			}

		}
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(rollCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rollCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rollCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
