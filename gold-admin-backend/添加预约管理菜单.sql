-- 添加预约管理菜单的SQL脚本
-- 使用方法：sqlite3 ./data/gold_admin.db < 添加预约管理菜单.sql

-- 1. 检查业务管理目录是否存在，不存在则创建
INSERT OR IGNORE INTO menus (id, parent_id, type, name, title, icon, path, component, sort, visible, status, created_at, updated_at)
VALUES (10, 0, 1, 'business', '业务管理', 'el-icon-s-management', '/business', 'Layout', 10, 1, 1, datetime('now'), datetime('now'));

-- 2. 添加预约管理菜单
INSERT OR IGNORE INTO menus (id, parent_id, type, name, title, icon, path, component, sort, visible, status, created_at, updated_at)
VALUES (12, 10, 2, 'appointment', '预约管理', 'el-icon-date', 'appointment', 'appointment/index', 2, 1, 1, datetime('now'), datetime('now'));

-- 3. 为超级管理员角色（ID=1）分配预约管理菜单权限
-- 先删除可能存在的旧记录，再插入
DELETE FROM role_menus WHERE role_id = 1 AND menu_id IN (10, 12);
INSERT INTO role_menus (role_id, menu_id, created_at)
VALUES 
  (1, 10, datetime('now')),
  (1, 12, datetime('now'));

-- 4. 显示添加结果
SELECT '✓ 菜单添加完成！' AS result;
SELECT * FROM menus WHERE id IN (10, 12);
SELECT '✓ 权限分配完成！' AS result;
SELECT * FROM role_menus WHERE menu_id IN (10, 12);

