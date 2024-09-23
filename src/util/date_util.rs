use std::time::{SystemTime, UNIX_EPOCH};

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