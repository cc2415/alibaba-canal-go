package main

import (
	"cc-go-canal/config"
	"cc-go-canal/syncEs"
	"cc-go-canal/table"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
)

func main() {
	// 192.168.199.17 替换成你的canal server的地址
	// example 替换成-e canal.destinations=example 你自己定义的名字
	connector := client.NewSimpleCanalConnector(config.AppConfig.AlibabaCanal.Address, config.AppConfig.AlibabaCanal.Port,
		config.AppConfig.AlibabaCanal.Username, config.AppConfig.AlibabaCanal.Password, config.AppConfig.AlibabaCanal.Destination, 60000, 60*60*1000)
	err := connector.Connect()
	fmt.Println("启动")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// https://github.com/alibaba/canal/wiki/AdminGuide
	//mysql 数据解析关注的表，Perl正则表达式.
	//
	//多个正则之间以逗号(,)分隔，转义符需要双斜杠(\\)
	//
	//常见例子：
	//
	//  1.  所有表：.*   or  .*\\..*
	//	2.  canal schema下所有表： canal\\..*
	//	3.  canal下的以canal打头的表：canal\\.canal.*
	//	4.  canal schema下的一张表：canal\\.test1
	//  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)

	//err = connector.Subscribe(".*\\..*")
	err = connector.Subscribe(config.AppConfig.AlibabaCanal.Database + "\\..*") //数据库\\所有表
	//err = connector.Subscribe("canal\\\\..")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println(connector.Connected)
	fmt.Println(connector.ClientIdentity)
	for {

		message, err := connector.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		batchId := message.Id

		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			//fmt.Println("===没有数据了=== " + time.DateTime)
			continue
		}
		//处理数据
		printEntry(message.Entries)

	}
}

func printEntry(entrys []pbe.Entry) {

	for _, entry := range entrys {
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println("")
			fmt.Println("")
			fmt.Println(fmt.Sprintf(" %s库.%s 表有变化", header.GetSchemaName(), header.GetTableName())) //GetSchemaName 数据库
			_, exit := syncEs.NeedInEsTableName[header.GetTableName()]
			if exit {
				syncEs.CreateChatMsgIndex(header.GetSchemaName()) // 创建索引
			}
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_DELETE {
					fmt.Println("……………………………… 数据删除显示开始 ………………………………")
					printColumn(rowData.GetBeforeColumns(), nil)
					fmt.Println("************ 数据删除显示结束 ************")
				} else if eventType == pbe.EventType_INSERT {
					fmt.Println("……………………………… 数据新增显示开始 ………………………………")
					rowDataRe := make(map[string]string) // 创建一个 map 用于存储结果
					printColumn(rowData.GetAfterColumns(), rowDataRe)
					if _, exit := syncEs.NeedInEsTableName[header.GetTableName()]; exit {
						switch header.GetTableName() {
						case table.TableChatMag:
							syncEs.CreateOrUpdateChatMsgDocument(header.GetTableName(), rowDataRe)
							break
						}
					}
					fmt.Println(rowDataRe)
					//table.CreateOrUpdateChatMsgDocument()
					fmt.Println("************ 数据新增显示结束 ************")
				} else {
					fmt.Println("……………………………… 数据更新显示开始 ………………………………")
					// 创建一个切片

					beforeData := []map[string]interface{}{}
					for _, col := range rowData.GetBeforeColumns() {
						beforeData = append(beforeData, map[string]interface{}{
							"column": col.GetName(),
							"value":  col.GetValue(),
						})
					}
					afterData := []map[string]interface{}{}
					for _, col := range rowData.GetAfterColumns() {
						afterData = append(afterData, map[string]interface{}{
							"column":  col.GetName(),
							"value":   col.GetValue(),
							"updated": col.GetUpdated(),
						})
					}
					fmt.Println("旧数据")
					//旧数据
					for _, datum := range beforeData {
						fmt.Print(fmt.Sprintf(" %s:%s ", datum["column"], datum["value"]))
					}
					fmt.Println("")

					rowDataRe := make(map[string]string)
					fmt.Println("新数据")
					updateData := []map[string]interface{}{}
					for i, datum := range afterData {
						if datum["updated"] == true {
							//变化的值
							updateData = append(updateData, map[string]interface{}{
								"column":   datum["column"],
								"oldValue": beforeData[i]["value"],
								"newValue": datum["value"],
							})
						}
						rowDataRe[datum["column"].(string)] = datum["value"].(string)
						fmt.Print(fmt.Sprintf(" %s:%s ", datum["column"], datum["value"]))
					}
					fmt.Println("")
					fmt.Println("变化的数据")
					for _, datum := range updateData {
						fmt.Println(fmt.Sprintf("字段:%s 从:%s 改为:%s", datum["column"], datum["oldValue"], datum["newValue"]))
					}
					if _, exit := syncEs.NeedInEsTableName[header.GetTableName()]; exit {
						switch header.GetTableName() {
						case table.TableChatMag:
							syncEs.CreateOrUpdateChatMsgDocument(header.GetTableName(), rowDataRe)
							break
						}
					}

					fmt.Println("……………………………… 数据更新显示结束 ………………………………")
				}
			}
		}
	}
}

func printColumn(columns []*pbe.Column, rowData map[string]string) {
	for _, col := range columns {
		fmt.Print(fmt.Sprintf(" %s:%s ", col.GetName(), col.GetValue()))
		rowData[col.GetName()] = col.GetValue()
	}
	fmt.Println("")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
