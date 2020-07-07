package main

import (
	"flag"
	"fmt"
	"opcdata-predict-client/pkg/predict"
	"os"
)

var (
	host     string
	auth     string
	serverIp string
)

func main() {
	flag.StringVar(&host, "host", "", "model domain name")
	flag.StringVar(&auth, "auth", "", "model auth code")
	flag.StringVar(&serverIp, "server", "", "server ip that will be request")
	flag.Parse()

	if host == "" || auth == "" || serverIp == "" {
		fmt.Println("--host or --auth or --server is missing")
		os.Exit(1)
	}
	/*
		url := "http://restapi3.apiary.io/notes"
		fmt.Println("URL:>", url)

		var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	*/
	data := `{"instances": [{"x1":6.2, "x2":2.2, "x3":1.1, "x4":1.2}]}`
	r, err := predict.Post2(serverIp, host, auth, []byte(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("response:\n%v", r)

}
