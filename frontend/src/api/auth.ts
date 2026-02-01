import {
    LoginResponse,
    RegisterResponse,
    RefreshResponse,
    CSRFTokenResponse,
    GetUser,
} from "@/types/response/user";
import { apiClient, apiClientJson } from "@/lib/apiClient";

export const authApi = {
    register: async (
        email: string,
        password: string,
        name: string,
    ): Promise<RegisterResponse> => {
        return apiClientJson<RegisterResponse>("/auth/register", {
            method: "POST",
            body: JSON.stringify({ email, password, name }),
        });
    },

    login: async (email: string, password: string): Promise<LoginResponse> => {
        return apiClientJson<LoginResponse>("/auth/login", {
            method: "POST",
            body: JSON.stringify({ email, password }),
        });
    },

    logout: async (): Promise<void> => {
        // Logout is a POST request and requires CSRF protection
        await apiClient("/auth/logout", {
            method: "POST",
        });
    },

    getCurrentUser: async (): Promise<GetUser | null> => {
        const response = await apiClient("/auth/me", {
            skipAuthRefresh: true, // Don't trigger refresh on initial auth check
        });

        if (!response.ok) {
            if (response.status === 401) {
                return null;
            }
            throw new Error("Failed to fetch current user");
        }

        return response.json();
    },

    refreshToken: async (): Promise<RefreshResponse> => {
        // Note: This is called internally by apiClient, but exposed for manual use
        return apiClientJson<RefreshResponse>("/auth/refresh", {
            method: "POST",
        });
    },

    getCsrfToken: async (): Promise<CSRFTokenResponse> => {
        return apiClientJson<CSRFTokenResponse>("/csrf");
    },
};
