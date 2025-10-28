import { Bars3Icon } from '@heroicons/react/24/outline';
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';

interface HeaderProps {
    onMenuClick: () => void;
}

export const Header = ({ onMenuClick }: HeaderProps) => {
    const navigate = useNavigate();
    const [showUserMenu, setShowUserMenu] = useState(false);

    const handleLogout = async () => {
        try {
            const response = await fetch('http://localhost:37285/api/logout', {
                method: 'POST',
                credentials: 'include',
            });

            if (response.ok) {
                // Clear auth token from localStorage
                localStorage.removeItem('auth_token');
                navigate('/login');
            }
        } catch (error) {
            console.error('Logout error:', error);
            // Even if API call fails, clear token and redirect
            localStorage.removeItem('auth_token');
            navigate('/login');
        }
    };

    return (
        <header className="sticky top-0 z-30 h-16 border-b border-gray-200 bg-white">
            <div className="flex h-full items-center justify-between px-4 sm:px-6 lg:px-8">
                {/* Left side - Menu button */}
                <div className="flex items-center">
                    {/* Mobile menu button */}
                    <button
                        onClick={onMenuClick}
                        className="rounded-md p-2 text-gray-600 hover:bg-gray-100 lg:hidden"
                    >
                        <Bars3Icon className="h-6 w-6" />
                    </button>
                </div>

                {/* Right side - User menu */}
                <div className="flex items-center gap-4">
                    {/* User menu */}
                    <div className="relative">
                        <button
                            onClick={() => setShowUserMenu(!showUserMenu)}
                            className="flex items-center gap-2 rounded-md p-2 text-gray-600 hover:bg-gray-100"
                        >
                            <div className="h-8 w-8 rounded-full bg-indigo-600 flex items-center justify-center">
                                <span className="text-sm font-medium text-white">U</span>
                            </div>
                            <div className="hidden md:block text-left">
                                <p className="text-sm font-medium text-gray-900">User</p>
                                <p className="text-xs text-gray-500">Admin</p>
                            </div>
                        </button>

                        {/* Dropdown menu */}
                        {showUserMenu && (
                            <>
                                <div
                                    className="fixed inset-0 z-40"
                                    onClick={() => setShowUserMenu(false)}
                                />
                                <div className="absolute right-0 mt-2 w-48 rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 z-50">
                                    <button
                                        onClick={() => {
                                            setShowUserMenu(false);
                                            navigate('/settings/profile');
                                        }}
                                        className="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
                                    >
                                        Profile
                                    </button>
                                    <button
                                        onClick={() => {
                                            setShowUserMenu(false);
                                            navigate('/settings');
                                        }}
                                        className="block w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
                                    >
                                        Settings
                                    </button>
                                    <hr className="my-1 border-gray-200" />
                                    <button
                                        onClick={handleLogout}
                                        className="block w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-gray-100"
                                    >
                                        Logout
                                    </button>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </div>
        </header>
    );
};
