use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use lazy_static::lazy_static;
use tokio::io::AsyncReadExt;
use tokio::net::{TcpListener, TcpStream};
use tokio::net::tcp::{OwnedReadHalf, OwnedWriteHalf};
use tokio::sync::{watch, Mutex};
use tokio::time::sleep;
use crate::{clientSessionMap, SessionInfo};
use crate::util::date_util;

pub struct ClientSignal {
    signal: watch::Sender<Arc<TcpStream>>, // 用于通知关闭任务的信号
}

lazy_static! {

/**
 * 客户端ID对应的Socket连接
 */
    pub static ref CLIENT_SESSION_MAP: Mutex<HashMap<i64, Arc<ClientSignal>>> = Mutex::new(HashMap::new());
}

pub async fn start(){

    // 绑定到指定地址和端口
    let tcp_listener = TcpListener::bind("0.0.0.0:3435").await.unwrap();
    loop {
        let (tcp_stream, _) = tcp_listener.accept().await.unwrap();
        tokio::spawn(async {
            handle_tcp(tcp_stream).await;
        });

        println!("-->reloop accept")
    }
}


///处理tcp连接
/// - tcp_stream
async fn handle_tcp(tcp_stream: TcpStream) {
    let tcp1 = Arc::new(tcp_stream);
    let tcp2 = tcp1.clone();
    let tcp3 = tcp1.clone();
    // let (r, mut w) = tcp1.into_split();
    let (signal, mut receiver) = watch::channel(tcp2);
    let cs = ClientSignal {
        signal
    };
    let csArc = Arc::new(cs);
    {
        let mut map = CLIENT_SESSION_MAP.lock().await;
        // if map.contains_key(&64) {
        //     let old = map.get_mut(&64).unwrap();
        //
        //     println!("-->001发送关闭信号");
        //     old.read_shutdown_signal.send(true);
        //     old.write_shutdown_signal.send(true);
        //     println!("-->002关闭信号完成");
        //     map.remove(&64);
        //     println!("-->003移除旧的session完成");
        // }
        println!("-->00添加新的session");
        map.insert(64, csArc);
    }

    tokio::spawn(async{
        let rs = receiver.changed();
        if let Ok(()) = rs{
            if *read_shutdown_receiver.borrow() {
                println!("接收到了关闭信号");
            }
        }
    });

    println!("-->11启动读数据任务");
    let read_task = tokio::spawn(async move {
        return crate::read_handle(r, read_shutdown_receiver).await;
    });
    println!("-->22启动写数据任务");
    let write_task = tokio::spawn(async move {
        return crate::write_handle(w, write_shutdown_receiver).await;
    });
    println!("-->33等待读写任务完成");

    //返回所有权
    let source_read = read_task.await;
    let source_write = write_task.await;
    println!("-->44读写任务完成，准备关闭");

    //最终关闭连接
    // source_read.unwrap().reunite(source_write.unwrap()).unwrap().shutdown().await;
    println!("-->55关闭了一个连接")
}

async fn read_handle(mut reader: OwnedReadHalf, mut shutdown_receiver: watch::Receiver<bool>) -> OwnedReadHalf {
    loop {

        // 创建一个长度为 head_len 的缓冲区
        let mut head_data_buf = [0u8; 1024];
        tokio::select! {

            //读取数据部分
            rs = reader.read(&mut head_data_buf) => {
                if let Ok(len) = rs{
                    println!("-->从客户端接收到数据:{}", String::from_utf8_lossy(&head_data_buf[..len]));
                }else{
                    break;
                }
            }

            // 如果收到关闭信号，则退出
            _ = shutdown_receiver.changed() => {
                if *shutdown_receiver.borrow() {
                    println!("接收到了关闭信号");
                break;
                }
            }
        }
    }
    reader
}

async fn write_handle(mut writer: OwnedWriteHalf, mut write_shutdown_receiver: watch::Receiver<bool>) -> OwnedWriteHalf {
    loop {
        println!("--->正在往客户端发送数据");

        let str = format!("{}:来自服务端的数据", date_util::timestamp());
        let send_data = str.as_bytes();
        tokio::select! {
            _ = writer.write(&send_data) =>{
            // _ = writer.write_u8(120) =>{
                writer.flush().await;
            }

            // 如果收到关闭信号，则退出
            _ = write_shutdown_receiver.changed() => {
                if *write_shutdown_receiver.borrow() {
                    writer.shutdown().await;
                    println!("接收到了关闭信号");
                    break;
                }
            }
        }
        sleep(Duration::from_secs_f32(1.0)).await;
    }
    println!("Shutting down connection with ");
    writer
}
