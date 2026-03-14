<template>
  <div class="task-management">
    <div class="header">
      <h2>{{ user?.role === 'teacher' ? '任务管理' : '我的任务' }}</h2>
      <el-button v-if="user?.role === 'teacher'" type="primary" @click="showCreateDialog = true">创建任务</el-button>
    </div>

    <!-- 任务Tab切换 -->
    <el-tabs v-model="activeTab" class="task-tabs">
      <el-tab-pane label="进行中" name="active"></el-tab-pane>
      <el-tab-pane label="已完成" name="completed"></el-tab-pane>
    </el-tabs>

    <!-- 任务列表 -->
    <el-card class="task-list">
      <el-table :data="filteredTasks" style="width: 100%">
        <el-table-column prop="name" label="任务名称" width="200"></el-table-column>
        <el-table-column prop="course_id" label="课程ID"></el-table-column>
        <el-table-column prop="start_date" label="开始日期" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.start_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="end_date" label="结束日期" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="student_count" label="学生数量" width="100"></el-table-column>
        <el-table-column label="进度" width="150">
          <template #default="scope">
            <el-progress 
              :percentage="getProgressPercentage(scope.row)" 
              :color="getProgressColor(getProgressPercentage(scope.row))"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="viewTaskDetails(scope.row.id)">查看</el-button>
            <el-button v-if="user?.role === 'teacher'" size="small" type="primary" @click="assignTask(scope.row.id)">分配</el-button>
            <el-button v-if="user?.role === 'teacher'" size="small" type="warning" @click="updateTaskStatus(scope.row.id, scope.row.status)">归档</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadTasks"
          @current-change="loadTasks"
        />
      </div>
    </el-card>

    <!-- 创建任务对话框 -->
    <el-dialog title="创建任务" v-model="showCreateDialog" width="500px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入任务名称"></el-input>
        </el-form-item>
        <el-form-item label="课程ID" prop="course_id">
          <el-input v-model="createForm.course_id" placeholder="请输入课程ID"></el-input>
        </el-form-item>
        <el-form-item label="任务描述" prop="description">
          <el-input v-model="createForm.description" placeholder="请输入任务描述" type="textarea"></el-input>
        </el-form-item>
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker v-model="createForm.start_date" type="datetime" placeholder="选择开始日期" style="width: 100%"></el-date-picker>
        </el-form-item>
        <el-form-item label="结束日期" prop="end_date">
          <el-date-picker v-model="createForm.end_date" type="datetime" placeholder="选择结束日期" style="width: 100%"></el-date-picker>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreateTask">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配任务对话框 -->
    <el-dialog title="分配任务" v-model="showAssignDialog" width="600px">
      <el-form :model="assignForm" ref="assignFormRef">
        <el-form-item label="任务名称">
          <el-input v-model="currentTaskName" disabled></el-input>
        </el-form-item>
        <el-form-item label="选择学生">
          <el-select v-model="assignForm.student_ids" multiple placeholder="请选择学生" style="width: 100%">
            <el-option 
              v-for="student in students" 
              :key="student.id" 
              :label="`${student.name} (${student.student_id})`" 
              :value="student.id"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAssignDialog = false">取消</el-button>
          <el-button type="primary" @click="handleAssignTask">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 任务详情对话框 -->
    <el-dialog title="任务详情" v-model="showTaskDialog" width="700px">
      <el-descriptions :column="2">
        <el-descriptions-item label="任务名称">{{ taskDetails?.name }}</el-descriptions-item>
        <el-descriptions-item label="课程ID">{{ taskDetails?.course_id }}</el-descriptions-item>
        <el-descriptions-item label="任务描述">{{ taskDetails?.description }}</el-descriptions-item>
        <el-descriptions-item label="教师">
          <template v-if="taskDetails?.teacher">
            {{ taskDetails.teacher.name }} ({{ taskDetails.teacher_id }})
          </template>
          <template v-else>
            {{ taskDetails?.teacher_id }}
          </template>
        </el-descriptions-item>
        <el-descriptions-item label="开始日期">{{ formatDateTime(taskDetails?.start_date) }}</el-descriptions-item>
        <el-descriptions-item label="结束日期">{{ formatDateTime(taskDetails?.end_date) }}</el-descriptions-item>
        <el-descriptions-item label="状态"><el-tag :type="getStatusType(taskDetails?.status || '')">{{ getStatusText(taskDetails?.status || '') }}</el-tag></el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(taskDetails?.created_at) }}</el-descriptions-item>
      </el-descriptions>

      <h3 style="margin-top: 20px">学生任务列表</h3>
      <el-table :data="studentTasks" style="width: 100%">
        <el-table-column prop="student_id" label="学生ID" width="150"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="150">
          <template #default="scope">
            <el-progress :percentage="scope.row.progress" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="scope">
            <el-button size="small" @click="viewStudentTask(scope.row.id)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { taskApi } from '../api/task'
