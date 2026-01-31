import { GetCounter } from "@/types/response/counter";

const API_BASE_URL = "http://localhost:8080/api";

export const counterApi = {
    getCounter: async (): Promise<GetCounter> => {
        const response = await fetch(`${API_BASE_URL}/counter`, {
            method: "GET",
        });
        if (!response.ok) {
            throw new Error(`Failed to fetch counter: ${response.statusText}`);
        }
        return response.json();
    },

    incrementCounter: async (): Promise<GetCounter> => {
        const response = await fetch(`${API_BASE_URL}/counter`, {
            method: "POST",
        });
        if (!response.ok) {
            throw new Error(
                `Failed to increment counter: ${response.statusText}`,
            );
        }
        return response.json();
    },
};
