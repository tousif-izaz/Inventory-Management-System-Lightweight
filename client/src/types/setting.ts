export interface Setting {
    setting_id: number;
    setting_key: string;
    setting_value: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateSettingRequest {
    setting_key: string;
    setting_value: string;
    description?: string;
}

export interface UpdateSettingRequest {
    setting_value: string;
}