import type { Task, CreateTaskRequest, AssignTaskRequest, StudentTask } from '../types/task'
import { formatDate, formatDateTime, getStatusType, getStatusText, getProgressPercentage, getProgressColor } from '../utils/format'

// 任务列表相关
const tasks = ref<Task[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('active')

// 去重函数
const deduplicateTasks = (tasks: any[]) => {
  const seen = new Set()
  return tasks.filter(task => {
    if (seen.has(task.id)) {
      return false
    }
    seen.add(task.id)
    return true
  })
}

// 计算属性：分开显示已完成和未完成任务
const filteredTasks = computed(() => {
  const uniqueTasks = deduplicateTasks(tasks.value)
  if (activeTab.value === 'active') {
    return uniqueTasks.filter(task => task.status !== 'completed' && task.status !== 'archived')
  } else {
    return uniqueTasks.filter(task => task.status === 'completed' || task.status === 'archived')
  }
})

// 对话框相关
const showCreateDialog = ref(false)
const showAssignDialog = ref(false)
const showTaskDialog = ref(false)
const createFormRef = ref()
const assignFormRef = ref()

// 表单数据
const createForm = reactive<CreateTaskRequest>({
  name: '',
  description: '',
  course_id: '',
  start_date: '',
  end_date: ''
})

const assignForm = reactive<AssignTaskRequest>({
  student_ids: []
})

// 验证规则
const createRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  course_id: [{ required: true, message: '请输入课程ID', trigger: 'blur' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'blur' }],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'blur' }]
}

// 学生列表
const students = ref<any[]>([])

// 任务详情
const taskDetails = ref<Task | null>(null)
const studentTasks = ref<StudentTask[]>([])
const currentTaskId = ref('')
const currentTaskName = ref('')
const user = ref<any>(null)

// 初始化用户信息
onMounted(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
  loadTasks()
})

