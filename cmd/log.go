// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/ajaxray/mqgit/db"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		dbPath, errFindDb := findRepoDb(pwd)

		if errFindDb == nil {
			dbhost := db.Read(dbPath, "settings", []byte("dbhost"))
			dbuser := db.Read(dbPath, "settings", []byte("dbuser"))
			dbpass := db.Read(dbPath, "settings", []byte("dbpass"))
			dbname := db.Read(dbPath, "settings", []byte("dbname"))

			fmt.Printf("Provided settings: server=%s;uid=%s;pwd=%s;database=%s\n", dbhost, dbuser, dbpass, dbname)
		} else {
			fmt.Errorf("No db found")
		}

		// db.Map(getDbPath(), "settings", func(k, v []byte) error {
		// 	fmt.Printf("key=%s, value=%s\n", k, v)
		// 	return nil
		// })
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
}
