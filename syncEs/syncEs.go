package syncEs

import (
	"cc-go-canal/config"
	"cc-go-canal/es"
	"cc-go-canal/table"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 获取索引名
func GetIndexName(tableName string) string {
	return config.AppConfig.DatabaseName + "." + tableName + "-000001"
}

// 获取别名
func GetAliasName(tableName string) string {
	return config.AppConfig.DatabaseName + "." + tableName
}

// 创建索引
func CreateIndexAndAlias(tableName string) {
	indexName := GetIndexName(tableName)
	err, e := es.CreateIndex(indexName) // 数据库名_表名
	if err != nil {
		panic(err)
	}
	// 创建别名
	es.SetAliasName(indexName, GetAliasName(tableName))
	fmt.Println(fmt.Sprintf("es初始化索引成功 %s  别名:%d", indexName, GetAliasName(tableName), e.StatusCode))
}

func CreateTableDocument(tableName string, data map[string]interface{}) {
	es.CreateDocument(GetAliasName(tableName), data)
}

// 根据id 更新数据
func UpdateTableDocument(tableName string, id string, data map[string]interface{}) {
	query, _ := json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": id,
			},
		},
	})
	es.UpdateDocumentByQuery(GetAliasName(tableName), query, data)
}

// 初始化表的数据
func InitTableData() {
	// 数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.AppConfig.AlibabaCanal.Username, config.AppConfig.AlibabaCanal.Password,
		config.AppConfig.AlibabaCanal.DbHost, config.AppConfig.AlibabaCanal.DbPort, config.AppConfig.DatabaseName)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %s\n", err)
	}

	for _, itemModel := range table.ModelList {
		itemTableName := itemModel.GetTableName()
		if itemModel.GetNeedInitData() && !es.CheckIndexExit(GetAliasName(itemTableName)) {
			CreateIndexAndAlias(itemTableName)
			es.SetIndexSettings(GetAliasName(itemTableName), 0, "-1")
			var results []map[string]interface{}
			db.Raw(fmt.Sprintf("select * from %s", itemTableName)).FindInBatches(&results, 1000, func(tx *gorm.DB, batch int) error {
				_, err := es.BulkCreate(GetAliasName(itemTableName), results)
				if err != nil {
					log.Println("批量插入数据失败", err)
				} else {
					log.Println("批量插入数据成功")
				}
				return nil
			})
			es.SetIndexSettings(GetAliasName(itemTableName), config.AppConfig.EsIndexShareReplicasNum, "-1")
		}

	}
}