// 加载任务列表
const loadTasks = async () => {
  try {
    console.log('开始加载任务...')
    // 获取用户信息
    const userStr = localStorage.getItem('user')
    const user = userStr ? JSON.parse(userStr) : null
    console.log('当前用户:', user)
    
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      console.log('使用模拟数据加载任务...')
      if (user && user.role === 'teacher') {
        // 教师用户显示所有任务
        tasks.value = [
          {
            id: '1',
            name: '研究方法课程作业',
            description: '完成研究方法课程的作业',
            course_id: 'RM2024',
            teacher_id: '1',
            start_date: new Date().toISOString(),
            end_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
            status: 'active',
            student_count: 10,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            teacher: {
              id: '1',
              name: '管理员',
              email: 'admin@example.com',
              role: 'teacher',
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString()
            }
          },
          {
            id: '2',
            name: '社会调查研究',
            description: '完成社会调查研究项目',
            course_id: 'SS2024',
            teacher_id: '1',
            start_date: new Date().toISOString(),
            end_date: new Date(Date.now() + 14 * 24 * 60 * 60 * 1000).toISOString(),
            status: 'active',
            student_count: 15,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            teacher: {
              id: '1',
              name: '管理员',
              email: 'admin@example.com',
              role: 'teacher',
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString()
            }
          }
        ]
      } else if (user && user.role === 'student') {
        // 学生用户只显示分配给他们的任务
        tasks.value = [
          {
            id: '1',
            name: '研究方法课程作业',
            description: '完成研究方法课程的作业',
            course_id: 'RM2024',
            teacher_id: '1',
            start_date: new Date().toISOString(),
            end_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
            status: 'active',
            student_count: 10,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            teacher: {
              id: '1',
              name: '管理员',
              email: 'admin@example.com',
              role: 'teacher',
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString()
            }
          }
        ]
      }
      total.value = tasks.value.length
      console.log('模拟数据加载成功:', tasks.value)
    } else {
      // 调用后端API
      if (user && user.role === 'teacher') {
        // 教师用户获取所有任务
        const response = await taskApi.getTasks()
        console.log('任务加载响应:', response)
        if (response.code === 200) {
          tasks.value = response.data
          total.value = response.data.length
          console.log('任务加载成功:', tasks.value)
        } else {
          console.error('任务加载失败:', response.message)
          ElMessage.error('加载任务失败: ' + response.message)
        }
      } else if (user && user.role === 'student') {
        // 学生用户获取分配给他们的任务
        // 这里需要调用学生任务列表API，假设API路径为/tasks/students/assigned
        try {
          const response = await taskApi.getAssignedTasks()
          console.log('学生任务加载响应:', response)
          if (response.code === 200) {
            tasks.value = response.data
            total.value = response.data.length
            console.log('学生任务加载成功:', tasks.value)
          } else {
            console.error('学生任务加载失败:', response.message)
            ElMessage.error('加载任务失败: ' + response.message)
          }
        } catch (error) {
          // 如果API不存在，使用模拟数据
          console.error('获取学生任务API不存在，使用模拟数据:', error)
          tasks.value = [
            {
              id: '1',
              name: '研究方法课程作业',
              description: '完成研究方法课程的作业',
              course_id: 'RM2024',
              teacher_id: '1',
              start_date: new Date().toISOString(),
              end_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
              status: 'active',
              student_count: 10,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
              teacher: {
                id: '1',
                name: '管理员',
                email: 'admin@example.com',
                role: 'teacher',
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString()
              }
            }
          ]
          total.value = tasks.value.length
          console.log('模拟数据加载成功:', tasks.value)
        }
      }
    }
  } catch (error) {
    console.error('加载任务失败:', error)
    ElMessage.error('加载任务失败: ' + error.message)
  }
}

// 加载学生列表
const loadStudents = async () => {
  try {
    console.log('开始加载学生...')
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      console.log('使用模拟数据加载学生...')
      students.value = [
        {
          id: '1',
          name: '学生张三',
          email: 'student1@example.com',
          role: 'student',
          student_id: '2024001',
          major: '计算机科学',
          grade: '大三',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        },
        {
          id: '2',
          name: '学生李四',
          email: 'student2@example.com',
          role: 'student',
          student_id: '2024002',
          major: '数据科学',
          grade: '大三',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        },
        {
          id: '3',
          name: '学生王五',
          email: 'student3@example.com',
          role: 'student',
          student_id: '2024003',
          major: '人工智能',
          grade: '大三',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
      ]
      console.log('模拟数据加载成功:', students.value)
    } else {
      // 调用后端API
      const response = await taskApi.getStudents()
      console.log('学生加载响应:', response)
      if (response.code === 200) {
        students.value = response.data
        console.log('学生加载成功:', students.value)
      } else {
        console.error('学生加载失败:', response.message)
        ElMessage.error('加载学生失败: ' + response.message)
      }
    }
  } catch (error) {
    console.error('加载学生失败:', error)
    ElMessage.error('加载学生失败: ' + error.message)
  }
}

