import {
    LoginResponse,
    RegisterResponse,
    RefreshResponse,
    CSRFTokenResponse,
    GetUser,
} from "@/types/response/user";

const API_BASE_URL = "http://localhost:8080/api";

export const authApi = {
    register: async (
        email: string,
        password: string,
        name: string,
    ): Promise<RegisterResponse> => {
        const response = await fetch(`${API_BASE_URL}/auth/register`, {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password, name }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Registration failed");
        }

        return response.json();
    },

    login: async (email: string, password: string): Promise<LoginResponse> => {
        const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Login failed");
        }

        return response.json();
    },

    logout: async (): Promise<void> => {
        await fetch(`${API_BASE_URL}/auth/logout`, {
            method: "POST",
            credentials: "include",
        });
    },

    getCurrentUser: async (): Promise<GetUser | null> => {
        const response = await fetch(`${API_BASE_URL}/auth/me`, {
            credentials: "include",
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
        const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
            method: "POST",
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Token refresh failed");
        }

        return response.json();
    },

    getCsrfToken: async (): Promise<CSRFTokenResponse> => {
        const response = await fetch(`${API_BASE_URL}/csrf`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch CSRF token");
        }

        return response.json();
    },
};
