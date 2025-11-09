package router

import (
	v1 "gold-admin-backend/api/v1"
	"gold-admin-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.Cors())

	// 操作日志中间件
	r.Use(middleware.OperationLogMiddleware())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API v1 路由组
	api := r.Group("/api")
	{
		// 公开路由（无需认证）
		api.POST("/login", v1.Login)

		openapi := api.Group("openapi")
		{
			openapi.GET("/prices", v1.GetPriceList)
		}

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// 认证相关
			auth.POST("/logout", v1.Logout)
			auth.GET("/user/info", v1.GetUserInfo)
			auth.GET("/user/menus", v1.GetUserMenus)

			// 首页看板
			auth.GET("/dashboard/stats", v1.GetDashboardStats)
			auth.GET("/dashboard/activities", v1.GetRecentActivities)
			auth.GET("/dashboard/trend", v1.GetAppointmentTrend)

			// 用户管理
			auth.GET("/users", v1.GetUserList)
			auth.GET("/users/:id", v1.GetUser)
			auth.POST("/users", v1.CreateUser)
			auth.PUT("/users/:id", v1.UpdateUser)
			auth.DELETE("/users/:id", v1.DeleteUser)
			auth.PUT("/users/:id/password", v1.UpdateUserPassword)

			// 角色管理
			auth.GET("/roles", v1.GetRoleList)
			auth.GET("/roles/all", v1.GetAllRoles)
			auth.GET("/roles/:id", v1.GetRole)
			auth.POST("/roles", v1.CreateRole)
			auth.PUT("/roles/:id", v1.UpdateRole)
			auth.DELETE("/roles/:id", v1.DeleteRole)

			// 菜单管理
			auth.GET("/menus", v1.GetMenuList)
			auth.GET("/menus/tree", v1.GetMenuTree)
			auth.GET("/menus/:id", v1.GetMenu)
			auth.POST("/menus", v1.CreateMenu)
			auth.PUT("/menus/:id", v1.UpdateMenu)
			auth.DELETE("/menus/:id", v1.DeleteMenu)

			// 价格管理
			auth.GET("/prices", v1.GetPriceList)
			auth.GET("/prices/:id", v1.GetPrice)
			auth.POST("/prices", v1.CreatePrice)
			auth.PUT("/prices/:id", v1.UpdatePrice)
			auth.DELETE("/prices/:id", v1.DeletePrice)
			auth.POST("/prices/sync", v1.SyncBasePrice)

			// 店铺管理
			auth.GET("/shops", v1.GetShopList)
			auth.GET("/shops/all", v1.GetAllShops)
			auth.GET("/shops/:id", v1.GetShop)
			auth.POST("/shops", v1.CreateShop)
			auth.PUT("/shops/:id", v1.UpdateShop)
			auth.DELETE("/shops/:id", v1.DeleteShop)

			// 预约管理
			auth.GET("/appointments", v1.GetAppointmentList)
			auth.GET("/appointments/stats", v1.GetAppointmentStats)
			auth.GET("/appointments/:id", v1.GetAppointment)
			auth.POST("/appointments", v1.CreateAppointment)
			auth.PUT("/appointments/:id", v1.UpdateAppointment)
			auth.DELETE("/appointments/:id", v1.DeleteAppointment)

			// 操作日志
			auth.GET("/logs", v1.GetOperationLogList)
			auth.GET("/logs/stats", v1.GetOperationLogStats)
			auth.GET("/logs/export", v1.ExportOperationLogs)
			auth.GET("/logs/:id", v1.GetOperationLog)
			auth.DELETE("/logs/:id", v1.DeleteOperationLog)
			auth.POST("/logs/clear", v1.ClearOperationLogs)
		}
	}

	return r
}
