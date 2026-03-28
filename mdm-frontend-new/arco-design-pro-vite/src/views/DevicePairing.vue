<template>
  <div class="pairing-container">
    <a-card>
      <template #title>
        <span>设备配对</span>
      </template>
      
      <a-steps :current="currentStep" style="margin-bottom: 32px;">
        <a-step title="选择设备类型" description="选择要配对的设备类型" />
        <a-step title="扫描设备" description="扫描二维码或手动输入" />
        <a-step title="验证配对码" description="在设备上确认配对码" />
        <a-step title="完成" description="配对成功" />
      </a-steps>

      <!-- Step 1: 选择设备类型 -->
      <div v-if="currentStep === 0" class="step-content">
        <a-row :gutter="[24, 24]">
          <a-col :span="8">
            <div class="device-type-card" :class="{ active: form.deviceType === 'M5Stack' }" @click="form.deviceType = 'M5Stack'">
              <icon-robot :size="48" />
              <div class="device-name">M5Stack</div>
              <div class="device-desc">主流设备系列</div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="device-type-card" :class="{ active: form.deviceType === 'ESP32' }" @click="form.deviceType = 'ESP32'">
              <icon-robot :size="48" />
              <div class="device-name">ESP32</div>
              <div class="device-desc">低成本解决方案</div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="device-type-card" :class="{ active: form.deviceType === 'RaspberryPi' }" @click="form.deviceType = 'RaspberryPi'">
              <icon-robot :size="48" />
              <div class="device-name">树莓派</div>
              <div class="device-desc">高级开发板</div>
            </div>
          </a-col>
        </a-row>
        <div style="text-align: center; margin-top: 32px;">
          <a-button type="primary" :disabled="!form.deviceType" @click="currentStep = 1">下一步</a-button>
        </div>
      </div>

      <!-- Step 2: 扫描设备 -->
      <div v-if="currentStep === 1" class="step-content">
        <a-row :gutter="24">
          <a-col :span="12">
            <a-card title="扫码配对">
              <div class="qr-placeholder">
                <icon-qrcode :size="120" />
                <div style="margin-top: 16px;">请使用设备扫描此二维码</div>
              </div>
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="或手动输入">
              <a-form :model="form" layout="vertical">
                <a-form-item label="设备码" required>
                  <a-input v-model="form.deviceCode" placeholder="请输入设备底部的设备码" />
                </a-form-item>
                <a-form-item label="设备名称(可选)">
                  <a-input v-model="form.deviceName" placeholder="为设备起个名字，如：客厅小黄" />
                </a-form-item>
              </a-form>
            </a-card>
          </a-col>
        </a-row>
        <div style="text-align: center; margin-top: 32px;">
          <a-button @click="currentStep = 0">上一步</a-button>
          <a-button type="primary" :disabled="!form.deviceCode" @click="handleSearchDevice" style="margin-left: 16px;">
            搜索设备
          </a-button>
        </div>
      </div>

      <!-- Step 3: 验证配对码 -->
      <div v-if="currentStep === 2" class="step-content">
        <a-result title="请确认配对码" sub-title="在设备屏幕上确认以下配对码">
          <template #icon>
            <div class="pairing-code">{{ pairingCode }}</div>
          </template>
          <template #extra>
            <a-button type="primary" @click="handleConfirmPairing">确认配对</a-button>
            <a-button @click="handleCancelPairing" style="margin-left: 16px;">取消</a-button>
          </template>
        </a-result>
      </div>

      <!-- Step 4: 完成 -->
      <div v-if="currentStep === 3" class="step-content">
        <a-result title="配对成功" sub-title="设备已成功添加到您的账户">
          <template #icon>
            <icon-check-circle :size="80" style="color: #00B42A;" />
          </template>
          <template #extra>
            <a-button type="primary" @click="handleFinish">完成</a-button>
            <a-button @click="handlePairAnother" style="margin-left: 16px;">继续配对</a-button>
          </template>
        </a-result>
      </div>

      <!-- 搜索到的设备 -->
      <a-modal v-model:visible="deviceFoundVisible" title="找到设备" @before-ok="handleSelectDevice">
        <a-descriptions bordered :column="2">
          <a-descriptions-item label="设备码">{{ foundDevice.deviceCode }}</a-descriptions-item>
          <a-descriptions-item label="设备类型">{{ foundDevice.type }}</a-descriptions-item>
          <a-descriptions-item label="固件版本">{{ foundDevice.firmware }}</a-descriptions-item>
          <a-descriptions-item label="设备名称">
            <a-input v-model="form.deviceName" placeholder="为设备起个名字" />
          </a-descriptions-item>
        </a-descriptions>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';

const currentStep = ref(0);
const deviceFoundVisible = ref(false);
const pairingCode = ref('123456');

const form = reactive({
  deviceType: '',
  deviceCode: '',
  deviceName: '',
});

const foundDevice = reactive({
  deviceCode: 'DEV001',
  type: 'M5Stack',
  firmware: 'v2.0.0',
});

const handleSearchDevice = () => {
  deviceFoundVisible.value = true;
};

const handleSelectDevice = (done: boolean) => {
  pairingCode.value = String(Math.floor(Math.random() * 900000) + 100000);
  currentStep.value = 2;
  done(true);
  deviceFoundVisible.value = false;
};

const handleConfirmPairing = () => {
  currentStep.value = 3;
};

const handleCancelPairing = () => {
  currentStep.value = 1;
};

const handleFinish = () => {
  // router back
};

const handlePairAnother = () => {
  currentStep.value = 0;
  form.deviceType = '';
  form.deviceCode = '';
  form.deviceName = '';
};
</script>

<style scoped>
.pairing-container { padding: 20px; max-width: 900px; margin: 0 auto; }
.step-content { min-height: 300px; }
.device-type-card {
  border: 2px solid #e5e6e8;
  border-radius: 8px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}
.device-type-card:hover { border-color: #165DFF; }
.device-type-card.active { border-color: #165DFF; background: #f2f3f5; }
.device-name { font-size: 18px; font-weight: 500; margin-top: 16px; }
.device-desc { font-size: 14px; color: #86909c; margin-top: 8px; }
.qr-placeholder {
  text-align: center;
  padding: 48px;
  color: #86909c;
}
.pairing-code {
  font-size: 48px;
  font-weight: bold;
  letter-spacing: 8px;
  color: #165DFF;
}
</style>
