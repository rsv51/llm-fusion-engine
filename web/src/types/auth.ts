export interface LoginRequest {
  username: string;
  password?: string;
  github_token?: string;
}

export interface User {
  id: number;
  username: string;
  name: string;
  email: string;
  role: number;
  status: number;
}

export interface LoginResponse {
  token: string;
  user: User;
}