<template>
  <div class="evidence-management">
    <div class="header">
      <h2>证据管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">上传证据</el-button>
    </div>

    <el-card class="evidence-list">
      <el-table :data="evidences" style="width: 100%">
        <el-table-column prop="id" label="证据ID" width="180"></el-table-column>
        <el-table-column prop="student_name" label="学生" width="120" v-if="isTeacher"></el-table-column>
        <el-table-column prop="task_name" label="任务" width="150" v-if="isTeacher"></el-table-column>
        <el-table-column prop="type" label="证据类型" width="120">
          <template #default="scope">
            {{ getTypeLabel(scope.row.type) }}
          </template>
        </el-table-column>
        <el-table-column prop="kbm_name" label="能力维度" width="150">
          <template #default="scope">
            {{ getKBMNameLabel(scope.row.kbm_name) }}
          </template>
        </el-table-column>
        <el-table-column prop="kbm_level" label="KBM级别" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.kbm_level > 0" :type="getLevelType(scope.row.kbm_level)">
              {{ scope.row.kbm_level }}级
            </el-tag>
            <span v-else>待分析</span>
          </template>
        </el-table-column>
        <el-table-column prop="file_name" label="文件" width="150">
          <template #default="scope">
            <el-link v-if="scope.row.file_name" type="primary" @click="handleDownloadFile(scope.row.id)">
              {{ scope.row.file_name }}
            </el-link>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewEvidenceDetails(scope.row)">查看</el-button>
            <el-button size="small" type="success" @click="analyzeEvidence(scope.row.id)" v-if="!scope.row.kbm_level || scope.row.kbm_level === 0">AI分析</el-button>
            <el-button size="small" type="info" @click="viewFeedback(scope.row.id)" v-if="scope.row.kbm_level && scope.row.kbm_level > 0">查看反馈</el-button>
            <el-button size="small" type="danger" @click="deleteEvidence(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog title="上传证据" v-model="showCreateDialog" width="700px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="120px">
        <el-form-item label="任务" prop="student_task_id">
          <el-select v-model="createForm.student_task_id" placeholder="请选择任务" style="width: 100%" filterable>
            <el-option 
              v-for="task in uniqueStudentTasks" 
              :key="task.id" 
              :label="getTaskLabel(task)"
              :value="task.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="证据类型" prop="type">
          <el-select v-model="createForm.type" placeholder="请选择证据类型" style="width: 100%">
            <el-option label="文档" value="document"></el-option>
            <el-option label="代码" value="code"></el-option>
            <el-option label="演示" value="presentation"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="能力维度" prop="kbm_name">
          <el-select v-model="createForm.kbm_name" placeholder="请选择能力维度" style="width: 100%">
            <el-option label="文献综述" value="literature_review"></el-option>
            <el-option label="研究设计" value="research_design"></el-option>
            <el-option label="数据分析" value="data_analysis"></el-option>
            <el-option label="批判性思维" value="critical_thinking"></el-option>
            <el-option label="学术写作" value="academic_writing"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="上传方式">
          <el-radio-group v-model="uploadMode">
            <el-radio label="text">文本输入</el-radio>
            <el-radio label="file">文件上传</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="内容" prop="content" v-if="uploadMode === 'text'">
          <el-input 
            v-model="createForm.content" 
            placeholder="请输入证据内容" 
            type="textarea" 
            :rows="4"
          ></el-input>
        </el-form-item>
        <el-form-item label="文件" prop="file" v-if="uploadMode === 'file'">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
            accept=".txt,.pdf,.doc,.docx,.md"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持 .txt, .pdf, .doc, .docx, .md 格式文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="resetCreateForm">取消</el-button>
          <el-button type="primary" @click="handleCreateEvidence">上传</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog title="证据详情" v-model="showEvidenceDialog" width="700px">
      <el-descriptions :column="2" v-if="evidenceDetails">
        <el-descriptions-item label="证据ID">{{ evidenceDetails.id }}</el-descriptions-item>
        <el-descriptions-item label="学生" v-if="evidenceDetails.student_name">{{ evidenceDetails.student_name }}</el-descriptions-item>
        <el-descriptions-item label="任务" v-if="evidenceDetails.task_name">{{ evidenceDetails.task_name }}</el-descriptions-item>
        <el-descriptions-item label="证据类型">{{ getTypeLabel(evidenceDetails.type) }}</el-descriptions-item>
        <el-descriptions-item label="能力维度">{{ getKBMNameLabel(evidenceDetails.kbm_name) }}</el-descriptions-item>
        <el-descriptions-item label="能力级别">
          <el-tag v-if="evidenceDetails.kbm_level > 0" :type="getLevelType(evidenceDetails.kbm_level)">
            {{ evidenceDetails.kbm_level }}级
          </el-tag>
          <span v-else>待分析</span>
        </el-descriptions-item>
        <el-descriptions-item label="文件" v-if="evidenceDetails.file_name">
          <el-link type="primary" @click="handleDownloadFile(evidenceDetails.id)">
            {{ evidenceDetails.file_name }}
          </el-link>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(evidenceDetails.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="内容" :span="2">
          <div style="white-space: pre-wrap; word-break: break-all;">
            {{ evidenceDetails.content || '(文件内容)' }}
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog title="AI反馈结果" v-model="showFeedbackDialog" width="700px">
      <el-descriptions :column="1" v-if="feedbackDetails">
        <el-descriptions-item label="KBM级别">
          <el-tag :type="getLevelType(feedbackDetails.kbm_level)">
            {{ feedbackDetails.kbm_level }}级
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="优点">
          <div style="white-space: pre-wrap;">{{ feedbackDetails.strengths || '暂无' }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="不足">
          <div style="white-space: pre-wrap;">{{ feedbackDetails.weaknesses || '暂无' }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="建议">
          <div style="white-space: pre-wrap;">{{ feedbackDetails.suggestions || '暂无' }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="总体评价">
          <div style="white-space: pre-wrap;">{{ feedbackDetails.content }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="反馈文件" v-if="feedbackDetails.file_name">
          <el-link type="primary">{{ feedbackDetails.file_name }}</el-link>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox, type UploadFile, type UploadUserFile } from 'element-plus'
import { evidenceApi } from '../api/evidence'
import { taskApi } from '../api/task'
import type { Evidence, CreateEvidenceRequest, Feedback } from '../types/evidence'
import type { StudentTask } from '../types/task'
import { formatDateTime } from '../utils/format'

const evidences = ref<Evidence[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showCreateDialog = ref(false)
const showEvidenceDialog = ref(false)
const showFeedbackDialog = ref(false)
const createFormRef = ref()
const uploadRef = ref()

const studentTasks = ref<StudentTask[]>([])
const uploadMode = ref<'text' | 'file'>('text')
const selectedFile = ref<UploadUserFile | null>(null)

const userRole = ref(localStorage.getItem('userRole') || 'student')
const isTeacher = computed(() => userRole.value === 'teacher')

const uniqueStudentTasks = computed(() => {
  const seen = new Set()
  return studentTasks.value.filter(task => {
    if (seen.has(task.task_id)) return false
    seen.add(task.task_id)
    return true
  })
})

const createForm = reactive<CreateEvidenceRequest>({
  student_task_id: '',
  type: '',
  content: '',
  kbm_name: ''
})

const createRules = {
  student_task_id: [{ required: true, message: '请选择任务', trigger: 'change' }],
  type: [{ required: true, message: '请选择证据类型', trigger: 'change' }],
  kbm_name: [{ required: true, message: '请选择能力维度', trigger: 'change' }],
  content: [{ required: true, message: '请输入证据内容', trigger: 'blur' }]
}

const evidenceDetails = ref<Evidence | null>(null)
const feedbackDetails = ref<Feedback | null>(null)

const getTaskLabel = (task: StudentTask) => {
  return task.task?.name || '未知任务'
}

const getTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    document: '文档',
    code: '代码',
    presentation: '演示',
    other: '其他'
  }
  return labels[type] || type
}

const getKBMNameLabel = (name: string) => {
  const labels: Record<string, string> = {
    literature_review: '文献综述',
    research_design: '研究设计',
    data_analysis: '数据分析',
    critical_thinking: '批判性思维',
    academic_writing: '学术写作'
  }
  return labels[name] || name
}

const getLevelType = (level: number) => {
  const types: Record<number, string> = {
    1: 'danger',
    2: 'warning',
    3: '',
    4: 'success',
    5: 'success'
  }
  return types[level] || ''
}

const loadStudentTasks = async () => {
  try {
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
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
        }
      ]
    } else {
      const response = await taskApi.getStudentTasks()
      if (response.code === 200 && response.data) {
        studentTasks.value = response.data
      }
    }
  } catch (error) {
    console.error('加载学生任务失败:', error)
    ElMessage.error('加载学生任务失败')
  }
}

