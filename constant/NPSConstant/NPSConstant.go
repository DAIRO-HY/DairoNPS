package NPSConstant

import (
	"math/rand/v2"
	"strconv"
)

// WEB管理端口
var WebPort = "1780"

// 服务端监听TCP端口,客户端通过此端口进行连接
var TcpPort = "1781"

// 服务端监听UDP端口,客户端通过此端口进行连接
var UDPPort = "1782"

// 关闭UDP连接池标记
const UDP_POOL_CLOSE_FLAG = ":[NPS-POOL-CLOSE]"

// 关闭UDP桥接标记
const UDP_BRIDIGE_CLOSE_FLAG = ":[NPS-BRIDGE-CLOSE]"

// 数据统计时间间隔（秒）
const STATISTICS_DATA_SIZE_TIMER = 60

/**
 * 因为UDP的不确定性,服务端无法检测存活状态,所以
 * 每个一段时间去检测过期的连接
 */
//const val RECYLE_UDP_TIME = 1 * 10 * 1000L
const RECYLE_UDP_TIME = 1 * 60 * 1000

// 每隔一段时间回收长时间不用的连接池（毫秒）
const RECYLE_POOL_TIME_OUT = 3 * 60 * 1000

//const RECYLE_POOL_TIME_OUT = 5 * 1000

// UDP桥接连接会话超时(毫秒)
const UDP_BRIDGE_TIMEOUT = 1 * 60 * 1000

//const UDP_BRIDGE_TIMEOUT = 5 * 1000

/**
 * 心跳间隔时间
 */
const HEART_TIME = 3000

/**
 * 读取数据缓存大小
 */
const READ_UDP_CACHE_SIZE = 1500

/**
 * 读取数据缓存大小
 */
const READ_CACHE_SIZE = 32 * 1024

/**
 * 连接池最大数量
 */
const MAX_POOL_COUNT = 6

/**
 * 连接池最低数量
 * 连接池中的Socket在一段时间内无任何操作
 */
const MIN_POOL_COUNT = 1

/**
 * 连接池不足时,一次性创建连接数
 */
const ADD_POOL_COUNT = 3

/**
 * 系统配置
 */
//var systemConfig = SystemConfigDao.SelectOne()

// 管理员用户名
var LoginName = "admin"

// 管理员登录密码 默认随机6位数
var LoginPwd = strconv.Itoa(rand.IntN(900000) + 100000)

// 是否开发模式
var IsDev bool
