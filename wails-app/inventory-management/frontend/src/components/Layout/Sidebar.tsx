import { Link, useLocation } from 'react-router-dom';
import {
    ChartBarIcon,
    CubeIcon,
    ArrowsRightLeftIcon,
    UsersIcon,
    DocumentChartBarIcon,
    CogIcon,
    ChevronDownIcon,
    ChevronRightIcon,
} from '@heroicons/react/24/outline';
import { useState, useEffect } from 'react';

interface NavItem {
    name: string;
    path?: string;
    icon: any;
    children?: NavItem[];
}

const navigationItems: NavItem[] = [
    {
        name: 'Dashboard',
        path: '/dashboard',
        icon: ChartBarIcon,
    },
    {
        name: 'Inventory',
        icon: CubeIcon,
        children: [
            { name: 'All Products', path: '/inventory/products', icon: CubeIcon },
            { name: 'Low Stock', path: '/inventory/low-stock', icon: CubeIcon },
            { name: 'Categories', path: '/inventory/categories', icon: CubeIcon },
        ],
    },
    {
        name: 'Transactions',
        icon: ArrowsRightLeftIcon,
        children: [
            { name: 'Sales', path: '/transactions/sales', icon: ArrowsRightLeftIcon },
            { name: 'Purchases', path: '/transactions/purchases', icon: ArrowsRightLeftIcon },
            { name: 'Adjustments', path: '/transactions/adjustments', icon: ArrowsRightLeftIcon },
            { name: 'Transfers', path: '/transactions/transfers', icon: ArrowsRightLeftIcon },
            { name: 'History', path: '/transactions/history', icon: ArrowsRightLeftIcon },
        ],
    },
    {
        name: 'Parties',
        icon: UsersIcon,
        children: [
            { name: 'Suppliers', path: '/parties/suppliers', icon: UsersIcon },
            { name: 'Customers', path: '/parties/customers', icon: UsersIcon },
        ],
    },
    {
        name: 'Reports',
        icon: DocumentChartBarIcon,
        children: [
            { name: 'Stock Summary', path: '/reports/stock-summary', icon: DocumentChartBarIcon },
            { name: 'Inventory by Location', path: '/reports/inventory-by-location', icon: DocumentChartBarIcon },
            { name: 'Sales Summary', path: '/reports/sales-summary', icon: DocumentChartBarIcon },
            { name: 'Purchase Summary', path: '/reports/purchase-summary', icon: DocumentChartBarIcon },
            { name: 'Low Stock Alerts', path: '/reports/low-stock-alerts', icon: DocumentChartBarIcon },
        ],
    },
    {
        name: 'Settings',
        path: '/settings',
        icon: CogIcon,
    },
];

interface SidebarProps {
    isOpen: boolean;
    onClose: () => void;
}

export const Sidebar = ({ isOpen, onClose }: SidebarProps) => {
    const location = useLocation();
    const [openMenu, setOpenMenu] = useState<string | null>(null);

    // Automatically open the menu that contains the current route
    useEffect(() => {
        const currentPath = location.pathname;
        const activeParent = navigationItems.find((item) => {
            if (!item.children) return false;
            return item.children.some((child) => child.path === currentPath);
        });

        if (activeParent) {
            setOpenMenu(activeParent.name);
        }
    }, [location.pathname]);

    const toggleMenu = (menuName: string) => {
        // If clicking the currently open menu, close it
        // Otherwise, open the clicked menu and close all others (accordion behavior)
        setOpenMenu((prev) => (prev === menuName ? null : menuName));
    };

    const isActive = (path: string) => {
        return location.pathname === path;
    };

    const isParentActive = (children?: NavItem[]) => {
        if (!children) return false;
        return children.some((child) => child.path && location.pathname === child.path);
    };

    return (
        <>
            {/* Mobile overlay */}
            {isOpen && (
                <div
                    className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
                    onClick={onClose}
                />
            )}

            {/* Sidebar */}
            <aside
                className={`fixed top-0 left-0 z-50 h-screen w-64 bg-black text-white transition-transform duration-300 ease-in-out lg:translate-x-0 ${
                    isOpen ? 'translate-x-0' : '-translate-x-full'
                }`}
            >
                <div className="flex h-full flex-col">
                    {/* Logo/Brand */}
                    <div className="flex h-16 items-center justify-between px-6 border-b border-gray-800">
                        <Link to="/dashboard" className="text-xl font-bold">
                            Inventory MS
                        </Link>
                        <button
                            onClick={onClose}
                            className="lg:hidden text-gray-400 hover:text-white"
                        >
                            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>

                    {/* Navigation */}
                    <nav className="flex-1 overflow-y-auto py-4">
                        <ul className="space-y-1 px-3">
                            {navigationItems.map((item) => (
                                <li key={item.name}>
                                    {item.children ? (
                                        // Parent menu item with children
                                        <div>
                                            <button
                                                onClick={() => toggleMenu(item.name)}
                                                className={`flex w-full items-center justify-between rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                                                    isParentActive(item.children)
                                                        ? 'bg-indigo-600 text-white'
                                                        : 'text-gray-300 hover:bg-[#434343] hover:text-white'
                                                }`}
                                            >
                                                <div className="flex items-center">
                                                    <item.icon className="h-5 w-5 mr-3" />
                                                    {item.name}
                                                </div>
                                                {openMenu === item.name ? (
                                                    <ChevronDownIcon className="h-4 w-4" />
                                                ) : (
                                                    <ChevronRightIcon className="h-4 w-4" />
                                                )}
                                            </button>

                                            {/* Submenu with smooth animation */}
                                            <div
                                                className={`overflow-hidden transition-all duration-300 ease-in-out ${
                                                    openMenu === item.name ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'
                                                }`}
                                            >
                                                <ul className="mt-1 ml-4 space-y-1">
                                                    {item.children.map((child) => (
                                                        <li key={child.name}>
                                                            <Link
                                                                to={child.path!}
                                                                className={`flex items-center rounded-md px-3 py-2 text-sm transition-colors ${
                                                                    isActive(child.path!)
                                                                        ? 'bg-indigo-600 text-white'
                                                                        : 'text-gray-400 hover:bg-[#434343] hover:text-white'
                                                                }`}
                                                                onClick={onClose}
                                                            >
                                                                {child.name}
                                                            </Link>
                                                        </li>
                                                    ))}
                                                </ul>
                                            </div>
                                        </div>
                                    ) : (
                                        // Single menu item without children
                                        <Link
                                            to={item.path!}
                                            className={`flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                                                isActive(item.path!)
                                                    ? 'bg-indigo-600 text-white'
                                                    : 'text-gray-300 hover:bg-[#434343] hover:text-white'
                                            }`}
                                            onClick={onClose}
                                        >
                                            <item.icon className="h-5 w-5 mr-3" />
                                            {item.name}
                                        </Link>
                                    )}
                                </li>
                            ))}
                        </ul>
                    </nav>

                    {/* User info at bottom */}
                    <div className="border-t border-gray-800 p-4">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <div className="h-8 w-8 rounded-full bg-indigo-600 flex items-center justify-center">
                                    <span className="text-sm font-medium">U</span>
                                </div>
                            </div>
                            <div className="ml-3">
                                <p className="text-sm font-medium">User</p>
                                <p className="text-xs text-gray-400">View profile</p>
                            </div>
                        </div>
                    </div>
                </div>
            </aside>
        </>
    );
};
