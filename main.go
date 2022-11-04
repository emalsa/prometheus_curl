package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Topic struct {
	Timestamp string  `json:"timestamp"`
	Payload   Payload `json:"jsonPayload"`
}
type Payload struct {
	Message Message `json:"message"`
}
type Message struct {
	MessageId  string     `json:"messageId"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Location    string `json:"location"`
	Type        string `json:"type"`
	SitecheckId string `json:"sitecheck_id"`
}

func main() {
	port := "8080"
	log.Print("Starting server...")
	handler := http.HandlerFunc(curlExecute)
	http.Handle("/", handler)
	log.Printf("listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

func curlExecute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	var topic Topic
	return

	//var result map[string]interface{}
	//json.Unmarshal([]byte(body), &result)
	errors := json.Unmarshal([]byte(body), &topic)
	if errors != nil {
		//fmt.Println(err.Error())
	}

	cmd := exec.Command("/usr/bin/curl", "-s", "-w", "@curl-format.txt", "--location", "--include", "--request", "GET", "--compressed", "https://untrusted-root.badssl.com/", "-vI")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		//fmt.Println(err.Error())
	}

	log.Print(topic.Payload.Message.Attributes.SitecheckId)
	log.Print(topic.Payload.Message.Attributes.Type)
	//log.Print(string(body))

	// Print the output
	//fmt.Println(string(stdout))
	w.Write(stdout)

	return
}

//cmd.Stdout = os.Stdout
//cmd.Stderr = os.Stderr
//err := cmd.Run()
//if err != nil {
//	log.Fatalf("cmd.Run() failed with %s\n", err)
//}
//fmt.Println(string(os.Stdout.s))
//}
