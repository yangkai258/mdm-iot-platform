<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>开发工具</a-breadcrumb-item>
      <a-breadcrumb-item>弹窗风格测试</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 说明区域 -->
    <a-card class="intro-card">
      <template #title>弹窗/页面风格测试</template>
      <a-space direction="vertical" :size="16">
        <p>点击下方按钮预览不同风格，选择你喜欢的一种</p>
        <a-alert type="info">
          当前页面采用标准 page-container 结构，点击按钮后会弹出不同风格的模态框/页面
        </a-alert>
      </a-space>
    </a-card>

    <!-- 风格选项 -->
    <a-card class="style-card">
      <template #title>风格选项</template>
      <a-row :gutter="[24, 24]">
        <!-- 风格A: 全屏页面 -->
        <a-col :span="12" :xs="24" :sm="12">
          <div class="style-option">
            <h3>A. 全屏页面</h3>
            <p>点击后进入新页面，无侧边栏，适合复杂表单</p>
            <a-button type="primary" @click="showStyleA">预览风格A</a-button>
          </div>
        </a-col>

        <!-- 风格B: 居中模态框 -->
        <a-col :span="12" :xs="24" :sm="12">
          <div class="style-option">
            <h3>B. 居中模态框</h3>
            <p>半透明遮罩 + 居中对话框，点击遮罩可关闭</p>
            <a-button type="primary" @click="showStyleB">预览风格B</a-button>
          </div>
        </a-col>

        <!-- 风格C: 右侧抽屉 -->
        <a-col :span="12" :xs="24" :sm="12">
          <div class="style-option">
            <h3>C. 右侧抽屉</h3>
            <p>从右侧滑入的抽屉面板，不遮挡主内容</p>
            <a-button type="primary" @click="showStyleC">预览风格C</a-button>
          </div>
        </a-col>

        <!-- 风格D: 全屏模态 -->
        <a-col :span="12" :xs="24" :sm="12">
          <div class="style-option">
            <h3>D. 全屏模态</h3>
            <p>类似风格A但是用遮罩实现，点击遮罩关闭</p>
            <a-button type="primary" @click="showStyleD">预览风格D</a-button>
          </div>
        </a-col>
      </a-row>
    </a-card>

    <!-- 风格A: 全屏页面 -->
    <a-modal 
      v-model:visible="styleAVisible" 
      :footer="false" 
      :mask-closable="true"
      title="创建设备 - 全屏页面"
      :width="900"
    >
      <div class="full-page-form">
        <a-card title="基本信息" class="form-section">
          <a-form :model="formA" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="设备名称">
                  <a-input v-model="formA.name" placeholder="请输入设备名称" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="设备类型">
                  <a-select v-model="formA.type" placeholder="请选择">
                    <a-option value="sensor">传感器</a-option>
                    <a-option value="controller">控制器</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="固件版本">
                  <a-input v-model="formA.firmware" placeholder="请输入固件版本" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="设备ID">
                  <a-input v-model="formA.deviceId" placeholder="自动生成" disabled />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-card>

        <a-card title="配置信息" class="form-section">
          <a-form :model="formA" layout="vertical">
            <a-form-item label="备注">
              <a-textarea v-model="formA.remark" placeholder="请输入备注" :rows="3" />
            </a-form-item>
          </a-form>
        </a-card>

        <div class="form-actions">
          <a-button @click="styleAVisible = false">取消</a-button>
          <a-button type="primary">确定创建</a-button>
        </div>
      </div>
    </a-modal>

    <!-- 风格B: 居中模态框 -->
    <a-modal 
      v-model:visible="styleBVisible" 
      title="创建设备 - 居中弹窗" 
      :width="520"
      @ok="styleBVisible = false"
      @cancel="styleBVisible = false"
    >
      <a-form :model="formB" layout="vertical">
        <a-form-item label="设备名称" required>
          <a-input v-model="formB.name" placeholder="请输入设备名称" />
        </a-form-item>
        <a-form-item label="设备类型" required>
          <a-select v-model="formB.type" placeholder="请选择">
            <a-option value="sensor">传感器</a-option>
            <a-option value="controller">控制器</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件版本">
          <a-input v-model="formB.firmware" placeholder="请输入固件版本" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 风格C: 右侧抽屉 -->
    <a-drawer 
      v-model:visible="styleCVisible" 
      title="创建设备 - 右侧抽屉" 
      width="480px"
      @before-ok="styleCVisible = false"
    >
      <a-form :model="formC" layout="vertical">
        <a-form-item label="设备名称" required>
          <a-input v-model="formC.name" placeholder="请输入设备名称" />
        </a-form-item>
        <a-form-item label="设备类型" required>
          <a-select v-model="formC.type" placeholder="请选择">
            <a-option value="sensor">传感器</a-option>
            <a-option value="controller">控制器</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="固件版本">
          <a-input v-model="formC.firmware" placeholder="请输入固件版本" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model="formC.remark" placeholder="请输入备注" :rows="3" />
        </a-form-item>
      </a-form>
    </a-drawer>

    <!-- 风格D: 全屏模态 -->
    <a-modal 
      v-model:visible="styleDVisible" 
      :footer="false" 
      :mask-closable="true"
      title="创建设备 - 全屏模态"
      :width="900"
    >
      <div class="full-page-form">
        <a-card title="基本信息" class="form-section">
          <a-form :model="formD" layout="vertical">
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="设备名称">
                  <a-input v-model="formD.name" placeholder="请输入设备名称" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="设备类型">
                  <a-select v-model="formD.type" placeholder="请选择">
                    <a-option value="sensor">传感器</a-option>
                    <a-option value="controller">控制器</a-option>
                  </a-select>
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="16">
              <a-col :span="12">
                <a-form-item label="固件版本">
                  <a-input v-model="formD.firmware" placeholder="请输入固件版本" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="设备ID">
                  <a-input v-model="formD.deviceId" placeholder="自动生成" disabled />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-card>

        <div class="form-actions">
          <a-button @click="styleDVisible = false">取消</a-button>
          <a-button type="primary">确定创建</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const styleAVisible = ref(false)
const styleBVisible = ref(false)
const styleCVisible = ref(false)
const styleDVisible = ref(false)

const formA = reactive({
  name: '',
  type: '',
  firmware: '',
  deviceId: '自动生成',
  remark: ''
})

const formB = reactive({
  name: '',
  type: '',
  firmware: ''
})

const formC = reactive({
  name: '',
  type: '',
  firmware: '',
  remark: ''
})

const formD = reactive({
  name: '',
  type: '',
  firmware: '',
  deviceId: '自动生成'
})

const showStyleA = () => {
  styleAVisible.value = true
}

const showStyleB = () => {
  styleBVisible.value = true
}

const showStyleC = () => {
  styleCVisible.value = true
}

const showStyleD = () => {
  styleDVisible.value = true
}
</script>

<style scoped>
.page-container {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.breadcrumb {
  margin-bottom: 16px;
}

.intro-card {
  margin-bottom: 16px;
}

.style-card {
  margin-bottom: 16px;
}

.style-option {
  padding: 16px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  height: 100%;
}

.style-option h3 {
  margin-bottom: 8px;
  font-size: 16px;
}

.style-option p {
  color: #666;
  font-size: 14px;
  margin-bottom: 16px;
}

.full-page-form {
  max-height: 60vh;
  overflow-y: auto;
}

.form-section {
  margin-bottom: 16px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e5e5e5;
}
</style>
