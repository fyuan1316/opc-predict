package predict

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var urlFormat = "http://%s/v1/models/m:predict"

func Post(server, host, auth string, data []byte) (string, error) {

	fmt.Println("server:", server)
	fmt.Println("host:", host)
	fmt.Println("auth:", auth)
	predictUrl := fmt.Sprintf(urlFormat, server)

	req, err := http.NewRequest("POST", predictUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Host = host
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Header Host: [%s]\n", host)
	log.Printf("Header Authorization: [%s]\n", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("response Status: [%s]\n", resp.Status)
	log.Printf("response Headers: [%s]\n", resp.Header)
	log.Printf("response Body: [%s]\n", string(body))
	return string(body), err
}

type Payload struct {
	Instances []Instances `json:"instances"`
}
type Instances struct {
	X1 float64 `json:"x1"`
	X2 float64 `json:"x2"`
	X3 float64 `json:"x3"`
	X4 float64 `json:"x4"`
}

func Post2(server, host, auth string, payloadBytes []byte) (string, error) {
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://139.155.92.20/v1/models/m:predict", body)
	if err != nil {
		// handle err
	}
	req.Host = "cnvdxnmodel-service.tbdsversion.com"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "e70e0dc4c3f74c7380bc3dfba91bf1d3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))
	return string(respBody), err
}
