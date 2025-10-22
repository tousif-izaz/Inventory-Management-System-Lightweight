import { api } from './api';
import { Setting, CreateSettingRequest, UpdateSettingRequest } from '../types/setting';

class SettingService {
    /**
     * Get all settings
     */
    async getAllSettings(): Promise<Setting[]> {
        return api.get<Setting[]>('/settings');
    }

    /**
     * Get a single setting by key
     */
    async getSettingByKey(key: string): Promise<Setting> {
        return api.get<Setting>(`/settings/${key}`);
    }

    /**
     * Get setting value by key (convenience method)
     */
    async getSettingValue(key: string): Promise<string> {
        const setting = await this.getSettingByKey(key);
        return setting.setting_value;
    }

    /**
     * Create a new setting
     */
    async createSetting(settingData: CreateSettingRequest): Promise<void> {
        return api.post<void>('/settings', settingData);
    }

    /**
     * Update a setting value
     */
    async updateSetting(key: string, value: string): Promise<void> {
        const updateData: UpdateSettingRequest = {
            setting_value: value,
        };
        return api.put<void>(`/settings/${key}`, updateData);
    }

    /**
     * Delete a setting
     */
    async deleteSetting(key: string): Promise<void> {
        return api.delete<void>(`/settings/${key}`);
    }
}

export const settingService = new SettingService();
