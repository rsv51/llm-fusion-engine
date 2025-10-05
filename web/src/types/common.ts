// 通用类型定义

export interface ApiResponse<T = any> {
  data?: T
  message?: string
  error?: string
}

export interface PaginationParams {
  page: number
  pageSize: number
}

export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface SelectOption {
  label: string
  value: string | number
}

export type LoadingState = 'idle' | 'loading' | 'success' | 'error'

export type StatusType = 'success' | 'error' | 'warning' | 'info'