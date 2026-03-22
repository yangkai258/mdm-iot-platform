<template>
  <div class="page-container">
    <!-- 面包屑 -->
    <a-breadcrumb class="breadcrumb">
      <a-breadcrumb-item>首页</a-breadcrumb-item>
      <a-breadcrumb-item>第三方集成</a-breadcrumb-item>
      <a-breadcrumb-item>宠物用品</a-breadcrumb-item>
    </a-breadcrumb>

    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-left">
        <h2>宠物用品商城</h2>
      </div>
      <div class="header-right">
        <a-button type="primary" status="warning" @click="showCartDrawer = true">
          <template #icon><icon-shopping-cart /></template>
          购物车
          <a-badge v-if="cartItems.length > 0" :count="cartItems.length" :max-count="99" />
        </a-button>
      </div>
    </div>

    <!-- 搜索筛选栏 -->
    <div class="filter-bar">
      <a-space wrap>
        <a-input-search
          v-model="searchKeyword"
          placeholder="搜索商品..."
          style="width: 280px"
          search-button
          @search="loadProducts"
        />
        <a-select
          v-model="filterCategory"
          placeholder="商品分类"
          style="width: 140px"
          allow-clear
          @change="loadProducts"
        >
          <a-option value="food">宠物食品</a-option>
          <a-option value="toy">玩具用品</a-option>
          <a-option value="care">护理清洁</a-option>
          <a-option value="clothing">宠物服饰</a-option>
          <a-option value="accessories">配件用具</a-option>
        </a-select>
        <a-select
          v-model="filterTag"
          placeholder="商品标签"
          style="width: 140px"
          allow-clear
          @change="loadProducts"
        >
          <a-option value="hot">热销</a-option>
          <a-option value="new">新品</a-option>
          <a-option value="recommend">推荐</a-option>
          <a-option value="discount">折扣</a-option>
        </a-select>
        <a-radio-group v-model="sortBy" type="button" size="small" @change="loadProducts">
          <a-radio value="default">默认</a-radio>
          <a-radio value="price_asc">价格↑</a-radio>
          <a-radio value="price_desc">价格↓</a-radio>
          <a-radio value="sales">销量</a-radio>
        </a-radio-group>
      </a-space>
    </div>

    <a-row :gutter="16">
      <!-- 商品列表 -->
      <a-col :span="17">
        <div class="product-list">
          <div v-if="loading" class="loading-state">
            <a-spin size="large" />
          </div>
          <div v-else-if="filteredProducts.length === 0" class="empty-state">
            <icon-shopping-bg class="empty-icon" />
            <p>暂无商品</p>
          </div>
          <div v-else class="product-grid">
            <div v-for="product in filteredProducts" :key="product.id" class="product-card">
              <div class="product-image" @click="viewProductDetail(product)">
                <img v-if="product.image" :src="product.image" :alt="product.name" />
                <div v-else class="image-placeholder">
                  <icon-image :size="40" style="color: #ccc" />
                </div>
                <div class="product-tags">
                  <span v-if="product.tag === 'hot'" class="tag tag-hot">热销</span>
                  <span v-if="product.tag === 'new'" class="tag tag-new">新品</span>
                  <span v-if="product.tag === 'discount'" class="tag tag-discount">折扣</span>
                </div>
              </div>
              <div class="product-info">
                <div class="product-name" @click="viewProductDetail(product)">{{ product.name }}</div>
                <div class="product-desc">{{ product.description }}</div>
                <div class="product-rating">
                  <a-rate :value="product.rating" :readonly="true" :size="12" />
                  <span class="rating-count">{{ product.rating_count }}条评价</span>
                </div>
                <div class="product-price-row">
                  <div class="price-info">
                    <span class="current-price">¥{{ product.price }}</span>
                    <span v-if="product.original_price" class="original-price">¥{{ product.original_price }}</span>
                  </div>
                  <a-button type="primary" size="small" @click.stop="addToCart(product)">
                    <template #icon><icon-cart /></template>
                    加购
                  </a-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </a-col>

      <!-- 推荐商品侧边栏 -->
      <a-col :span="7">
        <div class="recommend-sidebar">
          <div class="sidebar-title">
            <icon-thumb-up :size="18" />
            <span>为你推荐</span>
          </div>
          <div class="recommend-list">
            <div
              v-for="item in recommendedProducts"
              :key="item.id"
              class="recommend-item"
              @click="viewProductDetail(item)"
            >
              <img v-if="item.image" :src="item.image" :alt="item.name" class="recommend-img" />
              <div v-else class="recommend-img-placeholder">
                <icon-image :size="20" style="color: #ccc" />
              </div>
              <div class="recommend-info">
                <div class="recommend-name">{{ item.name }}</div>
                <div class="recommend-price">¥{{ item.price }}</div>
              </div>
            </div>
          </div>
        </div>
      </a-col>
    </a-row>

    <!-- 商品详情弹窗 -->
    <a-modal
      v-model:visible="detailModalVisible"
      :title="currentProduct?.name"
      :width="680"
      :footer="null"
    >
      <div v-if="currentProduct" class="product-detail">
        <a-row :gutter="24">
          <a-col :span="10">
            <div class="detail-image">
              <img v-if="currentProduct.image" :src="currentProduct.image" :alt="currentProduct.name" />
              <div v-else class="image-placeholder large">
                <icon-image :size="60" style="color: #ccc" />
              </div>
            </div>
          </a-col>
          <a-col :span="14">
            <div class="detail-info">
              <div class="detail-name">{{ currentProduct.name }}</div>
              <div class="detail-desc">{{ currentProduct.description }}</div>
              <div class="detail-rating">
                <a-rate :value="currentProduct.rating" :readonly="true" />
                <span>{{ currentProduct.rating_count }}条评价</span>
              </div>
              <div class="detail-price-row">
                <span class="detail-current-price">¥{{ currentProduct.price }}</span>
                <span v-if="currentProduct.original_price" class="detail-original-price">
                  ¥{{ currentProduct.original_price }}
                </span>
                <span v-if="currentProduct.original_price" class="detail-discount">
                  {{ Math.round((1 - currentProduct.price / currentProduct.original_price) * 100) }}折
                </span>
              </div>
              <div class="detail-meta">
                <div class="meta-item">
                  <span class="meta-label">品牌:</span>
                  <span>{{ currentProduct.brand || '官方' }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">分类:</span>
                  <span>{{ getCategoryName(currentProduct.category) }}</span>
                </div>
                <div class="meta-item">
                  <span class="meta-label">销量:</span>
                  <span>{{ currentProduct.sales || 0 }}</span>
                </div>
              </div>
              <div class="detail-quantity">
                <span class="quantity-label">数量:</span>
                <a-input-number v-model="buyQuantity" :min="1" :max="99" size="medium" />
                <span class="quantity-stock">库存: {{ currentProduct.stock || 100 }}</span>
              </div>
              <div class="detail-actions">
                <a-button type="primary" size="large" @click="addToCartAndClose">
                  <template #icon><icon-cart /></template>
                  加入购物车
                </a-button>
                <a-button status="warning" size="large" @click="buyNow">
                  立即购买
                </a-button>
              </div>
            </div>
          </a-col>
        </a-row>

        <!-- 商品详情 -->
        <a-divider>商品详情</a-divider>
        <div class="detail-content">
          <p>{{ currentProduct.description || '暂无详细描述' }}</p>
        </div>
      </div>
    </a-modal>

    <!-- 购物车抽屉 -->
    <a-drawer
      v-model:visible="showCartDrawer"
      title="购物车"
      :width="460"
      @cancel="showCartDrawer = false"
    >
      <div v-if="cartItems.length === 0" class="cart-empty">
        <icon-shopping-cart :size="48" style="color: #ccc" />
        <p>购物车是空的</p>
        <a-button type="primary" @click="showCartDrawer = false">去逛逛</a-button>
      </div>
      <div v-else>
        <div class="cart-list">
          <div v-for="item in cartItems" :key="item.id" class="cart-item">
            <div class="cart-item-info">
              <div class="cart-item-name">{{ item.name }}</div>
              <div class="cart-item-price">¥{{ item.price }}</div>
            </div>
            <div class="cart-item-controls">
              <a-input-number
                v-model="item.quantity"
                :min="1"
                :max="99"
                size="small"
                @change="updateCartItem(item)"
              />
              <a-button type="text" status="danger" size="small" @click="removeFromCart(item)">
                <template #icon><icon-delete /></template>
              </a-button>
            </div>
          </div>
        </div>
        <a-divider />
        <div class="cart-summary">
          <div class="summary-row">
            <span>商品件数:</span>
            <span>{{ totalQuantity }}件</span>
          </div>
          <div class="summary-row total">
            <span>总计:</span>
            <span class="total-price">¥{{ totalPrice }}</span>
          </div>
          <a-button type="primary" long size="large" @click="checkout">
            结算
          </a-button>
        </div>
      </div>
    </a-drawer>

    <!-- 下单成功提示 -->
    <a-modal v-model:visible="orderSuccessVisible" title="下单成功" :footer="null" @ok="orderSuccessVisible = false">
      <div class="order-success">
        <icon-check-circle :size="64" style="color: #00b42a; margin-bottom: 16px" />
        <p>您的订单已提交成功！</p>
        <p class="order-info">订单编号: {{ orderId }}</p>
        <a-button type="primary" @click="orderSuccessVisible = false">确定</a-button>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import {
  IconShoppingCart,
  IconThumbUp,
  IconImage,
  IconCart,
  IconDelete,
  IconCheckCircle
} from '@arco-design/web-icons/vue'

const loading = ref(false)
const products = ref([])
const searchKeyword = ref('')
const filterCategory = ref('')
const filterTag = ref('')
const sortBy = ref('default')
const detailModalVisible = ref(false)
const showCartDrawer = ref(false)
const currentProduct = ref(null)
const buyQuantity = ref(1)
const orderSuccessVisible = ref(false)
const orderId = ref('')
const cartItems = ref([])

const recommendedProducts = ref([
  { id: 'r1', name: '智能宠物饮水机', price: 129, image: '' },
  { id: 'r2', name: '宠物自动喂食器', price: 259, image: '' },
  { id: 'r3', name: '宠物外出便携包', price: 89, image: '' },
  { id: 'r4', name: '宠物电动梳毛器', price: 69, image: '' },
  { id: 'r5', name: '宠物LED发光项圈', price: 39, image: '' }
])

const filteredProducts = computed(() => {
  let result = products.value
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    result = result.filter(p => p.name.toLowerCase().includes(kw) || p.description?.toLowerCase().includes(kw))
  }
  if (filterCategory.value) {
    result = result.filter(p => p.category === filterCategory.value)
  }
  if (filterTag.value) {
    result = result.filter(p => p.tag === filterTag.value)
  }
  if (sortBy.value === 'price_asc') {
    result = [...result].sort((a, b) => a.price - b.price)
  } else if (sortBy.value === 'price_desc') {
    result = [...result].sort((a, b) => b.price - a.price)
  } else if (sortBy.value === 'sales') {
    result = [...result].sort((a, b) => (b.sales || 0) - (a.sales || 0))
  }
  return result
})

