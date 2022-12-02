package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

type Topic struct {
	Url         string `json:"url"`
	Type        string `json:"type"`
	CheckItemId string `json:"check_item_id"`
	CloudUrl    string `json:"cloud_url"`
}

type Response struct {
	Success     bool   `json:"success"`
	CheckItemId string `json:"check_item_id"`
	Response    string `json:"response"`
	Type        string `json:"type"`
	CloudUrl    string `json:"cloud_url"`
	Url         string `json:"url"`
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
	body, err := ioutil.ReadAll(r.Body)
	//fmt.Println(string(body))
	var topic Topic
	//return
	fmt.Println("Start test")
	log.Print("Start test")
	//var result map[string]interface{}
	//json.Unmarshal([]byte(body), &result)
	errors := json.Unmarshal([]byte(body), &topic)
	if errors != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Ok"))
	go func() {
		fmt.Println("Start sleep")
		log.Print("Log Start sleep")
		//w.Write(stdout)
		//time.Sleep(15 * time.Second)
		//log.Print("End sleep")

		//testUrl := "http://www.nicastro.io"
		cmd := exec.Command("/usr/bin/curl", "-w", "@curl-format.txt", "--request", "GET", "--compressed", "-Lvs", "-o", "/dev/null", topic.Url)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		//log.Print(topic.Payload.Message.Attributes.SitecheckId)
		//log.Print(topic.Payload.Message.Attributes.Type)

		// Print the output

		sEnc := b64.StdEncoding.EncodeToString([]byte(stdout))
		response := Response{
			Success:     true,
			CheckItemId: topic.CheckItemId,
			Response:    sEnc,
			Type:        topic.Type,
			CloudUrl:    topic.CloudUrl,
			Url:         topic.Url,
		}

		responseJson, _ := json.Marshal(response)
		apiUrl := "http://localhost:61032/api/check_item/update?XDEBUG_SESSION_START=PHPSTORM"
		// Pass new buffer for request with URL to post.
		// This will make a post request and will share the JSON data
		fmt.Println("responseJson: ", string(responseJson))
		resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(responseJson))

		// An error is returned if something goes wrong
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			panic(err)
		}
		// Need to close the response stream, once response is read.
		// Hence, defer close. It will automatically take care of it.
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				panic(err)
			}
		}(resp.Body)

		// Check response code.
		if resp.StatusCode == http.StatusOK {
			_, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				//Failed to read response.
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				panic(err)
			}

			// Convert bytes to String and print
			//jsonStr := string(body)
			//fmt.Println("Response: ", jsonStr)

		} else {
			//The status is not Created. print the error.
			fmt.Println("Get failed with error: ", resp.Status)
		}

		return
	}()
}
