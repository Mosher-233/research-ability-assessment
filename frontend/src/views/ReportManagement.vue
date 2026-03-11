<template>
  <div class="report-management">
    <el-card shadow="hover" class="card">
      <template #header>
        <div class="card-header">
          <span>报告管理</span>
          <el-button type="primary" @click="generateReport">
            <el-icon><Plus /></el-icon> 生成报告
          </el-button>
        </div>
      </template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="报告列表" name="list">
          <el-table :data="reports" style="width: 100%">
            <el-table-column prop="id" label="报告ID" width="80" />
            <el-table-column prop="studentId" label="学生ID" width="120" />
            <el-table-column prop="studentName" label="学生姓名" width="120" />
            <el-table-column prop="taskId" label="任务ID" width="120" />
            <el-table-column prop="taskName" label="任务名称" />
            <el-table-column prop="createdAt" label="生成时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="small" type="primary" @click="viewReport(scope.row)">
                  查看
                </el-button>
                <el-button size="small" type="success" @click="downloadReport(scope.row)">
                  下载
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="生成报告" name="generate">
          <el-form :model="reportForm" label-width="120px">
            <el-form-item label="学生选择">
              <el-select v-model="reportForm.studentId" placeholder="请选择学生">
                <el-option v-for="student in students" :key="student.id" :label="student.name" :value="student.id" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="任务选择">
              <el-select v-model="reportForm.taskId" placeholder="请选择任务">
                <el-option v-for="task in tasks" :key="task.id" :label="task.title" :value="task.id" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="报告类型">
              <el-radio-group v-model="reportForm.reportType">
                <el-radio label="detailed">详细报告</el-radio>
                <el-radio label="summary">摘要报告</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="submitReportForm">生成报告</el-button>
              <el-button @click="resetReportForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
    
    <!-- 报告详情对话框 -->
    <el-dialog
      v-model="reportDialogVisible"
      title="报告详情"
      width="80%"
    >
      <div v-if="currentReport" class="report-detail">
        <h3>{{ currentReport.taskName }} - 研究能力评价报告</h3>
        <div class="report-header">
          <p><strong>学生:</strong> {{ currentReport.studentName }}</p>
          <p><strong>生成时间:</strong> {{ formatDate(currentReport.createdAt) }}</p>
        </div>
        
        <div class="report-section">
          <h4>能力分析</h4>
          <div class="ability-scores">
            <div v-for="(score, key) in currentReport.abilityScores" :key="key" class="ability-item">
              <span class="ability-name">{{ getAbilityName(key) }}:</span>
              <el-progress :percentage="score" :color="getScoreColor(score)" />
            </div>
          </div>
        </div>
        
        <div class="report-section">
          <h4>能力雷达图</h4>
          <div class="chart-container">
            <!-- 雷达图占位符 -->
            <div class="chart-placeholder">
              <el-icon class="placeholder-icon"><PieChart /></el-icon>
              <p>能力雷达图</p>
            </div>
          </div>
        </div>
        
        <div class="report-section">
          <h4>详细分析</h4>
          <div class="analysis-content">
            <p>{{ currentReport.analysis }}</p>
          </div>
        </div>
        
        <div class="report-section">
          <h4>改进建议</h4>
          <div class="suggestions">
            <ul>
              <li v-for="(suggestion, index) in currentReport.suggestions" :key="index">
                {{ suggestion }}
              </li>
            </ul>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="reportDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="downloadReport(currentReport)">
            下载报告
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, PieChart } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatDate } from '../utils/format'

// 类型定义
interface Report {
  id: number
  studentId: number
  studentName: string
  taskId: number
  taskName: string
  reportType: string
  abilityScores: Record<string, number>
  analysis: string
  suggestions: string[]
  createdAt: string
}

interface Student {
  id: number
  name: string
}

interface Task {
  id: number
  title: string
}

// 响应式数据
const activeTab = ref('list')
const reports = ref<Report[]>([])
const students = ref<Student[]>([])
const tasks = ref<Task[]>([])
const reportForm = ref({
  studentId: '',
  taskId: '',
  reportType: 'detailed'
})
const reportDialogVisible = ref(false)
const currentReport = ref<Report | null>(null)

