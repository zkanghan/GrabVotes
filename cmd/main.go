package main

import (
	"GrabVotes/internal/controller"
	"GrabVotes/internal/dao/mysql"
	"GrabVotes/internal/dao/redis"
	"GrabVotes/internal/logic"
	"GrabVotes/internal/pkg/snowid"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := mysql.InitMysql(); err != nil {
		panic("MySQL初始化失败：" + err.Error())
	}

	if err := redis.InitRedis(); err != nil {
		panic("Redis初始化失败：" + err.Error())
	}

	if err := snowid.Init(); err != nil {
		panic("snowflake初始化失败：" + err.Error())
	}

	// 初始化Redis缓存中的票数
	if err := logic.SetRedisTicketNum("zkh_mirror", 5000); err != nil {
		panic("初始化Redis缓存失败：" + err.Error())
	}
}

func main() {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	// 抢购订单
	r.POST("/GrabAction", controller.JWTAuth, controller.GrabAction)
	// 支付订单
	r.POST("/defrayAction", controller.JWTAuth, controller.DefrayAction)

	go func() {
		mq, err := logic.NewRabbitMQSimple(logic.MqName)
		if err != nil {
			panic("消息队列初始化失败：" + err.Error())
		}
		mq.ConsumeSimple() //simple模式下，mq中有消息会立即被处理
	}()

	err := r.Run(":8080")
	if err != nil {
		panic("启动异常：" + err.Error())
	}
}