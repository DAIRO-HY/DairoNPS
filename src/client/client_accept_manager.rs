use std::collections::HashMap;
use std::sync::Arc;
use lazy_static::lazy_static;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::Mutex;
use crate::client::header_util;

lazy_static! {
    pub static ref CLIENT_SOCKET_MAP: Arc<Mutex<HashMap<String, TcpStream>>> = Arc::new(Mutex::new(HashMap::new()));
}

///开始
pub async fn start() {
    tokio::spawn(async {
        accept()
    });
}

///监听客户端连接
async fn accept() {

    // 绑定到指定地址和端口
    let bind_result = TcpListener::bind("0.0.0.0:1681").await;
    if let Err(e) = bind_result {
        eprintln!("创建TcpListener失败: {}", e);
        return;
    }

    //等待客户端连接
    println!("-->监听客户端连接开始");
    let listener = bind_result.unwrap();
    loop {
        let accept_result = listener.accept().await;
        if let Err(e) = accept_result {
            eprintln!("监听客户端监听失败: {}", e);

            //tokio::time::sleep(tokio::time::Duration::from_secs_f32(10.0)).await;
            return;
        }
        println!("-->接收到客户端连接请求");
        let (tcp_stream, _) = accept_result.unwrap();

        tokio::spawn(async {
            handle(tcp_stream)
        });
    }

    // println!("New connection: {:?}", socket.peer_addr());



    // let mut global_stream = CLIENT_SOCKET_MAP.lock().await;
    // global_stream.insert(String::from("TT"), socket);


    println!("-->监听客户端连接结束")
    // //开启异步任务
    // tokio::spawn(async {
    //     let mut global_stream = PIPELINE_SOCKET.lock().await;
    //
    //     if let Some(ref mut stream) = *global_stream {
    //         handle(stream).await;
    //     } else {
    //         println!("TcpStream is not set");
    //     }
    // });
    // }
}


/**
 * 分配连接
 * @param socketClient 与客户端的连接
 */
async fn handle(mut tcp_stream: TcpStream){

//保持长连接
// socketClient.keepAlive = true

//得到输入流
// val clientIStream = socketClient.inputStream
//
// //读取连接的第一个数据,设置超时,避免恶意连接
// socketClient.soTimeout = 5000

//println("-->读取标记开始${System.currentTimeMillis()}")

    // 创建一个长度为 1 的缓冲区
    let mut buf = [0; 1];

    //读取第一个标记字节,通过该自己判断该连接类型
    tcp_stream.read(&mut buf).await;

    let flag = buf[0];
    match flag{

        //标记该连接为:服务器端往客户端发送指令的连接
        header_util::CLIENT_TO_SERVER_MAIN_CONNECTION => {

        }

        //创建客户端Socket连接池
        2 => {
            // TCPPoolManager.add(socketClient)
        }
        _ => {
            tcp_stream.shutdown().await;
        }
    }
}

// ///等待客户端连接
// async fn handle(socket: &mut TcpStream) {
//     loop {
//
//         // // 创建一个缓冲区来读取数据
//         // let mut buffer = [0; 8 * 1024];
//         //
//         // // 从客户端读取数据
//         // let len_result = socket.read(&mut buffer).await;
//         // if let Err(e) = len_result {
//         //     eprintln!("读取数据失败: {}", e);
//         //     break;
//         // }
//         // let len = len_result.unwrap();
//         // if len == 0{//客户端已经关闭了
//         //     break;
//         // }
//         //
//         // // 打印接收到的数据
//         // println!("Received: {}", String::from_utf8_lossy(&buffer[..len]));
//
//         tokio::time::sleep(tokio::time::Duration::from_secs_f32(2.0)).await;
//
//         // 回应客户端
//         let write_result = socket.write(b"Hello from server!").await;
//         if let Err(e) = write_result {
//             eprintln!("写入数据失败: {}", e);
//             break;
//         }
//     }
// }