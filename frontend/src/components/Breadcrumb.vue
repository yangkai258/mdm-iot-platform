<template>
  <a-breadcrumb class="container-breadcrumb">
    <a-breadcrumb-item v-for="(item, index) in items" :key="index">
      <a :href="item.href" v-if="item.href && index < items.length - 1" @click.prevent="navigate(item.href)">
        {{ item.label }}
      </a>
      <span v-else>{{ item.label }}</span>
    </a-breadcrumb-item>
  </a-breadcrumb>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

defineProps({
  items: {
    type: Array as () => Array<{ label: string; href?: string }>,
    default() {
      return []
    }
  }
})

const router = useRouter()
const navigate = (href: string) => {
  if (href?.startsWith('/')) {
    router.push(href)
  }
}
</script>

<style scoped>
.container-breadcrumb {
  margin-bottom: 16px;
}
</style>
