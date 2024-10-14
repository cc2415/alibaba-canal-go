package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var Client *elasticsearch.Client

func init() {
	fmt.Println("初始化esClient")
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200/", "http://127.0.0.1:9201/", "http://127.0.0.1:9202/"},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 5 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Error creating the client: %s", err)
	} else {
		fmt.Println("连接成功")
		//log.Println(es.Info())
	}
	Client = es
}

func CreateIndex(indexName string) (error, *esapi.Response) {
	create, err := Client.Indices.Create(indexName)
	if err != nil {
		return err, nil
	}
	return nil, create
}
func CreateOrUpdateDocument(indexName string, id string, data interface{}) (error, *esapi.Response) {

	data2, _ := json.Marshal(data)
	request := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: id,
		Body:       bytes.NewReader(data2),
	}

	do, err := request.Do(context.Background(), Client)
	if err != nil {
		return err, nil
	}
	return err, do
}
func GetByDocumentId(indexName string, id string) (error, *esapi.Response) {
	request := esapi.GetRequest{Index: indexName, DocumentID: id}
	do, err := request.Do(context.Background(), Client)
	return err, do
}
func Search(indexName string) (error, *esapi.Response) {
	// 查询条件，搜索 "title" 字段中包含 "Hello" 的文档
	query := `{
        "query": {
            "match": {
                "name": "Hello"
            }
        }
    }`
	res, err := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(query),
	}.Do(context.Background(), Client)
	return err, res
}
