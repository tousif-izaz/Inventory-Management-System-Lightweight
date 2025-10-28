export interface User {
    user_id: number;
    username: string;
    name: string;
    email?: string;
    role: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
    last_login?: string;
}

export interface UpdatePasswordRequest {
    old_password: string;
    new_password: string;
}
