package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/fortifi/go-provision"
)

func main() {
	http.HandleFunc("/success", success)
	http.HandleFunc("/failed", failed)
	http.HandleFunc("/processing", processing)
	http.ListenAndServe("0.0.0.0:9060", nil)
}

func success(w http.ResponseWriter, r *http.Request) {

	//time.Sleep(15 * time.Second)
	response := provisioning.Response{
		BaseTransport: provisioning.BaseTransport{
			Timestamp: provisioning.Time(),
		},
		Log: []provisioning.LogMessage{{
			Type:      "debug",
			Message:   "That worked :)",
			Timestamp: provisioning.Time(),
		}},
		Type:    "success",
		Message: "Setup complete",
	}

	jsonBytes, _ := json.Marshal(response)

	fmt.Fprint(w, string(jsonBytes))
}
func failed(w http.ResponseWriter, r *http.Request) {

	response := provisioning.Response{
		BaseTransport: provisioning.BaseTransport{
			Timestamp: provisioning.Time(),
		},
		Log: []provisioning.LogMessage{{
			Type:      "debug",
			Message:   "That failed :(",
			Timestamp: provisioning.Time(),
		}},
		Type:    "failed",
		Message: "Unable to provision, because its not right",
	}

	jsonBytes, _ := json.Marshal(response)

	fmt.Fprint(w, string(jsonBytes))
}

func processing(w http.ResponseWriter, r *http.Request) {

	response := provisioning.Response{
		BaseTransport: provisioning.BaseTransport{
			Timestamp: provisioning.Time(),
		},
		Log: []provisioning.LogMessage{{
			Type:      "debug",
			Message:   "Queued for provisioning",
			Timestamp: provisioning.Time(),
		}},
		Type:    "processing",
		Message: "Domain queued for provisioning",
	}

	jsonBytes, _ := json.Marshal(response)

	fmt.Fprint(w, string(jsonBytes))
}
