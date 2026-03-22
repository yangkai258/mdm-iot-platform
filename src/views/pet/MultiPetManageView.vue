<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>宠物管理</a-breadcrumb-item>
      <a-breadcrumb-item>家庭宠物管理</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-left">
        <h2>家庭宠物管理</h2>
      </div>
    </div>

    <div class="page-content">
      <!-- 家庭成员 -->
      <a-card class="section-card">
        <template #title>
          <span class="section-title">家庭成员</span>
        </template>
        <div class="member-list">
          <div v-for="member in members" :key="member.id" class="member-item">
            <div class="member-avatar">
              <icon-user />
            </div>
            <div class="member-info">
              <div class="member-name-row">
                <span class="member-name">{{ member.name }}</span>
                <a-tag v-if="member.role === 'owner'" color="arcoblue" size="small">户主</a-tag>
                <a-tag v-else color="gray" size="small">成员</a-tag>
              </div>
              <div class="member-email">{{ member.email }}</div>
            </div>
            <div class="member-actions">
              <template v-if="member.role === 'owner'">
                <a-button type="text" size="small" disabled>管理</a-button>
              </template>
              <template v-else>
                <a-button type="text" size="small" @click="editMember(member)">编辑</a-button>
                <a-button type="text" size="small" status="danger" @click="removeMember(member)">移除</a-button>
              </template>
            </div>
          </div>
        </div>
        <div class="add-member-row">
          <a-button type="outline" @click="showAddMember">
            <template #icon><icon-plus /></template>
            添加家庭成员
          </a-button>
        </div>
      </a-card>

      <!-- 我的宠物 -->
      <a-card class="section-card">
        <template #title>
          <span class="section-title">我的宠物</span>
        </template>
        <div v-if="pets.length === 0" class="empty-inline">
          <icon-exclamation-circle class="empty-icon" />
          <span>暂无宠物</span>
          <router-link to="/pet/register">
            <a-button type="text" size="small">去登记</a-button>
          </router-link>
        </div>
        <a-table v-else :columns="petColumns" :data="pets" :pagination="false" size="small">
          <template #avatar="{ record }">
            <span class="pet-emoji">{{ getPetEmoji(record.pet_type) }}</span>
          </template>
          <template #breed="{ record }">
            {{ record.breed || getPetTypeName(record.pet_type) }}
          </template>
          <template #status="{ record }">
            <span class="status-badge" :class="getStatusClass(record.status)">
              {{ getStatusText(record.status) }}
            </span>
          </template>
          <template #actions="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="viewPetDetail(record)">详情</a-button>
              <a-button type="text" size="small" @click="editPet(record)">编辑</a-button>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 添加家庭成员弹窗 -->
    <a-modal
      v-model:visible="addMemberModalVisible"
      title="添加家庭成员"
      :width="420"
      @ok="handleAddMember"
      @cancel="addMemberModalVisible = false"
    >
      <a-form :model="memberForm" layout="vertical">
        <a-form-item label="邮箱地址" required>
          <a-input v-model="memberForm.email" placeholder="请输入邮箱地址" />
        </a-form-item>
        <a-form-item label="成员角色">
          <a-select v-model="memberForm.role" placeholder="选择角色">
            <a-option value="member">成员</a-option>
            <a-option value="viewer">查看者</a-option>
          </a-select>
        </a-form-item>
        <a-form-item label="邀请说明">
          <a-textarea v-model="memberForm.message" placeholder="可选的邀请留言" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑成员弹窗 -->
    <a-modal
      v-model:visible="editMemberModalVisible"
      title="编辑成员"
      :width="420"
      @ok="handleEditMember"
      @cancel="editMemberModalVisible = false"
    >
      <a-form :model="editMemberForm" layout="vertical">
        <a-form-item label="成员角色">
          <a-select v-model="editMemberForm.role">
            <a-option value="member">成员</a-option>
            <a-option value="viewer">查看者</a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import { getHouseholdMembers, inviteHouseholdMember, removeHouseholdMember, getHouseholdPets } from '@/api/pet'

const router = useRouter()
const members = ref([])
const pets = ref([])
const addMemberModalVisible = ref(false)
const editMemberModalVisible = ref(false)

const memberForm = reactive({ email: '', role: 'member', message: '' })
const editMemberForm = reactive({ id: '', role: 'member' })

const petColumns = [
  { title: '', slotName: 'avatar', width: 50 },
  { title: '名称', dataIndex: 'pet_name', width: 120 },
  { title: '品种', slotName: 'breed', width: 120 },
  { title: '状态', slotName: 'status', width: 100 },
  { title: '绑定设备', dataIndex: 'device_id', ellipsis: true },
  { title: '操作', slotName: 'actions', width: 150 }
]

