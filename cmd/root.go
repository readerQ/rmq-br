/*
Copyright © 2022 https://github.com/readerQ/rmq-br
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rmq-br",
	Short: "RabbitMQ message dump/send tool",
	Long: `RabbitMQ message dump/send tool. For example:
>rmq-br dump

	.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	serverUrl  string
	dataFolder string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rmq-br.yaml)")
	rootCmd.PersistentFlags().StringVarP(&serverUrl, "url", "u", "ampq://localhost:5672", "rabbitMQ server")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))

	rootCmd.PersistentFlags().StringVarP(&dataFolder, "dir", "d", "data", "folder for messages")
	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
