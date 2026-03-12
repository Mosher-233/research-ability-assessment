<template>
  <div class="evidence-management">
    <div class="header">
      <h2>证据管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">上传证据</el-button>
    </div>

    <!-- 证据列表 -->
    <el-card class="evidence-list">
      <el-table :data="evidences" style="width: 100%">
        <el-table-column prop="id" label="证据ID" width="180"></el-table-column>
        <el-table-column prop="type" label="证据类型" width="120"></el-table-column>
        <el-table-column prop="kbm_name" label="KBM名称" width="150"></el-table-column>
        <el-table-column prop="kbm_level" label="KBM级别" width="100"></el-table-column>
        <el-table-column prop="content" label="内容" show-overflow-tooltip></el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" @click="viewEvidenceDetails(scope.row.id)">查看</el-button>
            <el-button size="small" type="danger" @click="deleteEvidence(scope.row.id)">删除</el-button>
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
          @size-change="loadEvidences"
          @current-change="loadEvidences"
        />
      </div>
    </el-card>

    <!-- 创建证据对话框 -->
    <el-dialog title="上传证据" v-model="showCreateDialog" width="600px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="120px">
        <el-form-item label="任务" prop="student_task_id">
          <el-select v-model="createForm.student_task_id" placeholder="请选择任务">
            <el-option v-for="task in studentTasks" :key="task.id" :label="task.task.name" :value="task.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="证据类型" prop="type">
          <el-select v-model="createForm.type" placeholder="请选择证据类型">
            <el-option label="文档" value="document"></el-option>
            <el-option label="代码" value="code"></el-option>
            <el-option label="演示" value="presentation"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="能力维度" prop="kbm_name">
          <el-select v-model="createForm.kbm_name" placeholder="请选择能力维度">
            <el-option label="文献综述" value="literature_review"></el-option>
            <el-option label="研究设计" value="research_design"></el-option>
            <el-option label="数据分析" value="data_analysis"></el-option>
            <el-option label="批判性思维" value="critical_thinking"></el-option>
            <el-option label="学术写作" value="academic_writing"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="内容" prop="content">
          <el-input v-model="createForm.content" placeholder="请输入证据内容" type="textarea" :rows="4"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreateEvidence">上传</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 证据详情对话框 -->
    <el-dialog title="证据详情" v-model="showEvidenceDialog" width="600px">
      <el-descriptions :column="2">
        <el-descriptions-item label="证据ID">{{ evidenceDetails?.id }}</el-descriptions-item>
        <el-descriptions-item label="任务">{{ evidenceDetails?.student_task_id }}</el-descriptions-item>
        <el-descriptions-item label="证据类型">{{ evidenceDetails?.type }}</el-descriptions-item>
        <el-descriptions-item label="能力维度">{{ evidenceDetails?.kbm_name }}</el-descriptions-item>
        <el-descriptions-item label="能力级别">{{ evidenceDetails?.kbm_level }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(evidenceDetails?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="内容" :span="2">{{ evidenceDetails?.content }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { evidenceApi } from '../api/evidence'
import { taskApi } from '../api/task'
import type { Evidence, CreateEvidenceRequest } from '../types/evidence'
import type { StudentTask } from '../types/task'
import { formatDateTime } from '../utils/format'

// 证据列表相关
const evidences = ref<Evidence[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框相关
const showCreateDialog = ref(false)
const showEvidenceDialog = ref(false)
const createFormRef = ref()

// 学生任务列表
const studentTasks = ref<StudentTask[]>([])

// 表单数据
const createForm = reactive<CreateEvidenceRequest>({
  student_task_id: '',
  type: '',
  content: '',
  kbm_name: ''
})

// 验证规则
const createRules = {
  student_task_id: [{ required: true, message: '请选择任务', trigger: 'blur' }],
  type: [{ required: true, message: '请选择证据类型', trigger: 'blur' }],
  kbm_name: [{ required: true, message: '请选择能力维度', trigger: 'blur' }],
  content: [{ required: true, message: '请输入证据内容', trigger: 'blur' }]
}

// 证据详情
const evidenceDetails = ref<Evidence | null>(null)

// 加载学生任务列表
const loadStudentTasks = async () => {
  try {
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      studentTasks.value = [
        {
          id: 'st1',
          student_id: 's1',
          task_id: 't1',
          status: 'active',
          progress: 50,
          created_at: '2026-03-01T10:00:00Z',
          updated_at: '2026-03-01T10:00:00Z',
          student: {
            id: 's1',
            name: '学生张三',
            email: 'student1@example.com',
            role: 'student',
            student_id: '2024001',
            major: '计算机科学',
            grade: '大三',
            created_at: '2026-01-01T00:00:00Z',
            updated_at: '2026-01-01T00:00:00Z'
          },
          task: {
            id: 't1',
            name: '研究方法课程作业',
            description: '完成研究方法课程的作业',
            course_id: 'RM2024',
            teacher_id: '1',
            start_date: '2026-03-01T00:00:00Z',
            end_date: '2026-03-15T23:59:59Z',
            status: 'active',
            student_count: 10,
            created_at: '2026-02-20T00:00:00Z',
            updated_at: '2026-02-20T00:00:00Z',
            teacher: {
              id: '1',
              name: '管理员',
              email: 'admin@example.com',
              role: 'teacher',
              created_at: '2026-01-01T00:00:00Z',
              updated_at: '2026-01-01T00:00:00Z'
            }
          }
        },
        {
          id: 'st2',
          student_id: 's1',
          task_id: 't2',
          status: 'active',
          progress: 30,
          created_at: '2026-03-02T14:30:00Z',
          updated_at: '2026-03-02T14:30:00Z',
          student: {
            id: 's1',
            name: '学生张三',
            email: 'student1@example.com',
            role: 'student',
            student_id: '2024001',
            major: '计算机科学',
            grade: '大三',
            created_at: '2026-01-01T00:00:00Z',
            updated_at: '2026-01-01T00:00:00Z'
          },
          task: {
            id: 't2',
            name: '社会调查研究',
            description: '完成社会调查研究项目',
            course_id: 'SS2024',
            teacher_id: '1',
            start_date: '2026-03-01T00:00:00Z',
            end_date: '2026-03-30T23:59:59Z',
            status: 'active',
            student_count: 15,
            created_at: '2026-02-25T00:00:00Z',
            updated_at: '2026-02-25T00:00:00Z',
            teacher: {
              id: '1',
              name: '管理员',
              email: 'admin@example.com',
              role: 'teacher',
              created_at: '2026-01-01T00:00:00Z',
              updated_at: '2026-01-01T00:00:00Z'
            }
          }
        }
      ]
    } else {
      // 调用后端API获取学生任务列表
      const response = await taskApi.getStudentTasks()
      if (response.code === 200) {
        studentTasks.value = response.data
      }
    }
  } catch (error) {
    console.error('加载学生任务失败:', error)
    ElMessage.error('加载学生任务失败')
  }
}

// 加载证据列表
const loadEvidences = async () => {
  try {
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      evidences.value = [
        {
          id: '1',
          student_task_id: 'st1',
          type: 'document',
          content: '这是一份文献综述报告，包含了相关领域的最新研究进展和分析。',
          kbm_name: 'literature_review',
          kbm_level: 4,
          created_at: '2026-03-01T10:00:00Z',
          updated_at: '2026-03-01T10:00:00Z'
        },
        {
          id: '2',
          student_task_id: 'st1',
          type: 'code',
          content: '这是一段数据分析代码，使用Python实现了数据清洗和可视化。',
          kbm_name: 'data_analysis',
          kbm_level: 3,
          created_at: '2026-03-02T14:30:00Z',
          updated_at: '2026-03-02T14:30:00Z'
        }
      ]
      total.value = evidences.value.length
    } else {
      // 调用后端API获取证据列表
      const response = await evidenceApi.getEvidences()
      if (response.code === 200) {
        evidences.value = response.data
        total.value = response.data.length
      }
    }
  } catch (error) {
    console.error('加载证据失败:', error)
    ElMessage.error('加载证据失败')
  }
}

// 上传证据
const handleCreateEvidence = async () => {
  if (!createFormRef.value) return
  
  await createFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        // 显示加载中
        ElMessage({ message: '正在上传证据并进行AI分析...', type: 'info' })
        
        const response = await evidenceApi.createEvidence(createForm)
        if (response.code === 200) {
          // 上传成功后获取AI分析结果
          const analysisResponse = await evidenceApi.analyzeEvidence(response.data.id)
          if (analysisResponse.code === 200) {
            // 显示AI反馈
            ElMessage.success({
              message: `上传证据成功！AI分析结果：能力级别为${analysisResponse.data.kbm_level}级，${analysisResponse.data.feedback}`,
              duration: 5000
            })
          } else {
            ElMessage.success('上传证据成功')
          }
          showCreateDialog.value = false
          loadEvidences()
          // 重置表单
          Object.assign(createForm, {
            student_task_id: '',
            type: '',
            content: '',
            kbm_name: ''
          })
        }
      } catch (error) {
        console.error('上传证据失败:', error)
        ElMessage.error('上传证据失败')
      }
    }
  })
}

