package ClientSessionManagerInterface

// 向外部暴露函数
type ClientSessionManagerInterface interface {
	SendTCPPoolRequest(clientId int, count int)
	SendUDPPoolRequest(clientId int, count int)
}
