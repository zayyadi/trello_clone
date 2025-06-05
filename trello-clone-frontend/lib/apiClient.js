import axios from 'axios';

// Base URL for the backend API (e.g., http://localhost:8080)
const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL;
const BASE_AUTH_URL = process.env.NEXT_PUBLIC_AUTH_BASE_URL;

// apiClient for protected routes (e.g., /api/boards)
const apiClient = axios.create({
  baseURL: `${BASE_URL}/api`,
});

// Add a request interceptor to include the JWT token for protected routes
apiClient.interceptors.request.use(
  (config) => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// authApiClient for public authentication routes (e.g., /auth/login)
const authApiClient = axios.create({
  baseURL: `${BASE_AUTH_URL}/auth`,
});

export { apiClient, authApiClient };
