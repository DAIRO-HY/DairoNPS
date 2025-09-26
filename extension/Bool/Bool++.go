package Bool

				import (
					"DairoNPS/DebugTimer"
				)

// 模拟三元表达式
func Is[T any](flag bool, trueValue T, falseValue T) T {
DebugTimer.Add108()
	if flag {
DebugTimer.Add109()
		return trueValue
	} else {
		return falseValue
	}
}
