// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ajaxray/mqgit/db"
	u "github.com/ajaxray/mqgit/util"
	"github.com/spf13/cobra"
)

var forceInit bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initiate MqGIT repository",
	Long:  `Initiate MqGIT for in this (and all decendent) directory`,
	Run: func(cmd *cobra.Command, args []string) {
		pwd, _ := os.Getwd()
		dbPath, errFindDb := findRepoDb(pwd)

		initiated := errFindDb == nil

		if initiated && !forceInit {
			fmt.Println("Already initiated! Use \"mqgit init --force\" to clear and re-initiate.")
			os.Exit(0)
		}

		if !initiated || (initiated && forceInit) {

			// Get current directory name, possibly the db name
			currentDir := filepath.Base(pwd)

			dbhost := u.Prompt("Database hostname [127.0.0.1]:", "127.0.0.1")
			dbport := u.Prompt("Database port [3306]:", "3306")
			dbuser := u.Prompt("Database user [root]:", "root")
			dbpass := u.Prompt(fmt.Sprintf("Password of %s at %s:", dbuser, dbhost), "")
			dbname := u.Prompt(fmt.Sprintf("Database name [%s]:", currentDir), currentDir)

			fmt.Println("Initializng MqGIT...")
			verifyInitConfig(dbhost, dbport, dbuser, dbpass, dbname)

			dbPath = filepath.Join(pwd, dbFileName)
			db.Write(dbPath, "settings", []byte("dbhost"), []byte(dbhost))
			db.Write(dbPath, "settings", []byte("dbport"), []byte(dbport))
			db.Write(dbPath, "settings", []byte("dbuser"), []byte(dbuser))
			db.Write(dbPath, "settings", []byte("dbpass"), []byte(dbpass))
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

func verifyInitConfig(dbhost, dbport, dbuser, dbpass, dbname string) {
	_, err := u.RunCommand(fmt.Sprintf("mysql -h%s -P%s -u%s -p\"%s\" -e 'use %s'", dbhost, dbport, dbuser, dbpass, dbname))

	if err != nil {
		fmt.Println("Connection test failed!")
		fmt.Printf("Provided settings: server=%s:%s;uid=%s;pwd=%s;database=%s\n", dbhost, dbport, dbuser, dbpass, dbname)
		fmt.Println("Aborting initialization.")
		os.Exit(0)
	} else {
		fmt.Println("Connection successfull!")
	}
}
