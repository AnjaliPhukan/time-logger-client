package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	url := flag.String("url", "localhost:8443", "Address of web server (example: 'localhost:8443').")
	certFile := flag.String("cert", "certs/server.crt", "Path to server certificate (for self-signed certificates).")
	flag.Parse()

	stdout := log.Logger{}
	stdout.SetOutput(os.Stdout)

	cert, err := os.ReadFile(*certFile)
	if err != nil {
		stdout.Fatalf("Could not load cert file from path %s: %v\n", *certFile, err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		stdout.Fatalf("Unable to parse certificate input as PEM from path '%s'.\n", *certFile)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Get("https://" + *url + "/health")
	text, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response from server: %v\n", err)
	}
	stdout.Printf("Server response: %s\n", string(text))
}
