import { useState, useEffect } from 'react';
import { User } from '../types/user';

export const useAuth = () => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Fetch current user from API
        const fetchUser = async () => {
            try {
                const token = localStorage.getItem('auth_token');
                if (!token) {
                    setLoading(false);
                    return;
                }

                const response = await fetch('http://localhost:37285/api/profile', {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                    },
                    credentials: 'include',
                });

                if (response.ok) {
                    const userData = await response.json();
                    setUser(userData);
                    // Store in localStorage for quick access
                    localStorage.setItem('user', JSON.stringify(userData));
                } else {
                    // Clear localStorage if request fails
                    localStorage.removeItem('user');
                    localStorage.removeItem('auth_token');
                }
            } catch (error) {
                console.error('Failed to fetch user profile:', error);
                localStorage.removeItem('user');
            } finally {
                setLoading(false);
            }
        };

        fetchUser();
    }, []);

    const hasRole = (roles: string[]): boolean => {
        if (!user) return false;
        return roles.includes(user.role);
    };

    const canAccessInventoryManagement = (): boolean => {
        // Only admin can add/remove inventory
        if (!user) return false;
        return user.role === 'admin';
    };

    const canAccessSquareSettings = (): boolean => {
        // Only admin can access Square settings
        if (!user) return false;
        return user.role === 'admin';
    };

    const canViewReports = (): boolean => {
        // Admin and manager can view reports
        return user?.role === 'admin' || user?.role === 'manager';
    };

    const canManageSettings = (): boolean => {
        // Only admin can manage settings
        return user?.role === 'admin';
    };

    return {
        user,
        loading,
        hasRole,
        canAccessInventoryManagement,
        canAccessSquareSettings,
        canViewReports,
        canManageSettings,
    };
};
