<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>设备监控</a-breadcrumb-item>
      <a-breadcrumb-item>远程调试</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 设备选择 -->
    <a-card class="device-select-card">
      <a-space wrap>
        <a-select
          v-model="selectedDevice"
          placeholder="选择要调试的设备"
          style="width: 260px"
          filterable
          allow-search
          @change="onDeviceChange"
        >
          <a-option v-for="d in deviceList" :key="d.device_id" :value="d.device_id">
            {{ d.device_id }} ({{ d.is_online ? '在线' : '离线' }})
          </a-option>
        </a-select>
        <a-button @click="clearTerminal">清空终端</a-button>
        <a-button @click="exportLogs">导出日志</a-button>
        <a-badge :status="connected ? 'success' : 'error'" :text="connected ? '已连接' : '未连接'" />
      </a-space>
    </a-card>

    <!-- 快捷命令 -->
    <a-card class="cmd-card" title="快捷命令">
      <a-space wrap>
        <a-button
          v-for="cmd in quickCommands"
          :key="cmd.name"
          size="small"
          @click="sendCommand(cmd.cmd)"
        >
          {{ cmd.name }}
        </a-button>
        <a-button size="small" @click="sendCommand('help')">帮助</a-button>
      </a-space>
    </a-card>

    <!-- 终端输出区 -->
    <a-card class="terminal-card" :body-style="{ padding: 0 }">
      <div class="terminal-output" ref="terminalRef">
        <div class="terminal-line system">MDM Remote Debug Console v1.0.0</div>
        <div class="terminal-line system">输入 help 查看可用命令，输入 exit 退出连接</div>
        <div class="terminal-line system" v-if="!selectedDevice">请先选择要调试的设备...</div>
        <div
          v-for="(line, idx) in terminalLines"
          :key="idx"
          class="terminal-line"
          :class="line.type"
        >
          <span class="prompt" v-if="line.type === 'input'">&gt;&gt; </span>
          <span class="prompt" v-else-if="line.type === 'error'">❌ </span>
          <span class="prompt" v-else-if="line.type === 'success'">✅ </span>
          <span class="prompt" v-else-if="line.type === 'system'">ℹ️ </span>
          {{ line.text }}
        </div>
      </div>

      <!-- 命令输入区 -->
      <div class="terminal-input-bar">
        <span class="input-prompt">&gt;&gt;</span>
        <a-input
          v-model="commandInput"
          placeholder="输入调试命令，按 Enter 发送"
          class="terminal-input"
          :disabled="!selectedDevice"
          @press-enter="sendCommand(commandInput)"
        />
        <a-button type="primary" :disabled="!selectedDevice || !commandInput" @click="sendCommand(commandInput)">
          发送
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Message } from '@arco-design/web-vue'
import monitorApi from '@/api/monitor'

const selectedDevice = ref('')
const commandInput = ref('')
const connected = ref(false)
const terminalRef = ref(null)
const terminalLines = ref([])

const deviceList = ref([
  { device_id: 'DEV001', is_online: true },
  { device_id: 'DEV002', is_online: true },
  { device_id: 'DEV003', is_online: false },
  { device_id: 'DEV004', is_online: true },
  { device_id: 'DEV005', is_online: true }
])

const quickCommands = [
  { name: '设备信息', cmd: 'device info' },
  { name: '系统状态', cmd: 'system status' },
  { name: '网络诊断', cmd: 'network diag' },
  { name: '重启设备', cmd: 'device restart' },
  { name: '固件版本', cmd: 'firmware version' },
  { name: '内存状态', cmd: 'memory stats' },
  { name: '传感器列表', cmd: 'sensor list' },
  { name: '日志级别', cmd: 'log level' },
  { name: '恢复出厂', cmd: 'factory reset' }
]

const responses = {
  'device info': 'Device: ' + (selectedDevice.value || 'N/A') + '\nHardware: MDM-Pro-200\nFirmware: v1.2.0\nMAC: 00:11:22:33:44:55\nUptime: 7d 14h 32m',
  'system status': 'CPU: 23%\nMemory: 67%\nTemperature: 38°C\nBattery: 85%\nNetwork: WiFi (signal: -45dBm)\nMQTT: Connected\nLast Heartbeat: 2026-03-22 10:30:00',
  'network diag': 'Ping test: OK (latency: 12ms)\nDNS: Resolved\nGateway: 192.168.1.1\nMQTT Broker: Connected (tls: true)',
  'device restart': 'Restarting device... Please wait 30 seconds.',
  'firmware version': 'Current: v1.2.0\nLatest: v1.2.1\nUpdate available: Yes',
  'memory stats': 'Total: 512KB\nUsed: 342KB\nFree: 170KB\nHeap: OK\nStack: 1.2KB / 4KB',
  'sensor list': '1. Temperature  -> 38.2°C\n2. Humidity     -> 65%\n3. Battery      -> 85%\n4. Accelerometer -> Normal',
  'log level': 'Current: INFO\nAvailable: DEBUG, INFO, WARN, ERROR',
  'factory reset': '⚠️ WARNING: This will reset all device settings!\nConfirm? (type: factory reset confirm)',
  'help': 'Available commands:\n  device info    - Show device info\n  system status  - Show system status\n  network diag   - Network diagnostics\n  device restart - Restart device\n  firmware version - Check firmware\n  memory stats   - Memory usage\n  sensor list    - List sensors\n  log level      - Log level config\n  factory reset  - Factory reset\n  exit           - Disconnect',
  'exit': 'Disconnected from device.'
}

