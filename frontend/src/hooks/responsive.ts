// hooks/responsive.ts - 响应式钩子
import { ref, onMounted, onUnmounted } from 'vue';
import { useAppStore } from '@/store';

const MOBILE_WIDTH = 992;

export default function useResponsive(immediate = false) {
  const appStore = useAppStore();
  const isMobile = ref(window.innerWidth < MOBILE_WIDTH);

  const handleResize = () => {
    isMobile.value = window.innerWidth < MOBILE_WIDTH;
    appStore.updateSettings({
      device: isMobile.value ? 'mobile' : 'desktop',
    });
  };

  onMounted(() => {
    window.addEventListener('resize', handleResize);
    handleResize(); // 初始化
  });

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
  });

  return {
    isMobile,
  };
}
