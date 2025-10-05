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

export interface PaginationInfo {
	page: number
	pageSize: number
	total: number
	totalPage: number
}

export interface PaginationResponse<T> {
	data: T[]
	pagination: PaginationInfo
}

export interface SelectOption {
  label: string
  value: string | number
}

export type LoadingState = 'idle' | 'loading' | 'success' | 'error'

export type StatusType = 'success' | 'error' | 'warning' | 'info'