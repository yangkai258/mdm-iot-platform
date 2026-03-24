// hooks/permission.ts - 权限钩子
import { useUserStore } from '@/store';

export default function usePermission() {
  const userStore = useUserStore();

  // 检查是否有权限访问路由
  const accessRouter = (route: { meta?: { roles?: string[] } }) => {
    const { roles } = route.meta || {};
    
    // 如果路由没有设置角色要求，直接通过
    if (!roles || roles.length === 0) {
      return true;
    }
    
    // 检查用户角色是否在允许列表中
    const userRole = userStore.role;
    return roles.includes(userRole);
  };

  return {
    accessRouter,
  };
}
