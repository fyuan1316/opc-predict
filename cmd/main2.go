package cmd

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"opcdata-predict-client/cmd/option"
	"opcdata-predict-client/pkg/server"
	"os"
)

//var (
//	// predict
//	host     string
//	auth     string
//	serverIp string
//	// kafka
//	//command           string
//	brokers           string
//	topic             string
//	partition         int
//	saslEnable        bool
//	username          string
//	password          string
//	saslSASLMechanism string
//)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		//{
		//	Name:     "add",
		//	Aliases:  []string{"a"},
		//	Usage:    "add a task to the list",
		//	Category: "Add",
		//	Action: func(c *cli.Context) error {
		//		fmt.Println("added task: ", c.Args().First())
		//		return nil
		//	},
		//},
		{
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "options for task templates",
			Category: "server",
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "start server",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "stop",
					Usage: "stop server",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Name = "opc-predict-server"
	app.Usage = "application usage"
	app.Description = "application description" // 描述
	app.Version = "0.0.1"                       // 版本

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	parseArgsFromCli()
	fmt.Println(host)
	webSocketServerStart()

}

//var lock sync.Mutex
//var consumeState bool
func webSocketServerStart() {
	wsmgr := server.NewWebSocketManager(option.Options{})
	http.HandleFunc("/echo", wsmgr.WsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	fmt.Println("Listening :8080")
	http.ListenAndServe(":8080", nil)
}

func parseArgsFromCli() {
	flag.StringVar(&host, "host", "", "model domain name")
	flag.StringVar(&auth, "auth", "", "model auth code")
	flag.StringVar(&serverIp, "server", "", "server ip that will be request")

	flag.StringVar(&brokers, "brokers", "localhost:9092", "Common separated kafka hosts")
	flag.StringVar(&topic, "topic", "test", "Kafka topic")
	flag.IntVar(&partition, "partition", 0, "Kafka topic partition")

	flag.BoolVar(&saslEnable, "sasl", false, "SASL enable")
	flag.StringVar(&username, "username", "", "SASL Username")
	flag.StringVar(&password, "password", "", "SASL Password")
	flag.StringVar(&saslMechanism, "mechanism", sarama.SASLTypePlaintext, "SASL mechanism")

	flag.Parse()

	//opts := &option.Options{}
	//return nil
}
