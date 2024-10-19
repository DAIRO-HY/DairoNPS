package cn.dairo.cls.forward

import cn.dairo.cls.dao.ForwardDao
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

/**
 * 数据转发流量统计工具类
 */
object ForwardStatisticsUtil {

    /**
     * 统计时间间隔
     */
    private val STATISTICS_TIME = 5 * 60 * 1000L
    fun start() = GlobalScope.launch {
        while (true) {
            delay(STATISTICS_TIME)

            //保存代理服务流量
            ForwardAcceptManager.forwardIdToForwardAccept.forEach { _, v ->

                //当前正在会话的代理服务
                ForwardDao.setDataLen(v.forwardDto)
            }
        }
    }
}