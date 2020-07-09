package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"opcdata-predict/cmd/option"
	serverv2 "opcdata-predict/pkg/server/v2"
	"os"
	"os/signal"
)

var (
	// predict
	host     string
	auth     string
	serverIp string
	timeOut  int //predict timeout(ms)
	// kafka
	//command           string
	brokers       string
	topic         string
	partition     int
	saslEnable    bool
	username      string
	password      string
	saslMechanism string

	//ws server
	wsPort           int
	gracefulshutdown int
)
var stopCh chan os.Signal
var startCommand = &cobra.Command{
	Use:   "start",
	Short: "start ws server",
	Long: `start a ws server, retrieve message from kafka and do predict with tiOne model serving. 
then push result to http://<ip>/echo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		fmt.Printf("auth: %v\n", auth)
		fmt.Printf("host: %v\n", host)
		fmt.Printf("serverIp: %v\n", serverIp)

		var opFns []option.OptionFunc

		opFns = append(opFns,
			//predict
			option.SetPredictDomain(host),
			option.SetPredictAuth(auth),
			option.SetServerIp(serverIp),
			option.SetTimeout(timeOut),
			//kafka
			option.SetBrokers(brokers),
			option.SetTopic(topic),
			option.SetPartition(partition),
			//sasl default
			option.SetSaslEnable(saslEnable),
		)

		if saslEnable {
			opFns = append(opFns,
				option.SetSaslUser(username),
				option.SetSaslPassword(password),
				option.SetSaslMechanism(saslMechanism),
			)
		}

		opts := option.NewOptions(opFns...)

		fmt.Printf("start test server\n")

		//wsServer := server.NewWsServer(wsPort, opts)
		wsServer := serverv2.NewWsServer(wsPort, opts)
		wsServer.Boot()

		stopCh = make(chan os.Signal, 1)
		signal.Notify(stopCh, os.Interrupt)
		<-stopCh
		fmt.Printf("stopped test server\n")

	},
}
var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "stop ws server",
	Long:  `stop a ws server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop called")
		//TODO
		//close(stopCh)
	},
}

func init() {
	rootCmd.AddCommand(startCommand)
	//rootCmd.AddCommand(stopCommand)
}
