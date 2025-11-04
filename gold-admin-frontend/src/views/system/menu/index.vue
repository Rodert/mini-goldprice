<template>
  <div class="app-container">
    <!-- 操作按钮 -->
    <el-row :gutter="10" style="margin-bottom: 15px;">
      <el-col :span="1.5">
        <el-button type="primary" icon="el-icon-plus" @click="handleCreate">新增菜单</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="info" icon="el-icon-sort" @click="toggleExpand">{{ expandAll ? '折叠' : '展开' }}</el-button>
      </el-col>
    </el-row>

    <!-- 菜单列表（树形表格） -->
    <el-table
      v-loading="loading"
      :data="menuList"
      row-key="id"
      :default-expand-all="expandAll"
      :tree-props="{children: 'children', hasChildren: 'hasChildren'}"
      border
    >
      <el-table-column label="菜单名称" prop="title" min-width="200" />
      <el-table-column label="图标" width="80" align="center">
        <template slot-scope="scope">
          <i v-if="scope.row.icon" :class="scope.row.icon" />
        </template>
      </el-table-column>
      <el-table-column label="类型" width="80" align="center">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.type === 1" type="success">目录</el-tag>
          <el-tag v-else-if="scope.row.type === 2">菜单</el-tag>
          <el-tag v-else type="info">按钮</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="路由路径" prop="path" width="150" show-overflow-tooltip />
      <el-table-column label="组件路径" prop="component" width="200" show-overflow-tooltip />
      <el-table-column label="权限标识" prop="permission" width="180" show-overflow-tooltip />
      <el-table-column label="排序" prop="sort" width="70" align="center" />
      <el-table-column label="操作" width="200" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">编辑</el-button>
          <el-button v-if="scope.row.type !== 3" size="mini" type="text" icon="el-icon-plus" @click="handleCreate(scope.row)">新增</el-button>
          <el-button size="mini" type="text" icon="el-icon-delete" style="color: #f56c6c;" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="700px">
      <el-form ref="dataForm" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="上级菜单">
              <el-cascader
                v-model="form.parent_id"
                :options="menuTreeData"
                :props="cascaderProps"
                placeholder="请选择上级菜单（可选）"
                clearable
                style="width: 100%;"
                @change="handleParentChange"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单类型" prop="type">
              <el-radio-group v-model="form.type">
                <el-radio :label="1">目录</el-radio>
                <el-radio :label="2">菜单</el-radio>
                <el-radio :label="3">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单图标">
              <el-input v-model="form.icon" placeholder="如：el-icon-setting">
                <i v-if="form.icon" slot="prefix" :class="form.icon" />
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单标题" prop="title">
              <el-input v-model="form.title" placeholder="请输入菜单标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="name">
              <el-input v-model="form.name" placeholder="路由名称，如：User" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="路由路径" prop="path">
              <el-input v-model="form.path" placeholder="如：/system" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组件路径">
              <el-input v-model="form.component" placeholder="如：system/user/index" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限标识">
              <el-input v-model="form.permission" placeholder="如：system:user:list" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="显示排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" :max="999" style="width: 100%;" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="是否可见">
              <el-radio-group v-model="form.visible">
                <el-radio :label="1">是</el-radio>
                <el-radio :label="0">否</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单状态">
              <el-radio-group v-model="form.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="0">禁用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getMenuList, createMenu, updateMenu, deleteMenu } from '@/api/menu'

export default {
  name: 'MenuManagement',
  data() {
    return {
      loading: false,
      menuList: [],
      menuTreeData: [],
      expandAll: false,
      dialogVisible: false,
      dialogStatus: '',
      dialogTitle: '',
      form: {
        parent_id: null,
        type: 1,
        name: '',
        title: '',
        icon: '',
        path: '',
        component: '',
        permission: '',
        sort: 0,
        visible: 1,
        status: 1
      },
      cascaderProps: {
        checkStrictly: true,
        value: 'id',
        label: 'title',
        children: 'children',
        emitPath: false
      },
      rules: {
        type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }],
        name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
        title: [{ required: true, message: '请输入菜单标题', trigger: 'blur' }],
        sort: [{ required: true, message: '请输入显示排序', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.loading = true
      getMenuList().then(response => {
        this.menuList = response.data || []
        this.menuTreeData = this.buildTreeData([{ id: 0, title: '根菜单', children: response.data || [] }])
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },
    buildTreeData(data) {
      return data.map(item => ({
        id: item.id,
        title: item.title,
        children: item.children && item.children.length > 0 ? this.buildTreeData(item.children) : undefined
      }))
    },
    toggleExpand() {
      this.expandAll = !this.expandAll
    },
    resetForm() {
      this.form = {
        parent_id: null,
        type: 1,
        name: '',
        title: '',
        icon: '',
        path: '',
        component: '',
        permission: '',
        sort: 0,
        visible: 1,
        status: 1
      }
    },
    handleParentChange(val) {
      this.form.parent_id = val || 0
    },
    handleCreate(row) {
      this.resetForm()
      if (row && row.id) {
        this.form.parent_id = row.id
      }
      this.dialogStatus = 'create'
      this.dialogTitle = '新增菜单'
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate()
      })
    },
    handleUpdate(row) {
      this.form = Object.assign({}, row)
      this.dialogStatus = 'update'
      this.dialogTitle = '编辑菜单 - ' + row.title
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dataForm.clearValidate()
      })
    },
    handleDelete(row) {
      this.$confirm('确定要删除菜单 "' + row.title + '" 吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteMenu(row.id).then(() => {
          this.$message.success('删除成功')
          this.getList()
        })
      })
    },
    submitForm() {
      this.$refs.dataForm.validate(valid => {
        if (valid) {
          const data = Object.assign({}, this.form)
          data.parent_id = data.parent_id || 0

          if (this.dialogStatus === 'create') {
            createMenu(data).then(() => {
              this.$message.success('创建成功')
              this.dialogVisible = false
              this.getList()
            })
          } else {
            updateMenu(this.form.id, data).then(() => {
              this.$message.success('更新成功')
              this.dialogVisible = false
              this.getList()
            })
          }
        }
      })
    }
  }
}
</script>

<style scoped>
</style>

