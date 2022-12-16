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
	exchange   string
	routingKey string
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:     "send",
	Short:   "reades files from --dir, send to exchange",
	Long:    ``,
	Example: "rmq-br -x / -k \"test\"",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.PersistentFlags().StringVarP(&queue, "exchange", "x", "", "rabbitMQ exchange name")
	viper.BindPFlag("exchange", dumpCmd.PersistentFlags().Lookup("exchange"))
	sendCmd.MarkFlagRequired("exchange")

	sendCmd.PersistentFlags().StringVarP(&queue, "routingKey", "k", "", "routing key")
	viper.BindPFlag("routingKey", dumpCmd.PersistentFlags().Lookup("routingKey"))
	sendCmd.MarkFlagRequired("routingKey")

}
