package es

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log"
)

var c *elasticsearch.Client

func init() {
	log.SetFlags(0)
	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	c = es

	// 1. Get cluster info
	//
	res, err := c.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer closeBody(res.Body)
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
}

func IndexRequest(magnet Magnet) {
	// Build the request body.
	data, err := json.Marshal(magnet)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index: "magnet",
		//DocumentID: strconv.Itoa(1),
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), c)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer closeBody(res.Body)

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), 1)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

func closeBody(b io.ReadCloser) {
	if err := b.Close(); err != nil {
		log.Println(err)
	}
}

type File struct {
	Name string
	Size string
}

type Magnet struct {
	Name     string
	InfoHash string
	Magnet   string
	Size     string
	Torrent  string
	Files    []File
}
