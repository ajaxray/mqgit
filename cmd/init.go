// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ajaxray/mqgit/db"
	"github.com/spf13/cobra"
)

var forceInit bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		dbPath, errFindDb := findRepoDb(pwd)

		initiated := errFindDb == nil

		if initiated && !forceInit {
			fmt.Println("Already initiated. Use \"mqgit init --force\" to clear and re-initiate.")
			os.Exit(0)
		}

		if !initiated || (initiated && forceInit) {

			// Get current directory name, possibly the db name
			currentDir := filepath.Base(pwd)

			dbhost := prompt("Database hostname [127.0.0.1]:", "127.0.0.1")
			dbuser := prompt("Database user [root]:", "root")
			dbpass := prompt(fmt.Sprintf("Password of %s at %s:", dbuser, dbhost), "")
			dbname := prompt(fmt.Sprintf("Database name [%s]:", currentDir), currentDir)

			fmt.Println("Initializng MqGIT...")
			verifyInitConfig(dbhost, dbuser, dbpass, dbname)

			db.Write(dbPath, "settings", []byte("dbuser"), []byte(dbuser))
			db.Write(dbPath, "settings", []byte("dbpass"), []byte(dbpass))
			db.Write(dbPath, "settings", []byte("dbhost"), []byte(dbhost))
			db.Write(dbPath, "settings", []byte("dbname"), []byte(dbname))
			fmt.Println(`Happy logging! Use "mqgit commit -m message" to log database state`)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	initCmd.Flags().BoolVarP(&forceInit, "force", "f", false, "Clear and re-initiate if required.")
}

func verifyInitConfig(dbhost, dbuser, dbpass, dbname string) {
	cmd := makeShellCmd(fmt.Sprintf("mysql -h%s -u%s -p\"%s\" -e 'use %s'", dbhost, dbuser, dbpass, dbname))
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("Connection test failed!")
		fmt.Printf("Provided settings: server=%s;uid=%s;pwd=%s;database=%s\n", dbhost, dbuser, dbpass, dbname)
		fmt.Println("Aborting initialization.")
		os.Exit(0)
		return
	} else {
		fmt.Println("Connection successfull!")
	}
}
