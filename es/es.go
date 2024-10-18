package es

import (
	"bytes"
	"cc-go-canal/config"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var Client *elasticsearch.Client

func init() {
	fmt.Println("初始化esClient")
	cfg := elasticsearch.Config{
		Addresses: config.AppConfig.EsAddress,
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
	}
	Client = es
}

// 创建索引
func CreateIndex(indexName string) (error, *esapi.Response) {
	create, err := Client.Indices.Create(indexName)
	if err != nil {
		return err, nil
	}
	return nil, create
}

// 插入或更新文档
func CreateDocument(indexName string, data interface{}) (error, *esapi.Response) {

	data2, _ := json.Marshal(data)
	request := esapi.IndexRequest{
		Index: indexName,
		Body:  bytes.NewReader(data2),
	}

	do, err := request.Do(context.Background(), Client)
	if err != nil {
		return err, nil
	}
	return err, do
}

//	根据查询条件更新文档
//
// query: {"query":{"term":{"id":19}}}
func UpdateDocumentByQuery(indexName string, query []byte, data map[string]interface{}) {
	response := Search(indexName, query)
	if response.IsArray() {
		data = map[string]interface{}{"doc": data}
		upd, _ := json.Marshal(data)
		responseList := response.Array()
		for i := 0; i < len(responseList); i++ {
			do, _ := esapi.UpdateRequest{Index: indexName, DocumentID: responseList[i].Get("_id").String(),
				Body: bytes.NewReader([]byte(upd)),
			}.Do(context.Background(), Client)
			fmt.Println("更新结果", do)
		}
	}

}

// 查找数据
func Search(indexName string, query []byte) gjson.Result {
	response, _ := esapi.SearchRequest{Index: []string{indexName},
		Body: bytes.NewReader(query),
	}.Do(context.Background(), Client)
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	res := gjson.Get(string(body), "hits.hits")
	return res
}

// 判断索引是否存在
func CheckIndexExit(IndexName string) bool {
	res, _ := esapi.IndicesExistsRequest{Index: []string{IndexName}}.Do(context.Background(), Client)
	if res.StatusCode == 200 {
		return true
	}
	return false
}

// 设置别名
func SetAliasName(indexName string, aliasName string) {
	res, err := esapi.IndicesPutAliasRequest{Index: []string{indexName}, Name: aliasName}.Do(context.Background(), Client)
	if res.StatusCode != 200 {
		log.Fatalln("设置别名失败", res)
	} else {
		log.Println("设置别名成功", err, res)
	}
}

// 批量添加数据
func BulkCreate(indexName string, data []map[string]interface{}) (*esapi.Response, error) {
	var body []byte
	for _, doc := range data {
		meta := []byte(fmt.Sprintf(`{"index":{"_index":"%s"}}%s`, indexName, "\n"))
		body = append(body, meta...)

		source, _ := json.Marshal(doc)
		body = append(body, source...)
		body = append(body, '\n')
	}
	return esapi.BulkRequest{Index: indexName,
		Body: bytes.NewReader(body),
	}.Do(context.Background(), Client)
}

// 设置索引的副本数和刷新间隔
func SetIndexSettings(indexName string, replicas int, refreshInterval string) {
	settings := map[string]interface{}{
		"index": map[string]interface{}{
			"number_of_replicas": replicas,
			//"refresh_interval":   refreshInterval,
		},
	}

	body, _ := json.Marshal(settings)

	req, _ := esapi.IndicesPutSettingsRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}.Do(context.Background(), Client)

	defer req.Body.Close()
	fmt.Printf("Index settings updated: replicas=%d, refresh_interval=%s\n", replicas, refreshInterval)
}

// rollov
func Rollover(aliasName string, maxAge string, maxDocs int, maxSize string) {
	rolloverConditions := map[string]interface{}{
		"aliases": map[string]interface{}{
			aliasName: map[string]interface{}{"is_write_index": true},
		},
		"conditions": map[string]interface{}{
			"max_age":  maxAge,
			"max_docs": maxDocs,
			"max_size": maxSize,
		},
	}
	body, _ := json.Marshal(rolloverConditions)

	esapi.IndicesRolloverRequest{Body: bytes.NewReader(body), Alias: aliasName}.Do(context.Background(), Client)

}
