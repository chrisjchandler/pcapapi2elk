package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const esEndpoint = "http://your-elasticsearch-endpoint:9200/packet-captures/pcap"

func main() {
	http.HandleFunc("/upload", handleUpload)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().UTC().Format(time.RFC3339)
	pcap := map[string]interface{}{
		"timestamp": timestamp,
		"data":      body,
	}

	b, err := json.Marshal(pcap)
	if err != nil {
		http.Error(w, "failed to encode the packet capture", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(esEndpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		http.Error(w, "failed to store the packet capture in Elasticsearch", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		http.Error(w, "failed to store the packet capture in Elasticsearch: "+string(body), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

