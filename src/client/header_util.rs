use std::error::Error;
use std::sync::Arc;
use tokio::io::AsyncReadExt;
use tokio::net::tcp::OwnedReadHalf;
use tokio::net::TcpStream;
use tokio::sync::Mutex;
use crate::clientSessionMap;

/**
 * NPS客户端头部标记
 */


/**
 * 客户端与服务器端通信连接标记
 */
pub const CLIENT_TO_SERVER_MAIN_CONNECTION: u8 = 0;

/**
 * 与客户端通信心跳标记
 */
pub const MAIN_HEART_BEAT: u8 = 1;

/**
 * 向客户端发送clientId
 */
pub const SERVER_TO_CLIENT_ID: u8 = 2;

/**
 * 向客户端申请TCP连接池请求
 */
pub const SERVER_TCP_POOL_REQUEST: u8 = 3;

/**
 * 向客户端申请UDP连接池请求
 */
pub const SERVER_UDP_POOL_REQUEST: u8 = 4;

/**
 * 服务器向客户端同步当前处于激活状态的UDP连接池端口
 */
pub const SYNC_ACTIVE_POOL_UDP_PORT: u8 = 5;

/**
 * 向客户端同步当前处于激活状态的UDP连接端口
 */
pub const SYNC_ACTIVE_BRIDGE_UDP_PORT: u8 = 6;

/**
 * 向客户端发送clientId
 */
pub const SECURITY_CLIENT_KEY: u8 = 7;

/**
 * 获取客户端Socket头部信息
   // TODO:读取数据时应该设置超时，避免恶意连接导致并发激增
 */
pub async fn get_header(client_id: &i64) -> Result<String, Box<dyn Error>> {
    let mut map = clientSessionMap.lock().await;
    let (reader,_) = map.get_mut(client_id).unwrap();

    //第一个字节值则是后面head消息的数据长度
    let head_len = reader.read_i8().await?;

    // 创建一个长度为 head_len 的缓冲区
    let mut head_data_buf = vec![0u8; head_len as usize];

    //读取数据部分
    reader.read_exact(&mut *head_data_buf).await?;
    Ok(String::from_utf8_lossy(&head_data_buf).parse()?)
}