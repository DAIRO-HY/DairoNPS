mod pipeline {
    pub mod pipeline_tcp_accept;
}
mod proxy {
    pub mod proxy_tcp_accept;
}
mod bridge {
    pub mod tcp_bridge;
}
mod client {
    pub mod client_accept_manager;
    pub mod client_session_manager;
    pub mod client_session;
    pub mod header_util;
}
mod util {
    pub mod date_util;
}
mod dao {
    pub mod client_dao;
    pub(crate) mod dto {
        pub mod client_dto;
        pub mod channel_dto;
    }
}

//如果想单线程执行加上参数flavor = "current_thread"
#[tokio::main(flavor = "current_thread")]
// #[tokio::main]
async fn main() {
    tokio::spawn(async {
        pipeline::pipeline_tcp_accept::start().await;         // 调用 sub_module 中的函数
    });
    tokio::spawn(async {
        proxy::proxy_tcp_accept::start().await;         // 调用 sub_module 中的函数
    });

    tokio::time::sleep(tokio::time::Duration::from_secs_f32(9999999999.0)).await;
}

// use std::sync::Arc;
// use tokio::time::Duration;
// use tokio::sync::Mutex;
// use tokio::io::AsyncWriteExt;
// use tokio::net::TcpListener;
// use tokio::time::sleep;
// use crate::client::header_util;
//
// mod client {
//     pub mod header_util;
// }
//
// #[tokio::main]
// async fn main() {
//
//     // 绑定到指定地址和端口
//     let tcp_listener = TcpListener::bind("0.0.0.0:3435").await.unwrap();
//     loop {
//         let (tcp_stream, _) = tcp_listener.accept().await.unwrap();
//         let arc_tcp1 =  Arc::new(Mutex::new(tcp_stream));
//         let arc_tcp2 =  arc_tcp1.clone();
//         tokio::spawn(async move{
//             loop{
//                 println!("--->正在往客户端发送数据");
//                 let mut writter = arc_tcp2.lock().await;
//                 writter.write("sadfsa".as_bytes()).await;
//                 writter.flush().await;
//                 sleep(Duration::from_secs_f32(1.0)).await
//             }
//         });
//         let sd = header_util::get_header(arc_tcp1).await;
//         if let Ok(wef) = sd{
//             println!("-->{}",wef)
//         }else if let Err(err) = sd{
//             println!("-->{}",err)
//         }
//     }
// }
//
