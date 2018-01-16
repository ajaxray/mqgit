// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajaxray/mqgit/db"
	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"
const dbFileName = ".mqgit.db"

var cfgFile string
var showVersion bool
var showVerbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mqgit",
	Short: "VCS for MySQL Database",
	Long:  `VCS for MySQL Database - git style commit, log, checkout of the schema+data`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			fmt.Printf("mqgit version %s\n", VERSION)
		} else {
			fmt.Println("mqgit is the VCS for MySQL Database.\nType \"mqgit help\" for usages instructions.")
		}
	},
}

func findRepoDb(dir string) (string, error) {
	if "" == dir {
		return "", errors.New("no db file found in current or any of parent directories. Use \"mqgit init\"")
	}

	dbPath := filepath.Join(dir, dbFileName)
	if _, err := os.Stat(dbPath); err == nil {
		return dbPath, nil
	} else {
		parent, _ := filepath.Split(dir)
		return findRepoDb(strings.TrimSuffix(parent, "/"))
	}
}

func getDbOrDie() string {
	if pwd, errWd := os.Getwd(); errWd == nil {
		if dbPath, errDbPath := findRepoDb(pwd); errDbPath == nil {
			return dbPath
		} else {
			log.Fatal(errDbPath)
		}
	} else {
		log.Fatal(errWd)
	}
	os.Exit(0)
	return "Died!"
}

func getSettings(dbPath string) map[string]string {
	return db.Dictionary(dbPath, "settings")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "Display application version")
	rootCmd.PersistentFlags().BoolVarP(&showVerbose, "verbose", "v", false, "Display more insights.")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mvvcs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
