# SQLite 迁移到 MySQL 指南

## 为什么先用 SQLite？

### SQLite 的优势（开发阶段）
- ✅ **零配置**：无需安装数据库服务
- ✅ **单文件存储**：`gold_admin.db` 一个文件
- ✅ **开发便捷**：修改、备份、恢复都很简单
- ✅ **便携性**：复制文件即可迁移
- ✅ **演示友好**：客户演示不需要配置数据库

### MySQL 的优势（生产阶段）
- ✅ **高并发**：支持更多并发连接
- ✅ **大数据量**：处理百万级数据性能更好
- ✅ **成熟稳定**：企业级应用首选
- ✅ **工具丰富**：Navicat、PHPMyAdmin等
- ✅ **备份恢复**：专业的备份方案

## 何时迁移？

### 建议保持 SQLite 的场景
- 单店铺小规模应用（< 1000条记录/天）
- 演示系统
- 内部管理工具
- 并发用户 < 10

### 建议迁移到 MySQL 的场景
- 多店铺连锁（10家以上）
- 数据量大（> 10万条记录）
- 高并发访问（> 50并发用户）
- 需要主从备份
- 需要读写分离

## 迁移方法

### 方法一：GORM 自动迁移（推荐）⭐

**步骤：**

1. **安装 MySQL 驱动**
```bash
go get gorm.io/driver/mysql
```

2. **修改配置文件**
```yaml
# config/config.yaml
database:
  type: mysql                    # 改为 mysql
  host: 127.0.0.1
  port: 3306
  username: root
  password: your_password
  database: gold_admin
  charset: utf8mb4
  # 删除或注释 SQLite 配置
  # path: ./data/gold_admin.db
```

3. **修改数据库初始化代码**
```go
// models/init.go
package models

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func InitDB(config DBConfig) error {
    var err error
    var dialector gorm.Dialector
    
    if config.Type == "mysql" {
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
            config.Username,
            config.Password,
            config.Host,
            config.Port,
            config.Database,
            config.Charset,
        )
        dialector = mysql.Open(dsn)
    } else {
        dialector = sqlite.Open(config.Path)
    }
    
    DB, err = gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        return err
    }
    
    // 自动迁移（会在 MySQL 中创建所有表）
    return DB.AutoMigrate(
        &AdminUser{},
        &Role{},
        // ... 其他模型
    )
}
```

4. **创建 MySQL 数据库**
```sql
CREATE DATABASE gold_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. **运行程序（自动创建表结构）**
```bash
go run main.go
```

6. **导入 SQLite 数据到 MySQL（可选）**

使用迁移脚本：

```go
// tools/migrate_sqlite_to_mysql.go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    // 连接 SQLite
    sqliteDB, _ := gorm.Open(sqlite.Open("./data/gold_admin.db"), &gorm.Config{})
    
    // 连接 MySQL
    mysqlDSN := "root:password@tcp(127.0.0.1:3306)/gold_admin?charset=utf8mb4&parseTime=True"
    mysqlDB, _ := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
    
    // 迁移数据（示例：用户表）
    var users []AdminUser
    sqliteDB.Find(&users)
    
    for _, user := range users {
        mysqlDB.Create(&user)
    }
    
    println("数据迁移完成！")
}
```

运行迁移：
```bash
go run tools/migrate_sqlite_to_mysql.go
```

### 方法二：导出 SQL 再导入

**步骤：**

1. **导出 SQLite 数据为 SQL**
```bash
sqlite3 ./data/gold_admin.db .dump > dump.sql
```

2. **转换 SQL 语法**

SQLite 和 MySQL 语法差异需要手动调整：

```sql
-- SQLite
CREATE TABLE admin_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  ...
);

-- 改为 MySQL
CREATE TABLE admin_users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

3. **导入 MySQL**
```bash
mysql -u root -p gold_admin < dump.sql
```

### 方法三：使用工具迁移

**推荐工具：**

1. **DB Browser for SQLite** + **MySQL Workbench**
   - 导出 CSV
   - 导入到 MySQL

2. **DBeaver**（通用数据库工具）
   - 支持 SQLite → MySQL 数据迁移
   - 图形化界面，简单易用

## 语法差异对照表

| 功能 | SQLite | MySQL |
|------|--------|-------|
| 主键自增 | `INTEGER PRIMARY KEY AUTOINCREMENT` | `INT PRIMARY KEY AUTO_INCREMENT` |
| 字符串类型 | `TEXT` / `VARCHAR(n)` | `VARCHAR(n)` |
| 日期时间 | `DATETIME` | `DATETIME` |
| 布尔类型 | `INTEGER (0/1)` | `TINYINT(1)` / `BOOLEAN` |
| 注释 | `--` | `-- ` 或 `/* */` 或 `COMMENT '...'` |
| 引擎 | 无 | `ENGINE=InnoDB` |
| 字符集 | 无 | `CHARSET=utf8mb4` |
| 外键 | 默认不启用 | 完全支持 |