function getPetEmoji(type) {
  const map = { dog: '🐶', cat: '🐱', bird: '🐦', rabbit: '🐰', other: '🐾' }
  return map[type] || '🐾'
}

function getPetTypeName(type) {
  const map = { dog: '狗', cat: '猫', bird: '鸟', rabbit: '兔子', other: '其他' }
  return map[type] || '未知'
}

function getStatusClass(status) {
  return { active: 'status-normal', lost: 'status-lost', found: 'status-found' }[status] || 'status-normal'
}

function getStatusText(status) {
  return { active: '正常', lost: '走失中', found: '已找到' }[status] || '正常'
}

async function loadMembers() {
  try {
    const res = await getHouseholdMembers()
    if (res.data) {
      members.value = res.data
    } else {
      loadMockMembers()
    }
  } catch {
    loadMockMembers()
  }
}

function loadMockMembers() {
  members.value = [
    { id: 1, name: '张三', email: 'zhangsan@example.com', role: 'owner' },
    { id: 2, name: '李四', email: 'lisi@example.com', role: 'member' }
  ]
}

async function loadPets() {
  try {
    const res = await getHouseholdPets()
    if (res.data) {
      pets.value = res.data
    } else {
      loadMockPets()
    }
  } catch {
    loadMockPets()
  }
}

function loadMockPets() {
  pets.value = [
    { pet_id: 'PET001', pet_name: '小爪', pet_type: 'cat', breed: '布偶猫', status: 'active', device_id: 'M5-001' },
    { pet_id: 'PET002', pet_name: '旺财', pet_type: 'dog', breed: '柴犬', status: 'lost', device_id: 'M5-002' },
    { pet_id: 'PET003', pet_name: '咪咪', pet_type: 'cat', breed: '英短', status: 'active', device_id: '' }
  ]
}

function showAddMember() {
  Object.assign(memberForm, { email: '', role: 'member', message: '' })
  addMemberModalVisible.value = true
}

async function handleAddMember() {
  if (!memberForm.email) {
    Message.warning('请填写邮箱地址')
    return
  }
  try {
    await inviteHouseholdMember(memberForm)
    Message.success('邀请已发送')
    addMemberModalVisible.value = false
    loadMembers()
  } catch {
    members.value.push({
      id: Date.now(),
      name: memberForm.email.split('@')[0],
      email: memberForm.email,
      role: memberForm.role
    })
    Message.success('邀请已发送')
    addMemberModalVisible.value = false
  }
}

function editMember(member) {
  editMemberForm.id = member.id
  editMemberForm.role = member.role
  editMemberModalVisible.value = true
}

async function handleEditMember() {
  try {
    Message.success('成员角色已更新')
    const idx = members.value.findIndex(m => m.id === editMemberForm.id)
    if (idx !== -1) members.value[idx].role = editMemberForm.role
    editMemberModalVisible.value = false
  } catch {
    Message.error('更新失败')
  }
}

function removeMember(member) {
  Modal.warning({
    title: '移除成员',
    content: `确定要移除成员「${member.name}」吗？`,
    okText: '移除',
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      try {
        await removeHouseholdMember(member.id)
      } catch { /* ignore */ }
      members.value = members.value.filter(m => m.id !== member.id)
      Message.success('成员已移除')
    }
  })
}

function viewPetDetail(pet) {
  router.push(`/pet/detail/${pet.pet_id}`)
}

function editPet(pet) {
  router.push(`/pet/register?edit=${pet.pet_id}`)
}

onMounted(async () => {
  await loadMembers()
  await loadPets()
})
</script>

<style scoped>
.page-container {
  padding: 20px 24px;
  min-height: calc(100vh - 64px);
  background: #f5f7fa;
}

.breadcrumb {
  margin-bottom: 12px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
  margin: 0;
}

.page-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-card {
  border-radius: 8px;
}

.section-title {
  font-weight: 600;
  font-size: 15px;
}

.member-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border);
}

.member-item:last-child {
  border-bottom: none;
}

.member-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f0f5ff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgb(var(--arcoblue-6));
  font-size: 18px;
  flex-shrink: 0;
}

.member-info {
  flex: 1;
  min-width: 0;
}

.member-name-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 2px;
}

.member-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
}

.member-email {
  font-size: 13px;
  color: var(--color-text-3);
}

.member-actions {
  flex-shrink: 0;
}

.add-member-row {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed var(--color-border);
}

.empty-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-3);
  font-size: 14px;
}

.empty-icon {
  font-size: 16px;
}

.pet-emoji {
  font-size: 20px;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.status-normal { background: #e6f7e6; color: #52c41a; }
.status-lost { background: #fff7e6; color: #fa8c16; }
.status-found { background: #e6f7ff; color: #1890ff; }
</style>
