import { useState, useEffect } from 'react';
import { Button } from '../UI/Button';
import { Input } from '../UI/Input';
import { useToast } from '../UI/Toast';
import { settingService } from '../../services/settingService';
import { Setting } from '../../types/setting';
import { PencilIcon, CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';

export const Settings = () => {
    const { showToast } = useToast();
    const [settings, setSettings] = useState<Setting[]>([]);
    const [loading, setLoading] = useState(true);
    const [editingKey, setEditingKey] = useState<string | null>(null);
    const [editValues, setEditValues] = useState<{ [key: string]: string }>({});

    useEffect(() => {
        loadSettings();
    }, []);

    const loadSettings = async () => {
        try {
            setLoading(true);
            const data = await settingService.getAllSettings();
            // Filter out Square-related settings
            const filteredSettings = data.filter(
                (setting) => !setting.setting_key.startsWith('square_')
            );
            setSettings(filteredSettings);
        } catch (error) {
            console.error('Failed to load settings:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to load settings';
            showToast('error', `Failed to load settings: ${errorMessage}`);
        } finally {
            setLoading(false);
        }
    };

    const handleEdit = (setting: Setting) => {
        setEditingKey(setting.setting_key);
        setEditValues({ ...editValues, [setting.setting_key]: setting.setting_value });
    };

    const handleCancel = () => {
        setEditingKey(null);
        setEditValues({});
    };

    const handleSave = async (key: string) => {
        try {
            const newValue = editValues[key];
            if (newValue === undefined || newValue.trim() === '') {
                showToast('error', 'Setting value cannot be empty');
                return;
            }

            await settingService.updateSetting(key, newValue);
            showToast('success', 'Setting updated successfully');
            setEditingKey(null);
            loadSettings();
        } catch (error) {
            console.error('Failed to update setting:', error);
            const errorMessage = error instanceof Error ? error.message : 'Failed to update setting';
            showToast('error', `Failed to update setting: ${errorMessage}`);
        }
    };

    const handleInputChange = (key: string, value: string) => {
        setEditValues({ ...editValues, [key]: value });
    };

    const getSettingLabel = (key: string): string => {
        const labels: { [key: string]: string } = {
            store_name: 'Store Name',
            store_address: 'Store Address',
            store_phone: 'Store Phone',
            store_policy: 'Store Policy',
            sales_tax_rate: 'Sales Tax Rate (%)',
        };
        return labels[key] || key;
    };

    const getSettingDescription = (key: string): string => {
        const descriptions: { [key: string]: string } = {
            store_name: 'The name of your store displayed on receipts and documents',
            store_address: 'Physical address of your store',
            store_phone: 'Contact phone number for your store',
            store_policy: 'Return and exchange policy displayed on receipts',
            sales_tax_rate: 'Default sales tax rate percentage applied to sales',
        };
        return descriptions[key] || '';
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center h-64">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Settings</h1>
                <p className="mt-1 text-sm text-gray-500">
                    Manage system settings and business preferences
                </p>
            </div>

            <div className="bg-white shadow rounded-lg">
                <div className="px-6 py-4 border-b border-gray-200">
                    <h2 className="text-lg font-medium text-gray-900">Store Settings</h2>
                    <p className="mt-1 text-sm text-gray-500">
                        Configure your store information and defaults
                    </p>
                </div>

                <div className="divide-y divide-gray-200">
                    {settings.map((setting) => (
                        <div key={setting.setting_key} className="px-6 py-4">
                            <div className="flex items-start justify-between">
                                <div className="flex-1 mr-4">
                                    <label className="block text-sm font-medium text-gray-700 mb-1">
                                        {getSettingLabel(setting.setting_key)}
                                    </label>
                                    <p className="text-xs text-gray-500 mb-2">
                                        {setting.description || getSettingDescription(setting.setting_key)}
                                    </p>

                                    {editingKey === setting.setting_key ? (
                                        <div className="mt-2">
                                            {setting.setting_key === 'store_policy' ? (
                                                <textarea
                                                    value={editValues[setting.setting_key] || ''}
                                                    onChange={(e) =>
                                                        handleInputChange(setting.setting_key, e.target.value)
                                                    }
                                                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                                    rows={3}
                                                />
                                            ) : (
                                                <Input
                                                    type={
                                                        setting.setting_key === 'sales_tax_rate'
                                                            ? 'number'
                                                            : 'text'
                                                    }
                                                    value={editValues[setting.setting_key] || ''}
                                                    onChange={(e) =>
                                                        handleInputChange(setting.setting_key, e.target.value)
                                                    }
                                                    step={
                                                        setting.setting_key === 'sales_tax_rate'
                                                            ? '0.01'
                                                            : undefined
                                                    }
                                                />
                                            )}
                                        </div>
                                    ) : (
                                        <div className="mt-2 text-sm text-gray-900">
                                            {setting.setting_key === 'store_policy' ? (
                                                <p className="whitespace-pre-wrap">{setting.setting_value}</p>
                                            ) : setting.setting_key === 'sales_tax_rate' ? (
                                                <p>{setting.setting_value}%</p>
                                            ) : (
                                                <p>{setting.setting_value}</p>
                                            )}
                                        </div>
                                    )}
                                </div>

                                <div className="flex items-center space-x-2">
                                    {editingKey === setting.setting_key ? (
                                        <>
                                            <Button
                                                variant="primary"
                                                size="sm"
                                                onClick={() => handleSave(setting.setting_key)}
                                            >
                                                <CheckIcon className="h-4 w-4 mr-1" />
                                                Save
                                            </Button>
                                            <Button
                                                variant="secondary"
                                                size="sm"
                                                onClick={handleCancel}
                                            >
                                                <XMarkIcon className="h-4 w-4 mr-1" />
                                                Cancel
                                            </Button>
                                        </>
                                    ) : (
                                        <Button
                                            variant="secondary"
                                            size="sm"
                                            onClick={() => handleEdit(setting)}
                                        >
                                            <PencilIcon className="h-4 w-4 mr-1" />
                                            Edit
                                        </Button>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                <div className="px-6 py-4 bg-gray-50 border-t border-gray-200">
                    <p className="text-xs text-gray-500">
                        Last updated: {new Date().toLocaleDateString()}
                    </p>
                </div>
            </div>
        </div>
    );
};