const totalQuantity = computed(() => {
  return cartItems.value.reduce((sum, item) => sum + item.quantity, 0)
})

const totalPrice = computed(() => {
  return cartItems.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
})

function getCategoryName(cat) {
  const map = {
    food: '宠物食品',
    toy: '玩具用品',
    care: '护理清洁',
    clothing: '宠物服饰',
    accessories: '配件用具'
  }
  return map[cat] || cat
}

async function loadProducts() {
  loading.value = true
  try {
    const res = await fetch('/api/v1/pet-shop/products', {
      headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` }
    })
    const data = await res.json()
    if (data.data) {
      products.value = data.data
    } else {
      loadMockProducts()
    }
  } catch {
    loadMockProducts()
  } finally {
    loading.value = false
  }
}

function loadMockProducts() {
  products.value = [
    {
      id: 'p1',
      name: '天然狗粮 成犬专用 10kg',
      description: '精选鸡肉配方，营养均衡',
      price: 268,
      original_price: 328,
      category: 'food',
      tag: 'hot',
      rating: 4.8,
      rating_count: 2356,
      sales: 5800,
      stock: 200,
      image: ''
    },
    {
      id: 'p2',
      name: '猫咪逗猫棒 套装',
      description: '多种款式混合套装',
      price: 29.9,
      original_price: null,
      category: 'toy',
      tag: 'new',
      rating: 4.5,
      rating_count: 890,
      sales: 1200,
      stock: 500,
      image: ''
    },
    {
      id: 'p3',
      name: '宠物沐浴露 500ml',
      description: '温和去味，滋润护毛',
      price: 49,
      original_price: 69,
      category: 'care',
      tag: 'discount',
      rating: 4.6,
      rating_count: 1500,
      sales: 3200,
      stock: 300,
      image: ''
    },
    {
      id: 'p4',
      name: '宠物外出航空箱',
      description: '大型犬适用，通风设计',
      price: 159,
      original_price: null,
      category: 'accessories',
      tag: 'recommend',
      rating: 4.7,
      rating_count: 680,
      sales: 890,
      stock: 80,
      image: ''
    },
    {
      id: 'p5',
      name: '猫咪小鱼干零食 100g×3包',
      description: '冻干工艺，美味营养',
      price: 39.9,
      original_price: 59,
      category: 'food',
      tag: 'hot',
      rating: 4.9,
      rating_count: 4200,
      sales: 9600,
      stock: 600,
      image: ''
    },
    {
      id: 'p6',
      name: '宠物电热保暖窝',
      description: '冬季保暖，可拆洗',
      price: 129,
      original_price: 169,
      category: 'accessories',
      tag: 'new',
      rating: 4.4,
      rating_count: 450,
      sales: 560,
      stock: 50,
      image: ''
    },
    {
      id: 'p7',
      name: '宠物指甲剪套装',
      description: '安全剪甲，不伤宠物',
      price: 25,
      original_price: null,
      category: 'care',
      tag: 'recommend',
      rating: 4.3,
      rating_count: 920,
      sales: 1800,
      stock: 400,
      image: ''
    },
    {
      id: 'p8',
      name: '狗狗四脚衣 春秋款',
      description: '防脏防水，时尚舒适',
      price: 69,
      original_price: 89,
      category: 'clothing',
      tag: 'discount',
      rating: 4.5,
      rating_count: 340,
      sales: 780,
      stock: 120,
      image: ''
    }
  ]
}

function viewProductDetail(product) {
  currentProduct.value = product
  buyQuantity.value = 1
  detailModalVisible.value = true
}

function addToCart(product) {
  const existing = cartItems.value.find(item => item.id === product.id)
  if (existing) {
    existing.quantity += 1
    Message.success('已增加数量')
  } else {
    cartItems.value.push({
      id: product.id,
      name: product.name,
      price: product.price,
      quantity: 1
    })
    Message.success('已加入购物车')
  }
}

function addToCartAndClose() {
  if (currentProduct.value) {
    addToCart(currentProduct.value)
    detailModalVisible.value = false
  }
}

function buyNow() {
  if (currentProduct.value) {
    const item = {
      id: currentProduct.value.id + '_buynow',
      product_id: currentProduct.value.id,
      name: currentProduct.value.name,
      price: currentProduct.value.price,
      quantity: buyQuantity.value
    }
    cartItems.value.push(item)
    detailModalVisible.value = false
    showCartDrawer.value = true
    Message.info('已加入购物车，请点击结算')
  }
}

function updateCartItem(item) {
  // 实时更新，已通过 v-model 绑定
}

function removeFromCart(item) {
  const index = cartItems.value.findIndex(i => i.id === item.id)
  if (index !== -1) {
    cartItems.value.splice(index, 1)
    Message.success('已从购物车移除')
  }
}

function checkout() {
  if (cartItems.value.length === 0) {
    Message.warning('购物车为空')
    return
  }
  orderId.value = `ORD${Date.now()}`
  cartItems.value = []
  orderSuccessVisible.value = true
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.page-container {
  padding: 24px;
  min-height: 100vh;
  background: #f5f6f7;
}

.breadcrumb {
  margin-bottom: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-1);
}

.filter-bar {
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.product-list {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 16px;
}

.empty-icon {
  font-size: 48px;
  color: #ccc;
}

.empty-state p {
  color: #999;
  font-size: 14px;
  margin: 0;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.product-card {
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: box-shadow 0.2s, border-color 0.2s;
}

.product-card:hover {
  border-color: #1650d8;
  box-shadow: 0 4px 12px rgba(22, 80, 216, 0.1);
}

.product-image {
  position: relative;
  height: 180px;
  background: #f7f8fa;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f1f2;
}

.image-placeholder.large {
  height: 280px;
}

.product-tags {
  position: absolute;
  top: 8px;
  left: 8px;
  display: flex;
  gap: 4px;
}

.tag {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
  color: #fff;
}

.tag-hot { background: #ff4d4f; }
.tag-new { background: #1650d8; }
.tag-discount { background: #ff7d00; }

.product-info {
  padding: 12px;
}

.product-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-desc {
  font-size: 12px;
  color: #86909c;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 8px;
}

.rating-count {
  font-size: 12px;
  color: #86909c;
}

.product-price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price-info {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.current-price {
  font-size: 16px;
  font-weight: 700;
  color: #ff4d4f;
}

.original-price {
  font-size: 12px;
  color: #c0c1c6;
  text-decoration: line-through;
}

/* 推荐侧边栏 */
.recommend-sidebar {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
}

.sidebar-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--color-text-1);
}

.recommend-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recommend-item {
  display: flex;
  gap: 12px;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: background 0.2s;
}

.recommend-item:hover {
  background: #f7f8fa;
}

.recommend-img,
.recommend-img-placeholder {
  width: 64px;
  height: 64px;
  border-radius: 6px;
  background: #f0f1f2;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.recommend-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.recommend-info {
  flex: 1;
  min-width: 0;
}

.recommend-name {
  font-size: 13px;
  color: var(--color-text-1);
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.recommend-price {
  font-size: 14px;
  font-weight: 600;
  color: #ff4d4f;
}

/* 商品详情弹窗 */
.product-detail {
  padding: 8px 0;
}

.detail-image {
  background: #f7f8fa;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 280px;
}

.detail-image img {
  width: 100%;
  height: 280px;
  object-fit: cover;
}

.detail-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.detail-desc {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 12px;
}

.detail-rating {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #86909c;
  margin-bottom: 12px;
}

.detail-price-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 16px;
}

.detail-current-price {
  font-size: 24px;
  font-weight: 700;
  color: #ff4d4f;
}

.detail-original-price {
  font-size: 14px;
  color: #c0c1c6;
  text-decoration: line-through;
}

.detail-discount {
  font-size: 12px;
  color: #ff4d4f;
  background: #fff1f0;
  padding: 2px 6px;
  border-radius: 4px;
}

.detail-meta {
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  font-size: 13px;
  margin-bottom: 4px;
}

.meta-label {
  color: #86909c;
  width: 60px;
}

.detail-quantity {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.quantity-label {
  font-size: 14px;
  color: var(--color-text-1);
}

.quantity-stock {
  font-size: 13px;
  color: #86909c;
}

.detail-actions {
  display: flex;
  gap: 12px;
}

/* 购物车 */
.cart-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 16px;
}

.cart-empty p {
  color: #999;
  margin: 0;
}

.cart-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 50vh;
  overflow-y: auto;
}

.cart-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f7f8fa;
  border-radius: 8px;
}

.cart-item-info {
  flex: 1;
  min-width: 0;
}

.cart-item-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-1);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 4px;
}

.cart-item-price {
  font-size: 14px;
  font-weight: 600;
  color: #ff4d4f;
}

.cart-item-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cart-summary {
  padding-top: 8px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: var(--color-text-1);
  margin-bottom: 8px;
}

.summary-row.total {
  font-weight: 600;
  font-size: 16px;
  margin-bottom: 16px;
}

.total-price {
  font-size: 20px;
  font-weight: 700;
  color: #ff4d4f;
}

/* 下单成功 */
.order-success {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 0;
  text-align: center;
}

.order-success p {
  font-size: 16px;
  margin: 0 0 8px;
}

.order-info {
  color: #86909c;
  font-size: 13px;
  margin-bottom: 16px;
}
</style>
