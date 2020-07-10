package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "opc-predict-server",
	Short: "a websocket server for opc metrics prediction",
	Long:  ``,
}

func init() {
	//cobra.OnInitialize(initConfig) // viper是cobra集成的配置文件读取的库，以后我们会专门说
	//
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobraDemo.yaml)") // 添加`--config`参数。它是全局的~

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	startCommand.Flags().StringVar(&host, "modelDomain", "cnvdxnmodel-service.tbdsversion.com", "the model's host for http request")
	startCommand.Flags().StringVar(&auth, "auth", "e70e0dc4c3f74c7380bc3dfba91bf1d3", "the model's auth code for http request")
	startCommand.Flags().StringVar(&serverIp, "ip", "139.155.92.20", "the ip address for http request")
	startCommand.Flags().IntVar(&timeOut, "timeout", 100, "the timeout (ms) of prediction")
	//kafka
	startCommand.Flags().StringVar(&brokers, "brokers", "localhost:9092", "kafka's brokers list by comma split")
	startCommand.Flags().StringVar(&topic, "topic", "test", "kafka topic name to consume")
	startCommand.Flags().IntVar(&partition, "partition", 0, "kafka partition")
	startCommand.Flags().BoolVar(&saslEnable, "sasl", false, "sasl switcher")
	startCommand.Flags().StringVar(&username, "username", "", "sasl username")
	startCommand.Flags().StringVar(&password, "passwd", "", "sasl password")
	startCommand.Flags().StringVar(&saslMechanism, "mechanism", "PLAIN", "sasl mechanism")
	//websocket server
	startCommand.Flags().IntVar(&wsPort, "wsport", 8080, "websocket server port")

	//stopCommand.Flags().IntVar(&gracefulshutdown, "after", 5, "shutdown after this seconds")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
