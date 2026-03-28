<template>
  <div class="container">
    <a-card>
      <template #title>
        <a-space><icon-robot /> 虚拟宠物仿真</a-space>
      </template>

      <a-row :gutter="16">
        <a-col :span="8">
          <a-card title="仿真控制">
            <a-form :model="simForm" layout="vertical">
              <a-form-item label="仿真场景">
                <a-select v-model="simForm.scene">
                  <a-option value="home">家庭环境</a-option>
                  <a-option value="outdoor">户外环境</a-option>
                </a-select>
              </a-form-item>
              <a-form-item label="运行时间">
                <a-input-number v-model="simForm.duration" :min="1" :max="60" /> 分钟
              </a-form-item>
              <a-form-item label="加速倍数">
                <a-select v-model="simForm.speed">
                  <a-option value="1">1x</a-option>
                  <a-option value="5">5x</a-option>
                  <a-option value="10">10x</a-option>
                </a-select>
              </a-form-item>
              <a-form-item>
                <a-button type="primary" @click="handleStart">开始仿真</a-button>
                <a-button @click="handlePause">暂停</a-button>
                <a-button @click="handleStop">停止</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :span="16">
          <a-card title="仿真视图">
            <div class="sim-view">
              <svg width="600" height="400" viewBox="0 0 600 400">
                <rect x="0" y="0" width="600" height="400" fill="#e8e8e8" />
                <rect x="50" y="50" width="200" height="150" fill="#b3d9ff" />
                <rect x="270" y="50" width="150" height="150" fill="#d9ffb3" />
                <rect x="50" y="220" width="200" height="150" fill="#ffe6b3" />
                <circle v-for="pet in pets" :key="pet.id" :cx="pet.x" :cy="pet.y" r="20" fill="#409EFF" />
              </svg>
            </div>
          </a-card>

          <a-card title="仿真指标" style="margin-top: 16px">
            <a-row :gutter="16">
              <a-col :span="6">
                <a-statistic title="运行时间" :value="simStats.runTime" suffix="秒" />
              </a-col>
              <a-col :span="6">
                <a-statistic title="帧率" :value="simStats.fps" suffix="fps" />
              </a-col>
              <a-col :span="6">
                <a-statistic title="碰撞次数" :value="simStats.collisions" />
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const simForm = reactive({ scene: 'home', duration: 10, speed: '5' })
const simStats = reactive({ runTime: 0, fps: 60, collisions: 0 })
const pets = ref([{ id: 1, x: 200, y: 200 }])

const handleStart = () => { }
const handlePause = () => { }
const handleStop = () => { }
</script>

<style scoped>
.container { padding: 16px; }
.sim-view { background: #f0f0f0; border-radius: 8px; display: flex; justify-content: center; padding: 16px; }
</style>
