export interface PageParams { page?: number; pageSize?: number; keyword?: string; }
export interface PageResult<T> { items: T[]; total: number; page: number; pageSize: number; }
