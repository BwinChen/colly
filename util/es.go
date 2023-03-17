package util

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
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	c = es

	// 1. Get cluster info
	var res *esapi.Response
	res, err = c.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer closeBody(res.Body)
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	var m map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", m["version"].(map[string]interface{})["number"])
}

func IndexRequest(m Magnet) {
	// Build the request body.
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "magnet",
		DocumentID: Checksum(InfoHash(m.Magnet)),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	var res *esapi.Response
	res, err = req.Do(context.Background(), c)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer closeBody(res.Body)

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), 1)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}

func Search(URL string) int {
	// 3. Search for the indexed documents
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"URL": URL,
			},
		},
		"_source": []string{
			"URL",
		},
	}
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	var res *esapi.Response
	res, err = c.Search(
		c.Search.WithContext(context.Background()),
		c.Search.WithIndex("magnet"),
		c.Search.WithBody(&buf),
		c.Search.WithTrackTotalHits(true),
		c.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer closeBody(res.Body)

	var r map[string]interface{}
	if res.IsError() {
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				r["error"].(map[string]interface{})["type"],
				r["error"].(map[string]interface{})["reason"],
			)
		}
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	i := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(), i, int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	//for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	//	log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	//}
	return i
}

func closeBody(b io.ReadCloser) {
	if err := b.Close(); err != nil {
		log.Println(err)
	}
}

type File struct {
	Name string
	Size int64
}

type Magnet struct {
	Name    string
	URL     string
	Magnet  string
	Size    int64
	Torrent string
	Files   []File
}
