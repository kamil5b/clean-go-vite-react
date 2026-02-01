import { GetMessage } from "@/types/response/message";

const API_BASE_URL = "http://localhost:8080/api";

export const messageApi = {
    getMessage: async (): Promise<GetMessage | null> => {
        const response = await fetch(`${API_BASE_URL}/message`);

        if (!response.ok) {
            if (response.status === 401) {
                return null;
            }
            throw new Error("Failed to fetch current user");
        }

        return response.json();
    },
};
