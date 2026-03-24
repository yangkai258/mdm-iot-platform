// hooks/locale.ts - 国际化钩子（简化版，不依赖 vue-i18n）
import { ref } from 'vue';

export default function useLocale() {
  // 从 localStorage 读取当前语言，默认中文
  const currentLocale = ref(localStorage.getItem('locale') || 'zh-CN');
  
  // 切换语言
  const setLocale = (newLocale: string) => {
    currentLocale.value = newLocale;
    localStorage.setItem('locale', newLocale);
  };
  
  // 切换到中文
  const localeZh = () => setLocale('zh-CN');
  
  // 切换到英文
  const localeEn = () => setLocale('en-US');
  
  return {
    currentLocale,
    setLocale,
    localeZh,
    localeEn,
  };
}
