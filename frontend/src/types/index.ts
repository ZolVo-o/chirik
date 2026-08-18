export interface User {
  id: number;
  username: string;
  email: string;
  name: string;
  bio: string;
  created_at: string;
}

export interface Tweet {
  id: number;
  user_id: number;
  username: string;
  content: string;
  likes: number;
  views: number;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
