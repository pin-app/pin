const API_BASE_URL = 'http://localhost:8080';

export interface HealthResponse {
  status: string;
  timestamp: string;
  service: string;
}

export interface ApiError {
  message: string;
  status?: number;
}

export class ApiService {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  async checkHealth(): Promise<HealthResponse> {
    try {
      const response = await fetch(`${this.baseUrl}/health`);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data: HealthResponse = await response.json();
      return data;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Health check failed: ${error.message}`);
      }
      throw new Error('Health check failed: Unknown error');
    }
  }
}

export const apiService = new ApiService();