package syncEs

import (
	"cc-go-canal/es"
	"cc-go-canal/table"
	"encoding/json"
	"fmt"
)

// 索引
var IndexName = ""

// 创建索引
func CreateChatMsgIndex(schemaName string) {
	IndexName = schemaName + "." + table.TableChatMag
	err, e := es.CreateIndex(IndexName) // 数据库名_表名
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("es初始化索引成功 %s %d", table.TableChatMag, e.StatusCode))
}

func CreateOrUpdateChatMsgDocument(tableName string, msg map[string]string) {
	jsonData, _ := json.Marshal(msg)
	chatMsg := NeedInEsTableName[tableName].(table.BaseModel)
	json.Unmarshal(jsonData, &chatMsg)
	//fmt.Println("写入es")
	es.CreateOrUpdateDocument(IndexName, chatMsg.Id, msg)
	//fmt.Println(err, res)
}