## GORM 兼容性

好消息：使用 GORM 的话，几乎不需要修改代码！

**GORM 自动处理的差异：**
- ✅ 数据类型映射
- ✅ 主键自增
- ✅ 索引创建
- ✅ 外键约束

**需要注意的地方：**

1. **自动生成列（Generated Column）**

SQLite 不支持，MySQL 支持：

```go
// 价格表中的计算字段
type Price struct {
    BasePrice    float64 `gorm:"not null"`
    BuyPriceDiff float64 `gorm:"not null;default:0"`
    
    // MySQL 支持，SQLite 需要在应用层计算
    BuyPrice     float64 `gorm:"->"` // 只读字段
}

// 应用层计算（兼容两者）
func (p *Price) GetBuyPrice() float64 {
    return p.BasePrice + p.BuyPriceDiff
}
```

2. **全文搜索**

SQLite 的 FTS5 和 MySQL 的 FULLTEXT 语法不同，建议用 GORM 的 `Where` 实现：

```go
// 兼容写法
DB.Where("name LIKE ?", "%"+keyword+"%").Find(&users)
```

## 配置动态切换

**支持运行时选择数据库：**

```go
// config/database.go
package config

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func NewDatabase(cfg DatabaseConfig) (*gorm.DB, error) {
    var dialector gorm.Dialector
    
    switch cfg.Type {
    case "mysql":
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
        dialector = mysql.Open(dsn)
        
    case "sqlite":
        dialector = sqlite.Open(cfg.Path)
        
    default:
        return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
    }
    
    return gorm.Open(dialector, &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true, // 兼容性更好
    })
}
```

**配置文件示例：**

```yaml
# 开发环境 (config/dev.yaml)
database:
  type: sqlite
  path: ./data/dev.db

# 生产环境 (config/prod.yaml)
database:
  type: mysql
  host: 127.0.0.1
  port: 3306
  username: root
  password: prod_password
  database: gold_admin_prod
```

## 性能对比

| 场景 | SQLite | MySQL |
|------|--------|-------|
| 单用户读写 | ⭐⭐⭐⭐⭐ 极快 | ⭐⭐⭐⭐ 快 |
| 10并发用户 | ⭐⭐⭐ 一般 | ⭐⭐⭐⭐⭐ 优秀 |
| 50并发用户 | ⭐⭐ 较慢 | ⭐⭐⭐⭐⭐ 优秀 |
| 100万条记录 | ⭐⭐⭐ 可用 | ⭐⭐⭐⭐⭐ 优秀 |
| 复杂查询 | ⭐⭐⭐ 一般 | ⭐⭐⭐⭐⭐ 优秀 |
| 部署难度 | ⭐⭐⭐⭐⭐ 极简 | ⭐⭐⭐ 中等 |

## 最佳实践建议

### 开发阶段（推荐 SQLite）
```bash
# 开发配置
database:
  type: sqlite
  path: ./data/dev.db
```

**理由：**
- 快速启动，无需配置
- 便于测试和调试
- 数据库文件可随项目一起管理

### 演示/小规模生产（SQLite 也可以）
```bash
# 演示配置
database:
  type: sqlite
  path: ./data/gold_admin.db
```

**适用场景：**
- 单店铺
- 日均 < 100 笔交易
- 并发用户 < 10

### 大规模生产（推荐 MySQL）
```bash
# 生产配置
database:
  type: mysql
  host: db.production.com
  port: 3306
  username: gold_admin
  password: strong_password
  database: gold_admin
```

**适用场景：**
- 多店铺连锁
- 日均 > 1000 笔交易
- 并发用户 > 50

## 常见问题

### Q1: SQLite 数据库文件在哪？
A: 默认在 `./data/gold_admin.db`，可以在配置文件中修改。

### Q2: 如何备份 SQLite 数据库？
A: 直接复制 `gold_admin.db` 文件即可。

### Q3: SQLite 会影响性能吗？
A: 小规模应用（< 10并发用户）完全够用，性能甚至比 MySQL 更好。

### Q4: 迁移到 MySQL 需要改代码吗？
A: 使用 GORM 的话，只需要改配置文件，代码几乎不用改。

### Q5: 能否同时支持 SQLite 和 MySQL？
A: 可以！通过配置文件动态切换，见上文"配置动态切换"部分。

## 总结

✅ **推荐路线：**
```
SQLite (开发) → SQLite (演示) → MySQL (生产)
                                  ↑
                          根据实际需求决定是否迁移
```

✅ **核心优势：**
- 开发阶段使用 SQLite：简单、快速
- GORM 保证兼容性：迁移几乎无痛
- 根据规模决定：不一定需要 MySQL

✅ **迁移建议：**
- 数据量 < 10万：保持 SQLite
- 并发用户 > 20：考虑 MySQL
- 多店铺连锁：建议 MySQL

---

**文档版本**: v1.0  
**创建日期**: 2025-11-04  
**适用项目**: 沪金汇后台管理系统



