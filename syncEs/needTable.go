package syncEs

import "cc-go-canal/table"

// 需要创建索引的表
var NeedInEsTableName = map[string]interface{}{
	table.TableChatMag: table.BaseModel{},
}
