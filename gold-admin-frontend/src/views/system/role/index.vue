<template>
  <div class="app-container">
    <!-- 搜索栏 -->
    <el-form :inline="true" :model="queryParams" class="demo-form-inline">
      <el-form-item label="角色名称">
        <el-input v-model="queryParams.name" placeholder="请输入角色名称" clearable @keyup.enter.native="handleQuery" />
      </el-form-item>
      <el-form-item label="角色代码">
        <el-input v-model="queryParams.code" placeholder="请输入角色代码" clearable @keyup.enter.native="handleQuery" />
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
        <el-button type="primary" icon="el-icon-plus" @click="handleCreate">新增角色</el-button>
      </el-col>
    </el-row>

    <!-- 角色列表 -->
    <el-table v-loading="loading" :data="roleList" border>
      <el-table-column label="ID" prop="id" width="60" align="center" />
      <el-table-column label="角色名称" prop="name" width="150" />
      <el-table-column label="角色代码" prop="code" width="150" />
      <el-table-column label="描述" prop="description" min-width="200" show-overflow-tooltip />
      <el-table-column label="排序" prop="sort" width="80" align="center" />
      <el-table-column label="状态" width="80" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
            {{ scope.row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="created_at" width="160" :formatter="formatTime" />
      <el-table-column label="操作" width="200" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">编辑</el-button>
          <el-button size="mini" type="text" icon="el-icon-setting" @click="handlePermission(scope.row)">分配权限</el-button>
          <el-button v-if="scope.row.code !== 'super_admin'" size="mini" type="text" icon="el-icon-delete" style="color: #f56c6c;" @click="handleDelete(scope.row)">删除</el-button>
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
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色代码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色代码（英文）" :disabled="dialogStatus === 'update'" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
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

    <!-- 分配权限对话框 -->
    <el-dialog title="分配菜单权限" :visible.sync="permissionDialogVisible" width="400px">
      <el-tree
        ref="menuTree"
        :data="menuTreeData"
        :props="menuTreeProps"
        :default-checked-keys="checkedMenuIds"
        node-key="id"
        show-checkbox
        default-expand-all
      />
      <div slot="footer">
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPermission">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getRoleList, createRole, updateRole, deleteRole, getRole } from '@/api/role'
import { getMenuTree } from '@/api/menu'

export default {
  name: 'RoleManagement',
  data() {
    return {
      loading: false,
      roleList: [],
      total: 0,
      queryParams: {
        page: 1,
        page_size: 10,
        name: '',
        code: '',
        status: null
      },
      dialogVisible: false,
      permissionDialogVisible: false,
      dialogStatus: '',
      dialogTitle: '',
      form: {
        name: '',
        code: '',
        description: '',
        sort: 0,
        status: 1
      },
      currentRoleId: null,
      menuTreeData: [],
      menuTreeProps: {
        children: 'children',
        label: 'title'
      },
      checkedMenuIds: [],
      rules: {
        name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
        code: [{ required: true, message: '请输入角色代码', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getList()
    this.getMenuTree()
  },
  methods: {
    getList() {
      this.loading = true
      getRoleList(this.queryParams).then(response => {
        this.roleList = response.data.list || []
        this.total = response.data.total || 0
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },
    getMenuTree() {
      getMenuTree().then(response => {
        this.menuTreeData = response.data || []
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
        name: '',
        code: '',
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
        name: '',
        code: '',
        description: '',
        sort: 0,
        status: 1
      }
    },
    handleCreate() {
      this.resetForm()
      this.dialogStatus = 'create'
      this.dialogTitle = '新增角色'
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate()
      })
    },
    handleUpdate(row) {
      this.form = Object.assign({}, row)
      this.dialogStatus = 'update'
      this.dialogTitle = '编辑角色 - ' + row.name
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate()
      })
    },
    handlePermission(row) {
      this.currentRoleId = row.id
      // 获取角色已分配的菜单
      getRole(row.id).then(response => {
        this.checkedMenuIds = response.data.menu_ids || []
        this.permissionDialogVisible = true
        this.$nextTick(() => {
          this.$refs.menuTree.setCheckedKeys(this.checkedMenuIds)
        })
      })
    },
    handleDelete(row) {
      this.$confirm('确定要删除角色 "' + row.name + '" 吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteRole(row.id).then(() => {
          this.$message.success('删除成功')
          this.getList()
        })
      })
    },
    submitForm() {
      this.$refs.dataForm.validate(valid => {
        if (valid) {
          if (this.dialogStatus === 'create') {
            createRole(this.form).then(() => {
              this.$message.success('创建成功')
              this.dialogVisible = false
              this.getList()
            })
          } else {
            updateRole(this.form.id, this.form).then(() => {
              this.$message.success('更新成功')
              this.dialogVisible = false
              this.getList()
            })
          }
        }
      })
    },
    submitPermission() {
      const checkedKeys = this.$refs.menuTree.getCheckedKeys()
      const halfCheckedKeys = this.$refs.menuTree.getHalfCheckedKeys()
      const menuIds = [...checkedKeys, ...halfCheckedKeys]

      updateRole(this.currentRoleId, { menu_ids: menuIds }).then(() => {
        this.$message.success('权限分配成功')
        this.permissionDialogVisible = false
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

