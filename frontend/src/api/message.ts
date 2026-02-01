const API_BASE_URL = "http://localhost:8080/api";

export const messageApi = {
    getMessage: async (): Promise<{ content: string }> => {
        const response = await fetch(`${API_BASE_URL}/message`, {
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Failed to fetch message");
        }

        return response.json();
    },
};
