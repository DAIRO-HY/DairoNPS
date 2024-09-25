use std::collections::HashMap;
use std::net::SocketAddr;
use std::sync::Arc;
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::{Mutex, watch};
use tokio::io::{AsyncReadExt, AsyncWriteExt, ReadHalf, WriteHalf};
use tokio::io::split;

type ClientMap = Arc<Mutex<HashMap<SocketAddr, Arc<ClientConnection>>>>;

// 客户端连接结构体，包含读写两部分
struct ClientConnection {
    read_half: Mutex<ReadHalf<TcpStream>>,
    write_half: Mutex<WriteHalf<TcpStream>>,
    shutdown_signal: watch::Sender<bool>, // 用于通知关闭任务的信号
}

impl ClientConnection {
    // 创建一个新的ClientConnection实例
    fn new(stream: TcpStream) -> (Self, watch::Receiver<bool>) {
        let (read_half, write_half) = split(stream);
        let (shutdown_signal, shutdown_receiver) = watch::channel(false);
        (
            ClientConnection {
                read_half: Mutex::new(read_half),
                write_half: Mutex::new(write_half),
                shutdown_signal,
            },
            shutdown_receiver,
        )
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let clients: ClientMap = Arc::new(Mutex::new(HashMap::new()));
    let listener = TcpListener::bind("127.0.0.1:8080").await?;

    loop {
        let (socket, addr) = listener.accept().await?;
        let clients = Arc::clone(&clients);

        tokio::spawn(async move {
            handle_client(socket, addr, clients).await;
        });
    }
}

async fn handle_client(socket: TcpStream, addr: SocketAddr, clients: ClientMap) {
    // 创建新的连接和关闭接收器
    let (client_conn, mut shutdown_receiver) = ClientConnection::new(socket);

    // 处理同一IP的旧连接
    {
        let mut clients_lock = clients.lock().await;
        if let Some(old_conn) = clients_lock.remove(&addr) {
            println!("Closing previous connection from {}", addr);
            // 发送关闭信号给旧连接的读取任务
            let _ = old_conn.shutdown_signal.send(true);
        }
        clients_lock.insert(addr, Arc::new(client_conn));
    }

    let clients = Arc::clone(&clients);
    let read_conn = Arc::clone(&clients.lock().await[&addr]);

    // 启动读取任务和写入任务
    let read_task = tokio::spawn(async move {
        handle_read(read_conn, addr, &mut shutdown_receiver, clients).await;
    });

    let write_conn = Arc::clone(&clients.lock().await[&addr]);
    let write_task = tokio::spawn(async move {
        handle_write(write_conn, addr).await;
    });

    let _ = tokio::join!(read_task, write_task);
}

async fn handle_read(
    client_conn: Arc<ClientConnection>,
    addr: SocketAddr,
    shutdown_receiver: &mut watch::Receiver<bool>,
    clients: ClientMap
) {
    let mut buffer = [0u8; 1024];

    loop {
        tokio::select! {
            // 尝试读取数据
            result = client_conn.read_half.lock().await.read(&mut buffer) => {
                match result {
                    Ok(0) => {
                        println!("Client {} disconnected", addr);
                        let mut clients_lock = clients.lock().await;
                        clients_lock.remove(&addr);
                        break;
                    }
                    Ok(n) => {
                        println!("Received from {}: {}", addr, String::from_utf8_lossy(&buffer[..n]));
                    }
                    Err(e) => {
                        eprintln!("Error reading from {}: {:?}", addr, e);
                        break;
                    }
                }
            }
            // 如果收到关闭信号，则退出
            _ = shutdown_receiver.changed() => {
                if *shutdown_receiver.borrow() {
                    println!("Shutting down connection with {}", addr);
                    break;
                }
            }
        }
    }
}

async fn handle_write(client_conn: Arc<ClientConnection>, addr: SocketAddr) {
    loop {
        let message = format!("Hello from server to {}\n", addr);
        let mut write_half = client_conn.write_half.lock().await;

        if let Err(e) = write_half.write_all(message.as_bytes()).await {
            eprintln!("Error sending to {}: {:?}", addr, e);
            break;
        }

        tokio::time::sleep(tokio::time::Duration::from_secs(5)).await;
    }
}
