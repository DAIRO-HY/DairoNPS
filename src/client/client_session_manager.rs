use std::collections::HashMap;
use std::sync::Arc;
use lazy_static::lazy_static;
use tokio::net::TcpStream;
use tokio::sync::{Mutex, RwLock};
use crate::client::client_session;
use crate::client::client_session::ClientSession;
use crate::dao::dto::client_dto::ClientDto;

lazy_static! {

/**
 * 客户端ID对应的Socket连接
 */
    pub static ref clientSessionMap: RwLock<HashMap<i64, ClientSession>> = RwLock::new(HashMap::new());
}

/**
 * 添加互斥锁
 */
// private val clientSessionMapLock = Mutex()

/**
 * 获取与客户端的会话
 */
pub async fn get_session(client_id: i64) -> Arc<ClientSession> {
    let map = clientSessionMap.read().await;
    let dto = map.get(&client_id).unwrap();
    dto.clone()
}

/**
 * 当前客户端数量
 */
pub async fn size() -> usize {
    clientSessionMap.read().await.len()
}

/**
 * 获取当前客户端会话列表
 */
// pub async fn getSessionList() -> List<ClientSession> {
//     return clientSessionMap.toList();
// }

/**
 * 获取当前客户端列表
 */
// pub async fn getClientList() -> List<ClientDto> {
//     return clientSessionMap.toList();
// }

/**
 * 添加客户端连接
 */
pub async fn validate(clientSocket: TcpStream) {

    //得到头部数据
    // let  header = HeaderUtil.getHeader(clientSocket) ?: return

    //得到客户端key
    let key = header.substring(0, header.lastIndexOf("|"));
    let client = ClientDao.selectByKey(key);
    if (client == null) {
        println("key:${key}不存在");
        clientSocket.close();
        return;
    }
    if client.enableState == 0 {
        println!("key:${key}的客户端已停止服务");
        clientSocket.shutdown();
        return;
    }
    holdOnClient(client, clientSocket);

    //客户端ip
    let ip = clientSocket.inetAddress.hostAddress;

    //从头部信息中得到客户端版本号
    let version = header.substring(header.lastIndexOf("|") + 1);
    let loginClientDto = ClientDto();
    loginClientDto.id = client.id;
    loginClientDto.ip = ip;
    loginClientDto.version = version;
    ClientDao.setClientInfo(loginClientDto);

    //将客户端ID返回给客户端
    sendClientId(client.id);

    //将加密秘钥发送到客户端
    sendClientSecurityKey(client.id);

    //开启该客户端下所有隧道监听
    ProxyAcceptManager.accept(client);
}

/**
 * 保持客户端连接
 */
async fn holdOnClient(client: ClientDto, clientSocket: TcpStream) {
    let clientID= client.id;

    //先移除之前的连接
    close(clientID);
    let session = ClientSession {
        client,
        clientSocket,
        lastHeartBeatTime: 0,
    };

    let sessionClone1 = Arc::new(Mutex::new(session));
    let sessionClone2 = sessionClone1.clone();
    {
        let mut sessionMap = clientSessionMap.write().await;
        sessionMap.insert(clientID, sessionClone1.lock());
    }
    client_session::start(sessionClone2);
}

/**
 * 将客户端ID返回给客户端
 */
async fn sendClientId(clientID: Int) {
    this.send(clientID, HeaderUtil.SERVER_TO_CLIENT_ID, clientID.toString());
}

/**
 * 将加密秘钥发送到客户端
 */
async fn sendClientSecurityKey(clientID: u64) {
    this.clientSessionMap[clientID]?.send(SecurityUtil.clientKeyArray);
}

/**
 * 向客户端申请TCP连接池请求
 * @param clientID 客户端ID
 * @param count 申请数量
 */
pub async fn sendTCPPoolRequest(clientID: Int, count: u8) {
    send(clientID, HeaderUtil.SERVER_TCP_POOL_REQUEST, count.toString());
}

/**
 * 向客户端申请UDP连接池请求
 * @param clientID 客户端ID
 * @param count 申请数量
 */
pub async fn sendUDPPoolRequest(clientID: Int, count: Int) {
    send(clientID, HeaderUtil.SERVER_UDP_POOL_REQUEST, count.toString())
}

/**
 * 往客户端发送数据
 * @param clientID 客户端ID
 * @param flag 头部标记
 * @param message 头部消息
 */
pub async fn send(clientID: Int, flag: Int, message: String) {
    clientSessionMap[clientID]?.send(flag, message);
}

/**
 * 关闭客户端
 */
pub async fn close(clientId: Int) {
    if clientSessionMap.containsKey(clientId) { //如果存在
        let clientSession = this.clientSessionMap[clientId];
        clientSession.close();
        ClientDao.setDataLen(clientSession.client);
        clientSessionMap.remove(clientId);
    }
}
