package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log"
	"strings"
)

var index = "magnet"
var ec *elasticsearch.Client

func init() {
	//log.SetFlags(0)
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:9200", IP),
		},
		Username: "elastic",
		Password: "Bwin@0913",
	}
	var err error
	ec, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	// 1. Get cluster info
	res, err := ec.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer closeBody(res.Body)
	if res.IsError() {
		log.Fatalf("[%s] Error: %s", res.Status(), res.String())
	}
	var m map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", m["version"].(map[string]interface{})["number"])

	// 获取Elasticsearch集群的健康状态
	res, err = ec.Cluster.Health(ec.Cluster.Health.WithPretty())
	if err != nil {
		log.Fatalf("Error getting cluster health: %s", err)
	}
	defer closeBody(res.Body)
	if res.IsError() {
		log.Printf("[%s] Error: %s", res.Status(), res.String())
	}
	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf("Cluster health: %s", m["status"])
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
	res, err := req.Do(context.Background(), ec)
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

func IndexTorrent(t *Torrent) (string, error) {
	// Build the request body.
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(b),
		Refresh: "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), ec)
	if err != nil {
		return "", err
	}
	defer closeBody(res.Body)
	if res.IsError() {
		return "", errors.New(fmt.Sprintf("Error indexing document: [%s]", res.Status()))
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return "", err
	}
	return r["_id"].(string), nil
}

func Search(url string) int {
	// 3. Search for the indexed documents
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"url": url,
			},
		},
		"_source": []string{
			"url",
		},
	}
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	var res *esapi.Response
	res, err = ec.Search(
		ec.Search.WithContext(context.Background()),
		ec.Search.WithIndex("magnet"),
		ec.Search.WithBody(&buf),
		ec.Search.WithTrackTotalHits(true),
		ec.Search.WithPretty(),
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

func SearchByInfoHash(ih string) (int, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"infoHash": ih,
			},
		},
		"_source": []string{
			"infoHash",
		},
	}
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return 0, err
	}
	res, err := ec.Search(
		ec.Search.WithContext(context.Background()),
		ec.Search.WithIndex(index),
		ec.Search.WithBody(&buf),
		ec.Search.WithTrackTotalHits(true),
		ec.Search.WithPretty(),
	)
	if err != nil {
		return 0, err
	}
	defer closeBody(res.Body)
	var r map[string]interface{}
	if res.IsError() {
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return 0, err
		} else {
			return 0, errors.New(fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				r["error"].(map[string]interface{})["type"],
				r["error"].(map[string]interface{})["reason"]))
		}
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return 0, err
	}
	h := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(), h, int(r["took"].(float64)),
	)
	return h, nil
}

func DeleteByInfoHash(ih string) (int, error) {
	// 构建查询
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"infoHash": ih,
			},
		},
	}
	// 将查询转换为 JSON
	b, err := json.Marshal(q)
	if err != nil {
		return 0, err
	}
	// 执行 DeleteByQuery 请求
	res, err := ec.DeleteByQuery(
		[]string{index},
		strings.NewReader(string(b)),
		ec.DeleteByQuery.WithPretty(),
	)
	if err != nil {
		return 0, err
	}
	defer closeBody(res.Body)
	// 检查响应状态码
	if res.IsError() {
		return 0, errors.New(fmt.Sprintf("[ERROR] %s", res.String()))
	} else {
		var r map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return 0, err
		}
		log.Printf("Documents deleted successfully")
		return int(r["deleted"].(float64)), nil
	}
}

func closeBody(b io.ReadCloser) {
	if err := b.Close(); err != nil {
		log.Println(err)
	}
}

type File struct {
	Path   string `json:"path"`
	Length int64  `json:"length"`
}

type Torrent struct {
	Name         string `json:"name"`
	Length       int64  `json:"totalLength"`
	InfoHash     string `json:"infoHash"`
	CreationDate string `json:"creationDate"`
	Files        []File `json:"torrentFiles"`
	FileNumber   int    `json:"totalFiles"`
}

type Magnet struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	AddedTime string `json:"addedTime"`
	Magnet    string `json:"magnet"`
	Size      int64  `json:"size"`
	Torrent   string `json:"torrent"`
	Files     []File `json:"files"`
}
