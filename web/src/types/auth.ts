export interface LoginRequest {
  username: string;
  password: string;
}

export interface User {
  id: number;
  username: string;
  isAdmin: boolean;
}

export interface LoginResponse {
  token: string;
  user: User;
}