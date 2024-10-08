package SecurityUtil

import (
	"math/rand"
)

/**
 * 服务器端加密秘钥
 */
var ServerKeyArray = make([]byte, 128)

/**
 * 客户端加密秘钥
 */
var ClientKeyArray = make([]byte, 128)

func Init() {
	for i := range ServerKeyArray {
		ServerKeyArray[i] = byte(i)
	}

	// 打乱数组
	rand.Shuffle(128, func(i, j int) {
		ServerKeyArray[i], ServerKeyArray[j] = ServerKeyArray[j], ServerKeyArray[i]
	})

	//for i, it := range serverKeyArray {
	//	fmt.Printf("%d->%d\n", i, it)
	//}

	//服务端数组的值是客户端数组的序号,对应的服务端数组的序号则是客户端数组的值
	for i, it := range ServerKeyArray {
		ClientKeyArray[it] = byte(i)
	}

	//for i, it := range clientKeyArray {
	//	fmt.Printf("%d->%d\n", i, it)
	//}
}

/**
 * 加密数据
 * @param data 要加密的数据
 * @param len 要加密的数据长度
 */
func mapping(data []byte, len int) {
	for i := 0; i < len; i++ {
		value := data[i]
		if value < 0 {
			data[i] = -ServerKeyArray[-value-1] - 1
		} else {
			data[i] = ServerKeyArray[value]
		}
	}
}