const loadEvidences = async () => {
  try {
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      evidences.value = [
        {
          id: 'EV2026010201001',
          student_task_id: 'st1',
          type: 'document',
          content: '这是一份文献综述报告，包含了相关领域的最新研究进展和分析。',
          kbm_name: 'literature_review',
          kbm_level: 4,
          created_at: '2026-03-01T10:00:00Z',
          updated_at: '2026-03-01T10:00:00Z'
        }
      ]
      total.value = evidences.value.length
    } else {
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

const handleFileChange = (file: UploadUserFile) => {
  selectedFile.value = file
}

const handleCreateEvidence = async () => {
  if (!createFormRef.value) return
  
  await createFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        ElMessage({ message: '正在上传证据...', type: 'info' })
        
        let response
        if (uploadMode.value === 'file' && selectedFile.value) {
          const formData = new FormData()
          formData.append('file', selectedFile.value.raw as File)
          formData.append('student_task_id', createForm.student_task_id)
          formData.append('type', createForm.type)
          formData.append('kbm_name', createForm.kbm_name)
          response = await evidenceApi.uploadEvidenceFile(formData)
        } else {
          response = await evidenceApi.createEvidence(createForm)
        }
        
        if (response.code === 200) {
          ElMessage({ message: '正在进行AI分析...', type: 'info' })
          try {
            const analysisResponse = await evidenceApi.analyzeEvidence(response.data.id)
            if (analysisResponse.code === 200) {
              ElMessage.success({
                message: '上传证据成功！AI分析已完成',
                duration: 5000
              })
            } else {
              ElMessage.success('上传证据成功')
            }
          } catch (analysisError) {
            ElMessage.success('上传证据成功，但AI分析失败')
          }
          showCreateDialog.value = false
          loadEvidences()
          resetCreateForm()
        }
      } catch (error) {
        console.error('上传证据失败:', error)
        ElMessage.error('上传证据失败')
      }
    }
  })
}