// 模拟数据
const mockReports: Report[] = [
  {
    id: 1,
    studentId: 1,
    studentName: '张三',
    taskId: 1,
    taskName: '研究方法课程作业',
    reportType: 'detailed',
    abilityScores: {
      literatureReview: 85,
      researchDesign: 78,
      dataAnalysis: 90,
      criticalThinking: 82,
      academicWriting: 75
    },
    analysis: '该学生在数据分析方面表现突出，能够熟练运用统计方法处理数据。文献综述能力较强，能够全面收集和整理相关研究。研究设计能力有待提高，需要加强研究方案的逻辑性和可行性。批判性思维能力良好，能够对研究问题进行深入分析。学术写作能力一般，需要提高论文的结构和表达。',
    suggestions: [
      '加强研究设计的训练，学习如何制定更加科学合理的研究方案',
      '提高学术写作能力，注重论文的结构和逻辑',
      '继续保持数据分析的优势，尝试学习更多高级分析方法',
      '在文献综述中更加注重批判性分析，而不仅仅是整理',
      '参与更多学术讨论，提高批判性思维能力'
    ],
    createdAt: '2024-01-15T10:30:00Z'
  },
  {
    id: 2,
    studentId: 2,
    studentName: '李四',
    taskId: 2,
    taskName: '社会调查研究',
    reportType: 'summary',
    abilityScores: {
      literatureReview: 70,
      researchDesign: 85,
      dataAnalysis: 75,
      criticalThinking: 80,
      academicWriting: 85
    },
    analysis: '该学生在研究设计和学术写作方面表现优秀，能够制定合理的研究方案并清晰表达研究结果。文献综述能力一般，需要加强对相关研究的系统梳理。数据分析能力有待提高，需要学习更多数据分析方法。批判性思维能力良好，能够对研究问题进行深入思考。',
    suggestions: [
      '加强文献综述的训练，学习如何系统梳理和分析相关研究',
      '提高数据分析能力，学习更多统计分析方法',
      '继续保持研究设计和学术写作的优势',
      '在研究过程中更加注重理论与实践的结合',
      '多参加学术交流活动，拓展研究视野'
    ],
    createdAt: '2024-01-16T14:20:00Z'
  }
]

const mockStudents: Student[] = [
  { id: 1, name: '张三' },
  { id: 2, name: '李四' },
  { id: 3, name: '王五' }
]

const mockTasks: Task[] = [
  { id: 1, title: '研究方法课程作业' },
  { id: 2, title: '社会调查研究' },
  { id: 3, title: '实验设计与分析' }
]

// 方法
const loadReports = () => {
  // 模拟API调用
  reports.value = mockReports
}

const loadStudents = () => {
  // 模拟API调用
  students.value = mockStudents
}

const loadTasks = () => {
  // 模拟API调用
  tasks.value = mockTasks
}

const generateReport = () => {
  activeTab.value = 'generate'
}

const submitReportForm = () => {
  // 模拟API调用
  ElMessage.success('报告生成成功！')
  resetReportForm()
  activeTab.value = 'list'
  loadReports()
}

const resetReportForm = () => {
  reportForm.value = {
    studentId: '',
    taskId: '',
    reportType: 'detailed'
  }
}

const viewReport = (report: Report) => {
  currentReport.value = report
  reportDialogVisible.value = true
}

const downloadReport = (report: Report | null) => {
  if (!report) return
  // 模拟下载
  ElMessage.success('报告下载成功！')
}

const getAbilityName = (key: string): string => {
  const abilityMap: Record<string, string> = {
    literatureReview: '文献综述',
    researchDesign: '研究设计',
    dataAnalysis: '数据分析',
    criticalThinking: '批判性思维',
    academicWriting: '学术写作'
  }
  return abilityMap[key] || key
}

const getScoreColor = (score: number): string => {
  if (score >= 90) return '#67c23a'
  if (score >= 75) return '#e6a23c'
  return '#f56c6c'
}

// 生命周期
onMounted(() => {
  loadReports()
  loadStudents()
  loadTasks()
})
</script>

<style scoped>
.report-management {
  padding: 20px;
}

.card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.report-detail {
  padding: 20px;
}

.report-header {
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eaeaea;
}

.report-section {
  margin-bottom: 30px;
}

.report-section h4 {
  margin-bottom: 15px;
  color: #303133;
}

.ability-scores {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.ability-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.ability-name {
  font-weight: 500;
  color: #606266;
}

.chart-container {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.chart-placeholder {
  text-align: center;
  color: #909399;
}

.placeholder-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.analysis-content {
  line-height: 1.6;
  color: #606266;
}

.suggestions ul {
  padding-left: 20px;
  line-height: 1.6;
  color: #606266;
}

.suggestions li {
  margin-bottom: 8px;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>