<template>
  <div class="dashboard-container">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2>研究能力评价系统</h2>
      </div>
      <nav class="sidebar-menu">
        <!-- 任务管理 - 教师和学生都可以看到 -->
        <router-link to="/dashboard/tasks" class="menu-item">
          <el-icon><i-ep-task /></el-icon>
          <span>{{ user?.role === 'teacher' ? '任务管理' : '我的任务' }}</span>
        </router-link>
        
        <!-- 证据管理 - 教师和学生都可以看到 -->
        <router-link to="/dashboard/evidences" class="menu-item">
          <el-icon><i-ep-document /></el-icon>
          <span>证据管理</span>
        </router-link>
        
        <!-- 结果管理 - 教师和学生都可以看到 -->
        <router-link to="/dashboard/results" class="menu-item">
          <el-icon><i-ep-data-analysis /></el-icon>
          <span>结果管理</span>
        </router-link>
        
        <!-- 报告管理 - 教师和学生都可以看到 -->
        <router-link to="/dashboard/reports" class="menu-item">
          <el-icon><i-ep-document-copy /></el-icon>
          <span>报告管理</span>
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <div class="user-info">
          <span>{{ user?.name }}</span>
          <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
      </div>
    </aside>

    <!-- 主内容区域 -->
    <main class="main-content">
      <header class="main-header">
        <el-button type="text" class="menu-toggle" @click="toggleMenu">
          <el-icon><i-ep-menu /></el-icon>
        </el-button>
        <h1>{{ currentTitle }}</h1>
      </header>
      <div class="content">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { User } from '../types/user'

const router = useRouter()
const route = useRoute()
const menuCollapsed = ref(false)
const user = ref<User | null>(null)

// 计算当前页面标题
const currentTitle = computed(() => {
  const titleMap: Record<string, string> = {
    '/dashboard/tasks': '任务管理',
    '/dashboard/evidences': '证据管理',
    '/dashboard/results': '结果管理',
    '/dashboard/reports': '报告管理'
  }
  return titleMap[route.path] || '研究能力评价系统'
})

// 切换菜单展开/收起
const toggleMenu = () => {
  menuCollapsed.value = !menuCollapsed.value
}

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
  ElMessage.success('退出登录成功')
}

// 初始化用户信息
onMounted(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
})
</script>

<style scoped>
.dashboard-container {
  display: flex;
  min-height: 100vh;
  background-color: #f5f5f5;
}

/* 侧边栏样式 */
.sidebar {
  width: 240px;
  background-color: #1f2937;
  color: white;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #374151;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.sidebar-menu {
  flex: 1;
  padding: 20px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: #e5e7eb;
  text-decoration: none;
  transition: background-color 0.2s ease;
}

.menu-item:hover {
  background-color: #374151;
}

.menu-item.router-link-active {
  background-color: #3b82f6;
  color: white;
}

.menu-item el-icon {
  margin-right: 12px;
  font-size: 18px;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid #374151;
}

.user-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 主内容区域样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-header {
  background-color: white;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.menu-toggle {
  margin-right: 20px;
}

.main-header h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 1000;
    transform: translateX(-100%);
  }

  .sidebar.open {
    transform: translateX(0);
  }

  .main-content {
    margin-left: 0;
  }
}
</style>