// 查看证据详情
const viewEvidenceDetails = async (evidenceId: string) => {
  try {
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      const evidence = evidences.value.find(e => e.id === evidenceId)
      if (evidence) {
        evidenceDetails.value = evidence
      }
    } else {
      // 调用后端API获取证据详情
      const response = await evidenceApi.getEvidenceById(evidenceId)
      if (response.code === 200) {
        evidenceDetails.value = response.data
      }
    }
    showEvidenceDialog.value = true
  } catch (error) {
    console.error('获取证据详情失败:', error)
    ElMessage.error('获取证据详情失败')
  }
}

// 删除证据
const deleteEvidence = async (evidenceId: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个证据吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      evidences.value = evidences.value.filter(e => e.id !== evidenceId)
      total.value = evidences.value.length
    } else {
      // 调用后端API删除证据
      const response = await evidenceApi.deleteEvidence(evidenceId)
      if (response.code === 200) {
        loadEvidences()
      }
    }
    ElMessage.success('删除证据成功')
  } catch (error) {
    console.error('删除证据失败:', error)
    ElMessage.error('删除证据失败')
  }
}

// 初始化
onMounted(() => {
  loadStudentTasks()
  loadEvidences()
})
</script>

<style scoped>
.evidence-management {
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

.evidence-list {
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