/*
Copyright © 2022 https://github.com/readerQ/rmq-br
*/
package cmd

import (
	"fmt"
	"log"

	localio "github.com/readerQ/rmq-br/local-io"
	"github.com/readerQ/rmq-br/rabbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	queue string
	max   int
	min   int
	wait  bool
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

	conn := rabbit.RmqConnecton{
		Url: serverUrl,
	}

	wr := localio.NewFileWrited(dataFolder)

	cons := rabbit.NewConsumer(queue, min, min, wait).WithConnection(&conn)
	cons = cons.WithWriter(wr)

	err := cons.Consume(1, 2)
	if err != nil {
		log.Println(err)
	}

	//	cons = cons.WithConnection(&conn)

}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.PersistentFlags().StringVarP(&queue, "queue", "q", "", "rabbitMQ queue name")
	viper.BindPFlag("queue", dumpCmd.PersistentFlags().Lookup("queue"))
	dumpCmd.MarkFlagRequired("queue")

	dumpCmd.Flags().IntVarP(&max, "max", "", 10, "maximum number of messages. 0 = infinity")
	dumpCmd.Flags().IntVarP(&min, "min", "", 1, "minimum number of messages")
	dumpCmd.Flags().BoolVarP(&wait, "wait", "", false, "wait messages till --max 'll be reached")

}
