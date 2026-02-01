import { GetCounter } from "@/types/response/counter";
import { apiClientJson } from "@/lib/apiClient";

export const counterApi = {
    getCounter: async (): Promise<GetCounter> => {
        return apiClientJson<GetCounter>("/counter", {
            method: "GET",
        });
    },

    incrementCounter: async (): Promise<GetCounter> => {
        // POST request - will automatically include CSRF token via apiClient
        return apiClientJson<GetCounter>("/counter", {
            method: "POST",
        });
    },
};
