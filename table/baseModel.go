package table

type BaseModel struct {
	BaseModelInterface
	Id string `json:"id"`
	// 表名
	TableName string
	// 是否需要初始化
	NeedInitData bool
	// 是否需要同步到es
	NeedInEs bool
}
type BaseModelInterface interface {
	GetTableName() string
	GetNeedInitData() bool
	GetNeedInEs() bool
}

func (m BaseModel) GetNeedInEs() bool {
	return m.NeedInEs
}
func (m BaseModel) GetNeedInitData() bool {
	return m.NeedInitData
}
func (m BaseModel) GetTableName() string {
	return m.TableName
}

var ModelList []BaseModelInterface

// 判断是否需要同步到es
func CheckTableNeedInEs(tableName string) bool {
	for _, item := range ModelList {
		if item.GetTableName() == tableName && item.GetNeedInEs() {
			return true
		}
	}
	return false
}

// 获取model
func GetModel(tableName string) BaseModelInterface {
	for _, item := range ModelList {
		if item.GetTableName() == tableName {
			return item
		}
	}
	return nil
}
