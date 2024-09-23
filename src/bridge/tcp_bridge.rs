use std::sync::Arc;
use rand::Rng;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpStream};
use tokio::sync::Mutex;
use crate::util::date_util;
use crate::dao::dto::channel_dto::ChannelDto;
use crate::dao::dto::client_dto::ClientDto;
use crate::pipeline::pipeline_tcp_accept::CLIENT_SOCKET_MAP;

#[derive(Clone, Copy)]
pub struct TcpBridge {
    /**
     * 连接关闭同步锁
     */
    // private val closeLock = Mutex()

    /**
     * 最后一次读取到数据的时间,用来判断Socket是否存活
     */
    last_session_time: u128,

    /**
     * 是否加密数据
     */
    is_encode_data: bool,

    /**
     * 本次连接入网总计
     */
    in_data_total: i64,

    /**
     * 本次连接出网总计
     */
    out_data_total: i64,

    /**
     * 代理连接入方向是否被关闭
     */
    proxy_in_is_closed: bool,

    /**
     * 客户端连接入方向是否被关闭
     */
    client_in_is_closed: bool,

    tag: u128, // = System.currentTimeMillis().toString() + (Math.random() * 1000).toInt()
}

impl TcpBridge {
    pub fn new(is_encode_data: bool) -> Self {
        let mut rng = rand::thread_rng();
        let random_number: u32 = rng.gen_range(100..1000); // 生成100到999之间的随机整数
        TcpBridge {
            last_session_time: date_util::timestamp(),
            is_encode_data,
            in_data_total: 0,
            out_data_total: 0,
            proxy_in_is_closed: false,
            client_in_is_closed: false,
            tag: date_util::timestamp() + random_number as u128,
        }
    }
}

// pub async fn bridge(mut proxy_socket: TcpStream, pipeline_socket: TcpStream) {
//     loop {
//
//         // 创建一个缓冲区来读取数据
//         let mut buffer = [0; 8 * 1024];
//
//         // 从客户端读取数据
//         let len_result = proxy_socket.read(&mut buffer).await;
//         if let Err(e) = len_result {
//             eprintln!("读取数据失败: {}", e);
//             break;
//         }
//         let len = len_result.unwrap();
//         if len == 0 { //客户端已经关闭了
//             break;
//         }
//
//         // 打印接收到的数据
//         println!("Received: {}", String::from_utf8_lossy(&buffer[..len]));
//         pipeline_socket.write(&buffer[..len]).await.unwrap();
//         tokio::time::sleep(tokio::time::Duration::from_secs_f32(1.0)).await;
//     }
// }


/**
 * 开始传输数据
 */
pub async fn start(bridge: TcpBridge, client: ClientDto,
                   channel: ChannelDto,
                   proxy_socket: Arc<Mutex<TcpStream>>,
                   client_socket: Arc<Mutex<TcpStream>>,
) {
    println!("{}", client.key);

    // 使用 Arc 和 Mutex 以便共享 TcpStream
    let proxy_socket_clone1 = proxy_socket.clone();

    // 使用 Arc 和 Mutex 以便共享 TcpStream
    let proxy_socket_clone2 = proxy_socket.clone();


    // 使用 Arc 和 Mutex 以便共享 TcpStream
    let client_socket_clone1 = client_socket.clone();
    let client_socket_clone2 = client_socket.clone();
    let client_socket_clone3 = client_socket.clone();

    // let client_socket2 = &client_socket;
    // let client_socket3 = &client_socket;
    // let client_socket4 = &client_socket;

    tokio::spawn(async move {

        //发送目标端口信息
        send_header_to_client(channel, client_socket_clone1).await;
        receive_by_proxy_send_to_client(bridge, Arc::clone(&proxy_socket_clone1), client_socket_clone2).await;
    });
    tokio::spawn(async move{
        receive_by_client_send_to_proxy(bridge, Arc::clone(&proxy_socket_clone2), client_socket_clone3).await;
    });
}

/**
 * 发送目标端口信息
 */
