// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ajaxray/mqgit/db"
	"github.com/ajaxray/mqgit/util"
	"github.com/spf13/cobra"
)

var cleanUp bool

// checkoutCmd represents the commit command
var checkoutCmd = &cobra.Command{
	Use:   "checkout CommitID",
	Short: "Checkout database snapshot by CommitID ",
	Long: `Checkout a snapshot (technically, restore a backup) by commit ID.
			It will restore database schema, data, triggers and stored routines 
			of configured database at comitted state.`,
	Aliases: []string{"restore"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("accepts (only) CommitID as argument")
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			return errors.New("non-numeric CommitID provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		CommitID := args[0]
		dbPath := getDbOrDie()

		commitData := db.Read(dbPath, "commits", []byte(CommitID))

		if commitData == "" {
			fmt.Println("Provided CommitID was not found! Use \"mqgit log\" to get the valid CommitID")
			return
		}

		var commit db.Commit
		commit.FromJSON(commitData)

		conf := getSettings(dbPath)
		sqlFileName := writeToFile(commit.Sql)
		defer os.Remove(sqlFileName)

		cmdStr := getRestoreCommand(conf, sqlFileName)
		//fmt.Println(cmdStr)

		_, err := util.RunCommand(cmdStr)
		if err != nil {
			fmt.Println(err)
			fmt.Errorf("restoring failed for Commit %s \n", CommitID)
		} else {
			if showVerbose {
				fmt.Println("Database state restored to " + CommitID)
			}
		}

	},
}

func getRestoreCommand(conf map[string]string, sqlFileName string) string {
	// mysqlimport -u [uname] -p[pass] --local [dbname] [backupfile.sql]
	connection := fmt.Sprintf("-h%s -P%s  -u%s -p%s", conf["dbhost"], conf["dbport"], conf["dbuser"], conf["dbpass"])

	cmdParts := []string{"mysql", connection, conf["dbname"], "<", sqlFileName}
	return strings.Join(cmdParts, " ")
}

func writeToFile(content string) string {
	tmpfile, err := ioutil.TempFile("", "mqgit_")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.WriteString(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func init() {

	rootCmd.AddCommand(checkoutCmd)

	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// @TODO : implement or remove
	checkoutCmd.Flags().BoolVarP(&cleanUp, "cleanup", "c", false, "Re-create database before restoring (requires permission)")
}
