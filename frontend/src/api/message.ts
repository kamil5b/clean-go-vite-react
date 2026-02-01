import { GetMessage } from "@/types/response/message";
import { apiClient } from "@/lib/apiClient";

export const messageApi = {
    getMessage: async (): Promise<GetMessage | null> => {
        const response = await apiClient("/message", {
            skipAuthRefresh: true, // Public endpoint - don't trigger auth refresh
        });

        if (!response.ok) {
            if (response.status === 401) {
                return null;
            }
            throw new Error("Failed to fetch message");
        }

        return response.json();
    },
};
