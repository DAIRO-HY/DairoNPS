package SecurityUtil

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"fmt"
	"math/rand"
)

// 服务器端加密秘钥
var ServerSecurityKey = [256]uint8{}

// 客户端加密秘钥
var ClientSecurityKey = [256]uint8{}

func init() {
DebugTimer.Add468()
	for i := range ServerSecurityKey {
DebugTimer.Add469()
		ServerSecurityKey[i] = uint8(i)
	}

	// 打乱数组
	rand.Shuffle(256, func(i, j int) {
		ServerSecurityKey[i], ServerSecurityKey[j] = ServerSecurityKey[j], ServerSecurityKey[i]
	})
	//for i, it := range ServerSecurityKey {
	//	fmt.Printf("%d->%d\n", i, it)
	//}
	fmt.Println("----------------------------------------------------------------------------------------------")

	//服务端数组的值是客户端数组的序号,对应的服务端数组的序号则是客户端数组的值
	for i, it := range ServerSecurityKey {
DebugTimer.Add470()
		ClientSecurityKey[it] = uint8(i)
	}
	//for i, it := range ClientSecurityKey {
	//	fmt.Printf("%d->%d\n", i, it)
	//}
}

/**
 * 加密数据
 * @param data 要加密的数据
 * @param len 要加密的数据长度
 */
func Mapping(data []uint8, len int) {
DebugTimer.Add471()
	for i := 0; i < len; i++ {
DebugTimer.Add472()
		value := data[i]
		data[i] = ServerSecurityKey[value]
	}
}
