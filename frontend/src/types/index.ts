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

export interface Conversation {
  id: number;
  other_user: User;
  last_message?: Message;
  updated_at: string;
}

export interface Message {
  id: number;
  conversation_id: number;
  sender_id: number;
  sender_username: string;
  content: string;
  created_at: string;
  updated_at: string;
  deleted: boolean;
}
