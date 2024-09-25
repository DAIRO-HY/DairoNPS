use std::collections::HashMap;
use std::time::{SystemTime, UNIX_EPOCH};
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
    let sdf = tokio::spawn(async{
        test1().await
    });
    let sdf2 = tokio::spawn(async{
        sleep(tokio::time::Duration::from_secs_f32(10.0)).await;
        test1().await
    });

    let _ = tokio::join!(sdf,sdf2);
}

async fn test1(){
    let tcp_stream = TcpStream::connect("127.0.0.1:3435").await.unwrap();
    {
        let ( reader,  writer) = tcp_stream.into_split();
        clientSessionMap.lock().await.insert(64, (reader, writer));
    }
    tokio::spawn(async move {
        sleep(tokio::time::Duration::from_secs_f32(2.0)).await;
        let mut buf = [0; 1024];
        loop {
            let mut map = clientSessionMap.lock().await;
            let (reader,_) = map.get_mut(&64).unwrap();
            let len = reader.read(&mut buf).await.unwrap();
            println!("-->接收到数据：{}", String::from_utf8_lossy(&buf[..len]))
        }
    });
    {
        loop{
            let mut map = clientSessionMap.lock().await;
            let (_, writer) = map.get_mut(&64).unwrap();
            // let head = "1234567890";
            let head = format!("now:{}", timestamp());
            writer.write_i8((head.len()) as i8).await;
            writer.write_all(head.as_bytes()).await;
            writer.flush().await;
            sleep(tokio::time::Duration::from_secs_f32(100.0)).await;
        }
    }

    println!("-->finish");
    sleep(tokio::time::Duration::from_secs_f32(9999.0)).await;
}

pub fn timestamp() -> u128 {
    // 获取当前时间
    let start = SystemTime::now();

    // 计算从 UNIX_EPOCH 到当前时间的持续时间
    let duration = start.duration_since(UNIX_EPOCH)
        .expect("时间回溯到 UNIX_EPOCH");

    // 获取时间戳（以秒为单位）
    let timestamp = duration.as_millis();
    timestamp
}