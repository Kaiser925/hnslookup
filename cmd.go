package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hnslookup name [, name2]",
	Short: "hnslookup query DNS record from cloudflare.",
	Long:  `A DNS query tool use DoH.`,
	Args:  argsHandle,
	RunE:  handleCmd,
}

var recordType string

func init() {
	rootCmd.PersistentFlags().StringVar(&recordType, "type", "", "The type of DNS record, default A")
}

func argsHandle(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("need arguments for domain name")
	}
	return nil
}

func handleCmd(cmd *cobra.Command, args []string) error {
	recordType := cmd.Flag("type").Value.String()

	questions := make([]*Question, len(args))
	for i, name := range args {
		question := &Question{
			Type: recordType,
			Name: name,
		}
		questions[i] = question
	}

	return HandleQuestions(questions)
}

func ExecuteCmd() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
