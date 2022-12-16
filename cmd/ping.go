/*
Copyright © 2022 https://github.com/readerQ/rmq-br
*/
package cmd

import (
	"log"
	"os"

	"github.com/readerQ/rmq-br/rabbit"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "check connection",
	Long: `create connection & channel (without queue). 
claims that user, password, host, post and vhost are correct
retrun 1 on error`,
	Run: func(cmd *cobra.Command, args []string) {

		conn := rabbit.RmqConnecton{
			Url: serverUrl,
		}

		cons := rabbit.NewConsumer(queue, max, wait).WithConnection(&conn)

		err := cons.Ping()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		} else {
			log.Println("seems to be ok")
		}

	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
