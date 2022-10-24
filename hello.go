package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"os/exec"
)

func main() {

	log.Print("starting server...")
	http.HandleFunc("/", curlExecute)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func curlExecute(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	cmd := exec.Command("/usr/bin/curl", "-v", "-w", "@curl-format.txt", "--location", "--include", "--request", "GET", "--compressed", "http://nicastro.io", "-vI")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

//cmd.Stdout = os.Stdout
//cmd.Stderr = os.Stderr
//err := cmd.Run()
//if err != nil {
//	log.Fatalf("cmd.Run() failed with %s\n", err)
//}
//fmt.Println(string(os.Stdout.s))
//}
