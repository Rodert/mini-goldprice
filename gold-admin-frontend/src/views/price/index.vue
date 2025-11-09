<template>
  <div class="app-container">
    <!-- 操作按钮 -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="24">
        <el-button type="primary" icon="el-icon-plus" @click="handleCreate">新增品种</el-button>
        <el-button type="warning" icon="el-icon-refresh" @click="handleSync">同步基础价</el-button>
      </el-col>
    </el-row>

    <!-- 价格列表 -->
    <el-row :gutter="20">
      <el-col v-for="item in list" :key="item.id" :xs="24" :sm="12" :md="8" :lg="6">
        <el-card shadow="hover" class="price-card">
          <div class="price-header">
            <div class="price-icon" :style="{backgroundColor: item.icon_color || '#FFD700'}">
              {{ item.icon || 'Au' }}
            </div>
            <div class="price-info">
              <div class="price-name">{{ item.name }}</div>
              <div class="price-subtitle">{{ item.subtitle }}</div>
              <div class="price-code">{{ item.code }}</div>
            </div>
          </div>

          <el-divider />

          <div class="price-detail">
            <div class="price-row base-price">
              <span class="label">基础价</span>
              <span class="value">¥{{ item.base_price }}/克</span>
            </div>
            <div class="price-row">
              <span class="label">回购差价</span>
              <span class="diff" :class="item.buy_price_diff < 0 ? 'negative' : 'positive'">
                {{ item.buy_price_diff >= 0 ? '+' : '' }}{{ item.buy_price_diff }}
              </span>
            </div>
            <div class="price-row buy-price">
              <span class="label">回购价</span>
              <span class="value">¥{{ item.buy_price }}/克</span>
            </div>
            <div class="price-row">
              <span class="label">销售差价</span>
              <span class="diff" :class="item.sell_price_diff < 0 ? 'negative' : 'positive'">
                {{ item.sell_price_diff >= 0 ? '+' : '' }}{{ item.sell_price_diff }}
              </span>
            </div>
            <div class="price-row sell-price">
              <span class="label">销售价</span>
              <span class="value">¥{{ item.sell_price }}/克</span>
            </div>
          </div>

          <div class="price-time">
            更新于 {{ item.updated_at | formatTime }}
          </div>

          <div class="price-actions">
            <el-button type="primary" size="small" @click="handleUpdate(item)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(item)">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px">
      <el-form ref="dataForm" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="品种名称" prop="name">
          <el-input v-model="form.name" placeholder="例如：黄金9999" />
        </el-form-item>
        <el-form-item label="唯一标识" prop="code">
          <el-input v-model="form.code" placeholder="例如：gold_9999（仅英文、数字、下划线）" :disabled="dialogStatus === 'update'" />
        </el-form-item>
        <el-form-item label="副标题" prop="subtitle">
          <el-input v-model="form.subtitle" placeholder="例如：Au9999 · 千足金" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="例如：Au" />
        </el-form-item>
        <el-form-item label="图标颜色" prop="icon_color">
          <el-color-picker v-model="form.icon_color" />
        </el-form-item>
        <el-divider />
        <el-form-item label="基础价格" prop="base_price">
          <el-input-number v-model="form.base_price" :precision="2" :step="1" :min="0" @change="calculatePrices" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">元/克</span>
        </el-form-item>
        <el-form-item label="回购差价" prop="buy_price_diff">
          <el-input-number v-model="form.buy_price_diff" :precision="2" :step="1" @change="calculatePrices" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">元/克（通常为负数）</span>
        </el-form-item>
        <el-form-item label="回购价">
          <span style="font-size: 16px; font-weight: bold; color: #f56c6c;">¥{{ calculatedBuyPrice }}/克</span>
        </el-form-item>
        <el-form-item label="销售差价" prop="sell_price_diff">
          <el-input-number v-model="form.sell_price_diff" :precision="2" :step="1" @change="calculatePrices" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">元/克（通常为正数）</span>
        </el-form-item>
        <el-form-item label="销售价">
          <span style="font-size: 16px; font-weight: bold; color: #67c23a;">¥{{ calculatedSellPrice }}/克</span>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getPriceList, createPrice, updatePrice, deletePrice, syncBasePrice } from '@/api/price'

