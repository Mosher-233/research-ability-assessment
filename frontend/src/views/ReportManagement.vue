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
            <el-table-column prop="student_id" label="学生ID" width="120" />
            <el-table-column prop="student_name" label="学生姓名" width="120">
              <template #default="scope">
                {{ scope.row.student_name || '未知' }}
              </template>
            </el-table-column>
            <el-table-column prop="task_id" label="任务ID" width="120" />
            <el-table-column prop="task_name" label="任务名称">
              <template #default="scope">
                {{ scope.row.task_name || '未知' }}
              </template>
            </el-table-column>
            <el-table-column prop="overall_score" label="总体得分" width="120">
              <template #default="scope">
                {{ formatScore(scope.row.overall_score) }}
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="生成时间" width="180">
              <template #default="scope">
                {{ formatDateTime(scope.row.created_at) }}
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
          
          <div v-if="reports.length === 0" class="empty-state">
            <el-empty description="暂无报告数据" />
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="生成报告" name="generate">
          <el-form :model="reportForm" label-width="120px">
            <el-form-item label="学生选择">
              <el-select v-model="reportForm.studentId" placeholder="请选择学生" style="width: 350px">
                <el-option v-for="student in students" :key="student.id" :value="student.id">
                  <span>{{ student.name }}</span>
                  <span style="color: #909399; margin-left: 10px; font-size: 12px;">{{ student.id }}</span>
                </el-option>
              </el-select>
            </el-form-item>
            
            <el-form-item label="任务选择">
              <el-select v-model="reportForm.taskId" placeholder="请选择任务" style="width: 400px">
                <el-option v-for="task in tasks" :key="task.id" :value="task.id">
                  <span>{{ task.title }}</span>
                  <span style="color: #909399; margin-left: 10px; font-size: 12px;">{{ task.id }}</span>
                </el-option>
              </el-select>
            </el-form-item>
            
            <el-form-item label="报告类型">
              <el-radio-group v-model="reportForm.reportType">
                <el-radio label="detailed">详细报告</el-radio>
                <el-radio label="summary">摘要报告</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="submitReportForm" :loading="generating">生成报告</el-button>
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
        <h3>研究能力评价报告</h3>
        <div class="report-header">
          <div class="info-bar">
            <div class="info-item">
              <span class="info-label">报告编号:</span>
              <span class="info-value">{{ currentReport.id }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">学生ID:</span>
              <span class="info-value">{{ currentReport.student_id }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">学生姓名:</span>
              <span class="info-value">{{ currentReport.student_name || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">任务ID:</span>
              <span class="info-value">{{ currentReport.task_id }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">任务名称:</span>
              <span class="info-value">{{ currentReport.task_name || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">生成时间:</span>
              <span class="info-value">{{ formatDateTime(currentReport.created_at) }}</span>
            </div>
          </div>
        </div>
        
        <div class="report-section">
          <h4>综合评价</h4>
          <div class="overall-score">
            <div class="score-circle">
              <div class="score-number">{{ formatScore(currentReport.overall_score) }}</div>
              <div class="score-level">{{ currentReport.overall_level }}</div>
            </div>
            <div class="rank-info" v-if="currentReport.rank > 0">
              <p>班级排名: 第 {{ currentReport.rank }} 名</p>
              <p>超越比例: {{ currentReport.percentile.toFixed(1) }}%</p>
            </div>
          </div>
        </div>
        
        <div class="report-section">
          <h4>能力分析</h4>
          <div class="ability-scores" v-if="abilityScores.length > 0">
            <div v-for="score in abilityScores" :key="score.name" class="ability-item">
              <span class="ability-name">{{ score.name }}:</span>
              <el-progress :percentage="score.score" :color="getScoreColor(score.score)" />
            </div>
          </div>
          <div v-else class="empty-data">
            <el-empty description="暂无能力分析数据" :image-size="80" />
          </div>
        </div>
        
        <div class="report-section">
          <h4>能力雷达图</h4>
          <div class="chart-container">
            <div ref="radarChartRef" class="radar-chart"></div>
          </div>
        </div>
        
        <div class="report-section">
          <h4>优势分析</h4>
          <div class="strengths" v-if="strengths.length > 0">
            <ul>
              <li v-for="(strength, index) in strengths" :key="index">
                {{ strength }}
              </li>
            </ul>
          </div>
          <div v-else class="empty-data">
            <el-empty description="暂无优势分析数据" :image-size="80" />
          </div>
        </div>
        
        <div class="report-section">
          <h4>待提升方向</h4>
          <div class="weaknesses" v-if="weaknesses.length > 0">
            <ul>
              <li v-for="(weakness, index) in weaknesses" :key="index">
                {{ weakness }}
              </li>
            </ul>
          </div>
          <div v-else class="empty-data">
            <el-empty description="暂无待提升方向数据" :image-size="80" />
          </div>
        </div>
        
        <div class="report-section">
          <h4>改进建议</h4>
          <div class="suggestions" v-if="suggestions.length > 0">
            <div v-for="suggestion in suggestions" :key="suggestion.id" class="suggestion-item">
              <h5>【{{ suggestion.dimension_name }}】(优先级: {{ suggestion.priority }})</h5>
              <p>当前得分: {{ formatScore(suggestion.current_score) }} → 目标得分: {{ formatScore(suggestion.target_score) }}</p>
              <p class="suggestion-text">{{ suggestion.suggestion }}</p>
              <div class="action-items">
                <strong>行动项:</strong>
                <ul>
                  <li v-for="(item, index) in suggestion.action_items" :key="index">
                    {{ item }}
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div v-else class="empty-data">
            <el-empty description="暂无改进建议数据" :image-size="80" />
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
import { ref, onMounted, watch } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatDateTime, formatScore } from '../utils/format'
import * as echarts from 'echarts'
import { generateRadarData } from '../utils/chart'
import { resultApi, type Report } from '../api/result'
import { taskApi } from '../api/task'

// 类型定义
interface Student {
  id: number
  name: string
}

interface Task {
  id: number
  title: string
}

interface AbilityScore {
  name: string
  score: number
}

interface Suggestion {
  id: string
  dimension: string
  dimension_name: string
  current_score: number
  target_score: number
  suggestion: string
  action_items: string[]
  resources: any[]
  priority: number
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
const generating = ref(false)

// 报告详情数据
const abilityScores = ref<AbilityScore[]>([])
const strengths = ref<string[]>([])
const weaknesses = ref<string[]>([])
const suggestions = ref<Suggestion[]>([])

// 雷达图相关
const radarChartRef = ref<HTMLElement | null>(null)
let radarChart: echarts.ECharts | null = null

// 维度名称映射
const getDimensionName = (dimId: string): string => {
  const dimMap: Record<string, string> = {
    'literature_review': '文献综述',
    'research_design': '研究设计',
    'data_analysis': '数据分析',
    'critical_thinking': '批判性思维',
    'literature': '文献综述',
    'experiment_design': '研究设计',
    'data_processing': '数据分析',
    'innovation': '创新能力'
  }
  return dimMap[dimId] || dimId
}

// 方法
const loadReports = async () => {
  try {
    const userStr = localStorage.getItem('user')
    const user = userStr ? JSON.parse(userStr) : null
    
    if (user && user.role === 'student') {
      const response = await resultApi.getStudentReports()
      if (response.code === 200) {
        reports.value = response.data
      }
    } else {
      const response = await resultApi.getReports()
      if (response.code === 200) {
        reports.value = response.data
      }
    }
  } catch (error) {
    ElMessage.error('获取报告列表失败')
    console.error('Error loading reports:', error)
  }
}

const loadStudents = async () => {
  try {
    const response = await taskApi.getStudents()
    if (response.code === 200) {
      students.value = response.data.map((student: any) => ({
        id: student.student_id || student.id,
        name: student.name
      }))
    }
  } catch (error) {
    ElMessage.error('获取学生列表失败')
    console.error('Error loading students:', error)
  }
}

const loadTasks = async () => {
  try {
    const response = await taskApi.getTasks()
    if (response.code === 200) {
      tasks.value = response.data.map((task: any) => ({
        id: task.id,
        title: task.name
      }))
    }
  } catch (error) {
    ElMessage.error('获取任务列表失败')
    console.error('Error loading tasks:', error)
  }
}

const generateReport = () => {
  const userStr = localStorage.getItem('user')
  const user = userStr ? JSON.parse(userStr) : null
  
  if (user && user.role === 'student') {
    ElMessage.info('您只能查看自己的报告')
  } else {
    activeTab.value = 'generate'
  }
}

const submitReportForm = async () => {
  try {
    if (!reportForm.value.studentId || !reportForm.value.taskId) {
      ElMessage.warning('请选择学生和任务')
      return
    }
    
    generating.value = true
    const response = await resultApi.generateStudentReport(
      reportForm.value.studentId,
      reportForm.value.taskId
    )
    
    if (response.code === 200) {
      ElMessage.success('报告生成成功！')
      resetReportForm()
      activeTab.value = 'list'
      await loadReports()
    }
  } catch (error) {
    ElMessage.error('生成报告失败')
    console.error('Error generating report:', error)
  } finally {
    generating.value = false
  }
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
  
  abilityScores.value = []
  strengths.value = []
  weaknesses.value = []
  suggestions.value = []
  
  if (report.dimension_scores) {
    try {
      const dimScores = typeof report.dimension_scores === 'string' 
        ? JSON.parse(report.dimension_scores) 
        : report.dimension_scores
      
      Object.entries(dimScores).forEach(([dimId, score]: [string, any]) => {
        const normalizedScore = score.score <= 1 ? score.score * 100 : score.score
        const displayName = getDimensionName(dimId)
        abilityScores.value.push({
          name: displayName,
          score: Math.round(normalizedScore)
        })
      })
    } catch (e) {
      console.error('解析维度得分失败:', e)
    }
  }
  
  if (report.strengths) {
    try {
      strengths.value = typeof report.strengths === 'string' 
        ? JSON.parse(report.strengths) 
        : report.strengths
    } catch (e) {
      console.error('解析优势分析失败:', e)
    }
  }
  
  if (report.weaknesses) {
    try {
      weaknesses.value = typeof report.weaknesses === 'string' 
        ? JSON.parse(report.weaknesses) 
        : report.weaknesses
    } catch (e) {
      console.error('解析待提升方向失败:', e)
    }
  }
  
  if (report.suggestions) {
    try {
      suggestions.value = typeof report.suggestions === 'string' 
        ? JSON.parse(report.suggestions) 
        : report.suggestions
    } catch (e) {
      console.error('解析改进建议失败:', e)
    }
  }
  
  reportDialogVisible.value = true
}

const downloadReport = (report: Report | null) => {
  if (!report) return
  ElMessage.success('报告下载成功！')
}

const getScoreColor = (score: number): string => {
  if (score >= 90) return '#67c23a'
  if (score >= 75) return '#e6a23c'
  return '#f56c6c'
}

const initRadarChart = () => {
  if (!radarChartRef.value) return
  
  radarChart = echarts.init(radarChartRef.value)
  updateRadarChart()
}

const updateRadarChart = () => {
  if (!radarChart || !currentReport.value) return
  
  const normalizedScores: Record<string, { score: number }> = {}
  abilityScores.value.forEach(score => {
    normalizedScores[score.name] = { score: score.score }
  })
  
  if (Object.keys(normalizedScores).length > 0) {
    const chartData = generateRadarData(normalizedScores)
    const option = {
      radar: chartData.radar,
      series: [{
        type: 'radar',
        data: chartData.series[0].data,
        areaStyle: {
          opacity: 0.2
        },
        lineStyle: {
          width: 2
        },
        itemStyle: {
          color: '#409EFF'
        }
      }]
    }
    
    radarChart.setOption(option)
  }
}

watch(reportDialogVisible, (newVal) => {
  if (newVal) {
    setTimeout(() => {
      initRadarChart()
    }, 100)
  } else {
    if (radarChart) {
      radarChart.dispose()
      radarChart = null
    }
  }
})

watch(currentReport, () => {
  updateRadarChart()
}, { deep: true })

onMounted(async () => {
  await loadReports()
  await loadStudents()
  await loadTasks()
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

.report-detail h3 {
  text-align: center;
  margin-bottom: 20px;
  color: #303133;
}

.info-bar {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.info-label {
  font-size: 12px;
  color: #909399;
}

.info-value {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.report-section {
  margin-bottom: 30px;
}

.report-section h4 {
  margin-bottom: 15px;
  color: #303133;
  padding-bottom: 10px;
  border-bottom: 2px solid #409EFF;
}

.overall-score {
  display: flex;
  align-items: center;
  gap: 40px;
  padding: 20px;
}

.score-circle {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.score-number {
  font-size: 32px;
  font-weight: bold;
}

.score-level {
  font-size: 14px;
  margin-top: 5px;
}

.rank-info {
  flex: 1;
}

.rank-info p {
  margin: 8px 0;
  font-size: 16px;
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

.radar-chart {
  width: 100%;
  height: 100%;
}

.strengths ul,
.weaknesses ul {
  padding-left: 20px;
  line-height: 1.8;
  color: #606266;
}

.strengths li,
.weaknesses li {
  margin-bottom: 8px;
}

.suggestion-item {
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 15px;
}

.suggestion-item h5 {
  margin: 0 0 10px 0;
  color: #409EFF;
}

.suggestion-text {
  color: #606266;
  margin: 10px 0;
}

.action-items {
  margin-top: 15px;
}

.action-items ul {
  padding-left: 20px;
  margin-top: 5px;
}

.action-items li {
  margin: 5px 0;
  color: #606266;
}

.empty-state,
.empty-data {
  padding: 40px 0;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