const resetCreateForm = () => {
  showCreateDialog.value = false
  Object.assign(createForm, {
    student_task_id: '',
    type: '',
    content: '',
    kbm_name: ''
  })
  uploadMode.value = 'text'
  selectedFile.value = null
  if (uploadRef.value) {
    uploadRef.value.clearFiles()
  }
}

const viewEvidenceDetails = async (evidence: Evidence) => {
  try {
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      evidenceDetails.value = evidence
    } else {
      const response = await evidenceApi.getEvidenceById(evidence.id)
      if (response.code === 200) {
        evidenceDetails.value = { ...response.data, ...evidence }
      }
    }
    showEvidenceDialog.value = true
  } catch (error) {
    console.error('获取证据详情失败:', error)
    ElMessage.error('获取证据详情失败')
  }
}

const analyzeEvidence = async (evidenceId: string) => {
  try {
    await ElMessageBox.confirm('确定要对该证据进行AI分析吗？', '分析确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    ElMessage({ message: '正在进行AI分析...', type: 'info' })
    const response = await evidenceApi.analyzeEvidence(evidenceId)
    if (response.code === 200) {
      ElMessage.success('AI分析完成！')
      loadEvidences()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('AI分析失败:', error)
      ElMessage.error('AI分析失败')
    }
  }
}

const viewFeedback = async (evidenceId: string) => {
  try {
    const response = await evidenceApi.getFeedbackByEvidenceId(evidenceId)
    if (response.code === 200) {
      feedbackDetails.value = response.data
      showFeedbackDialog.value = true
    }
  } catch (error) {
    console.error('获取反馈失败:', error)
    ElMessage.error('获取反馈失败')
  }
}

const handleDownloadFile = async (evidenceId: string) => {
  try {
    await evidenceApi.downloadEvidenceFile(evidenceId)
    ElMessage.success('下载成功')
  } catch (error) {
    console.error('下载文件失败:', error)
    ElMessage.error('下载文件失败')
  }
}

const deleteEvidence = async (evidenceId: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个证据吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      evidences.value = evidences.value.filter(e => e.id !== evidenceId)
      total.value = evidences.value.length
    } else {
      const response = await evidenceApi.deleteEvidence(evidenceId)
      if (response.code === 200) {
        loadEvidences()
      }
    }
    ElMessage.success('删除证据成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除证据失败:', error)
      ElMessage.error('删除证据失败')
    }
  }
}

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

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>
