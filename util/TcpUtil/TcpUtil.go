package TcpUtil

				import (
					"DairoNPS/DebugTimer"
				)

import "net"

// WriteAll net.Conn写入所有数据
func WriteAll(tcp net.Conn, data []uint8) error {
DebugTimer.Add473()

	//要发送的数据总长度
	total := len(data)

	//已经大宋的数据长度
	sendedTotal := 0

	for {
DebugTimer.Add474()
		//n:当前发送的数据长度
		n, err := tcp.Write(data[sendedTotal:])
		if err != nil {
DebugTimer.Add475()
			return err
		}
		sendedTotal += n
		if sendedTotal == total {
DebugTimer.Add476()
			break
		}
	}
	return nil
}

// ReadNByte 读取指定长度数据
func ReadNByte(tcp net.Conn, len int) ([]uint8, error) {
DebugTimer.Add477()

	//记录已经读取到的数据大小
	var readLen = 0
	data := make([]uint8, len)
	for {
DebugTimer.Add478()
		buffer := make([]uint8, len-readLen)
		n, readErr := tcp.Read(buffer)
		if n > 0 {
DebugTimer.Add479()
			copy(data[readLen:readLen+n], buffer[:n])
			readLen += n
			if readLen == len {
DebugTimer.Add480()
				break
			}
		}
		if readErr != nil {
DebugTimer.Add481()
			return nil, readErr
		}
	}
	return data, nil
}
