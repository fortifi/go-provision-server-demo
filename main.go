package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/fortifi/go-provision"
	"io/ioutil"
	"bytes"
	"log"
)

func main() {
	http.HandleFunc("/success", success)
	http.HandleFunc("/failed", failed)
	http.HandleFunc("/processing", processing)
	http.ListenAndServe("0.0.0.0:9060", nil)
}

func getRequest(w http.ResponseWriter, r *http.Request) *provisioning.Request {
	request := &provisioning.Request{}
	inBytes, byteErr := ioutil.ReadAll(r.Body)
	log.Print(string(inBytes))
	if byteErr != nil {
		log.Print(byteErr)
		http.Error(w, byteErr.Error(), 500)
		return nil
	}
	jsErr := json.Unmarshal(inBytes, request)
	if jsErr != nil {
		log.Print(jsErr)
		http.Error(w, jsErr.Error(), 500)
		return nil
	}
	return request
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

	request := getRequest(w, r)

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
		Message: "Service queued for provisioning",
	}

	jsonBytes, jsErr := json.Marshal(response)
	log.Print(jsErr)

	fmt.Fprint(w, string(jsonBytes))

	if request != nil {
		go postUpdate(request)
	}
}

func postUpdate(r *provisioning.Request) {

	response := provisioning.Response{}
	response.TransportKey = r.TransportKey
	response.SubscriptionFid = r.SubscriptionFid
	response.Timestamp = provisioning.Time()
	response.Message = "All finished"
	response.CustomerFid = r.CustomerFid
	response.Type = provisioning.RESPONSE_SUCCESS
	response.Log = []provisioning.LogMessage{{
		Type:      "debug",
		Message:   "Pulled from provisioning queue",
		Timestamp: provisioning.Time(),
	}, {
		Type:      "info",
		Message:   "Setting up service",
		Timestamp: provisioning.Time(),
	}, {
		Type:      "info",
		Message:   "Service setup with ID 34i634hr",
		Timestamp: provisioning.Time(),
	}}
	response.Properties = []provisioning.TransportProperty{{
		Key:         "serviceId",
		Type:        provisioning.TRANSPROP_TYPE_STRING,
		StringValue: "34i634hr",
	}, {
		Key:       "provisioningComplete",
		Type:      provisioning.TRANSPROP_TYPE_FLAG,
		FlagValue: true,
	}}

	jsonB, _ := json.Marshal(response)
	payload := bytes.NewBuffer(jsonB)
	resp, err := http.Post(r.UpdateUrl, "application/json", payload)
	if err != nil {
		log.Print(err)
	}
	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	} else {
		log.Print(string(respB))
	}
}
