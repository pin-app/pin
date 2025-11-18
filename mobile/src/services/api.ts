export const API_BASE_URL =
  process.env.EXPO_PUBLIC_API_BASE ??
  (__DEV__ ? "http://localhost:8080" : "todo, actual link");

export const DEV_MODE = process.env.EXPO_PUBLIC_DEV_MODE === 'true' || __DEV__;

export interface HealthResponse {
  status: string;
  timestamp: string;
  service: string;
  count: number;
}

export interface ApiError {
  message: string;
  status?: number;
}

export interface User {
  id: string;
  email: string;
  username?: string;
  display_name?: string;
  bio?: string;
  location?: string;
  pfp_url?: string;
  created_at: string;
  updated_at: string;
}

export interface PostImage {
  id: string;
  post_id: string;
  image_url: string;
  caption?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: string;
  user_id: string;
  place_id: string;
  description?: string;
  images: PostImage[];
  created_at: string;
  updated_at: string;
  user?: User;
  place?: {
    id: string;
    name: string;
    geometry: string;
    properties: Record<string, any>;
    created_at: string;
    updated_at: string;
  };
  likes_count: number;
  comments_count: number;
  liked_by_user: boolean;
}

export interface Place {
  id: string;
  name: string;
  geometry: string;
  properties: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: string;
  post_id: string;
  user_id: string;
  parent_id?: string;
  content: string;
  created_at: string;
  updated_at: string;
  user?: User;
  replies?: Comment[];
}

export interface OAuthResponse {
  session_token: string;
  user: User;
  expires_at: string;
}

export class ApiService {
  private baseUrl: string;
  private sessionToken: string | null = null;
  private devUserId: string | null = null;
  private isDevMode: boolean = DEV_MODE;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  setSessionToken(token: string | null) {
    this.sessionToken = token;
  }

  setDevMode(enabled: boolean) {
    this.isDevMode = enabled;
  }

  setDevUserId(userId: string | null) {
    this.devUserId = userId;
  }

