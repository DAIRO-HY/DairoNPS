use std::collections::HashMap;
use std::error::Error;
use std::fmt::format;
use std::io;
use std::io::ErrorKind;
use std::sync::Arc;
use lazy_static::lazy_static;
use tokio::net::TcpStream;
use tokio::sync::{Mutex, RwLock};
use crate::client::client_session::ClientSession;
use crate::client::{client_session, header_util};
use crate::dao::client_dao;
use crate::dao::dto::client_dto::ClientDto;
// use crate::client::{client_session, header_util};
// use crate::client::client_session::ClientSession;
// use crate::dao::client_dao;
// use crate::dao::dto::client_dto::ClientDto;

lazy_static! {

/**
 * 客户端ID对应的Socket连接
 */
    pub static ref clientSessionMap: RwLock<HashMap<i64, Arc<ClientSession>>> = RwLock::new(HashMap::new());
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
pub async fn validate(tcp: TcpStream) -> Result<(), Box<dyn Error>> {

    //得到头部数据
    let (tcp,head) = header_util::get_header(tcp).await?;
    let heads: Vec<_> = head.split("|").collect();

    //得到客户端key
    let key = heads[0];
    let client_opt = client_dao::selectByKey(key);
    if client_opt.is_none() {
        Err::<ErrorKind, io::Error>(io::Error::new(ErrorKind::Other, "key".to_owned() + &key + "不存在"));
        // clientSocket.close();
    }
    let client = client_opt.unwrap();
    if client.enable_state == 0 {
        Err::<ErrorKind, io::Error>(io::Error::new(ErrorKind::Other, "key".to_owned() + &key + "的客户端已停止服务"));
    }

    let client_id = client.id;

    //客户端ip
    // let ip = clientSocket.inetAddress.hostAddress;
    let ip = "";

    //从头部信息中得到客户端版本号
    let version = heads[1];

    client_dao::setClientInfo(client_id, ip.to_string(), version.to_string());

    //将客户端ID返回给客户端
    // sendClientId(client_id);


    holdOnClient(tcp,client);
    //
    // //将加密秘钥发送到客户端
    // sendClientSecurityKey(client.id);
    //
    // //开启该客户端下所有隧道监听
    // ProxyAcceptManager.accept(client);

    Ok(())
}

/**
 * 保持客户端连接
 */
async fn holdOnClient(tcp: TcpStream,client: ClientDto) {
    let client_id = client.id;

    //先移除之前的连接
    close(client_id);
    let session = ClientSession {
        client,
        clientSocket: Mutex::new(tcp),
        lastHeartBeatTime: 0,
    };
    {
        //将会话添加到map
        clientSessionMap.write().await.insert(client_id, Arc::new(session));
    }
    let session = clientSessionMap.read().await.get(&client_id).unwrap();
    client_session::start(session);
}

/**
 * 将客户端ID返回给客户端
 */
async fn sendClientId(clientID: i64) {
    send(clientID, header_util::SERVER_TO_CLIENT_ID as i8, clientID.to_string()).await;
}

// /**
//  * 将加密秘钥发送到客户端
//  */
// async fn sendClientSecurityKey(clientID: u64) {
//     this.clientSessionMap[clientID]?.send(SecurityUtil.clientKeyArray);
// }
//
// /**
//  * 向客户端申请TCP连接池请求
//  * @param clientID 客户端ID
//  * @param count 申请数量
//  */
// pub async fn sendTCPPoolRequest(clientID: Int, count: u8) {
//     send(clientID, HeaderUtil.SERVER_TCP_POOL_REQUEST, count.toString());
// }
//
// /**
//  * 向客户端申请UDP连接池请求
//  * @param clientID 客户端ID
//  * @param count 申请数量
//  */
// pub async fn sendUDPPoolRequest(clientID: Int, count: Int) {
//     send(clientID, HeaderUtil.SERVER_UDP_POOL_REQUEST, count.toString())
// }

/**
 * 往客户端发送数据
 * @param clientID 客户端ID
 * @param flag 头部标记
 * @param message 头部消息
 */
pub async fn send(clientID: i64, flag: i8, message: String) {
    clientSessionMap.read().await.get(&clientID).unwrap().sendHeader(flag, message).await;
}

/**
 * 关闭客户端
 */
pub async fn close(clientId: i64) {
    let mut map = clientSessionMap.read().await;
    if map.contains_key(&clientId) {
        let session = map.get(&clientId);
        session.unwrap().close().await;
        map.remove(&clientId);
    }
}
