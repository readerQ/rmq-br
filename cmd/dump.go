/*
Copyright © 2022 https://github.com/readerQ/rmq-br
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	localio "github.com/readerQ/rmq-br/local-io"
	"github.com/readerQ/rmq-br/rabbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	queue string
	max   int
	wait  bool
	ext   string
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
	wr.Extension = ext

	cons := rabbit.NewConsumer(queue, max, wait).WithConnection(&conn)
	cons = cons.WithWriter(wr)

	err := cons.Consume(1, 2)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	//	cons = cons.WithConnection(&conn)

}

func init() {
	rootCmd.AddCommand(dumpCmd)

	dumpCmd.PersistentFlags().StringVarP(&queue, "queue", "q", "", "rabbitMQ queue name")
	viper.BindPFlag("queue", dumpCmd.PersistentFlags().Lookup("queue"))
	dumpCmd.MarkFlagRequired("queue")

	dumpCmd.Flags().IntVarP(&max, "max", "", 0, "maximum number of messages. 0 = infinity")

	dumpCmd.Flags().BoolVarP(&wait, "wait", "", false, "wait messages till --max 'll be reached")

	dumpCmd.Flags().StringVarP(&ext, "ext", "", "json", "file extentions")

}
