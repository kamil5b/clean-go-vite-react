const API_BASE_URL = "http://localhost:8080/api";

type Item = {
    id: string;
    title: string;
    description: string;
    user_id: string;
    created_at: string;
    updated_at: string;
};

export const itemsApi = {
    createItem: async (title: string, description: string): Promise<Item> => {
        const response = await fetch(`${API_BASE_URL}/items`, {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, description }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to create item");
        }

        return response.json();
    },

    getItems: async (): Promise<Item[]> => {
        const response = await fetch(`${API_BASE_URL}/items`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch items");
        }

        const data = await response.json();
        return data || [];
    },

    getItem: async (id: string): Promise<Item> => {
        const response = await fetch(`${API_BASE_URL}/items/${id}`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch item");
        }

        return response.json();
    },

    updateItem: async (
        id: string,
        title: string,
        description: string,
    ): Promise<Item> => {
        const response = await fetch(`${API_BASE_URL}/items/${id}`, {
            method: "PUT",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, description }),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to update item");
        }

        return response.json();
    },

    deleteItem: async (id: string): Promise<void> => {
        const response = await fetch(`${API_BASE_URL}/items/${id}`, {
            method: "DELETE",
            credentials: "include",
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to delete item");
        }
    },
};

export type { Item };
