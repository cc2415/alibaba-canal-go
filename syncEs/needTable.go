package syncEs

import "cc-go-canal/table"

// 需要创建索引的表
var NeedInEsTableName = map[string]interface{}{
	table.TableChatMag: table.ChatMsgModel{},
}

// 需要初始化数据的表
var NeedInitDataTableName = []string{
	table.TableChatMag,
}