  private getHeaders(contentType: string | null = 'application/json'): Record<string, string> {
    const headers: Record<string, string> = {};
    if (contentType) {
      headers['Content-Type'] = contentType;
    }

    if (this.isDevMode && this.devUserId) {
      headers['X-Dev-User-ID'] = this.devUserId;
      console.log('🔧 Dev Mode: Using dev user ID:', this.devUserId);
    } else if (this.sessionToken) {
      headers['Authorization'] = `Bearer ${this.sessionToken}`;
      console.log('🔐 Auth: Using session token');
    } else {
      console.log('⚠️ No authentication headers');
    }

    return headers;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers = { ...this.getHeaders(), ...options.headers };

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (!response.ok) {
        let errorData: any = {};
        try {
          const text = await response.text();
          if (text) {
            errorData = JSON.parse(text);
          }
        } catch {
          // If parsing fails, use empty object
        }
        const error = new Error(errorData.error || `HTTP error! status: ${response.status}`);
        (error as any).status = response.status;
        throw error;
      }

      // Handle empty responses (like 204 No Content)
      const text = await response.text();
      if (!text) {
        return null as T;
      }
      
      return JSON.parse(text);
    } catch (error) {
      if (error instanceof Error) {
        // Preserve the original error with status code if it exists
        if ((error as any).status) {
          throw error;
        }
        throw new Error(`API request failed: ${error.message}`);
      }
      throw new Error('API request failed: Unknown error');
    }
  }

  // Health check
  async checkHealth(): Promise<HealthResponse> {
    return this.request<HealthResponse>('/health');
  }

  // OAuth endpoints
  async googleAuth(redirectUrl?: string): Promise<string> {
    const params = new URLSearchParams();
    if (redirectUrl) {
      params.append('redirect_url', redirectUrl);
    }
    const url = `/api/auth/google?${params.toString()}`;
    return `${this.baseUrl}${url}`;
  }

  async appleAuth(redirectUrl?: string): Promise<string> {
    const params = new URLSearchParams();
    if (redirectUrl) {
      params.append('redirect_url', redirectUrl);
    }
    const url = `/api/auth/apple?${params.toString()}`;
    return `${this.baseUrl}${url}`;
  }

  async logout(sessionToken: string): Promise<void> {
    await this.request('/api/auth/logout', {
      method: 'POST',
    });
  }

  // User endpoints
  async getUser(userId: string): Promise<User> {
    return this.request<User>(`/api/users/${userId}`);
  }

  async createUser(userData: Partial<User>): Promise<User> {
    return this.request<User>('/api/users', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async updateUser(userId: string, userData: Partial<User>): Promise<User> {
    return this.request<User>(`/api/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  }

  async searchUsers(query: string, limit = 20, offset = 0): Promise<User[]> {
    const params = new URLSearchParams({
      q: query,
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{users: User[]}>(`/api/users/search?${params.toString()}`);
    return response.users;
  }

  // Following endpoints
  async followUser(userId: string): Promise<void> {
    await this.request<void>(`/api/users/${userId}/follow`, {
      method: 'POST',
    });
  }

  async unfollowUser(userId: string): Promise<void> {
    await this.request<void>(`/api/users/${userId}/follow`, {
      method: 'DELETE',
    });
  }

  async getFollowing(userId: string, limit = 20, offset = 0): Promise<User[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{users: User[]}>(`/api/users/${userId}/following?${params.toString()}`);
    return response.users;
  }

  async getFollowers(userId: string, limit = 20, offset = 0): Promise<User[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{users: User[]}>(`/api/users/${userId}/followers?${params.toString()}`);
    return response.users;
  }

  async isFollowing(userId: string): Promise<boolean> {
    const response = await this.request<{is_following: boolean}>(`/api/users/${userId}/follow-status`);
    return response.is_following;
  }

  // User stats endpoint
  async getUserStats(userId: string): Promise<{
    user_id: string;
    posts_count: number;
    following_count: number;
    followers_count: number;
  }> {
    return this.request(`/api/users/${userId}/stats`);
  }

  // Place endpoints
  async getPlace(placeId: string): Promise<Place> {
    return this.request<Place>(`/api/places/${placeId}`);
  }

  async createPlace(placeData: {
    name: string;
    geometry: string;
    properties: Record<string, any>;
  }): Promise<Place> {
    return this.request<Place>('/api/places', {
      method: 'POST',
      body: JSON.stringify(placeData),
    });
  }

  async searchPlaces(query: string, limit = 20, offset = 0): Promise<Place[]> {
    const params = new URLSearchParams({
      q: query,
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{places: Place[]}>(`/api/places/search?${params.toString()}`);
    return response.places;
  }

  async listPlaces(limit = 20, offset = 0): Promise<Place[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{places: Place[]}>(`/api/places?${params.toString()}`);
    return response.places;
  }

  async searchNearbyPlaces(
    lat: number,
    lng: number,
    radius = 10,
    limit = 20
  ): Promise<Place[]> {
    const params = new URLSearchParams({
      lat: lat.toString(),
      lng: lng.toString(),
      radius: radius.toString(),
      limit: limit.toString(),
    });
    return this.request<Place[]>(`/api/places/nearby?${params.toString()}`);
  }

  // Post endpoints
  async getPosts(limit = 20, offset = 0): Promise<Post[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{posts: Post[]}>(`/api/posts?${params.toString()}`);
    return response.posts;
  }

  async getPost(postId: string): Promise<Post> {
    return this.request<Post>(`/api/posts/${postId}`);
  }

  async createPost(postData: {
    place_id: string;
    description: string;
    images: string[];
  }): Promise<Post> {
    return this.request<Post>('/api/posts', {
      method: 'POST',
      body: JSON.stringify(postData),
    });
  }

  async uploadPostImage(uri: string, fileName?: string): Promise<{ url: string }> {
    const name = fileName || uri.split('/').pop() || `upload-${Date.now()}.jpg`;
    const formData = new FormData();
    formData.append('file', {
      uri,
      name,
      type: this.getMimeType(name),
    } as any);

    const headers = this.getHeaders(null);
    const response = await fetch(`${this.baseUrl}/api/uploads`, {
      method: 'POST',
      headers,
      body: formData,
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(text || 'Failed to upload image');
    }

    return response.json();
  }

  async likePost(postId: string): Promise<{ post_id: string; likes_count: number; liked: boolean }> {
    return this.request<{ post_id: string; likes_count: number; liked: boolean }>(`/api/posts/${postId}/likes`, {
      method: 'POST',
    });
  }

  async unlikePost(postId: string): Promise<{ post_id: string; likes_count: number; liked: boolean }> {
    return this.request<{ post_id: string; likes_count: number; liked: boolean }>(`/api/posts/${postId}/likes`, {
      method: 'DELETE',
    });
  }

  async getPostsByUser(userId: string, limit = 20, offset = 0): Promise<Post[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    console.log('API: Getting posts for user:', userId, 'with params:', params.toString());
    const response = await this.request<{posts: Post[]}>(`/api/users/${userId}/posts?${params.toString()}`);
    console.log('API: Response received:', response);
    return response.posts;
  }

  async getPostsByPlace(placeId: string, limit = 20, offset = 0): Promise<Post[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{posts: Post[]}>(`/api/places/${placeId}/posts?${params.toString()}`);
    return response.posts;
  }

  // Comment methods
  async getCommentsByPost(postId: string, limit = 20, offset = 0): Promise<Comment[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await this.request<{comments: Comment[]}>(`/api/posts/${postId}/comments?${params.toString()}`);
    return response.comments;
  }

  async createComment(postId: string, content: string, parentId?: string): Promise<Comment> {
    const commentData: any = {
      post_id: postId,
      content: content,
    };
    
    if (parentId) {
      commentData.parent_id = parentId;
    }

    return this.request<Comment>('/api/comments', {
      method: 'POST',
      body: JSON.stringify(commentData),
    });
  }

  async updateComment(commentId: string, content: string): Promise<Comment> {
    return this.request<Comment>(`/api/comments/${commentId}`, {
      method: 'PUT',
      body: JSON.stringify({ content }),
    });
  }

  async deleteComment(commentId: string): Promise<void> {
    return this.request<void>(`/api/comments/${commentId}`, {
      method: 'DELETE',
    });
  }

  private getMimeType(fileName: string): string {
    const ext = fileName.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'png':
        return 'image/png';
      case 'jpg':
      case 'jpeg':
        return 'image/jpeg';
      case 'gif':
        return 'image/gif';
      case 'heic':
        return 'image/heic';
      default:
        return 'application/octet-stream';
    }
  }
}

export const apiService = new ApiService();
