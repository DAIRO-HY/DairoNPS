// 声明 my_module 模块
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
}
mod util {
    pub mod date_util;
}
mod dao {
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

// mod dao {
//     pub(crate) mod dto {
//         pub mod client_dto;
//     }
// }
// use std::collections::HashMap;
// use std::sync::Arc;
// use lazy_static::lazy_static;
// use tokio::sync::RwLock;
// use crate::dao::dto::client_dto::ClientDto;
//
// lazy_static! {
//
// /**
//  * 客户端ID对应的Socket连接
//  */
//     pub static ref clientSessionMap: RwLock<HashMap<u64, Arc<ClientDto>>> = RwLock::new(HashMap::new());
// }
//
// #[tokio::main]
// async fn main() {
//     let client_dto = ClientDto {
//         // id
//         id: 12,
//
//         // 名称
//         name: String::from("String"),
//
//         // 版本号
//         version: String::from("String"),
//
//         // 连接认证秘钥
//         key: String::from("String"),
//
//         // ip地址
//         ip: String::from("String"),
//
//         // 入网流量
//         in_data_total: 12,
//
//         // 出网流量
//         out_data_total: 12,
//
//         // 在线状态,0:离线 1:在线
//         online_state: 12,
//
//         // 启用状态
//         enable_state: 12,
//
//         // 最后一次连接时间
//         last_login_date: 12,
//
//         // 创建时间
//         create_date: 12,
//
//         // 最后一次更新时间戳
//         update_date: 12,
//
//         // 一些备注信息,错误信息等
//         remark: String::from("String"),
//     };
//
//     {
//         let mut map = clientSessionMap.write().await;
//         map.insert(23, Arc::new(client_dto));
//     }
//
//     let size = size().await;
//     println!("-->{}", size);
// }
//
//
// /**
//  * 当前客户端数量
//  */
// pub async fn size() -> usize {
//     clientSessionMap.read().await.len()
// }