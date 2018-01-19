// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"fmt"
	"time"

	"github.com/ajaxray/mqgit/db"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var msgLen uint

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show list of available commits",
	Long: `Displays list of available commits. 
	       This list can be filtered and styled in various ways using flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := getDbOrDie()
		table := uitable.New()
		table.Separator = "  "
		table.AddRow("CommitID", "CreatedAt", "Commit Message")
		table.AddRow("-----------", "-------------------------", "------------------------------")
		table.MaxColWidth = msgLen

		db.Map(dbPath, "commits", func(k, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			var commit db.Commit
			commit.FromJSON(string(v))
			at := time.Unix(commit.ID, 0)
			table.AddRow(commit.ID, at.Format("Mon Jan _2 15:04:05 2006"), commit.Message)

			return nil
		})
		fmt.Println(table)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	logCmd.Flags().UintVarP(&msgLen, "msg-len", "l", 40, "Length of message to show in default view")
}
