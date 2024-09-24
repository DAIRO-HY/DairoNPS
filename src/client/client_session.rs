use std::sync::Arc;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::tcp::{OwnedReadHalf, OwnedWriteHalf};
use tokio::net::TcpStream;
use tokio::sync::Mutex;
use crate::client::header_util;
use crate::dao::dto::client_dto::ClientDto;
use crate::util::date_util;

/**
 * 服务端与客户端通信类
 * @param client 客户端DTO
 * @param clientSocket 与客户端的连接
 */
pub struct ClientSession {
    pub client: ClientDto,
    pub clientSocket: TcpStream,

    ///读数据
    pub(crate) reader:Mutex<OwnedReadHalf>,

    ///写数据
    pub(crate) writer:Mutex<OwnedWriteHalf>,

    /**
     * 最后一次收到客户端心跳时间
     */
    pub lastHeartBeatTime: u128,
}

/**
 * 开始
 */
pub async fn start(session: &Arc<ClientSession>) {
    let mut proxy_socket_clone1 = session.clone();
    tokio::spawn(async move {
        proxy_socket_clone1.await.receive().await;
    });
}

impl ClientSession {

    /**
     * 接收从客户端发来的数据
     */
    async fn receive(&mut self) {
        let mut reader = self.reader.lock().await;
        loop {

            // 创建一个长度为 1 的缓冲区
            let mut buf = [0; 1];

            //读取第一个标记字节,通过该自己判断该连接类型
            let lenResult = reader.read(&mut buf).await;
            if let Err(e) = lenResult {
                break;
            }
            if lenResult.unwrap() == 0 { //可能对方已经关闭
                break;
            }

            //读取到标记
            let flag = buf[0];
            self.handle(flag).await
        }

        //显示丢弃所有权
        drop(reader);
        self.clientSocket.shutdown().await;
    }

    /**
     * 处理从客户端收到的消息
     */
    async fn handle(&mut self, flag: u8) {
        match flag {

            //客户端心跳
            header_util::MAIN_HEART_BEAT => {
                //println("-->接收到客户端的心跳数据${Date()}")

                //记录与客户端最后一次心跳时间戳
                self.lastHeartBeatTime = date_util::timestamp();

                //回复客户端心跳
                self.clientSocket.write_i8(1).await;
            }
            _ => {}
        }
    }


    /**
     * 往客户端发送数据
     * @param flag 头部标记
     * @param message 头部消息
     */
    pub async fn sendHeader(&mut self, flag: i8, message: String) {
        if message.is_empty() {
            return;
        }

        // 要发送的消息
        let data_array = message.as_bytes();

        let mut writer = self.writer.lock().await;

        // 发送消息到服务器
        writer.write_all(data_array).await.expect("TODO: panic message");
        writer.flush().await.expect("TODO: panic message");

        if data_array.len() > 127 {
            //throw RuntimeException("一次发送数据长度不能超过${Byte.MAX_VALUE}字节")
            return;
        }
        writer.write_i8(flag).await;
        writer.write_i8(data_array.len() as i8).await;
        writer.write_all(data_array).await;
        writer.flush().await;
    }


    /**
     * 往客户端发送数据
     * @param data 要发送的数据
     * @param len 数据长度
     */
    pub async fn send(&mut self, data: &[u8], len: usize) {
        let mut writer = self.writer.lock().await;
        writer.write_all(&data[..len]).await;
        writer.flush().await;
    }


    /**
     * 关闭与内网穿透客户端的会话连接
     */
    pub async fn close(&mut self) {

        //关闭所有TCP连接池
        // try {
        // TCPPoolManager.closeByClient(this.client.id!!)
        // } catch (e: Exception) {
        // e.printStackTrace()
        // }
        //
        // //关闭所有UDP连接池
        // try {
        // UDPPoolManager.closeByClient(this.client.id!!)
        // } catch (e: Exception) {
        // e.printStackTrace()
        // }
        //
        // try {
        // //关闭正在通信的UDP连接
        // UDPBridgeManager.closeByClient(this.client.id!!)
        // } catch (e: Exception) {
        // e.printStackTrace()
        // }
        //
        // try {
        // //关闭代理监听
        // ProxyAcceptManager.closeByClient(this.client.id!!)
        // } catch (e: Exception) {
        // e.printStackTrace()
        // }
        //
        // try {
        // //关闭客户端所有正在通信的连接
        // TCPBridgeManager.closeByClient(this.client.id!!)
        // } catch (e: Exception) {
        // e.printStackTrace()
        // }

        //关闭连接
        self.clientSocket.shutdown().await;
    }
}

