package router

// 路由文件
import (
	con "backend/controller"
	"backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	// 解决跨域问题
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 暴露的响应头
		AllowCredentials: true,                                                // 允许传递凭据（例如 Cookie）
		MaxAge:           12 * time.Hour,                                      // 预检请求的有效期
	}))
	// 设置静态文件目录
	r.Static("/static", "./dist/static")
	r.LoadHTMLGlob("dist/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// 测试GET请求
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//注册
	r.POST("/register", con.Register)
	//登录
	r.POST("/login", con.Login)
	//登出
	r.POST("/logout", con.Logout)
	//查询用户的类型
	r.POST("/getInfo", middleware.JWTAuthMiddleware(), con.GetInfo)
	//农产品上链
	r.POST("/uplink", middleware.JWTAuthMiddleware(), con.Uplink)
	// 获取农产品的上链信息
	r.POST("/getFruitInfo", con.GetFruitInfo)
	// 获取用户的农产品ID列表
	r.POST("/getFruitList", middleware.JWTAuthMiddleware(), con.GetFruitList)
	// 获取所有的农产品信息
	r.POST("/getAllFruitInfo", middleware.JWTAuthMiddleware(), con.GetAllFruitInfo)
	// 获取农产品上链历史(溯源)
	r.POST("/getFruitHistory", middleware.JWTAuthMiddleware(), con.GetFruitHistory)
	
	// 异常报警与自动召回相关接口
	// 异常检测
	r.POST("/detectAnomalies", middleware.JWTAuthMiddleware(), con.DetectAnomalies)
	// 创建异常报警
	r.POST("/createAlert", middleware.JWTAuthMiddleware(), con.CreateAlert)
	// 创建召回记录
	r.POST("/createRecall", middleware.JWTAuthMiddleware(), con.CreateRecall)
	// 更新报警状态
	r.POST("/updateAlertStatus", middleware.JWTAuthMiddleware(), con.UpdateAlertStatus)
	// 更新召回状态
	r.POST("/updateRecallStatus", middleware.JWTAuthMiddleware(), con.UpdateRecallStatus)
	// 获取产品的报警记录
	r.POST("/getFruitAlerts", con.GetFruitAlerts)
	// 获取产品的召回记录
	r.POST("/getFruitRecalls", con.GetFruitRecalls)
	// 获取所有待处理报警
	r.POST("/getPendingAlerts", middleware.JWTAuthMiddleware(), con.GetPendingAlerts)
	// 获取所有进行中的召回
	r.POST("/getActiveRecalls", middleware.JWTAuthMiddleware(), con.GetActiveRecalls)
	
	r.GET("/getImg/:filename", con.GetImg) // 获取图片
	return r
}
