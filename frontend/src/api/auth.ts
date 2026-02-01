const API_BASE_URL = "http://localhost:8080/api";

// Simple inline types - no over-engineering
type User = {
    id: string;
    email: string;
    name: string;
};

type AuthResponse = {
    user: User;
    token: string;
};

export const authApi = {
    register: async (
        email: string,
        password: string,
        name: string,
    ): Promise<AuthResponse> => {
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

    login: async (email: string, password: string): Promise<AuthResponse> => {
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

    getCurrentUser: async (): Promise<User | null> => {
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
};

// Export types for use in other files
export type { User, AuthResponse };
