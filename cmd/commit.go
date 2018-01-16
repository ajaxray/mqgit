// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ajaxray/mqgit/db"
	"github.com/ajaxray/mqgit/util"
	"github.com/labstack/gommon/bytes"
	"github.com/spf13/cobra"
)

var commitMsg string
var amend bool

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit current state of Database",
	Long: `Commit a snapshot (technically, a compressed backup with metadata) 
			of current schema, data, triggers and stored routines of configured database.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := getDbOrDie()
		conf := getSettings(dbPath)
		cmdStr := getDmpCommand(conf)

		output, err := util.RunCommand(cmdStr)
		if err != nil {
			fmt.Println(err)
		}

		commitID := time.Now().Unix()
		data := makeCommitJSON(commitID, output)

		db.Write(dbPath, "commits", []byte(strconv.FormatInt(commitID, 10)), []byte(data))
		if showVerbose {
			fmt.Printf("Commit saved. Compressed dump size: %s\n", bytes.Format(int64(len(data))))
		}

	},
}

func getDmpCommand(conf map[string]string) string {
	connection := fmt.Sprintf("-h%s -P%s  -u%s -p%s", conf["dbhost"], conf["dbport"], conf["dbuser"], conf["dbpass"])

	options := "--skip-comments --routines --triggers"

	cmdParts := []string{"mysqldump", connection, options, conf["dbname"]}
	return strings.Join(cmdParts, " ")
}

func makeCommitJSON(commitID int64, sql []byte) string {
	commit := db.Commit{
		ID:         commitID,
		Message:    commitMsg,
		GitHash:    util.CurrentGitHash(),
		GitMessage: util.LastGitMessage(),
		Sql:        string(sql),
	}

	commitJSON, _ := commit.ToJSON()
	return commitJSON
}

func init() {

	rootCmd.AddCommand(commitCmd)

	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	commitCmd.Flags().BoolVarP(&amend, "amend", "u", false, "Update latest commit")
	commitCmd.Flags().StringVarP(&commitMsg, "message", "m", "", "Commit message")
}
