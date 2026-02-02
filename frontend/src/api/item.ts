import { apiClient, apiClientJson } from "@/lib/apiClient";
import { CreateItemRequest, UpdateItemRequest } from "@/types/request/item";
import { ItemResponse, ItemPaginationResponse } from "@/types/response/item";

export const itemApi = {
  create: async (data: CreateItemRequest): Promise<ItemResponse> => {
    return apiClientJson<ItemResponse>("/items", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  getById: async (id: number): Promise<ItemResponse> => {
    return apiClientJson<ItemResponse>(`/items/${id}`);
  },

  update: async (id: number, data: UpdateItemRequest): Promise<ItemResponse> => {
    return apiClientJson<ItemResponse>(`/items/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  delete: async (id: number): Promise<void> => {
    await apiClient(`/items/${id}`, {
      method: "DELETE",
    });
  },

  getAll: async (
    page: number = 1,
    limit: number = 10,
    search: string = ""
  ): Promise<ItemPaginationResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    if (search) {
      params.append("search", search);
    }
    return apiClientJson<ItemPaginationResponse>(`/items?${params.toString()}`);
  },
};
