import { User, Tweet, AuthResponse, Conversation, Message } from '../types';

const API_URL = '/api';

let token: string | null = null;

export const api = {
  setToken: (newToken: string | null) => {
    token = newToken;
  },

  async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers: {
        ...headers,
        ...options.headers,
      },
    });

    if (response.status === 401) {
      throw new Error('Unauthorized');
    }

    const text = await response.text();
    if (!text) {
      return null as T;
    }

    let data: unknown;
    try {
      data = JSON.parse(text);
    } catch {
      throw new Error('Invalid JSON response');
    }

    if (!response.ok) {
      const message = typeof data === 'object' && data !== null && 'error' in data
        ? String(data.error)
        : 'Request failed';
      throw new Error(message);
    }

    return data as T;
  },

  async register(username: string, email: string, password: string, name: string) {
    return this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password, name }),
    });
  },

  async login(email: string, password: string) {
    return this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  },

  async getProfile() {
    return this.request<User>('/users/profile');
  },

  async updateProfile(name: string, bio: string) {
    return this.request<User>('/users/update', {
      method: 'POST',
      body: JSON.stringify({ name, bio }),
    });
  },

  async getTweets() {
    const result = await this.request<Tweet[]>('/tweets');
    return result || [];
  },

  async getTweetsByUser(userId: number) {
    const result = await this.request<Tweet[]>(`/tweets/user/${userId}`);
    return result || [];
  },

  async createTweet(content: string) {
    return this.request<Tweet>('/tweets', {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  },

  async likeTweet(tweetId: number) {
    return this.request(`/tweets/like/${tweetId}`, {
      method: 'POST',
    });
  },

  async viewTweet(tweetId: number) {
    return this.request(`/tweets/view/${tweetId}`, {
      method: 'POST',
    });
  },

  async follow(userId: number) {
    return this.request(`/follow/${userId}`, {
      method: 'POST',
    });
  },

  async unfollow(userId: number) {
    return this.request(`/unfollow/${userId}`, {
      method: 'POST',
    });
  },

  async getFollowing(userId: number) {
    const result = await this.request<any[]>(`/following/${userId}`);
    return result || [];
  },

  async getFollowers(userId: number) {
    const result = await this.request<any[]>(`/followers/${userId}`);
    return result || [];
  },

  async searchUsers(query: string) {
    const result = await this.request<User[]>(`/users/search?q=${encodeURIComponent(query)}`);
    return result || [];
  },

  async getConversations() {
    const result = await this.request<Conversation[]>('/conversations');
    return result || [];
  },

  async createConversation(userId: number) {
    return this.request<Conversation>('/conversations', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    });
  },

  async getMessages(conversationId: number) {
    const result = await this.request<Message[]>(`/conversations/${conversationId}/messages`);
    return result || [];
  },

  async sendMessage(conversationId: number, content: string) {
    return this.request<Message>(`/conversations/${conversationId}/messages`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  },

  async updateMessage(messageId: number, content: string) {
    return this.request<Message>(`/messages/${messageId}`, {
      method: 'PUT',
      body: JSON.stringify({ content }),
    });
  },

  async deleteMessage(messageId: number) {
    return this.request<Message>(`/messages/${messageId}`, {
      method: 'DELETE',
    });
  },
};
