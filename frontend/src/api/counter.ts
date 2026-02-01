const API_BASE_URL = "http://localhost:8080/api";

export const counterApi = {
    getCounter: async (): Promise<{ value: number }> => {
        const response = await fetch(`${API_BASE_URL}/counter`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch counter");
        }

        return response.json();
    },

    incrementCounter: async (): Promise<{ value: number }> => {
        const response = await fetch(`${API_BASE_URL}/counter`, {
            method: "POST",
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to increment counter");
        }

        return response.json();
    },
};
