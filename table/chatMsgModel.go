package table

type ChatMsgModel struct {
	BaseModel
}

var chatMsgModel = ChatMsgModel{
	BaseModel{
		TableName:    "fa_chat_msg",
		NeedInitData: true,
		NeedInEs:     true,
	}}

func init() {
	ModelList = append(ModelList, chatMsgModel)
}
