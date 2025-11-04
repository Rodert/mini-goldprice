<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <el-form :inline="true" :model="queryParams" class="demo-form-inline">
      <el-form-item label="姓名">
        <el-input v-model="queryParams.name" placeholder="请输入姓名" clearable @keyup.enter.native="handleQuery" />
      </el-form-item>
      <el-form-item label="电话">
        <el-input v-model="queryParams.phone" placeholder="请输入电话" clearable @keyup.enter.native="handleQuery" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
          <el-option label="待确认" value="pending" />
          <el-option label="已确认" value="confirmed" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </el-form-item>
      <el-form-item label="预约时间">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="yyyy-MM-dd"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="el-icon-search" @click="handleQuery">搜索</el-button>
        <el-button icon="el-icon-refresh" @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 操作按钮 -->
    <el-row :gutter="10" style="margin-bottom: 15px;">
      <el-col :span="1.5">
        <el-button type="primary" icon="el-icon-plus" @click="handleCreate">新增预约</el-button>
      </el-col>
    </el-row>

    <!-- 预约列表 -->
    <el-table v-loading="loading" :data="appointmentList" border>
      <el-table-column label="ID" prop="id" width="60" align="center" />
      <el-table-column label="姓名" prop="name" width="100" />
      <el-table-column label="电话" prop="phone" width="120" />
      <el-table-column label="品种" prop="metal_type" width="100" />
      <el-table-column label="服务类型" width="100" align="center">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.service_type === 'store'" type="success">到店</el-tag>
          <el-tag v-else type="warning">上门</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="预约时间" width="180" align="center">
        <template slot-scope="scope">
          {{ formatDateTime(scope.row.appointment_time) }}
        </template>
      </el-table-column>
      <el-table-column label="地址" prop="address" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="100" align="center">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.status === 'pending'" type="warning">待确认</el-tag>
          <el-tag v-else-if="scope.row.status === 'confirmed'" type="primary">已确认</el-tag>
          <el-tag v-else-if="scope.row.status === 'completed'" type="success">已完成</el-tag>
          <el-tag v-else type="info">已取消</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="客户备注" prop="note" width="150" show-overflow-tooltip />
      <el-table-column label="创建时间" width="180" align="center">
        <template slot-scope="scope">
          {{ formatDateTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button size="mini" type="text" icon="el-icon-view" @click="handleDetail(scope.row)">详情</el-button>
          <el-button
            v-if="scope.row.status === 'pending'"
            size="mini"
            type="text"
            icon="el-icon-check"
            @click="handleConfirm(scope.row)"
          >确认</el-button>
          <el-button
            v-if="scope.row.status === 'confirmed'"
            size="mini"
            type="text"
            icon="el-icon-circle-check"
            @click="handleComplete(scope.row)"
          >完成</el-button>
          <el-button size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">编辑</el-button>
          <el-button size="mini" type="text" icon="el-icon-delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      :current-page="queryParams.page"
      :page-sizes="[10, 20, 50, 100]"
      :page-size="queryParams.page_size"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />

    <!-- 新增/编辑预约对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px" :close-on-click-modal="false">
      <el-form ref="form" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="品种" prop="metal_type">
          <el-input v-model="form.metal_type" placeholder="请输入品种" />
        </el-form-item>
        <el-form-item label="服务类型" prop="service_type">
          <el-radio-group v-model="form.service_type">
            <el-radio label="store">到店</el-radio>
            <el-radio label="home">上门</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="预约时间" prop="appointment_time">
          <el-date-picker
            v-model="form.appointment_time"
            type="datetime"
            placeholder="选择预约时间"
            value-format="yyyy-MM-dd HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item v-if="form.service_type === 'home'" label="地址" prop="address">
          <el-input v-model="form.address" type="textarea" :rows="3" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="客户备注">
          <el-input v-model="form.note" type="textarea" :rows="3" placeholder="请输入客户备注" />
        </el-form-item>
        <el-form-item v-if="isEdit" label="管理员备注">
          <el-input v-model="form.admin_remark" type="textarea" :rows="3" placeholder="请输入管理员备注" />
        </el-form-item>
        <el-form-item v-if="isEdit" label="状态">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option label="待确认" value="pending" />
            <el-option label="已确认" value="confirmed" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </div>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog title="预约详情" :visible.sync="detailVisible" width="600px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag v-if="detailData.status === 'pending'" type="warning">待确认</el-tag>
          <el-tag v-else-if="detailData.status === 'confirmed'" type="primary">已确认</el-tag>
          <el-tag v-else-if="detailData.status === 'completed'" type="success">已完成</el-tag>
          <el-tag v-else type="info">已取消</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="姓名">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="电话">{{ detailData.phone }}</el-descriptions-item>
        <el-descriptions-item label="品种">{{ detailData.metal_type }}</el-descriptions-item>
        <el-descriptions-item label="服务类型">
          <el-tag v-if="detailData.service_type === 'store'" type="success">到店</el-tag>
          <el-tag v-else type="warning">上门</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="预约时间" :span="2">
          {{ formatDateTime(detailData.appointment_time) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailData.address" label="地址" :span="2">
          {{ detailData.address }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailData.note" label="客户备注" :span="2">
          {{ detailData.note }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailData.admin_remark" label="管理员备注" :span="2">
          {{ detailData.admin_remark }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatDateTime(detailData.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailData.confirmed_at" label="确认时间" :span="2">
          {{ formatDateTime(detailData.confirmed_at) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailData.completed_at" label="完成时间" :span="2">
          {{ formatDateTime(detailData.completed_at) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="detailData.cancelled_at" label="取消时间" :span="2">
          {{ formatDateTime(detailData.cancelled_at) }}
        </el-descriptions-item>
      </el-descriptions>
      <div slot="footer" class="dialog-footer">
        <el-button @click="detailVisible = false">关闭</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getAppointmentList, createAppointment, updateAppointment, deleteAppointment } from '@/api/appointment'

export default {
  name: 'Appointment',
  data() {
    return {
      loading: false,
      appointmentList: [],
      total: 0,
      queryParams: {
        page: 1,
        page_size: 10,
        name: '',
        phone: '',
        status: '',
        start_date: '',
        end_date: ''
      },
      dateRange: [],
      dialogVisible: false,
      dialogTitle: '',
      isEdit: false,
      form: {
        id: null,
        name: '',
        phone: '',
        metal_type: '',
        service_type: 'store',
        appointment_time: '',
        address: '',
        note: '',
        admin_remark: '',
        status: 'pending'
      },
      rules: {
        name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
        phone: [{ required: true, message: '请输入电话', trigger: 'blur' }],
        metal_type: [{ required: true, message: '请输入品种', trigger: 'blur' }],
        service_type: [{ required: true, message: '请选择服务类型', trigger: 'change' }],
        appointment_time: [{ required: true, message: '请选择预约时间', trigger: 'change' }],
        address: [{ required: true, message: '请输入地址', trigger: 'blur' }]
      },
      detailVisible: false,
      detailData: {}
    }
  },
  created() {
    this.getList()
  },
  methods: {
    // 获取预约列表
    getList() {
      this.loading = true
      // 处理日期范围
      if (this.dateRange && this.dateRange.length === 2) {
        this.queryParams.start_date = this.dateRange[0]
        this.queryParams.end_date = this.dateRange[1]
      } else {
        this.queryParams.start_date = ''
        this.queryParams.end_date = ''
      }
      
      getAppointmentList(this.queryParams).then(response => {
        this.appointmentList = response.data.list || []
        this.total = response.data.total || 0
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },

    // 搜索
    handleQuery() {
      this.queryParams.page = 1
      this.getList()
    },

    // 重置搜索
    handleReset() {
      this.queryParams = {
        page: 1,
        page_size: 10,
        name: '',
        phone: '',
        status: '',
        start_date: '',
        end_date: ''
      }
      this.dateRange = []
      this.getList()
    },

    // 分页大小变化
    handleSizeChange(val) {
      this.queryParams.page_size = val
      this.getList()
    },

    // 页码变化
    handlePageChange(val) {
      this.queryParams.page = val
      this.getList()
    },

    // 新增预约
    handleCreate() {
      this.resetForm()
      this.dialogTitle = '新增预约'
      this.isEdit = false
      this.dialogVisible = true
    },

    // 编辑预约
    handleUpdate(row) {
      this.resetForm()
      this.form = { ...row }
      this.dialogTitle = '编辑预约'
      this.isEdit = true
      this.dialogVisible = true
    },

    // 查看详情
    handleDetail(row) {
      this.detailData = { ...row }
      this.detailVisible = true
    },

    // 确认预约
    handleConfirm(row) {
      this.$confirm('确认该预约？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        updateAppointment(row.id, { status: 'confirmed' }).then(() => {
          this.$message.success('确认成功')
          this.getList()
        })
      })
    },

    // 完成预约
    handleComplete(row) {
      this.$confirm('确认完成该预约？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        updateAppointment(row.id, { status: 'completed' }).then(() => {
          this.$message.success('操作成功')
          this.getList()
        })
      })
    },

    // 删除预约
    handleDelete(row) {
      this.$confirm('确定删除该预约吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteAppointment(row.id).then(() => {
          this.$message.success('删除成功')
          this.getList()
        })
      })
    },

    // 提交表单
    submitForm() {
      this.$refs.form.validate(valid => {
        if (!valid) {
          return
        }

        if (this.isEdit) {
          // 更新
          updateAppointment(this.form.id, this.form).then(() => {
            this.$message.success('更新成功')
            this.dialogVisible = false
            this.getList()
          })
        } else {
          // 新增
          createAppointment(this.form).then(() => {
            this.$message.success('创建成功')
            this.dialogVisible = false
            this.getList()
          })
        }
      })
    },

    // 重置表单
    resetForm() {
      this.form = {
        id: null,
        name: '',
        phone: '',
        metal_type: '',
        service_type: 'store',
        appointment_time: '',
        address: '',
        note: '',
        admin_remark: '',
        status: 'pending'
      }
      if (this.$refs.form) {
        this.$refs.form.clearValidate()
      }
    },

    // 格式化日期时间
    formatDateTime(dateTime) {
      if (!dateTime) return ''
      const date = new Date(dateTime)
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      const seconds = String(date.getSeconds()).padStart(2, '0')
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>