const scrollToBottom = () => {
  nextTick(() => {
    if (terminalRef.value) {
      terminalRef.value.scrollTop = terminalRef.value.scrollHeight
    }
  })
}

const addLine = (text, type = 'output') => {
  const lines = text.split('\n')
  lines.forEach(line => {
    terminalLines.value.push({ text: line, type })
  })
  scrollToBottom()
}

const sendCommand = async (cmd) => {
  const trimmed = cmd.trim()
  if (!trimmed) return
  if (!selectedDevice.value) {
    Message.warning('请先选择设备')
    return
  }

  // 显示输入
  addLine(trimmed, 'input')
  commandInput.value = ''

  // 处理退出
  if (trimmed === 'exit') {
    connected.value = false
    addLine('已断开设备连接', 'system')
    return
  }

  // 模拟响应
  const response = responses[trimmed]
  if (response) {
    await new Promise(r => setTimeout(r, 300 + Math.random() * 500))
    if (trimmed === 'device restart') {
      addLine(response, 'success')
      setTimeout(() => addLine('Device rebooting...', 'system'), 1500)
    } else if (trimmed === 'factory reset confirm') {
      addLine('Factory reset initiated. Device will reboot in 10s...', 'error')
    } else if (trimmed === 'exit') {
      addLine(response, 'system')
    } else {
      addLine(response, 'output')
    }
  } else {
    await new Promise(r => setTimeout(r, 200))
    addLine(`Command not found: ${trimmed}`, 'error')
  }
}

const onDeviceChange = (deviceId) => {
  terminalLines.value = []
  if (deviceId) {
    connected.value = true
    addLine(`正在连接到设备 ${deviceId}...`, 'system')
    setTimeout(() => {
      addLine(`已连接到 ${deviceId}`, 'success')
      addLine('输入 help 查看可用命令', 'system')
    }, 800)
  } else {
    connected.value = false
  }
}

const clearTerminal = () => {
  terminalLines.value = []
  addLine('终端已清空', 'system')
}

const exportLogs = () => {
  const content = terminalLines.value.map(l => `[${l.type}] ${l.text}`).join('\n')
  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `debug-${selectedDevice.value || 'session'}-${Date.now()}.log`
  a.click()
  URL.revokeObjectURL(url)
  Message.success('日志已导出')
}

onMounted(() => {
  // 从 URL 参数获取设备
  const hash = window.location.hash
  const match = hash.match(/device=([^&]+)/)
  if (match) {
    selectedDevice.value = decodeURIComponent(match[1])
    onDeviceChange(selectedDevice.value)
  }
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.breadcrumb { margin-bottom: 0; }

.device-select-card,
.cmd-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}

.terminal-card {
  background: #1e1e1e;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
  flex: 1;
  min-height: 480px;
  display: flex;
  flex-direction: column;
}

.terminal-output {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
  line-height: 1.6;
  min-height: 400px;
  max-height: calc(100vh - 380px);
  scrollbar-width: thin;
  scrollbar-color: #555 #1e1e1e;
}

.terminal-output::-webkit-scrollbar { width: 6px; }
.terminal-output::-webkit-scrollbar-track { background: #1e1e1e; }
.terminal-output::-webkit-scrollbar-thumb { background: #555; border-radius: 3px; }

.terminal-line {
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-all;
}

.terminal-line.input { color: #4fc3f7; }
.terminal-line.error { color: #f53f3f; }
.terminal-line.success { color: #00b42a; }
.terminal-line.system { color: #b392f0; }
.terminal-line .prompt { font-weight: bold; }

.terminal-input-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #333;
  background: #252526;
}

.input-prompt {
  color: #4fc3f7;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  font-weight: bold;
}

.terminal-input {
  flex: 1;
  background: #2d2d2d;
  border-color: #3c3c3c;
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
}

.terminal-input :deep(.arco-input) {
  background: transparent;
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
}
</style>
