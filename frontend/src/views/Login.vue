<template>
  <div class="login-container">
    <div class="login-form">
      <h2>登录</h2>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="loginForm.email" placeholder="请输入邮箱" type="email"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" placeholder="请输入密码" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading">登录</el-button>
          <el-button @click="navigateToRegister">注册</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { authApi } from '../api/auth'
import type { LoginRequest } from '../types/user'

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive<LoginRequest>({
  email: '',
  password: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        // 尝试调用后端API
        const response = await authApi.login(loginForm)
        if (response.code === 200) {
          localStorage.setItem('token', response.data.token)
          localStorage.setItem('user', JSON.stringify(response.data.user))
          router.push('/dashboard/tasks')
        }
      } catch (error) {
        // 如果后端不可用，使用模拟数据进行演示
        console.log('后端不可用，使用模拟数据演示')
        if (loginForm.email === 'admin@example.com' && loginForm.password === 'password123') {
          const mockUser = {
            id: '1',
            name: '管理员',
            email: 'admin@example.com',
            role: 'teacher',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          }
          const mockToken = 'mock-token-' + Date.now()
          localStorage.setItem('token', mockToken)
          localStorage.setItem('user', JSON.stringify(mockUser))
          ElMessage.success('登录成功（演示模式）')
          router.push('/dashboard/tasks')
        } else if (loginForm.email === 'student@example.com' && loginForm.password === 'password123') {
          const mockUser = {
            id: '2',
            name: '学生张三',
            email: 'student@example.com',
            role: 'student',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          }
          const mockToken = 'mock-token-' + Date.now()
          localStorage.setItem('token', mockToken)
          localStorage.setItem('user', JSON.stringify(mockUser))
          ElMessage.success('登录成功（演示模式）')
          router.push('/dashboard/tasks')
        } else {
          ElMessage.error('登录失败，请检查邮箱和密码')
        }
      } finally {
        loading.value = false
      }
    }
  })
}

const navigateToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.login-form {
  width: 400px;
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-form h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
}

.el-form-item {
  margin-bottom: 20px;
}

.el-form-item__content {
  width: 100%;
}

.el-button {
  width: 100%;
  margin-bottom: 10px;
}
</style>