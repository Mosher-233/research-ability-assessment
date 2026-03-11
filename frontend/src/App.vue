<template>
  <div class="app">
    <!-- 登录和注册页面不需要布局 -->
    <template v-if="$route.path === '/login' || $route.path === '/register'">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </template>
    
    <!-- 主布局 -->
    <template v-else>
      <div class="main-layout">
        <!-- 侧边栏 -->
        <aside class="sidebar">
          <div class="sidebar-header">
            <h2>研究能力评价系统</h2>
          </div>
          <nav class="sidebar-nav">
            <el-menu
              :default-active="activeMenu"
              class="el-menu-vertical-demo"
              router
            >
              <el-menu-item index="/dashboard/tasks">
                <el-icon><Check /></el-icon>
                <span>任务管理</span>
              </el-menu-item>
              <el-menu-item index="/dashboard/evidences">
                <el-icon><Document /></el-icon>
                <span>证据管理</span>
              </el-menu-item>
              <el-menu-item index="/dashboard/results">
                <el-icon><DataAnalysis /></el-icon>
                <span>结果管理</span>
              </el-menu-item>
              <el-menu-item index="/dashboard/reports">
                <el-icon><Reading /></el-icon>
                <span>报告管理</span>
              </el-menu-item>
            </el-menu>
          </nav>
          <div class="sidebar-footer">
            <el-button type="text" @click="logout">
              <el-icon><SwitchButton /></el-icon> 退出登录
            </el-button>
          </div>
        </aside>
        
        <!-- 主内容区域 -->
        <main class="main-content">
          <header class="content-header">
            <div class="header-left">
              <h1>{{ pageTitle }}</h1>
            </div>
            <div class="header-right">
              <el-dropdown>
                <span class="user-info">
                  <el-avatar :size="32" :src="userAvatar"></el-avatar>
                  <span class="user-name">{{ userName }}</span>
                  <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item>个人中心</el-dropdown-item>
                    <el-dropdown-item>设置</el-dropdown-item>
                    <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </header>
          <div class="content-body">
            <router-view v-slot="{ Component }">
              <transition name="fade" mode="out-in">
                <component :is="Component" />
              </transition>
            </router-view>
          </div>
        </main>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Check, Document, DataAnalysis, Reading, SwitchButton, ArrowDown } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

// 计算当前激活的菜单
const activeMenu = computed(() => {
  return route.path
})

// 计算页面标题
const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/dashboard/tasks': '任务管理',
    '/dashboard/evidences': '证据管理',
    '/dashboard/results': '结果管理',
    '/dashboard/reports': '报告管理'
  }
  return titles[route.path] || '研究能力评价系统'
})

// 用户信息
const userName = ref('管理员')
const userAvatar = ref('https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png')

// 退出登录
const logout = () => {
  localStorage.removeItem('token')
  ElMessage.success('退出登录成功！')
  router.push('/login')
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Arial, sans-serif;
  background-color: #f5f5f5;
  color: #333;
}

.app {
  min-height: 100vh;
}

/* 主布局 */
.main-layout {
  display: flex;
  min-height: 100vh;
}

/* 侧边栏 */
.sidebar {
  width: 240px;
  background-color: #2c3e50;
  color: #fff;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #34495e;
}

.sidebar-header h2 {
  font-size: 18px;
  font-weight: 600;
}

.sidebar-nav {
  flex: 1;
  padding: 20px 0;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid #34495e;
}

/* 主内容区域 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
}

.content-header {
  background-color: #fff;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.content-header h1 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.user-name {
  font-size: 14px;
  color: #333;
}

.content-body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 60px;
  }
  
  .sidebar-header h2,
  .el-menu-item span {
    display: none;
  }
  
  .sidebar-footer {
    text-align: center;
  }
  
  .sidebar-footer .el-button span {
    display: none;
  }
}
</style>