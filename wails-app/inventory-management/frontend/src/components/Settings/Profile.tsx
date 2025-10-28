import { useState, useEffect } from 'react';
import { UserIcon, KeyIcon } from '@heroicons/react/24/outline';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { useToast } from '../UI/Toast';
import { api } from '../../services/api';
import { User, UpdatePasswordRequest } from '../../types/user';

export const Profile = () => {
    const { showToast } = useToast();
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const [showPasswordChange, setShowPasswordChange] = useState(false);
    const [passwordData, setPasswordData] = useState<UpdatePasswordRequest>({
        old_password: '',
        new_password: '',
    });
    const [confirmPassword, setConfirmPassword] = useState('');
    const [updating, setUpdating] = useState(false);

    useEffect(() => {
        fetchProfile();
    }, []);

    const fetchProfile = async () => {
        try {
            setLoading(true);
            const data = await api.get<User>('/profile');
            setUser(data);
        } catch (error) {
            console.error('Failed to fetch profile:', error);
            showToast('error', 'Failed to load profile');
        } finally {
            setLoading(false);
        }
    };

    const handlePasswordChange = async (e: React.FormEvent) => {
        e.preventDefault();

        if (passwordData.new_password !== confirmPassword) {
            showToast('error', 'New password and confirm password do not match');
            return;
        }

        if (passwordData.new_password.length < 8) {
            showToast('error', 'New password must be at least 8 characters');
            return;
        }

        try {
            setUpdating(true);
            await api.put('/profile/password', passwordData);
            showToast('success', 'Password updated successfully');
            setShowPasswordChange(false);
            setPasswordData({ old_password: '', new_password: '' });
            setConfirmPassword('');
        } catch (error) {
            console.error('Failed to update password:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to update password';
            showToast('error', errorMessage);
        } finally {
            setUpdating(false);
        }
    };

    const formatDate = (dateString?: string) => {
        if (!dateString) return 'Never';
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
            </div>
        );
    }

    if (!user) {
        return (
            <div className="bg-white shadow rounded-lg p-12 text-center">
                <h3 className="text-lg font-medium text-gray-900 mb-2">Profile not found</h3>
                <p className="text-gray-500">Unable to load user profile</p>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Profile</h1>
                <p className="mt-1 text-sm text-gray-500">
                    View your account information and change your password
                </p>
            </div>

            {/* User Information Card */}
            <div className="bg-white shadow rounded-lg p-6 mb-6">
                <div className="flex items-center mb-6">
                    <div className="h-16 w-16 rounded-full bg-indigo-600 flex items-center justify-center">
                        <UserIcon className="h-8 w-8 text-white" />
                    </div>
                    <div className="ml-4">
                        <h2 className="text-2xl font-bold text-gray-900">{user.name}</h2>
                        <p className="text-sm text-gray-500">@{user.username}</p>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Email</label>
                        <p className="mt-1 text-sm text-gray-900">{user.email || 'Not provided'}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Role</label>
                        <p className="mt-1">
                            <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-indigo-100 text-indigo-800">
                                {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                            </span>
                        </p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Account Status</label>
                        <p className="mt-1">
                            <span
                                className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                                    user.is_active
                                        ? 'bg-green-100 text-green-800'
                                        : 'bg-red-100 text-red-800'
                                }`}
                            >
                                {user.is_active ? 'Active' : 'Inactive'}
                            </span>
                        </p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Last Login</label>
                        <p className="mt-1 text-sm text-gray-900">{formatDate(user.last_login)}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Account Created</label>
                        <p className="mt-1 text-sm text-gray-900">{formatDate(user.created_at)}</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-500">Last Updated</label>
                        <p className="mt-1 text-sm text-gray-900">{formatDate(user.updated_at)}</p>
                    </div>
                </div>
            </div>

            {/* Password Change Card */}
            <div className="bg-white shadow rounded-lg p-6">
                <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center">
                        <KeyIcon className="h-6 w-6 text-gray-400 mr-2" />
                        <h2 className="text-lg font-semibold text-gray-900">Change Password</h2>
                    </div>
                    {!showPasswordChange && (
                        <Button onClick={() => setShowPasswordChange(true)}>Change Password</Button>
                    )}
                </div>

                {showPasswordChange ? (
                    <form onSubmit={handlePasswordChange} className="space-y-4">
                        <Input
                            label="Current Password"
                            type="password"
                            value={passwordData.old_password}
                            onChange={(e) =>
                                setPasswordData({ ...passwordData, old_password: e.target.value })
                            }
                            required
                        />
                        <Input
                            label="New Password"
                            type="password"
                            value={passwordData.new_password}
                            onChange={(e) =>
                                setPasswordData({ ...passwordData, new_password: e.target.value })
                            }
                            required
                            minLength={8}
                            placeholder="At least 8 characters"
                        />
                        <Input
                            label="Confirm New Password"
                            type="password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            required
                            minLength={8}
                        />
                        <div className="flex gap-3 pt-4">
                            <Button
                                type="button"
                                variant="secondary"
                                onClick={() => {
                                    setShowPasswordChange(false);
                                    setPasswordData({ old_password: '', new_password: '' });
                                    setConfirmPassword('');
                                }}
                            >
                                Cancel
                            </Button>
                            <Button type="submit" disabled={updating}>
                                {updating ? 'Updating...' : 'Update Password'}
                            </Button>
                        </div>
                    </form>
                ) : (
                    <p className="text-sm text-gray-500">
                        Keep your account secure by using a strong password
                    </p>
                )}
            </div>
        </div>
    );
};
