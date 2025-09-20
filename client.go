package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Note      string    `json:"note"`
}

func main() {
	url := flag.String("url", "localhost:8443", "Address of web server (example: 'localhost:8443').")
	certFile := flag.String("cert", "certs/server.crt", "Path to server certificate (for self-signed certificates).")
	infoFlag := flag.Bool("info", false, "Print the API information.")
	testFlag := flag.Bool("test", false, "POST a test log entry to the server.")
	flag.Parse()

	// for printing to stdout in a thread safe way
	stdout := log.New(os.Stdout, "", log.Lmsgprefix)

	// set up client using the server's certificate as part of a CA pool
	cert, err := os.ReadFile(*certFile)
	if err != nil {
		stdout.Fatalf("Could not load cert file from path %s: %v\n", *certFile, err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		stdout.Fatalf("Unable to parse certificate input as PEM from path '%s'.\n", *certFile)
	}

	// create a Client object for managing connections
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}

	// if info flag is set, query to API instructions using a GET request
	if *infoFlag {
		resp, err := client.Get("https://" + *url + "/")
		if err != nil {
			stdout.Printf("Error while getting response from server: %v\n", err)
		}
		defer resp.Body.Close()
		text, err := io.ReadAll(resp.Body)
		if err != nil {
			stdout.Printf("Error reading response from server: %v\n", err)
		}
		stdout.Print(string(text))
		os.Exit(0) // setting info flag should only print instructions, then exit with success
	}
	if *testFlag {
		testEntry := LogEntry{
			StartTime: time.Now().Add(-time.Hour),
			EndTime:   time.Now(),
			Note:      "Test Data",
		}
		testData, err := json.Marshal(testEntry)
		if err != nil {
			stdout.Fatalf("%v\n", err)
		}
		resp, err := client.Post("https://"+*url+"/logs", "application/json", bytes.NewReader(testData))
		if err != nil {
			stdout.Fatalf("%v\n", err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			stdout.Fatalf("Unable to parse server response: %v\n", err)
		}
		stdout.Printf("Server response: %s\n", string(respBody))
	}

}
