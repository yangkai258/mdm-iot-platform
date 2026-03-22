<template>
  <div class="test-container">
    <div class="test-header">
      <h1>弹窗/页面风格测试</h1>
      <p>点击下方按钮预览不同风格，选择你喜欢的一种</p>
    </div>

    <div class="test-section">
      <h2>当前风格问题</h2>
      <p>点击按钮后，弹窗/页面依然显示侧边栏，不够清爽</p>
    </div>

    <div class="test-section">
      <h2>风格选项</h2>
      
      <div class="style-cards">
        <!-- 风格A: 全屏页面 -->
        <div class="style-card">
          <h3>A. 全屏页面</h3>
          <p>点击后进入新页面，无侧边栏，适合复杂表单</p>
          <div class="preview-box style-a-preview">
            <div class="preview-sidebar">侧边栏</div>
            <div class="preview-content">
              <div class="preview-header">页面A</div>
              <div class="preview-form">表单内容</div>
            </div>
          </div>
          <button class="arco-btn arco-btn-primary" @click="showStyleA">预览风格A</button>
        </div>

        <!-- 风格B: 居中模态框 -->
        <div class="style-card">
          <h3>B. 居中模态框</h3>
          <p>半透明遮罩 + 居中对话框，点击遮罩可关闭</p>
          <div class="preview-box style-b-preview">
            <div class="preview-overlay">
              <div class="preview-dialog">对话框</div>
            </div>
          </div>
          <button class="arco-btn arco-btn-primary" @click="showStyleB">预览风格B</button>
        </div>

        <!-- 风格C: 右侧抽屉 -->
        <div class="style-card">
          <h3>C. 右侧抽屉</h3>
          <p>从右侧滑入的抽屉面板，不遮挡主内容</p>
          <div class="preview-box style-c-preview">
            <div class="preview-content-with-drawer">
              <div class="preview-drawer">抽屉</div>
            </div>
          </div>
          <button class="arco-btn arco-btn-primary" @click="showStyleC">预览风格C</button>
        </div>

        <!-- 风格D: 全屏模态 -->
        <div class="style-card">
          <h3>D. 全屏模态</h3>
          <p>类似风格A但是用遮罩实现，点击遮罩关闭</p>
          <div class="preview-box style-d-preview">
            <div class="preview-overlay-full">全屏模态</div>
          </div>
          <button class="arco-btn arco-btn-primary" @click="showStyleD">预览风格D</button>
        </div>
      </div>
    </div>

    <!-- 风格A: 全屏页面 -->
    <div v-if="activeStyle === 'A'" class="full-page-modal">
      <div class="full-page-header">
        <span>创建设备 - 全屏页面</span>
        <button class="arco-btn arco-btn-text" @click="activeStyle = null">关闭</button>
      </div>
      <div class="full-page-content">
        <div class="form-section">
          <h3>基本信息</h3>
          <div class="form-row">
            <label>设备名称</label>
            <input type="text" placeholder="请输入设备名称" />
          </div>
          <div class="form-row">
            <label>设备类型</label>
            <select>
              <option>请选择</option>
              <option>传感器</option>
              <option>控制器</option>
            </select>
          </div>
        </div>
        <div class="form-section">
          <h3>配置信息</h3>
          <div class="form-row">
            <label>固件版本</label>
            <input type="text" placeholder="请输入固件版本" />
          </div>
        </div>
        <div class="form-actions">
          <button class="arco-btn arco-btn-secondary">取消</button>
          <button class="arco-btn arco-btn-primary">确定</button>
        </div>
      </div>
    </div>

    <!-- 风格B: 居中模态框 -->
    <div v-if="activeStyle === 'B'" class="modal-overlay" @click.self="activeStyle = null">
      <div class="modal-dialog">
        <div class="modal-header">
          <span>创建设备 - 居中弹窗</span>
          <button class="arco-btn arco-btn-text" @click="activeStyle = null">×</button>
        </div>
        <div class="modal-body">
          <div class="form-row">
            <label>设备名称</label>
            <input type="text" placeholder="请输入设备名称" />
          </div>
          <div class="form-row">
            <label>设备类型</label>
            <select>
              <option>请选择</option>
              <option>传感器</option>
              <option>控制器</option>
            </select>
          </div>
          <div class="form-row">
            <label>固件版本</label>
            <input type="text" placeholder="请输入固件版本" />
          </div>
        </div>
        <div class="modal-footer">
          <button class="arco-btn arco-btn-secondary" @click="activeStyle = null">取消</button>
          <button class="arco-btn arco-btn-primary">确定</button>
        </div>
      </div>
    </div>

    <!-- 风格C: 右侧抽屉 -->
    <div v-if="activeStyle === 'C'" class="drawer-overlay" @click.self="activeStyle = null">
      <div class="drawer-panel">
        <div class="drawer-header">
          <span>创建设备 - 右侧抽屉</span>
          <button class="arco-btn arco-btn-text" @click="activeStyle = null">×</button>
        </div>
        <div class="drawer-body">
          <div class="form-row">
            <label>设备名称</label>
            <input type="text" placeholder="请输入设备名称" />
          </div>
          <div class="form-row">
            <label>设备类型</label>
            <select>
              <option>请选择</option>
              <option>传感器</option>
              <option>控制器</option>
            </select>
          </div>
          <div class="form-row">
            <label>固件版本</label>
            <input type="text" placeholder="请输入固件版本" />
          </div>
          <div class="form-row">
            <label>备注</label>
            <textarea placeholder="请输入备注"></textarea>
          </div>
        </div>
        <div class="drawer-footer">
          <button class="arco-btn arco-btn-secondary" @click="activeStyle = null">取消</button>
          <button class="arco-btn arco-btn-primary">确定</button>
        </div>
      </div>
    </div>

    <!-- 风格D: 全屏模态 -->
    <div v-if="activeStyle === 'D'" class="full-modal-overlay" @click.self="activeStyle = null">
      <div class="full-modal-content">
        <div class="full-modal-header">
          <span>创建设备 - 全屏模态</span>
          <button class="arco-btn arco-btn-text" @click="activeStyle = null">关闭</button>
        </div>
        <div class="full-modal-body">
          <div class="form-section">
            <h3>基本信息</h3>
            <div class="form-grid">
              <div class="form-row">
                <label>设备名称</label>
                <input type="text" placeholder="请输入设备名称" />
              </div>
              <div class="form-row">
                <label>设备类型</label>
                <select>
                  <option>请选择</option>
                  <option>传感器</option>
                  <option>控制器</option>
                </select>
              </div>
              <div class="form-row">
                <label>固件版本</label>
                <input type="text" placeholder="请输入固件版本" />
              </div>
              <div class="form-row">
                <label>设备ID</label>
                <input type="text" placeholder="自动生成" disabled />
              </div>
            </div>
          </div>
          <div class="form-section">
            <h3>其他信息</h3>
            <div class="form-row">
              <label>备注</label>
              <textarea placeholder="请输入备注"></textarea>
            </div>
          </div>
        </div>
        <div class="full-modal-footer">
          <button class="arco-btn arco-btn-secondary" @click="activeStyle = null">取消</button>
          <button class="arco-btn arco-btn-primary">确定创建</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const activeStyle = ref(null)

