<template>
  <div class="app-container">
    <!-- 用户数量统计卡片 -->
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="8" :lg="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon">👥</div>
          <div class="stat-title">管理员数量</div>
          <div class="stat-value">{{ stats.total_users || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getDashboardStats } from '@/api/dashboard'

export default {
  name: 'Dashboard',
  data() {
    return {
      stats: {}
    }
  },
  created() {
    this.getStats()
  },
  methods: {
    getStats() {
      getDashboardStats().then(response => {
        this.stats = response.data || {}
      }).catch(error => {
        // 错误已经在 request.js 中处理了，这里只做静默失败
        console.error('获取统计数据失败:', error)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.stat-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  height: 160px;

  ::v-deep .el-card__body {
    padding: 25px;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  .stat-icon {
    font-size: 40px;
    margin-bottom: 10px;
  }

  .stat-title {
    font-size: 14px;
    opacity: 0.9;
    margin-bottom: 10px;
  }

  .stat-value {
    font-size: 32px;
    font-weight: bold;
    margin-bottom: 5px;
  }

  .stat-desc {
    font-size: 12px;
    opacity: 0.8;
  }
}
</style>
