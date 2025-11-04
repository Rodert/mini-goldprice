<template>
  <div class="app-container">
    <!-- 搜索栏 -->
    <el-form :inline="true" :model="queryParams" class="demo-form-inline">
      <el-form-item label="用户名">
        <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable @keyup.enter.native="handleQuery" />
      </el-form-item>
      <el-form-item label="真实姓名">
        <el-input v-model="queryParams.real_name" placeholder="请输入真实姓名" clearable @keyup.enter.native="handleQuery" />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="queryParams.status" placeholder="状态" clearable>
          <el-option label="启用" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="el-icon-search" @click="handleQuery">搜索</el-button>
        <el-button icon="el-icon-refresh" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 操作按钮 -->
    <el-row :gutter="10" style="margin-bottom: 15px;">
      <el-col :span="1.5">
        <el-button type="primary" icon="el-icon-plus" @click="handleCreate">新增用户</el-button>
      </el-col>
    </el-row>

    <!-- 用户列表 -->
    <el-table v-loading="loading" :data="userList" border>
      <el-table-column label="ID" prop="id" width="60" align="center" />
      <el-table-column label="用户名" prop="username" width="120" />
      <el-table-column label="真实姓名" prop="real_name" width="120" />
      <el-table-column label="手机号" prop="phone" width="120" />
      <el-table-column label="邮箱" prop="email" min-width="150" />
      <el-table-column label="状态" width="80" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="created_at" width="160" :formatter="formatTime" />
      <el-table-column label="操作" width="240" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">编辑</el-button>
          <el-button size="mini" type="text" icon="el-icon-key" @click="handlePassword(scope.row)">重置密码</el-button>
          <el-button v-if="scope.row.id !== 1" size="mini" type="text" icon="el-icon-delete" style="color: #f56c6c;" @click="handleDelete(scope.row)">删除</el-button>
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
      style="margin-top: 20px;"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px">
      <el-form ref="dataForm" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" :disabled="dialogStatus === 'update'" />
        </el-form-item>
        <el-form-item v-if="dialogStatus === 'create'" label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="真实姓名" prop="real_name">
          <el-input v-model="form.real_name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="角色" prop="role_ids">
          <el-select v-model="form.role_ids" multiple placeholder="请选择角色" style="width: 100%;">
            <el-option v-for="role in roleOptions" :key="role.id" :label="role.name" :value="role.id" />
          </el-select>
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

    <!-- 重置密码对话框 -->
    <el-dialog title="重置密码" :visible.sync="passwordDialogVisible" width="400px">
      <el-form ref="passwordForm" :model="passwordForm" :rules="passwordRules" label-width="80px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="passwordForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPassword">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getUserList, getUser, createUser, updateUser, deleteUser, updateUserPassword } from '@/api/user'
import { getAllRoles } from '@/api/role'

export default {
  name: 'UserManagement',
  data() {
    return {
      loading: false,
      userList: [],
      roleOptions: [],
      total: 0,
      queryParams: {
        page: 1,
        page_size: 10,
        username: '',
        real_name: '',
        status: null
      },
      dialogVisible: false,
      passwordDialogVisible: false,
      dialogStatus: '',
      dialogTitle: '',
      form: {
        username: '',
        password: '',
        real_name: '',
        phone: '',
        email: '',
        status: 1,
        role_ids: []
      },
      passwordForm: {
        user_id: null,
        password: ''
      },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }, { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }],
        real_name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }]
      },
      passwordRules: {
        password: [{ required: true, message: '请输入新密码', trigger: 'blur' }, { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getList()
    this.getRoles()
  },
  methods: {
    getList() {
      this.loading = true
      getUserList(this.queryParams).then(response => {
        this.userList = response.data.list || []
        this.total = response.data.total || 0
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },
    getRoles() {
      getAllRoles().then(response => {
        this.roleOptions = response.data || []
      })
    },
    handleQuery() {
      this.queryParams.page = 1
      this.getList()
    },
    resetQuery() {
      this.queryParams = {
        page: 1,
        page_size: 10,
        username: '',
        real_name: '',
        status: null
      }
      this.getList()
    },
    handleSizeChange(val) {
      this.queryParams.page_size = val
      this.getList()
    },
    handleCurrentChange(val) {
      this.queryParams.page = val
      this.getList()
    },
    resetForm() {
      this.form = {
        username: '',
        password: '',
        real_name: '',
        phone: '',
        email: '',
        status: 1,
        role_ids: []
      }
    },
    handleCreate() {
      this.resetForm()
      this.dialogStatus = 'create'
      this.dialogTitle = '新增用户'
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate()
      })
    },
    handleUpdate(row) {
      // 获取用户详情（包括角色信息）
      getUser(row.id).then(response => {
        const userData = response.data
        this.form = Object.assign({}, userData.user)
        this.form.role_ids = userData.role_ids || []
        this.dialogStatus = 'update'
        this.dialogTitle = '编辑用户 - ' + row.username
        this.dialogVisible = true
        this.$nextTick(() => {
          this.$refs.dataForm.clearValidate()
        })
      })
    },
    handlePassword(row) {
      this.passwordForm = {
        user_id: row.id,
        password: ''
      }
      this.passwordDialogVisible = true
      this.$nextTick(() => {
        this.$refs.passwordForm && this.$refs.passwordForm.clearValidate()
      })
    },
    handleDelete(row) {
      this.$confirm('确定要删除用户 "' + row.username + '" 吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser(row.id).then(() => {
          this.$message.success('删除成功')
          this.getList()
        })
      })
    },
    submitForm() {
      this.$refs.dataForm.validate(valid => {
        if (valid) {
          if (this.dialogStatus === 'create') {
            createUser(this.form).then(() => {
              this.$message.success('创建成功')
              this.dialogVisible = false
              this.getList()
            })
          } else {
            updateUser(this.form.id, this.form).then(() => {
              this.$message.success('更新成功')
              this.dialogVisible = false
              this.getList()
            })
          }
        }
      })
    },
    submitPassword() {
      this.$refs.passwordForm.validate(valid => {
        if (valid) {
          updateUserPassword(this.passwordForm.user_id, { password: this.passwordForm.password }).then(() => {
            this.$message.success('密码重置成功')
            this.passwordDialogVisible = false
          })
        }
      })
    },
    formatTime(row, column, cellValue) {
      if (!cellValue) return ''
      return new Date(cellValue).toLocaleString('zh-CN', { hour12: false })
    }
  }
}
</script>

<style scoped>
.demo-form-inline {
  margin-bottom: 0;
}
</style>

