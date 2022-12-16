/*
Copyright © 2022 https://github.com/readerQ/rmq-br
*/
package cmd

import (
	"fmt"
	"os"

	localio "github.com/readerQ/rmq-br/local-io"
	"github.com/readerQ/rmq-br/rabbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	exchange             string
	routingKey           string
	mandatory, immediate bool
	contentType          string
	mask                 string
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:     "send",
	Short:   "reades files from --dir, send to exchange",
	Long:    ``,
	Example: "rmq-br -x / -k \"test\"",
	Run:     execSend,
}

func execSend(cmd *cobra.Command, args []string) {
	fmt.Println("send called")

	conn := rabbit.RmqConnecton{
		Url: serverUrl,
	}

	rdr := localio.NewMessageReader(dataFolder, mask, contentType)

	cons := rabbit.NewSender(exchange, routingKey, mandatory, immediate, quite).WithConnection(&conn)
	cons = cons.WithReader(rdr)

	err := cons.Send()

	if err != nil {
		os.Exit(1)
	}

}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.PersistentFlags().StringVarP(&exchange, "exchange", "x", "", "rabbitMQ exchange name")
	viper.BindPFlag("exchange", sendCmd.PersistentFlags().Lookup("exchange"))
	sendCmd.MarkFlagRequired("exchange")

	sendCmd.PersistentFlags().StringVarP(&routingKey, "routingKey", "k", "", "routing key")
	viper.BindPFlag("routingKey", sendCmd.PersistentFlags().Lookup("routingKey"))
	sendCmd.MarkFlagRequired("routingKey")

	sendCmd.PersistentFlags().StringVarP(&dataFolder, "dir", "d", "data", "folder with messages")
	viper.BindPFlag("dir", sendCmd.PersistentFlags().Lookup("dir"))

	sendCmd.PersistentFlags().StringVarP(&mask, "mask", "m", "*.*", "file search mask")
	viper.BindPFlag("mask", sendCmd.PersistentFlags().Lookup("mask"))

	sendCmd.PersistentFlags().StringVarP(&contentType, "contentType", "", "application/json", "ContentType header")
	viper.BindPFlag("contentType", sendCmd.PersistentFlags().Lookup("contentType"))

	sendCmd.Flags().BoolVarP(&quite, "quite", "", false, "less mesages")
}
