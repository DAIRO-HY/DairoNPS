package Number

import "fmt"

// 数据流量单位换算
func ToDataSize(input any) string {
	var inputFloat64 float64
	switch v := input.(type) {
	case int:
		inputFloat64 = float64(v)
	case int8:
		inputFloat64 = float64(v)
	case int16:
		inputFloat64 = float64(v)
	case int32:
		inputFloat64 = float64(v)
	case int64:
		inputFloat64 = float64(v)
	case uint:
		inputFloat64 = float64(v)
	case uint8:
		inputFloat64 = float64(v)
	case uint16:
		inputFloat64 = float64(v)
	case uint32:
		inputFloat64 = float64(v)
	case uint64:
		inputFloat64 = float64(v)
	case float32:
		inputFloat64 = float64(v)
	case float64:
		inputFloat64 = v
	default:
		inputFloat64 = 0.0
	}

	var dataSize float64
	var unit string
	if inputFloat64 >= 1024*1024*1024*1024 {
		dataSize = inputFloat64 / (1024 * 1024 * 1024 * 1024)
		unit = "TB"
	} else if inputFloat64 >= 1024*1024*1024 {
		dataSize = inputFloat64 / (1024 * 1024 * 1024)
		unit = "GB"
	} else if inputFloat64 >= 1024*1024 {
		dataSize = inputFloat64 / (1024 * 1024)
		unit = "MB"
	} else if inputFloat64 >= 1024 {
		dataSize = inputFloat64 / 1024
		unit = "KB"
	} else {
		dataSize = inputFloat64
		unit = "B"
	}
	dataSizeStr := fmt.Sprintf("%.2f\n", dataSize)
	return dataSizeStr + unit
}
