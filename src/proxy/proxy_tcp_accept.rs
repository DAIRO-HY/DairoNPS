use std::collections::HashMap;
use std::sync::Arc;
use lazy_static::lazy_static;
// use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::Mutex;
// use crate::pipeline::pipeline_tcp_accept;
use crate::bridge::tcp_bridge;
use crate::bridge::tcp_bridge::TcpBridge;
use crate::dao::dto::channel_dto::ChannelDto;
use crate::dao::dto::client_dto::ClientDto;


lazy_static! {
    pub static ref PROXY_SOCKET_MAP: Arc<Mutex<HashMap<String, TcpStream>>> = Arc::new(Mutex::new(HashMap::new()));
}

///等待客户端连接
pub async fn start() {

    // 绑定到指定地址和端口
    let bind_result = TcpListener::bind("0.0.0.0:3435").await;
    if let Err(e) = bind_result {
        eprintln!("创建TcpListener失败: {}", e);
        return;
    }
    println!("Server listening on port 3435");
    let listener = bind_result.unwrap();
    loop {
        let accept_result = listener.accept().await;
        if let Err(e) = accept_result {
            eprintln!("接受连接失败: {}", e);
            break;
        }
        let (proxy_socket, _) = accept_result.unwrap();
        println!("New connection: {:?}", proxy_socket.peer_addr());


        // let mut global_stream = PROXY_SOCKET_MAP.lock().await;
        // global_stream.insert(String::from("PROXY"),proxy_socket);

        //开启异步任务
        tokio::spawn(async {
            // handle(socket).await;

            let bridge = TcpBridge::new(false);
            let client_dto = ClientDto {
                // id
                id: 12,

                // 名称
                name: String::from("String"),

                // 版本号
                version: String::from("String"),

                // 连接认证秘钥
                key: String::from("String"),

                // ip地址
                ip: String::from("String"),

                // 入网流量
                in_data_total: 12,

                // 出网流量
                out_data_total: 12,

                // 在线状态,0:离线 1:在线
                online_state: 12,

                // 启用状态
                enable_state: 12,

                // 最后一次连接时间
                last_login_date: 12,

                // 创建时间
                create_date: 12,

                // 最后一次更新时间戳
                update_date: 12,

                // 一些备注信息,错误信息等
                remark: String::from("String"),
            };

            let channel_dto = ChannelDto {
                // 隧道ID
                id: 34,

                // 客户端id
                client_id: 345,

                // 隧道名
                name: String::from("String"),

                // 隧道模式, 1:TCP  2:UDP
                mode: 34,

                // 服务端端口
                server_port: 12,

                // 目标端口(ip:端口)
                target_port: String::from("String"),

                // 入网流量
                in_data_total: 12,

                // 出网流量
                out_data_total: 12,

                // 启用状态 1:开启  0:停止
                enable_state: 12,

                // 是否加密传输
                security_state: 12,

                // 黑白名单开启状态 0:关闭 1:白名单 2:黑名单
                acl_state: 12,

                // 创建时间
                create_date: 12,

                // 最后一次更新时间
                update_date: 12,

                // 一些备注信息,错误信息等
                remark: String::from("String"),
            };

            let proxy_socket_clone1 = Arc::new(Mutex::new(proxy_socket));
            let proxy_socket_clone2 = proxy_socket_clone1.clone();

            // 将可变引用转换为不可变引用

            tcp_bridge::start(bridge, client_dto, channel_dto, proxy_socket_clone1, proxy_socket_clone2).await;
        });
    }
}

// ///等待客户端连接
// async fn handle(mut socket: TcpStream) {
//     loop {
//
//         // 创建一个缓冲区来读取数据
//         let mut buffer = [0; 8 * 1024];
//
//         // 从客户端读取数据
//         let len_result = socket.read(&mut buffer).await;
//         if let Err(e) = len_result {
//             eprintln!("读取数据失败: {}", e);
//             break;
//         }
//         let len = len_result.unwrap();
//         if len == 0{//客户端已经关闭了
//             break;
//         }
//
//         // 打印接收到的数据
//         println!("Received: {}", String::from_utf8_lossy(&buffer[..len]));
//
//         tokio::time::sleep(tokio::time::Duration::from_secs_f32(5.0)).await;
//
//         // 回应客户端
//         let write_result = socket.write(b"Hello from server!").await;
//         if let Err(e) = write_result {
//             eprintln!("写入数据失败: {}", e);
//             break;
//         }
//     }
// }