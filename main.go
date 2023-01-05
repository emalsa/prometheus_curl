package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type Topic struct {
	Url         string `json:"url"`
	Type        string `json:"type"`
	CheckItemId string `json:"check_item_id"`
	CloudUrl    string `json:"cloud_url"`
	FromHost    string `json:"from_host"`
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
	fmt.Println("Starting server...")
	handler := http.HandlerFunc(curlExecute)
	http.Handle("/", handler)
	fmt.Println("listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

func curlExecute(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	var topic Topic
	fmt.Println("Start test")

	errors := json.Unmarshal([]byte(body), &topic)
	if errors != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusRequestHeaderFieldsTooLarge)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Ok"))
	go func() {

		cmd := exec.Command("/usr/bin/curl", "-w", "@curl-format.txt", "--request", "GET", "--compressed", "-Lvs", "-o", "/dev/null", topic.Url)
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusFailedDependency)
			w.Write([]byte(err.Error()))
		}

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
		//apiUrl := "http://localhost:61032/api/check_item/update?XDEBUG_SESSION_START=PHPSTORM"
		apiUrl := topic.FromHost
		// Pass new buffer for request with URL to post.
		// This will make a post request and will share the JSON data
		fmt.Println("responseJson: ", string(responseJson))
		resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(responseJson))

		// An error is returned if something goes wrong
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte(err.Error()))
			panic(err)
		}

		// Need to close the response stream, once response is read.
		// Hence, defer close. It will automatically take care of it.
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusConflict)
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

		} else {
			//The status is not Created. print the error.
			fmt.Println("Get failed with error: ", resp.Status)
		}

		return
	}()
}
