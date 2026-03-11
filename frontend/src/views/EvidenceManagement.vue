<template>
  <div class="evidence-management">
    <div class="header">
      <h2>证据管理</h2>
      <el-button type="primary" @click="showCreateDialog = true">创建证据</el-button>
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
    <el-dialog title="创建证据" v-model="showCreateDialog" width="600px">
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="120px">
        <el-form-item label="学生任务ID" prop="student_task_id">
          <el-input v-model="createForm.student_task_id" placeholder="请输入学生任务ID"></el-input>
        </el-form-item>
        <el-form-item label="证据类型" prop="type">
          <el-select v-model="createForm.type" placeholder="请选择证据类型">
            <el-option label="文档" value="document"></el-option>
            <el-option label="代码" value="code"></el-option>
            <el-option label="演示" value="presentation"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="KBM名称" prop="kbm_name">
          <el-input v-model="createForm.kbm_name" placeholder="请输入KBM名称"></el-input>
        </el-form-item>
        <el-form-item label="KBM级别" prop="kbm_level">
          <el-select v-model="createForm.kbm_level" placeholder="请选择KBM级别">
            <el-option label="1" value="1"></el-option>
            <el-option label="2" value="2"></el-option>
            <el-option label="3" value="3"></el-option>
            <el-option label="4" value="4"></el-option>
            <el-option label="5" value="5"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="createForm.content" placeholder="请输入证据内容" type="textarea" :rows="4"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="handleCreateEvidence">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 证据详情对话框 -->
    <el-dialog title="证据详情" v-model="showEvidenceDialog" width="600px">
      <el-descriptions :column="2">
        <el-descriptions-item label="证据ID">{{ evidenceDetails?.id }}</el-descriptions-item>
        <el-descriptions-item label="学生任务ID">{{ evidenceDetails?.student_task_id }}</el-descriptions-item>
        <el-descriptions-item label="证据类型">{{ evidenceDetails?.type }}</el-descriptions-item>
        <el-descriptions-item label="KBM名称">{{ evidenceDetails?.kbm_name }}</el-descriptions-item>
        <el-descriptions-item label="KBM级别">{{ evidenceDetails?.kbm_level }}</el-descriptions-item>
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
import type { Evidence, CreateEvidenceRequest } from '../types/evidence'
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

// 表单数据
const createForm = reactive<CreateEvidenceRequest>({
  student_task_id: '',
  type: '',
  content: '',
  kbm_name: '',
  kbm_level: 0
})

// 验证规则
const createRules = {
  student_task_id: [{ required: true, message: '请输入学生任务ID', trigger: 'blur' }],
  type: [{ required: true, message: '请选择证据类型', trigger: 'blur' }],
  kbm_name: [{ required: true, message: '请输入KBM名称', trigger: 'blur' }],
  kbm_level: [{ required: true, message: '请选择KBM级别', trigger: 'blur' }],
  content: [{ required: true, message: '请输入证据内容', trigger: 'blur' }]
}

// 证据详情
const evidenceDetails = ref<Evidence | null>(null)

// 加载证据列表
const loadEvidences = async () => {
  try {
    // 这里简化处理，实际应该根据学生任务ID获取证据列表
    // 暂时使用模拟数据
    evidences.value = [
      {
        id: '1',
        student_task_id: 'st1',
        type: 'document',
        content: '这是一份文献综述报告',
        kbm_name: '文献综述质量',
        kbm_level: 4,
        created_at: '2026-03-01T10:00:00Z',
        updated_at: '2026-03-01T10:00:00Z'
      },
      {
        id: '2',
        student_task_id: 'st1',
        type: 'code',
        content: '这是一段数据分析代码',
        kbm_name: '数据分析方法选择',
        kbm_level: 3,
        created_at: '2026-03-02T14:30:00Z',
        updated_at: '2026-03-02T14:30:00Z'
      }
    ]
    total.value = evidences.value.length
  } catch (error) {
    console.error('加载证据失败:', error)
    ElMessage.error('加载证据失败')
  }
}

// 创建证据
const handleCreateEvidence = async () => {
  if (!createFormRef.value) return
  
  await createFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        const response = await evidenceApi.createEvidence(createForm)
        if (response.code === 200) {
          ElMessage.success('创建证据成功')
          showCreateDialog.value = false
          loadEvidences()
          // 重置表单
          Object.assign(createForm, {
            student_task_id: '',
            type: '',
            content: '',
            kbm_name: '',
            kbm_level: 0
          })
        }
      } catch (error) {
        console.error('创建证据失败:', error)
        ElMessage.error('创建证据失败')
      }
    }
  })
}

// 查看证据详情
const viewEvidenceDetails = async (evidenceId: string) => {
  try {
    // 这里简化处理，实际应该调用API获取证据详情
    // 暂时使用模拟数据
    evidenceDetails.value = {
      id: evidenceId,
      student_task_id: 'st1',
      type: 'document',
      content: '这是一份详细的文献综述报告，包含了相关领域的最新研究进展和分析。',
      kbm_name: '文献综述质量',
      kbm_level: 4,
      created_at: '2026-03-01T10:00:00Z',
      updated_at: '2026-03-01T10:00:00Z'
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
    
    // 这里简化处理，实际应该调用API删除证据
    ElMessage.success('删除证据成功')
    loadEvidences()
  } catch (error) {
    console.error('删除证据失败:', error)
    ElMessage.error('删除证据失败')
  }
}

// 初始化
onMounted(() => {
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