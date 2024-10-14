package table

// 表名
const TableChatMag = "fa_chat_msg"

type StructChatMsg struct {
	BaseModel
	Content    map[string]interface{} `json:"content"`
	ToUserId   int64                  `json:"to_user_id"`
	FromUserId int64                  `json:"from_user_id"`
	CreateAt   int64                  `json:"create_at"`
	UpdateAt   int64                  `json:"update_at"`
}
