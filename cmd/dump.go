/*
Copyright © 2022 https://github.com/readerQ/rmq-br
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	queue string
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:              "dump",
	Short:            "dumps all messages",
	Long:             `dumps all messages`,
	TraverseChildren: true,
	Example:          "rmq-br dump -q test", // TODO

	Run: execDump,
}

func execDump(cmd *cobra.Command, args []string) {
	fmt.Println("dump called " + serverUrl)
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.PersistentFlags().StringVarP(&queue, "queue", "q", "", "rabbitMQ queue name")
	viper.BindPFlag("queue", dumpCmd.PersistentFlags().Lookup("queue"))
	dumpCmd.MarkFlagRequired("queue")

}
