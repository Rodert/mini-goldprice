package main

import (
	"gold-admin-backend/config"
	"gold-admin-backend/models"
	"gold-admin-backend/utils"
	"log"
	"time"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("./config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建超级管理员角色
	role := &models.Role{
		Name:        "超级管理员",
		Code:        "super_admin",
		Description: "拥有所有权限",
		Sort:        1,
		Status:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 检查角色是否已存在
	var existingRole models.Role
	err := models.DB.Where("code = ?", role.Code).First(&existingRole).Error
	if err == nil {
		log.Printf("角色 [%s] 已存在，跳过创建", role.Name)
	} else {
		if err := models.DB.Create(role).Error; err != nil {
			log.Fatalf("创建角色失败: %v", err)
		}
		log.Printf("✓ 创建角色: %s", role.Name)
		existingRole = *role
	}

	// 创建管理员账号
	hashedPassword, _ := utils.HashPassword("admin123")
	admin := &models.AdminUser{
		Username:  "admin",
		Password:  hashedPassword,
		RealName:  "系统管理员",
		Phone:     "",
		Email:     "",
		Avatar:    "",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 检查管理员是否已存在
	var existingAdmin models.AdminUser
	err = models.DB.Where("username = ?", admin.Username).First(&existingAdmin).Error
	if err == nil {
		log.Printf("管理员 [%s] 已存在，跳过创建", admin.Username)
	} else {
		if err := models.DB.Create(admin).Error; err != nil {
			log.Fatalf("创建管理员失败: %v", err)
		}
		log.Printf("✓ 创建管理员: %s", admin.Username)
		existingAdmin = *admin
	}

	// 分配角色
	userRole := &models.UserRole{
		UserID:    existingAdmin.ID,
		RoleID:    existingRole.ID,
		CreatedAt: time.Now(),
	}

	// 检查是否已分配角色
	var existingUserRole models.UserRole
	err = models.DB.Where("user_id = ? AND role_id = ?", userRole.UserID, userRole.RoleID).First(&existingUserRole).Error
	if err == nil {
		log.Printf("用户角色关系已存在，跳过创建")
	} else {
		if err := models.DB.Create(userRole).Error; err != nil {
			log.Fatalf("分配角色失败: %v", err)
		}
		log.Printf("✓ 分配角色成功")
	}

	// 创建基础菜单
	createMenus()

	log.Println("\n========================================")
	log.Println("✓ 初始化完成！")
	log.Println("========================================")
	log.Println("管理员账号: admin")
	log.Println("管理员密码: admin123")
	log.Println("========================================")
}

// createMenus 创建基础菜单
func createMenus() {
	menus := []models.Menu{
		// 1. 系统管理（目录）
		{ParentID: 0, Type: 1, Name: "system", Title: "系统管理", Icon: "el-icon-setting", Path: "/system", Component: "Layout", Sort: 1, Visible: 1, Status: 1},

		// 1.1 用户管理（菜单）
		{ParentID: 1, Type: 2, Name: "user", Title: "用户管理", Icon: "el-icon-user", Path: "user", Component: "system/user/index", Permission: "system:user:list", Sort: 1, Visible: 1, Status: 1},

		// 1.2 角色管理（菜单）
		{ParentID: 1, Type: 2, Name: "role", Title: "角色管理", Icon: "el-icon-s-custom", Path: "role", Component: "system/role/index", Permission: "system:role:list", Sort: 2, Visible: 1, Status: 1},

		// 1.3 菜单管理（菜单）
		{ParentID: 1, Type: 2, Name: "menu", Title: "菜单管理", Icon: "el-icon-menu", Path: "menu", Component: "system/menu/index", Permission: "system:menu:list", Sort: 3, Visible: 1, Status: 1},

		// 2. 价格管理（目录）
		{ParentID: 0, Type: 1, Name: "price", Title: "价格管理", Icon: "el-icon-s-finance", Path: "/price", Component: "Layout", Sort: 2, Visible: 1, Status: 1},

		// 2.1 价格列表（菜单）
		{ParentID: 5, Type: 2, Name: "price-list", Title: "价格列表", Icon: "el-icon-price-tag", Path: "list", Component: "price/index", Permission: "price:list", Sort: 1, Visible: 1, Status: 1},
	}

	for i, menu := range menus {
		var existing models.Menu
		err := models.DB.Where("name = ? AND parent_id = ?", menu.Name, menu.ParentID).First(&existing).Error
		if err == nil {
			log.Printf("菜单 [%s] 已存在，跳过创建", menu.Title)
			continue
		}

		menu.CreatedAt = time.Now()
		menu.UpdatedAt = time.Now()

		if err := models.DB.Create(&menu).Error; err != nil {
			log.Printf("创建菜单 [%s] 失败: %v", menu.Title, err)
			continue
		}

		// 更新 menus 中的 ID（用于后续的 ParentID 引用）
		menus[i].ID = menu.ID
		log.Printf("✓ 创建菜单: %s", menu.Title)
	}

	// 为超级管理员分配所有菜单权限
	var role models.Role
	if err := models.DB.Where("code = ?", "super_admin").First(&role).Error; err != nil {
		log.Printf("查询超级管理员角色失败: %v", err)
		return
	}

	var allMenus []models.Menu
	models.DB.Find(&allMenus)

	for _, menu := range allMenus {
		var existing models.RoleMenu
		err := models.DB.Where("role_id = ? AND menu_id = ?", role.ID, menu.ID).First(&existing).Error
		if err == nil {
			continue
		}

		roleMenu := &models.RoleMenu{
			RoleID:    role.ID,
			MenuID:    menu.ID,
			CreatedAt: time.Now(),
		}

		if err := models.DB.Create(roleMenu).Error; err != nil {
			log.Printf("分配菜单权限失败: %v", err)
		}
	}

	log.Printf("✓ 分配菜单权限成功")
}



