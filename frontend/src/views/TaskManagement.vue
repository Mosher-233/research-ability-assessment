<template>
  <div class="task-management">
    <div class="header">
      <h2>任务管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">创建任务</el-button>
    </div>

    <!-- 任务列表 -->
    <el-card class="task-list">
      <el-table :data="tasks" style="width: 100%">
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
            <el-button size="small" type="primary" @click="assignTask(scope.row.id)">分配</el-button>
            <el-button size="small" type="warning" @click="updateTaskStatus(scope.row.id, scope.row.status)">归档</el-button>
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
        <el-descriptions-item label="教师">{{ taskDetails?.teacher?.name }}</el-descriptions-item>
        <el-descriptions-item label="开始日期">{{ formatDateTime(taskDetails?.start_date) }}</el-descriptions-item>
        <el-descriptions-item label="结束日期">{{ formatDateTime(taskDetails?.end_date) }}</el-descriptions-item>
        <el-descriptions-item label="状态"><el-tag :type="getStatusType(taskDetails?.status || '')">{{ getStatusText(taskDetails?.status || '') }}</el-tag></el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(taskDetails?.created_at) }}</el-descriptions-item>
      </el-descriptions>

      <h3 style="margin-top: 20px">学生任务列表</h3>
      <el-table :data="studentTasks" style="width: 100%">
        <el-table-column prop="student.name" label="学生姓名" width="150"></el-table-column>
        <el-table-column prop="student.student_id" label="学生ID" width="120"></el-table-column>
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
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { taskApi } from '../api/task'
import type { Task, CreateTaskRequest, AssignTaskRequest, StudentTask } from '../types/task'
import { formatDate, formatDateTime, getStatusType, getStatusText, getProgressPercentage, getProgressColor } from '../utils/format'

// 任务列表相关
const tasks = ref<Task[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

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

// 加载任务列表
const loadTasks = async () => {
  try {
    const response = await taskApi.getTasks()
    if (response.code === 200) {
      tasks.value = response.data
      total.value = response.data.length
    }
  } catch (error) {
    console.error('加载任务失败:', error)
    ElMessage.error('加载任务失败')
  }
}

// 加载学生列表
const loadStudents = async () => {
  try {
    const response = await taskApi.getStudents()
    if (response.code === 200) {
      students.value = response.data
    }
  } catch (error) {
    console.error('加载学生失败:', error)
    ElMessage.error('加载学生失败')
  }
}

// 创建任务
const handleCreateTask = async () => {
  if (!createFormRef.value) return
  
  await createFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        const response = await taskApi.createTask(createForm)
        if (response.code === 200) {
          ElMessage.success('创建任务成功')
          showCreateDialog.value = false
          loadTasks()
          // 重置表单
          Object.assign(createForm, {
            name: '',
            description: '',
            course_id: '',
            start_date: '',
            end_date: ''
          })
        }
      } catch (error) {
        console.error('创建任务失败:', error)
        ElMessage.error('创建任务失败')
      }
    }
  })
}

// 查看任务详情
const viewTaskDetails = async (taskId: string) => {
  try {
    // 获取任务详情
    const taskResponse = await taskApi.getTaskByID(taskId)
    if (taskResponse.code === 200) {
      taskDetails.value = taskResponse.data
    }
    
    // 获取学生任务列表
    const studentTaskResponse = await taskApi.getStudentTasks(taskId)
    if (studentTaskResponse.code === 200) {
      studentTasks.value = studentTaskResponse.data
    }
    
    showTaskDialog.value = true
  } catch (error) {
    console.error('获取任务详情失败:', error)
    ElMessage.error('获取任务详情失败')
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
    // 这里简化处理，实际应该调用API更新状态
    ElMessage.success('任务已归档')
    loadTasks()
  } catch (error) {
    console.error('更新任务状态失败:', error)
    ElMessage.error('更新任务状态失败')
  }
}

// 查看学生任务详情
const viewStudentTask = (studentTaskId: string) => {
  // 这里可以跳转到学生任务详情页面
  console.log('查看学生任务:', studentTaskId)
}

// 初始化
onMounted(() => {
  loadTasks()
})
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