async fn send_header_to_client(channel: ChannelDto, mut client_socket1: Arc<Mutex<TcpStream>>) {

    //将加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
    //1:加密  0:不加密
    let header = channel.security_state.to_string() + "|" + &*channel.target_port.to_string();

    // 要发送的消息
    let header_data = header.as_bytes();


    // 获取锁并保存到局部变量
    let mut map = CLIENT_SOCKET_MAP.lock().await;

    // 访问条目
    if let Some(stream) = map.get_mut("TT") { // 使用 get_mut 以获取可变引用
        // 使用 stream
        stream.write_all(header_data).await.expect("Failed to write to stream");
        stream.flush().await.expect("Failed to flush stream");
    } else {
        println!("No stream found for TT");
    }
    // client_socket.lock().await.write_all(header_data).await.expect("TODO: panic message");
    // client_socket.lock().await.flush().await.expect("TODO: panic message");
}

/**
 * 从代理服务接收数据发送到客户端
 */
async fn receive_by_proxy_send_to_client(mut bridge: TcpBridge,
                                         proxy_socket: Arc<Mutex<TcpStream>>,
                                         client_socket: Arc<Mutex<TcpStream>>) {


    // 创建一个缓冲区来读取响应
    let mut buffer = [0; 8 * 1024];
    loop {
        let len = proxy_socket.lock().await.read(&mut buffer).await.unwrap();
        if len == 0 {
            break;
        }

        //标记最后一次读取到数据的时间,TODO:频繁设置事件可能影响性能,待整改
        bridge.last_session_time = date_util::timestamp();


        //TODO:
        // //入网统计
        // this.in_data_total = this.in_data_total.plus(len)
        // this.channel.in_data_total = this.channel.in_data_total!!.plus(len)
        // this.client.in_data_total = this.client.in_data_total!!.plus(len)
        // CLSConfig.systemConfig.in_data_total = CLSConfig.systemConfig.in_data_total!!.plus(len)

        if bridge.is_encode_data { //加密数据TODO:

        }

        //将从代理服务器读取到的数据写入到隧道客户端
        client_socket.lock().await.write(&buffer[..len]).await.unwrap();
        client_socket.lock().await.flush().await.unwrap();
    }
    // this.proxy_socket.shutdownInput()
    // this.client_socket.shutdownOutput()

    //标记代理连接入方向是否被关闭
    bridge.proxy_in_is_closed = true;
    recyle(proxy_socket, client_socket).await;
}

/**
 * 从客户端接收发送到代理服务器
 */
async fn receive_by_client_send_to_proxy(mut bridge: TcpBridge,
                                         proxy_socket: Arc<Mutex<TcpStream>>,
                                         client_socket: Arc<Mutex<TcpStream>>) {

    // 创建一个缓冲区来读取响应
    let mut buffer = [0; 8 * 1024];
    loop {
        let len = client_socket.lock().await.read(&mut buffer).await.unwrap();
        if len == 0 {
            break;
        }

        //标记最后一次读取到数据的时间,TODO:频繁设置事件可能影响性能,待整改
        bridge.last_session_time = date_util::timestamp();


        //TODO:
        //出网统计
        // this.out_data_total = this.out_data_total.plus(len)
        // this.channel.out_data_total = this.channel.out_data_total!!.plus(len)
        // this.client.out_data_total = this.client.out_data_total!!.plus(len)
        // CLSConfig.systemConfig.out_data_total = CLSConfig.systemConfig.out_data_total!!.plus(len)

        if bridge.is_encode_data { //解密数据 TODO:

        }

        //将从代理服务器读取到的数据写入到隧道客户端
        proxy_socket.lock().await.write(&buffer[..len]).await.expect("TODO: panic message");
        proxy_socket.lock().await.flush().await.expect("TODO: panic message");
    }


    // client_socket.shutdownInput()
    // proxy_socket.shutdownOutput()

    //标记客户端连接入方向是否被关闭
    bridge.client_in_is_closed = true;
    recyle(proxy_socket, client_socket).await
}

/**
 * 资源回收
 */
async fn recyle(
    proxy_socket: Arc<Mutex<TcpStream>>,
    client_socket: Arc<Mutex<TcpStream>>) {
    proxy_socket.lock().await.shutdown().await.expect("TODO: panic message");
    client_socket.lock().await.shutdown().await.expect("TODO: panic message");
    // TCPBridgeManager.removeBridgeList(this)
}

/**
 * 关闭连接
 */
async fn close(proxy_socket: Arc<Mutex<TcpStream>>,
               client_socket: Arc<Mutex<TcpStream>>) {
    proxy_socket.lock().await.shutdown().await.expect("TODO: panic message");
    client_socket.lock().await.shutdown().await.expect("TODO: panic message");
}