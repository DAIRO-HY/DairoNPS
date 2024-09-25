use std::collections::HashMap;
use std::sync::Arc;
use lazy_static::lazy_static;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::tcp::{OwnedReadHalf, OwnedWriteHalf};
use tokio::net::TcpStream;
use tokio::sync::Mutex;
use tokio::time::sleep;

lazy_static! {

/**
 * 客户端ID对应的Socket连接
 */
    pub static ref clientSessionMap: Mutex<HashMap<i64, (OwnedReadHalf,OwnedWriteHalf)>> = Mutex::new(HashMap::new());
}

#[tokio::main]
async fn main() {
    let tcp_stream = TcpStream::connect("127.0.0.1:3435").await.unwrap();
    {
        let ( reader,  writer) = tcp_stream.into_split();
        clientSessionMap.lock().await.insert(64, (reader, writer));
    }
    tokio::spawn(async move {
        sleep(tokio::time::Duration::from_secs_f32(2.0)).await;
        let mut buf = [0; 1024];
        let mut map = clientSessionMap.lock().await;
        let (reader,_) = map.get_mut(&64).unwrap();
        loop {
            let len = reader.read(&mut buf).await.unwrap();
            println!("-->接收到数据：{}", String::from_utf8_lossy(&buf[..len]))
        }
    });
    {
        let mut map = clientSessionMap.lock().await;
        let (_, writer) = map.get_mut(&64).unwrap();
        loop{
            println!("--->正在往服务端发送数据");
            let head = "1234567890";
            writer.write_i8((head.len()) as i8).await;
            writer.write_all(head.as_bytes()).await;
            writer.flush().await;
            sleep(tokio::time::Duration::from_secs_f32(1.0)).await;
        }
    }

    println!("-->finish");
    sleep(tokio::time::Duration::from_secs_f32(9999.0)).await;
}