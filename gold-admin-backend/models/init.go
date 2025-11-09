package models

import (
	"fmt"
	"gold-admin-backend/config"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	var err error
	var dialector gorm.Dialector

	cfg := config.AppConfig.Database

	// 根据配置选择数据库驱动
	switch cfg.Type {
	case "sqlite":
		// 确保数据库文件目录存在
		dbDir := filepath.Dir(cfg.Path)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("创建数据库目录失败: %w", err)
		}
		dialector = sqlite.Open(cfg.Path)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}

	// 打开数据库连接
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// 禁用外键约束（SQLite兼容性更好）
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// 自动迁移所有表
	log.Println("开始自动迁移数据库表...")
	err = DB.AutoMigrate(
		&AdminUser{},
		&Role{},
		&Menu{},
		&UserRole{},
		&RoleMenu{},
		&Price{},
		&Shop{},
		&Appointment{},
		&OperationLog{},
		// 后续添加其他模型...
	)
	if err != nil {
		return fmt.Errorf("自动迁移失败: %w", err)
	}

	log.Println("数据库初始化完成")
	return nil
}

// InitData 初始化基础数据
func InitData() error {
	// 检查是否已有管理员用户
	var count int64
	DB.Model(&AdminUser{}).Count(&count)
	if count > 0 {
		log.Println("数据库已有数据，跳过初始化")
		return nil
	}

	log.Println("开始初始化基础数据...")

	// 开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建默认角色
	roles := []Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "拥有所有权限",
			Sort:        1,
			Status:      1,
		},
		{
			Name:        "总部店长",
			Code:        "head_manager",
			Description: "管理所有店铺",
			Sort:        2,
			Status:      1,
		},
		{
			Name:        "单店店长",
			Code:        "shop_manager",
			Description: "管理单个店铺",
			Sort:        3,
			Status:      1,
		},
		{
			Name:        "店员",
			Code:        "shop_staff",
			Description: "处理日常业务（只读）",
			Sort:        4,
			Status:      1,
		},
		{
			Name:        "财务",
			Code:        "finance",
			Description: "查看数据、导出报表",
			Sort:        5,
			Status:      1,
		},
	}

	for _, role := range roles {
		if err := tx.Create(&role).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建角色失败: %w", err)
		}
	}
	log.Println("✓ 创建默认角色完成")

	// 2. 创建默认菜单
	menus := []Menu{
		// 首页
		{ID: 1, ParentID: 0, Type: 2, Name: "dashboard", Title: "首页", Icon: "el-icon-s-home", Path: "/dashboard", Component: "dashboard/index", Sort: 1, Visible: 1, Status: 1},

		// 业务管理
		{ID: 10, ParentID: 0, Type: 1, Name: "business", Title: "业务管理", Icon: "el-icon-s-management", Path: "/business", Component: "Layout", Sort: 10, Visible: 1, Status: 1},
		{ID: 11, ParentID: 10, Type: 2, Name: "price", Title: "价格管理", Icon: "el-icon-money", Path: "price", Component: "price/index", Sort: 1, Visible: 1, Status: 1},
		{ID: 12, ParentID: 10, Type: 2, Name: "appointment", Title: "预约管理", Icon: "el-icon-date", Path: "appointment", Component: "appointment/index", Sort: 2, Visible: 1, Status: 1},
		{ID: 13, ParentID: 10, Type: 2, Name: "shop", Title: "店铺管理", Icon: "el-icon-office-building", Path: "shop", Component: "shop/index", Sort: 3, Visible: 1, Status: 1},

		// 系统管理
		{ID: 20, ParentID: 0, Type: 1, Name: "system", Title: "系统管理", Icon: "el-icon-setting", Path: "/system", Component: "Layout", Sort: 20, Visible: 1, Status: 1},
		{ID: 21, ParentID: 20, Type: 2, Name: "user", Title: "用户管理", Icon: "el-icon-user", Path: "user", Component: "system/user/index", Sort: 1, Visible: 1, Status: 1},
		{ID: 22, ParentID: 20, Type: 2, Name: "role", Title: "角色管理", Icon: "el-icon-s-custom", Path: "role", Component: "system/role/index", Sort: 2, Visible: 1, Status: 1},
		{ID: 23, ParentID: 20, Type: 2, Name: "menu", Title: "菜单管理", Icon: "el-icon-menu", Path: "menu", Component: "system/menu/index", Sort: 3, Visible: 1, Status: 1},
	}

	for _, menu := range menus {
		if err := tx.Create(&menu).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建菜单失败: %w", err)
		}
	}
	log.Println("✓ 创建默认菜单完成")

	// 3. 为超级管理员角色分配所有菜单
	var menuIDs []uint
	tx.Model(&Menu{}).Pluck("id", &menuIDs)
	for _, menuID := range menuIDs {
		roleMenu := RoleMenu{
			RoleID: 1, // 超级管理员
			MenuID: menuID,
		}
		if err := tx.Create(&roleMenu).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("分配菜单失败: %w", err)
		}
	}
	log.Println("✓ 为超级管理员分配菜单完成")

	// 4. 创建默认管理员账号
	hashedPassword, err := HashPassword("admin123")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("密码加密失败: %w", err)
	}

	adminUser := AdminUser{
		Username: "admin",
		Password: hashedPassword,
		RealName: "系统管理员",
		Status:   1,
	}
	if err := tx.Create(&adminUser).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建管理员失败: %w", err)
	}
	log.Println("✓ 创建默认管理员完成")

	// 5. 为管理员分配超级管理员角色
	userRole := UserRole{
		UserID: adminUser.ID,
		RoleID: 1, // 超级管理员
	}
	if err := tx.Create(&userRole).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("分配角色失败: %w", err)
	}
	log.Println("✓ 为管理员分配角色完成")

	// 6. 创建默认店铺数据
	shops := []Shop{
		{
			Name:        "沪金汇总店",
			Code:        "shop1",
			Address:     "上海市黄浦区南京东路XXX号",
			Phone:       "021-12345678",
			Mobile:      "13800138000",
			Hours:       "周一至周日 09:00 - 18:00",
			Latitude:    31.2304,
			Longitude:   121.4737,
			Description: "沪金汇总店，提供专业的贵金属回收服务",
			Status:      1,
			Sort:        1,
		},
		{
			Name:        "沪金汇浦东店",
			Code:        "shop2",
			Address:     "上海市浦东新区陆家嘴XXX号",
			Phone:       "021-87654321",
			Mobile:      "13900139000",
			Hours:       "周一至周日 09:00 - 18:00",
			Latitude:    31.2397,
			Longitude:   121.4994,
			Description: "沪金汇浦东店，方便浦东地区客户",
			Status:      1,
			Sort:        2,
		},
	}

	for _, shop := range shops {
		if err := tx.Create(&shop).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建店铺数据失败: %w", err)
		}
	}
	log.Println("✓ 创建默认店铺数据完成")

	// 7. 创建默认价格数据（示例）
	prices := []Price{
		{
			Code:          "gold_9999",
			Name:          "黄金9999",
			Subtitle:      "Au9999 · 千足金",
			Icon:          "Au",
			IconColor:     "#FFD700",
			BasePrice:     560.00,
			BuyPriceDiff:  -10.00,
			SellPriceDiff: 15.00,
			Sort:          1,
			Status:        1,
		},
		{
			Code:          "gold_999",
			Name:          "黄金999",
			Subtitle:      "Au999 · 足金",
			Icon:          "Au",
			IconColor:     "#FFD700",
			BasePrice:     555.00,
			BuyPriceDiff:  -12.00,
			SellPriceDiff: 15.00,
			Sort:          2,
			Status:        1,
		},
		{
			Code:          "silver_999",
			Name:          "白银999",
			Subtitle:      "Ag999 · 纯银",
			Icon:          "Ag",
			IconColor:     "#C0C0C0",
			BasePrice:     6.50,
			BuyPriceDiff:  -0.50,
			SellPriceDiff: 1.00,
			Sort:          3,
			Status:        1,
		},
		{
			Code:          "london_gold",
			Name:          "伦敦金",
			Subtitle:      "London Gold · 国际现货",
			Icon:          "Au",
			IconColor:     "#FFD700",
			BasePrice:     2050.00,
			BuyPriceDiff:  -5.00,
			SellPriceDiff: 8.00,
			Sort:          12,
			Status:        1,
		},
		{
			Code:          "shanghai_gold",
			Name:          "上海金",
			Subtitle:      "SGE Gold · 上海金交所",
			Icon:          "Au",
			IconColor:     "#FFD700",
			BasePrice:     562.00,
			BuyPriceDiff:  -10.00,
			SellPriceDiff: 15.00,
			Sort:          13,
			Status:        1,
		},
		{
			Code:          "london_silver",
			Name:          "伦敦银",
			Subtitle:      "London Silver · 国际现货",
			Icon:          "Ag",
			IconColor:     "#C0C0C0",
			BasePrice:     24.50,
			BuyPriceDiff:  -0.30,
			SellPriceDiff: 0.50,
			Sort:          14,
			Status:        1,
		},
		{
			Code:          "shanghai_silver",
			Name:          "上海银",
			Subtitle:      "SGE Silver · 上海金交所",
			Icon:          "Ag",
			IconColor:     "#C0C0C0",
			BasePrice:     6.52,
			BuyPriceDiff:  -0.50,
			SellPriceDiff: 1.00,
			Sort:          15,
			Status:        1,
		},
		{
			Code:          "comex_gold",
			Name:          "纽约金",
			Subtitle:      "COMEX Gold · 期货",
			Icon:          "Au",
			IconColor:     "#FFD700",
			BasePrice:     2048.00,
			BuyPriceDiff:  -5.00,
			SellPriceDiff: 8.00,
			Sort:          17,
			Status:        1,
		},
		{
			Code:          "comex_silver",
			Name:          "纽约银",
			Subtitle:      "COMEX Silver · 期货",
			Icon:          "Ag",
			IconColor:     "#C0C0C0",
			BasePrice:     24.30,
			BuyPriceDiff:  -0.30,
			SellPriceDiff: 0.50,
			Sort:          18,
			Status:        1,
		},
		{
			Code:          "t_d_gold",
			Name:          "T+D黄金",
			Subtitle:      "T+D Gold · 延期交收",
			Icon:          "Au",
			IconColor:     "#FFD700",
			BasePrice:     561.50,
			BuyPriceDiff:  -10.00,
			SellPriceDiff: 15.00,
			Sort:          19,
			Status:        1,
		},
	}

	for _, price := range prices {
		if err := tx.Create(&price).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建价格数据失败: %w", err)
		}
	}
	log.Println("✓ 创建默认价格数据完成")

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	log.Println("========================================")
	log.Println("初始化数据完成！")
	log.Println("默认管理员账号：admin")
	log.Println("默认密码：admin123")
	log.Println("========================================")

	return nil
}

// HashPassword 加密密码（初始化时使用）
func HashPassword(password string) (string, error) {
	// 导入 bcrypt
	// import "golang.org/x/crypto/bcrypt"
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