const showStyleA = () => { activeStyle.value = 'A' }
const showStyleB = () => { activeStyle.value = 'B' }
const showStyleC = () => { activeStyle.value = 'C' }
const showStyleD = () => { activeStyle.value = 'D' }
</script>

<style scoped>
.test-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.test-header {
  text-align: center;
  margin-bottom: 32px;
}

.test-header h1 {
  font-size: 24px;
  margin-bottom: 8px;
}

.test-header p {
  color: #666;
}

.test-section {
  margin-bottom: 32px;
}

.test-section h2 {
  font-size: 18px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
}

.style-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.style-card {
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  padding: 20px;
}

.style-card h3 {
  font-size: 16px;
  margin-bottom: 8px;
}

.style-card p {
  color: #666;
  font-size: 14px;
  margin-bottom: 16px;
}

.preview-box {
  height: 160px;
  border: 1px dashed #ccc;
  border-radius: 4px;
  margin-bottom: 16px;
  display: flex;
  overflow: hidden;
}

.style-card button {
  width: 100%;
}

/* Preview styles */
.preview-sidebar {
  width: 60px;
  background: #001529;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.preview-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.preview-header {
  background: #fff;
  border-bottom: 1px solid #eee;
  padding: 8px 12px;
  font-size: 12px;
}

.preview-form {
  flex: 1;
  background: #f5f5f5;
  padding: 12px;
  font-size: 12px;
}

.style-b-preview,
.style-d-preview {
  position: relative;
}

.preview-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-dialog {
  width: 70%;
  height: 70%;
  background: #fff;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.preview-overlay-full {
  position: absolute;
  inset: 0;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.style-c-preview {
  position: relative;
}

.preview-content-with-drawer {
  flex: 1;
  position: relative;
  background: #f5f5f5;
}

.preview-drawer {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 50%;
  background: #fff;
  border-left: 1px solid #e5e5e5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

/* ========== 风格A: 全屏页面 ========== */
.full-page-modal {
  position: fixed;
  inset: 0;
  background: #f5f5f5;
  z-index: 1000;
  display: flex;
  flex-direction: column;
}

.full-page-header {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e5e5e5;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 16px;
  font-weight: 500;
}

.full-page-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.form-section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
}

.form-section h3 {
  font-size: 16px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.form-row {
  margin-bottom: 16px;
}

.form-row:last-child {
  margin-bottom: 0;
}

.form-row label {
  display: block;
  margin-bottom: 6px;
  color: #666;
  font-size: 14px;
}

.form-row input,
.form-row select,
.form-row textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
}

.form-row textarea {
  min-height: 100px;
  resize: vertical;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

/* ========== 风格B: 居中模态框 ========== */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-dialog {
  width: 520px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e5e5e5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #e5e5e5;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* ========== 风格C: 右侧抽屉 ========== */
.drawer-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 1000;
}

.drawer-panel {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 480px;
  background: #fff;
  display: flex;
  flex-direction: column;
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.15);
}

.drawer-header {
  height: 56px;
  padding: 0 20px;
  border-bottom: 1px solid #e5e5e5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.drawer-body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.drawer-footer {
  padding: 16px 20px;
  border-top: 1px solid #e5e5e5;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* ========== 风格D: 全屏模态 ========== */
.full-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.full-modal-content {
  width: 900px;
  max-height: 90vh;
  background: #fff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

.full-modal-header {
  height: 56px;
  padding: 0 24px;
  border-bottom: 1px solid #e5e5e5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.full-modal-body {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

.full-modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e5e5e5;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
