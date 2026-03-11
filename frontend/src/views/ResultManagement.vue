<template>
  <div class="result-management">
    <div class="header">
      <h2>结果管理</h2>
    </div>

    <!-- 结果列表 -->
    <el-card class="result-list">
      <el-table :data="results" style="width: 100%">
        <el-table-column prop="id" label="结果ID" width="180"></el-table-column>
        <el-table-column prop="student.name" label="学生姓名" width="150"></el-table-column>
        <el-table-column prop="student.student_id" label="学生ID" width="120"></el-table-column>
        <el-table-column prop="task.name" label="任务名称" width="200"></el-table-column>
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
    <el-dialog title="结果详情" v-model="showResultDialog" width="800px">
      <el-descriptions :column="2">
        <el-descriptions-item label="结果ID">{{ resultDetails?.id }}</el-descriptions-item>
        <el-descriptions-item label="学生姓名">{{ resultDetails?.student?.name }}</el-descriptions-item>
        <el-descriptions-item label="学生ID">{{ resultDetails?.student?.student_id }}</el-descriptions-item>
        <el-descriptions-item label="任务名称">{{ resultDetails?.task?.name }}</el-descriptions-item>
        <el-descriptions-item label="总体得分">{{ formatScore(resultDetails?.overall_score) }}</el-descriptions-item>
        <el-descriptions-item label="总体等级"><el-tag :type="getLevelType(resultDetails?.overall_level || '')">{{ resultDetails?.overall_level }}</el-tag></el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDateTime(resultDetails?.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="推理理由">{{ resultDetails?.reasoning }}</el-descriptions-item>
      </el-descriptions>

      <h3 style="margin-top: 20px">维度得分</h3>
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

      <h3 style="margin-top: 20px">能力雷达图</h3>
      <div class="chart-container">
        <!-- 这里可以集成ECharts等图表库来展示雷达图 -->
        <div class="placeholder-chart">
          <p>能力雷达图将在此显示</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { resultApi } from '../api/result'
import type { InferenceResult, DimensionScore } from '../types/result'
import { formatDateTime, formatScore } from '../utils/format'

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

// 获取等级类型
const getLevelType = (level: string): string => {
  const typeMap: Record<string, string> = {
    '优秀': 'success',
    '良好': 'primary',
    '中等': 'warning',
    '待提高': 'danger'
  }
  return typeMap[level] || 'info'
}

// 加载结果列表
const loadResults = async () => {
  try {
    // 这里简化处理，实际应该调用API获取结果列表
    // 暂时使用模拟数据
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
    total.value = results.value.length
  } catch (error) {
    console.error('加载结果失败:', error)
    ElMessage.error('加载结果失败')
  }
}

// 查看结果详情
const viewResultDetails = async (resultId: string) => {
  try {
    // 这里简化处理，实际应该调用API获取结果详情
    // 暂时使用模拟数据
    const result = results.value.find(r => r.id === resultId)
    if (result) {
      resultDetails.value = result
      // 转换维度得分格式
      dimensionScores.value = Object.values(result.dimension_scores)
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

// 初始化
onMounted(() => {
  loadResults()
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

.chart-container {
  margin-top: 20px;
  height: 400px;
  background-color: #f9f9f9;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-chart {
  text-align: center;
  color: #999;
}
</style>