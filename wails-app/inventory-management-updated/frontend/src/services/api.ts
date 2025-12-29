// API helper functions
const API_BASE = 'http://localhost:37285/api';

interface ApiError {
    error_message: string;
}

class ApiService {
    private async handleResponse<T>(response: Response): Promise<T> {
        if (!response.ok) {
            const error: ApiError = await response.json().catch(() => ({
                error_message: 'An unexpected error occurred',
            }));

            // Enhanced error logging for debugging
            console.error('API Error:', {
                url: response.url,
                status: response.status,
                statusText: response.statusText,
                errorMessage: error.error_message,
            });

            throw new Error(error.error_message);
        }

        // Handle 204 No Content or 201 Created with no body
        if (response.status === 204 || response.status === 201) {
            // Check if there's actually content
            const contentType = response.headers.get('content-type');
            if (!contentType || !contentType.includes('application/json')) {
                return {} as T;
            }
        }

        // Try to parse JSON, return empty object if no content
        const text = await response.text();
        if (!text || text.trim() === '') {
            return {} as T;
        }

        try {
            return JSON.parse(text);
        } catch (error) {
            console.error('JSON parse error:', error, 'Response text:', text);
            return {} as T;
        }
    }

    private getHeaders(): HeadersInit {
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
        };
        // Primary auth: cookies (credentials: 'include')
        // Fallback: Authorization header from localStorage
        const token = localStorage.getItem('auth_token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        return headers;
    }

    async get<T>(endpoint: string): Promise<T> {
        console.log(`API GET: ${API_BASE}${endpoint}`);
        const response = await fetch(`${API_BASE}${endpoint}`, {
            headers: this.getHeaders(),
            credentials: 'include',
        });
        return this.handleResponse<T>(response);
    }

    async post<T>(endpoint: string, data: unknown): Promise<T> {
        console.log(`API POST: ${API_BASE}${endpoint}`, data);
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: 'POST',
            headers: this.getHeaders(),
            credentials: 'include',
            body: JSON.stringify(data),
        });
        return this.handleResponse<T>(response);
    }

    async put<T>(endpoint: string, data: unknown): Promise<T> {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: 'PUT',
            headers: this.getHeaders(),
            credentials: 'include',
            body: JSON.stringify(data),
        });
        return this.handleResponse<T>(response);
    }

    async delete<T>(endpoint: string): Promise<T> {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: 'DELETE',
            headers: this.getHeaders(),
            credentials: 'include',
        });
        return this.handleResponse<T>(response);
    }
}

export const api = new ApiService();
