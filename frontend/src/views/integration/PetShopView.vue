<template>
  <div class="pet-shop-view">
    <a-card title="宠物用品商城">
      <a-alert type="info" style="margin-bottom: 16px">
        第三方宠物用品商城集成，支持宠物食品、玩具、服装等商品的浏览和购买。
      </a-alert>

      <a-row :gutter="16" style="margin-bottom: 24px">
        <a-col :span="6">
          <a-statistic title="本月订单" :value="stats.orders" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="消费金额" :value="stats.spent" prefix="¥" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="收藏商品" :value="stats.favorites" />
        </a-col>
        <a-col :span="6">
          <a-statistic title="优惠券" :value="stats.coupons" />
        </a-col>
      </a-row>

      <a-tabs>
        <a-tab-pane key="products" title="商品列表">
          <a-space style="margin-bottom: 12px">
            <a-input-search v-model="searchKeyword" placeholder="搜索商品" style="width: 300px" @search="loadProducts" />
            <a-select v-model="category" placeholder="分类" style="width: 140px" allow-clear>
              <a-option value="food">食品</a-option>
              <a-option value="toy">玩具</a-option>
              <a-option value="clothing">服装</a-option>
              <a-option value="health">保健</a-option>
            </a-select>
          </a-space>
          <a-table :columns="productColumns" :data="products" :loading="loading" :pagination="pagination" @page-change="onPageChange">
            <template #image="{ record }">
              <a-image width="60" height="60" :src="record.image" />
            </template>
            <template #price="{ record }">
              <span style="color: #f53f3f; font-weight: bold">¥{{ record.price }}</span>
            </template>
            <template #actions="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="addToCart(record)">加入购物车</a-button>
                <a-button type="text" size="small" @click="buyNow(record)">立即购买</a-button>
              </a-space>
            </template>
          </a-table>
        </a-tab-pane>
        <a-tab-pane key="cart" title="购物车">
          <a-empty v-if="cartItems.length === 0" description="购物车是空的" />
          <a-list v-else :data="cartItems" style="max-width: 600px">
            <template #render="{ item }">
              <a-list-item>
                <a-list-item-meta :title="item.name" :description="'¥' + item.price + ' x ' + item.quantity" />
                <template #actions>
                  <a-button type="text" status="danger" size="small" @click="removeCart(item)">删除</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane key="orders" title="我的订单">
          <a-table :columns="orderColumns" :data="orders" :loading="loading" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const searchKeyword = ref('')
const category = ref('')
const products = ref([])
const cartItems = ref([])
const orders = ref([])

const stats = reactive({ orders: 12, spent: 2580, favorites: 8, coupons: 3 })

const productColumns = [
  { title: '商品', slotName: 'image', dataIndex: 'image' },
  { title: '名称', dataIndex: 'name' },
  { title: '分类', dataIndex: 'category' },
  { title: '价格', slotName: 'price' },
  { title: '库存', dataIndex: 'stock' },
  { title: '操作', slotName: 'actions' }
]

const orderColumns = [
  { title: '订单号', dataIndex: 'order_no' },
  { title: '商品', dataIndex: 'product' },
  { title: '金额', dataIndex: 'amount' },
  { title: '状态', dataIndex: 'status' },
  { title: '时间', dataIndex: 'created_at' }
]

const pagination = reactive({ total: 0, current: 1, pageSize: 12 })

async function loadProducts() {
  loading.value = true
  try {
    products.value = [
      { id: 1, name: '天然狗粮 5kg', category: '食品', price: 299, stock: 50, image: '' },
      { id: 2, name: '智能逗猫棒', category: '玩具', price: 129, stock: 30, image: '' },
      { id: 3, name: '宠物冬季外套', category: '服装', price: 199, stock: 25, image: '' },
      { id: 4, name: '宠物益生菌', category: '保健', price: 89, stock: 100, image: '' }
    ]
    pagination.total = 4
  } finally {
    loading.value = false
  }
}

function onPageChange(page) {
  pagination.current = page
  loadProducts()
}

function addToCart(product) {
  const exist = cartItems.value.find(i => i.id === product.id)
  if (exist) exist.quantity++
  else cartItems.value.push({ ...product, quantity: 1 })
  Message.success('已加入购物车')
}

function buyNow(product) {
  Message.info('立即购买: ' + product.name)
}

function removeCart(item) {
  cartItems.value = cartItems.value.filter(i => i.id !== item.id)
}

onMounted(() => loadProducts())
</script>
