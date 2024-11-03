package ClientSessionManagerInterface

// 向外部暴露函数
type ClientSessionManagerInterface interface {
	SendTCPPoolRequest(clientId int, count int)
	SendUDPPoolRequest(clientId int, count int)

	/**
	 * 向客户端当前激活的UDP端口
	 * @param clientID 客户端ID
	 * @param count 申请数量
	 */
	SendActiveUDPBridge(clientId int, ports string)
}
