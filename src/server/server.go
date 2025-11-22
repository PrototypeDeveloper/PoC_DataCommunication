package main

import (
	"PoC_DataCommunication/src/models"
	"io"
	"log"
	"net/http"

	"github.com/fxamacker/cbor/v2"
	"github.com/quic-go/quic-go/http3"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/cbor", cborHandler)

	addr := "localhost:1234"
	if err := http3.ListenAndServeTLS(
		addr,
		"../../cert/cert.pem",
		"../../cert/key.pem",
		mux,
	); err != nil {
		log.Fatal(err)
	}
}

func cborHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body error", http.StatusBadRequest)
		return
	}

	var reqMsg models.Signal
	if err := cbor.Unmarshal(body, &reqMsg); err != nil {
		http.Error(w, "invalid CBOR", http.StatusBadRequest)
		return
	}
	log.Printf("received: %+v", reqMsg)

	resMsg := models.Signal{Type: "Response", Message: "Connection"}
	respBytes, err := cbor.Marshal(resMsg)
	if err != nil {
		http.Error(w, "encode CBOR error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/cbor")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(respBytes); err != nil {
		log.Printf("write resp error: %v", err)
	}
}
