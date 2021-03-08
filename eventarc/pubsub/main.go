// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START eventarc_pubsub_handler]

// Sample pubsub is a Cloud Run service which handles Pub/Sub messages.
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pubsub "github.com/googleapis/google-cloudevents-go/cloud/pubsub/v1"
)

// HelloEventsPubSub receives and processes a Pub/Sub push message.
func HelloEventsPubSub(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad HTTP Request", http.StatusBadRequest)
		log.Printf("Bad HTTP Request: %v", http.StatusBadRequest)
		return
	}
	e, err := pubsub.UnmarshalMessagePublishedData(body)
	if err != nil {
		http.Error(w, "Bad Pub/Sub Request", http.StatusBadRequest)
		log.Printf("Bad Pub/Sub Request: %v", http.StatusBadRequest)
		return
	}
	nameBytes, err := base64.URLEncoding.DecodeString(*e.Message.Data)
	if err != nil {
		http.Error(w, "Bad Pub/Sub message", http.StatusBadRequest)
		log.Printf("Bad Pub/Sub message: %v", http.StatusBadRequest)
		return
	}
	name := string(nameBytes)
	if name == "" {
		name = "World"
	}
	s := fmt.Sprintf("Hello, %s! ID: %s", name, string(r.Header.Get("Ce-Id")))
	log.Printf(s)
	fmt.Fprintln(w, s)
}

// [END eventarc_pubsub_handler]
// [START eventarc_pubsub_server]

func main() {
	http.HandleFunc("/", HelloEventsPubSub)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// [END eventarc_pubsub_server]
