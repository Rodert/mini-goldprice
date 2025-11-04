<template>
  <div class="dashboard-container">
    <div class="app-container">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-card">
              <div class="stat-icon" style="background: #409eff;">
                <i class="el-icon-user" />
              </div>
              <div class="stat-content">
                <div class="stat-title">用户数</div>
                <div class="stat-value">128</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-card">
              <div class="stat-icon" style="background: #67c23a;">
                <i class="el-icon-s-finance" />
              </div>
              <div class="stat-content">
                <div class="stat-title">今日价格更新</div>
                <div class="stat-value">5</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-card">
              <div class="stat-icon" style="background: #e6a23c;">
                <i class="el-icon-s-order" />
              </div>
              <div class="stat-content">
                <div class="stat-title">预约数</div>
                <div class="stat-value">32</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card shadow="hover">
            <div class="stat-card">
              <div class="stat-icon" style="background: #f56c6c;">
                <i class="el-icon-s-data" />
              </div>
              <div class="stat-content">
                <div class="stat-title">回收记录</div>
                <div class="stat-value">156</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :span="24">
          <el-card>
            <div slot="header">
              <span>欢迎使用沪金汇管理系统</span>
            </div>
            <div class="welcome-content">
              <h3>系统功能</h3>
              <ul>
                <li>✅ 用户权限管理（RBAC）</li>
                <li>✅ 价格管理（基础价 + 差价模式）</li>
                <li>⏳ 预约管理（待开发）</li>
                <li>⏳ 回收记录（待开发）</li>
                <li>⏳ 数据统计（待开发）</li>
              </ul>

              <h3 style="margin-top: 30px;">快速开始</h3>
              <el-alert
                title="1. 系统管理 -> 用户管理：管理系统用户"
                type="info"
                :closable="false"
                style="margin: 10px 0;"
              />
              <el-alert
                title="2. 系统管理 -> 角色管理：配置用户角色和权限"
                type="info"
                :closable="false"
                style="margin: 10px 0;"
              />
              <el-alert
                title="3. 价格管理：配置贵金属回购和销售价格"
                type="info"
                :closable="false"
                style="margin: 10px 0;"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Dashboard'
}
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  min-height: 100vh;
}

.stat-card {
  display: flex;
  align-items: center;

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
    color: #fff;
  }

  .stat-content {
    margin-left: 20px;
    flex: 1;

    .stat-title {
      font-size: 14px;
      color: #909399;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 24px;
      font-weight: bold;
      color: #303133;
    }
  }
}

.welcome-content {
  h3 {
    color: #303133;
    margin-bottom: 15px;
  }

  ul {
    list-style: none;
    padding: 0;

    li {
      padding: 8px 0;
      font-size: 14px;
      color: #606266;
    }
  }
}
</style>



