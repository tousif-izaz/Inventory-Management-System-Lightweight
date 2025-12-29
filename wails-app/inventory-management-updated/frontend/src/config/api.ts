// API Configuration
// In production (Wails app), we need to call the full localhost URL
// In development, Vite proxy handles /api requests
const API_BASE_URL = import.meta.env.PROD 
    ? 'http://localhost:8080' 
    : '/api';

export const getApiUrl = (path: string): string => {
    // Remove leading slash if present
    const cleanPath = path.startsWith('/') ? path.slice(1) : path;
    
    if (import.meta.env.PROD) {
        // In production, construct full URL
        return `${API_BASE_URL}/${cleanPath}`;
    } else {
        // In development, use /api prefix for Vite proxy
        return `/api/${cleanPath}`;
    }
};

export default API_BASE_URL;
