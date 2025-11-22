package main

import (
	"PoC_DataCommunication/src/models"
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/quic-go/quic-go/http3"
)

func main() {

	url := "https://localhost:1234/cbor"
	transport := &http3.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS13,
		},
	}
	defer transport.Close()

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	reqMsg := models.Signal{Type: "Request", Message: "Hellow World"}
	reqBytes, err := cbor.Marshal(reqMsg)
	if err != nil {
		log.Fatalf("cbor.Marshal: %v", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBytes))
	if err != nil {
		log.Fatalf("NewRequest: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/cbor")

	res, err := client.Do(httpReq)
	if err != nil {
		log.Fatalf("client.Do: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}

	var resMsg models.Signal
	if err := cbor.Unmarshal(body, &resMsg); err != nil {
		log.Fatalf("cbor.Unmarshal: %v", err)
	}
	log.Printf("resMsg : %+v", resMsg)
}