// 创建任务
const handleCreateTask = async () => {
  if (!createFormRef.value) return
  
  await createFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        // 转换日期格式为ISO字符串
        const taskData = {
          ...createForm,
          start_date: createForm.start_date ? new Date(createForm.start_date).toISOString() : '',
          end_date: createForm.end_date ? new Date(createForm.end_date).toISOString() : ''
        }
        
        const response = await taskApi.createTask(taskData)
        console.log('创建任务响应:', response)
        if (response.code === 200) {
          ElMessage.success('创建任务成功')
          showCreateDialog.value = false
          await loadTasks()
          // 重置表单
          Object.assign(createForm, {
            name: '',
            description: '',
            course_id: '',
            start_date: '',
            end_date: ''
          })
        } else {
          ElMessage.error('创建任务失败: ' + response.message)
        }
      } catch (error) {
        console.error('创建任务失败:', error)
        ElMessage.error('创建任务失败: ' + (error.message || '未知错误'))
      }
    }
  })
}

// 查看任务详情
const viewTaskDetails = async (taskId: string) => {
  try {
    // 获取任务详情
    const taskResponse = await taskApi.getTaskByID(taskId)
    console.log('任务详情响应:', taskResponse)
    if (taskResponse.code === 200) {
      taskDetails.value = taskResponse.data
    }
    
    // 获取学生任务列表并去重
    const studentTaskResponse = await taskApi.getStudentTasks(taskId)
    console.log('学生任务列表响应:', studentTaskResponse)
    if (studentTaskResponse.code === 200) {
      // 根据student_id去重
      const seen = new Set()
      studentTasks.value = studentTaskResponse.data.filter((st: any) => {
        if (seen.has(st.student_id)) return false
        seen.add(st.student_id)
        return true
      })
    }
    
    showTaskDialog.value = true
  } catch (error) {
    console.error('获取任务详情失败:', error)
    ElMessage.error('获取任务详情失败: ' + (error.message || '未知错误'))
  }
}

// 分配任务
const assignTask = async (taskId: string) => {
  currentTaskId.value = taskId
  
  // 获取任务详情
  try {
    const response = await taskApi.getTaskByID(taskId)
    if (response.code === 200) {
      currentTaskName.value = response.data.name
    }
  } catch (error) {
    console.error('获取任务详情失败:', error)
    ElMessage.error('获取任务详情失败')
  }
  
  // 加载学生列表
  await loadStudents()
  
  // 重置分配表单
  assignForm.student_ids = []
  
  showAssignDialog.value = true
}

// 处理分配任务
const handleAssignTask = async () => {
  if (assignForm.student_ids.length === 0) {
    ElMessage.warning('请选择学生')
    return
  }
  
  try {
    const response = await taskApi.assignTask(currentTaskId.value, assignForm)
    if (response.code === 200) {
      ElMessage.success('分配任务成功')
      showAssignDialog.value = false
      loadTasks()
    }
  } catch (error) {
    console.error('分配任务失败:', error)
    ElMessage.error('分配任务失败')
  }
}

// 更新任务状态
const updateTaskStatus = async (taskId: string, currentStatus: string) => {
  try {
    await ElMessageBox.confirm(
      '确定要归档此任务吗？归档后任务将移至"已完成"列表。',
      '确认归档',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await taskApi.updateTaskStatus(taskId, 'archived')
    if (response.code === 200) {
      ElMessage.success('任务已归档')
      await loadTasks()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新任务状态失败:', error)
      ElMessage.error('更新任务状态失败')
    }
  }
}

// 查看学生任务详情
const viewStudentTask = (studentTaskId: string) => {
  // 这里可以跳转到学生任务详情页面
  console.log('查看学生任务:', studentTaskId)
}


</script>

<style scoped>
.task-management {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
}

.task-tabs {
  margin-bottom: 20px;
}

.task-list {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>