export default {
  name: 'PriceList',
  filters: {
    formatTime(time) {
      if (!time) return ''
      return new Date(time).toLocaleString('zh-CN', { hour12: false })
    }
  },
  data() {
    return {
      list: [],
      dialogVisible: false,
      dialogStatus: '',
      dialogTitle: '',
      form: {
        name: '',
        code: '',
        subtitle: '',
        icon: 'Au',
        icon_color: '#FFD700',
        base_price: 560,
        buy_price_diff: -10,
        sell_price_diff: 15,
        sort: 0,
        status: 1
      },
      rules: {
        name: [{ required: true, message: '请输入品种名称', trigger: 'blur' }],
        code: [
          { required: true, message: '请输入唯一标识', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: '只能包含英文、数字和下划线', trigger: 'blur' }
        ],
        base_price: [{ required: true, message: '请输入基础价格', trigger: 'change' }]
      },
      calculatedBuyPrice: 0,
      calculatedSellPrice: 0
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      getPriceList().then(response => {
        this.list = response.data
      })
    },
    calculatePrices() {
      this.calculatedBuyPrice = (parseFloat(this.form.base_price) + parseFloat(this.form.buy_price_diff)).toFixed(2)
      this.calculatedSellPrice = (parseFloat(this.form.base_price) + parseFloat(this.form.sell_price_diff)).toFixed(2)
    },
    resetForm() {
      this.form = {
        name: '',
        code: '',
        subtitle: '',
        icon: 'Au',
        icon_color: '#FFD700',
        base_price: 560,
        buy_price_diff: -10,
        sell_price_diff: 15,
        sort: 0,
        status: 1
      }
      this.calculatePrices()
    },
    handleCreate() {
      this.resetForm()
      this.dialogStatus = 'create'
      this.dialogTitle = '新增贵金属品种'
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleUpdate(row) {
      this.form = Object.assign({}, row)
      this.dialogStatus = 'update'
      this.dialogTitle = '编辑价格 - ' + row.name
      this.dialogVisible = true
      this.calculatePrices()
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    handleDelete(row) {
      this.$confirm('确定要删除该品种吗？删除后将无法恢复', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deletePrice(row.id).then(() => {
          this.$message.success('删除成功')
          this.fetchData()
        })
      })
    },
    handleSync() {
      this.$confirm('确定要同步基础价格吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        syncBasePrice().then(() => {
          this.$message.success('同步成功')
          this.fetchData()
        })
      })
    },
    submitForm() {
      this.$refs['dataForm'].validate(valid => {
        if (valid) {
          if (this.dialogStatus === 'create') {
            createPrice(this.form).then(() => {
              this.$message.success('创建成功')
              this.dialogVisible = false
              this.fetchData()
            })
          } else {
            updatePrice(this.form.id, this.form).then(() => {
              this.$message.success('更新成功')
              this.dialogVisible = false
              this.fetchData()
            })
          }
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.price-card {
  margin-bottom: 20px;

  .price-header {
    display: flex;
    align-items: center;
    margin-bottom: 10px;

    .price-icon {
      width: 50px;
      height: 50px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 20px;
      font-weight: bold;
      color: #fff;
      margin-right: 15px;
    }

    .price-info {
      flex: 1;

      .price-name {
        font-size: 16px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 4px;
      }

      .price-subtitle {
        font-size: 12px;
        color: #909399;
        margin-bottom: 2px;
      }

      .price-code {
        font-size: 11px;
        color: #c0c4cc;
        font-family: monospace;
      }
    }
  }

  .price-detail {
    .price-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;

      &.base-price {
        background: #f5f7fa;
        margin: -10px -20px 10px;
        padding: 10px 20px;
        border-radius: 4px;

        .value {
          color: #409eff;
          font-weight: bold;
        }
      }

      &.buy-price .value {
        color: #f56c6c;
        font-weight: bold;
      }

      &.sell-price .value {
        color: #67c23a;
        font-weight: bold;
      }

      .label {
        font-size: 14px;
        color: #606266;
      }

      .value {
        font-size: 14px;
        color: #303133;
      }

      .diff {
        font-size: 14px;
        font-weight: bold;

        &.negative {
          color: #f56c6c;
        }

        &.positive {
          color: #67c23a;
        }
      }
    }
  }

  .price-time {
    font-size: 12px;
    color: #909399;
    margin-top: 10px;
    padding-top: 10px;
    border-top: 1px solid #ebeef5;
  }

  .price-actions {
    margin-top: 15px;
    display: flex;
    gap: 10px;

    .el-button {
      flex: 1;
    }
  }
}
</style>














