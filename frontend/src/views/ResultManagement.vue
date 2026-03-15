<template>
  <div class="result-management">
    <div class="header">
      <h2>结果管理</h2>
    </div>

    <!-- 学生任务选择区域 -->
    <el-card class="task-selection" v-if="isStudent">
      <template #header>
        <div class="card-header">
          <span>选择任务并生成推理结果</span>
        </div>
      </template>
      <el-form :inline="true" :model="taskSelectionForm">
        <el-form-item label="选择任务">
          <el-select v-model="taskSelectionForm.selectedTaskId" placeholder="请选择任务" style="width: 450px">
            <el-option
              v-for="task in assignedTasks"
              :key="task.id"
              :value="task.id"
            >
              <span>{{ task.name }}</span>
              <span style="color: #909399; margin-left: 10px; font-size: 12px;">{{ task.id }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="generateInference" :loading="generating">
            生成推理结果
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 结果列表 -->
    <el-card class="result-list">
      <el-table :data="results" style="width: 100%">
        <el-table-column prop="id" label="结果ID" width="180"></el-table-column>
        <el-table-column label="学生姓名" width="150">
          <template #default="scope">
            {{ scope.row.student?.name || scope.row.student_name || '未知' }}
          </template>
        </el-table-column>
        <el-table-column label="学生ID" width="150">
          <template #default="scope">
            {{ scope.row.student?.student_id || scope.row.student_id || '未知' }}
          </template>
        </el-table-column>
        <el-table-column label="任务名称" width="200">
          <template #default="scope">
            {{ scope.row.task?.name || scope.row.task_name || '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="overall_score" label="总体得分" width="100">
          <template #default="scope">
            {{ formatScore(scope.row.overall_score) }}
          </template>
        </el-table-column>
        <el-table-column prop="overall_level" label="总体等级" width="100">
          <template #default="scope">
            <el-tag :type="getLevelType(scope.row.overall_level)">{{ scope.row.overall_level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" @click="viewResultDetails(scope.row.id)">查看</el-button>
            <el-button size="small" type="primary" @click="generateReport(scope.row.id)">报告</el-button>
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
          @size-change="loadResults"
          @current-change="loadResults"
        />
      </div>
    </el-card>

    <!-- 结果详情对话框 -->
    <el-dialog title="结果详情" v-model="showResultDialog" width="900px">
      <div class="result-detail">
        <h3>研究能力评估结果</h3>
        <div class="info-bar">
          <div class="info-item">
            <span class="info-label">结果ID:</span>
            <span class="info-value">{{ resultDetails?.id }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">学生姓名:</span>
            <span class="info-value">{{ resultDetails?.student?.name || '未知' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">学生ID:</span>
            <span class="info-value">{{ resultDetails?.student?.student_id || '未知' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">任务名称:</span>
            <span class="info-value">{{ resultDetails?.task?.name || '未知' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">创建时间:</span>
            <span class="info-value">{{ formatDateTime(resultDetails?.created_at || '') }}</span>
          </div>
        </div>
        
        <div class="overall-section">
          <div class="score-display">
            <div class="score-circle">
              <div class="score-number">{{ formatScore(resultDetails?.overall_score || 0) }}</div>
              <div class="score-level">
                <el-tag :type="getLevelType(resultDetails?.overall_level || '')">{{ resultDetails?.overall_level || '未评级' }}</el-tag>
              </div>
            </div>
          </div>
        </div>
        
        <div class="reasoning-section">
          <h4>推理理由</h4>
          <p>{{ resultDetails?.reasoning || '暂无推理理由' }}</p>
        </div>

        <div class="dimension-section">
          <h4>维度得分</h4>
          <el-table :data="dimensionScores" style="width: 100%">
            <el-table-column prop="name" label="维度名称" width="150"></el-table-column>
            <el-table-column prop="score" label="得分" width="100">
              <template #default="scope">
                {{ formatScore(scope.row.score) }}
              </template>
            </el-table-column>
            <el-table-column prop="level" label="等级" width="100">
              <template #default="scope">
                <el-tag :type="getLevelType(scope.row.level)">{{ scope.row.level }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="details" label="详情" show-overflow-tooltip></el-table-column>
          </el-table>
        </div>

        <div class="radar-section">
          <h4>能力雷达图</h4>
          <div class="chart-container">
            <div ref="radarChartRef" class="radar-chart"></div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { resultApi } from '../api/result'
import { taskApi } from '../api/task'
import type { InferenceResult, DimensionScore } from '../types/result'
import type { Task } from '../types/task'
import { formatDateTime, formatScore } from '../utils/format'
import * as echarts from 'echarts'
import { generateRadarData } from '../utils/chart'

// 结果列表相关
const results = ref<InferenceResult[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框相关
const showResultDialog = ref(false)

// 结果详情
const resultDetails = ref<InferenceResult | null>(null)
const dimensionScores = ref<DimensionScore[]>([])

// 雷达图相关
const radarChartRef = ref<HTMLElement | null>(null)
let radarChart: echarts.ECharts | null = null

// 用户和任务相关
const user = ref<any>(null)
const isStudent = ref(false)
const assignedTasks = ref<Task[]>([])
const taskSelectionForm = ref({
  selectedTaskId: ''
})
const generating = ref(false)

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

// 维度详情补充
const getDimensionDetails = (dimId: string, details: string): string => {
  if (details) return details
  
  const defaultDetails: Record<string, string> = {
    'literature_review': '文献检索、综述撰写能力评估',
    'research_design': '研究方案设计、实验规划能力评估',
    'data_analysis': '数据处理、统计分析能力评估',
    'critical_thinking': '批判性思考、创新思维能力评估',
    'literature': '文献检索、综述撰写能力评估',
    'experiment_design': '实验方案设计、变量控制能力评估',
    'data_processing': '数据处理、分析方法应用能力评估',
    'innovation': '问题提出、解决方案原创性评估'
  }
  return defaultDetails[dimId] || '该维度的综合能力评估'
}

// 获取等级类型
const getLevelType = (level: string): string => {
  const typeMap: Record<string, string> = {
    '优秀': 'success',
    '良好': 'primary',
    '中等': 'warning',
    '待提高': 'danger',
    '及格': 'info',
    '不及格': 'danger'
  }
  return typeMap[level] || 'info'
}

// 加载用户信息
const loadUserInfo = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
    isStudent.value = user.value.role === 'student'
  }
}

// 加载学生分配的任务
const loadAssignedTasks = async () => {
  if (!isStudent.value) return
  
  try {
    const response = await taskApi.getAssignedTasks()
    if (response.code === 200) {
      assignedTasks.value = response.data
    }
  } catch (error) {
    console.error('加载任务失败:', error)
  }
}

// 加载结果列表
const loadResults = async () => {
  try {
    // 检查是否使用模拟数据
    const token = localStorage.getItem('token')
    if (token && token.startsWith('mock-token-')) {
      // 使用模拟数据
      if (isStudent.value) {
        // 学生用户只显示自己的结果
        results.value = [
          {
            id: '1',
            student_id: 's1',
            task_id: 't1',
            overall_score: 0.85,
            overall_level: '优秀',
            dimension_scores: {
              literature: {
                name: 'literature',
                score: 0.9,
                level: '优秀',
                details: '文献检索策略合理，综述质量高',
                evidence_ids: ['e1', 'e2']
              },
              experiment_design: {
                name: 'experiment_design',
                score: 0.8,
                level: '优秀',
                details: '实验方案合理，变量控制较好',
                evidence_ids: ['e3']
              },
              data_processing: {
                name: 'data_processing',
                score: 0.85,
                level: '优秀',
                details: '数据分析方法选择恰当，结果解释准确',
                evidence_ids: ['e4']
              },
              innovation: {
                name: 'innovation',
                score: 0.75,
                level: '良好',
                details: '问题提出有一定新颖性，解决方案有一定原创性',
                evidence_ids: ['e5']
              }
            },
            reasoning: '基于收集到的证据，对学生的研究能力进行了综合评估。总体得分为0.85，等级为优秀。各维度表现均衡，文献能力和实验设计能力突出。',
            created_at: '2026-03-01T10:00:00Z',
            updated_at: '2026-03-01T10:00:00Z',
            student: {
              id: 's1',
              name: user.value?.name,
              student_id: user.value?.student_id
            },
            task: {
              id: 't1',
              name: '研究方法课程作业'
            }
          }
        ]
      } else {
        // 教师用户显示所有结果
        results.value = [
          {
            id: '1',
            student_id: 's1',
            task_id: 't1',
            overall_score: 0.85,
            overall_level: '优秀',
            dimension_scores: {
              literature: {
                name: 'literature',
                score: 0.9,
                level: '优秀',
                details: '文献检索策略合理，综述质量高',
                evidence_ids: ['e1', 'e2']
              },
              experiment_design: {
                name: 'experiment_design',
                score: 0.8,
                level: '优秀',
                details: '实验方案合理，变量控制较好',
                evidence_ids: ['e3']
              },
              data_processing: {
                name: 'data_processing',
                score: 0.85,
                level: '优秀',
                details: '数据分析方法选择恰当，结果解释准确',
                evidence_ids: ['e4']
              },
              innovation: {
                name: 'innovation',
                score: 0.75,
                level: '良好',
                details: '问题提出有一定新颖性，解决方案有一定原创性',
                evidence_ids: ['e5']
              }
            },
            reasoning: '基于收集到的证据，对学生的研究能力进行了综合评估。总体得分为0.85，等级为优秀。各维度表现均衡，文献能力和实验设计能力突出。',
            created_at: '2026-03-01T10:00:00Z',
            updated_at: '2026-03-01T10:00:00Z',
            student: {
              id: 's1',
              name: '张三',
              student_id: '2022001'
            },
            task: {
              id: 't1',
              name: '2026春季学期研究能力评价'
            }
          },
          {
            id: '2',
            student_id: 's2',
            task_id: 't1',
            overall_score: 0.72,
            overall_level: '良好',
            dimension_scores: {
              literature: {
                name: 'literature',
                score: 0.75,
                level: '良好',
                details: '文献检索策略基本合理，综述质量一般',
                evidence_ids: ['e6']
              },
              experiment_design: {
                name: 'experiment_design',
                score: 0.7,
                level: '良好',
                details: '实验方案基本合理，变量控制一般',
                evidence_ids: ['e7']
              },
              data_processing: {
                name: 'data_processing',
                score: 0.75,
                level: '良好',
                details: '数据分析方法选择基本恰当，结果解释基本准确',
                evidence_ids: ['e8']
              },
              innovation: {
                name: 'innovation',
                score: 0.65,
                level: '良好',
                details: '问题提出缺乏新颖性，解决方案缺乏原创性',
                evidence_ids: ['e9']
              }
            },
            reasoning: '基于收集到的证据，对学生的研究能力进行了综合评估。总体得分为0.72，等级为良好。各维度表现基本均衡，但创新能力有待提高。',
            created_at: '2026-03-02T14:30:00Z',
            updated_at: '2026-03-02T14:30:00Z',
            student: {
              id: 's2',
              name: '李四',
              student_id: '2022002'
            },
            task: {
              id: 't1',
              name: '2026春季学期研究能力评价'
            }
          }
        ]
      }
      total.value = results.value.length
    } else {
      // 调用后端API获取结果列表
      if (isStudent.value) {
        // 学生用户获取自己的结果
        const response = await resultApi.getStudentResults()
        if (response.code === 200) {
          results.value = response.data
          total.value = response.data.length
        }
      } else {
        // 教师用户获取所有结果
        const response = await resultApi.getResults()
        if (response.code === 200) {
          results.value = response.data
          total.value = response.data.length
        }
      }
    }
  } catch (error) {
    console.error('加载结果失败:', error)
    ElMessage.error('加载结果失败')
  }
}

// 生成推理结果
const generateInference = async () => {
  if (!taskSelectionForm.value.selectedTaskId) {
    ElMessage.warning('请选择任务')
    return
  }

  generating.value = true
  try {
    const response = await resultApi.generateStudentInference(taskSelectionForm.value.selectedTaskId)
    if (response.code === 200) {
      ElMessage.success(response.message)
      // 重新加载结果列表
      await loadResults()
    } else {
      ElMessage.error(response.message)
    }
  } catch (error: any) {
    console.error('生成推理结果失败:', error)
    ElMessage.error(error.response?.data?.message || '生成推理结果失败')
  } finally {
    generating.value = false
  }
}

// 查看结果详情
const viewResultDetails = async (resultId: string) => {
  try {
    const result = results.value.find(r => r.id === resultId)
    if (result) {
      resultDetails.value = result
      
      // 转换维度得分格式，确保名称和详情正确
      const dimScores = Object.entries(result.dimension_scores).map(([dimId, dimScore]: [string, any]) => {
        return {
          ...dimScore,
          name: getDimensionName(dimId),
          details: getDimensionDetails(dimId, dimScore.details)
        }
      })
      dimensionScores.value = dimScores
    }
    showResultDialog.value = true
  } catch (error) {
    console.error('获取结果详情失败:', error)
    ElMessage.error('获取结果详情失败')
  }
}

// 生成报告
const generateReport = async (resultId: string) => {
  try {
    // 这里简化处理，实际应该调用API生成报告
    ElMessage.success('报告生成成功')
  } catch (error) {
    console.error('生成报告失败:', error)
    ElMessage.error('生成报告失败')
  }
}

// 初始化雷达图
const initRadarChart = () => {
  if (!radarChartRef.value) return
  
  radarChart = echarts.init(radarChartRef.value)
  updateRadarChart()
}

// 更新雷达图
const updateRadarChart = () => {
  if (!radarChart || !resultDetails.value) return
  
  const chartData = generateRadarData(resultDetails.value.dimension_scores)
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

// 监听对话框显示状态
watch(showResultDialog, (newVal) => {
  if (newVal) {
    // 延迟初始化，确保DOM已经渲染
    setTimeout(() => {
      initRadarChart()
    }, 100)
  } else {
    // 对话框关闭时销毁雷达图
    if (radarChart) {
      radarChart.dispose()
      radarChart = null
    }
  }
})

// 监听结果详情变化
watch(resultDetails, () => {
  updateRadarChart()
}, { deep: true })

// 初始化
onMounted(() => {
  loadUserInfo()
  loadResults()
  loadAssignedTasks()
})
</script>

<style scoped>
.result-management {
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

.result-list {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.result-detail {
  padding: 20px;
}

.result-detail h3 {
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

.overall-section {
  display: flex;
  justify-content: center;
  margin-bottom: 25px;
}

.score-display {
  display: flex;
  justify-content: center;
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

.reasoning-section,
.dimension-section,
.radar-section {
  margin-bottom: 25px;
}

.reasoning-section h4,
.dimension-section h4,
.radar-section h4 {
  margin-bottom: 15px;
  color: #303133;
  padding-bottom: 10px;
  border-bottom: 2px solid #409EFF;
}

.reasoning-section p {
  line-height: 1.8;
  color: #606266;
}

.chart-container {
  height: 400px;
  background-color: #f5f7fa;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.radar-chart {
  width: 100%;
  height: 100%;
}

.placeholder-chart {
  text-align: center;
  color: #999;
}
</style>