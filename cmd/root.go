// Copyright © 2018 Anis Uddin Ahmad <anis.programmer@gmail.com>

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func getDbPath() string {
	if pwd, err := os.Getwd(); err == nil {
		return filepath.Join(pwd, dbFileName)
	} else {
		log.Fatal(err)
		return ""
	}
}

func prompt(question string, defaultAnswer string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	if text == "\n" {
		return defaultAnswer
	}

	return strings.TrimSuffix(text, "\n")
}

// Make a exec.Cmd from string
// Then you can use it for example: output, err := cmd.Output()
func makeShellCmd(cmd string) *exec.Cmd {
	return exec.Command("sh", "-c", cmd)
}

func findRepoDb(dir string) (string, error) {
	if "" == dir {
		return "", errors.New("no db file found in current or any of parent directories")
	}

	dbPath := filepath.Join(dir, dbFileName)
	if _, err := os.Stat(dbPath); err == nil {
		return dbPath, nil
	} else {
		parent, _ := filepath.Split(dir)
		return findRepoDb(strings.TrimSuffix(parent, "/"))
	}
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
