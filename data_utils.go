package mdast

// DataKey 定义 Data 中的键
type DataKey string

// DataTable 定义 Data 类型
type DataTable map[DataKey]any

// 定义 Data 中的常量
const (
	NDK_Alt           DataKey = "alt"
	NDK_URL           DataKey = "url"
	NDK_Title         DataKey = "title"
	NDK_Identifier    DataKey = "identifier"
	NDK_Label         DataKey = "label"
	NDK_ReferenceType DataKey = "referenceType"
	NDK_Depth         DataKey = "depth"
	NDK_Lang          DataKey = "lang"
	NDK_Meta          DataKey = "meta"
	NDK_Ordered       DataKey = "ordered"
	NDK_Start         DataKey = "start"
	NDK_Spread        DataKey = "spread"
	NDK_Checked       DataKey = "checked"
	NDK_Align         DataKey = "align"
)

// GetString 从 DataTable 中获取字符串值
func (dt DataTable) GetString(key DataKey) (string, bool) {
	value, ok := dt[key]
	if !ok {
		return "", false
	}
	strValue, ok := value.(string)
	return strValue, ok
}

// GetInt 从 DataTable 中获取整数值
func (dt DataTable) GetInt(key DataKey) (int, bool) {
	value, ok := dt[key]
	if !ok {
		return 0, false
	}
	intValue, ok := value.(int)
	return intValue, ok
}

// GetBool 从 DataTable 中获取布尔值
func (dt DataTable) GetBool(key DataKey) (bool, bool) {
	value, ok := dt[key]
	if !ok {
		return false, false
	}
	boolValue, ok := value.(bool)
	return boolValue, ok
}

// GetAlignType 从 DataTable 中获取 AlignType 值
func (dt DataTable) GetAlignType(key DataKey) (AlignType, bool) {
	value, ok := dt[key]
	if !ok {
		return AlignNone, false
	}
	alignValue, ok := value.(AlignType)
	return alignValue, ok
}

// GetReferenceType 从 DataTable 中获取 ReferenceType 值
func (dt DataTable) GetReferenceType(key DataKey) (ReferenceType, bool) {
	value, ok := dt[key]
	if !ok {
		return "", false
	}
	refValue, ok := value.(ReferenceType)
	return refValue, ok
}
