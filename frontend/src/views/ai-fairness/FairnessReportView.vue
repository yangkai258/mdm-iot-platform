<template>
  <div class="fairness-report">
    <a-card title="公平性测试报告">
      <a-spin :loading="loading">
        <a-row :gutter="16">
          <a-col :span="8">
            <a-card title="总体评分">
              <a-statistic :value="report.overall_score" :suffix="'/' + report.max_score">
                <template #prefix>
                  <icon-check-circle style="color: green" v-if="report.overall_score >= 80" />
                  <icon-close-circle style="color: red" v-else />
                </template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="测试用例">
              <a-statistic :value="report.total_tests" />
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="通过率">
              <a-statistic :value="report.pass_rate" suffix="%" />
            </a-card>
          </a-col>
        </a-row>
        
        <a-divider>偏见检测结果</a-divider>
        <a-table :columns="columns" :data="report.bias_results" :pagination="false" />
        
        <a-divider>建议措施</a-divider>
        <a-list :data="report.recommendations">
          <template #item="{ item }">
            <a-list-item>
              <icon-check-circle style="color: green; margin-right: 8px" />
              {{ item }}
            </a-list-item>
          </template>
        </a-list>
        
        <template #extra>
          <a-space>
            <a-button @click="downloadPDF">下载PDF</a-button>
            <a-button type="primary" @click="shareReport">分享报告</a-button>
          </a-space>
        </template>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Message } from '@arco-design/web-vue'

const API_BASE = '/api/v1'
const route = useRoute()
const loading = ref(false)
const report = ref({ overall_score: 0, max_score: 100, total_tests: 0, pass_rate: 0, bias_results: [], recommendations: [] })

const columns = [
  { title: '检测项', dataIndex: 'check_item' },
  { title: '测试组', dataIndex: 'test_group' },
  { title: '通过', dataIndex: 'passed' },
  { title: '失败', dataIndex: 'failed' },
  { title: '状态', slotName: 'status' },
]

const loadReport = async () => {
  loading.value = true
  try {
    const id = route.params.id
    const res = await fetch(`${API_BASE}/ai-fairness/reports/${id}`)
    report.value = await res.json()
  } catch (e) {
    Message.error('加载报告失败')
  } finally {
    loading.value = false
  }
}

const downloadPDF = () => Message.info('PDF下载开发中')
const shareReport = () => Message.info('分享功能开发中')

onMounted(loadReport)
</script>
