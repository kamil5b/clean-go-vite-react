import { apiClient, apiClientJson } from "@/lib/apiClient";
import { CreateTagRequest, UpdateTagRequest } from "@/types/request/tag";
import { TagResponse, TagPaginationResponse } from "@/types/response/tag";

export const tagApi = {
  create: async (data: CreateTagRequest): Promise<TagResponse> => {
    return apiClientJson<TagResponse>("/tags", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  getById: async (id: number): Promise<TagResponse> => {
    return apiClientJson<TagResponse>(`/tags/${id}`);
  },

  update: async (id: number, data: UpdateTagRequest): Promise<TagResponse> => {
    return apiClientJson<TagResponse>(`/tags/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  delete: async (id: number): Promise<void> => {
    await apiClient(`/tags/${id}`, {
      method: "DELETE",
    });
  },

  getAll: async (
    page: number = 1,
    limit: number = 10,
    search: string = ""
  ): Promise<TagPaginationResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    if (search) {
      params.append("search", search);
    }
    return apiClientJson<TagPaginationResponse>(`/tags?${params.toString()}`);
  },
};
