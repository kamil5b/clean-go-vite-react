import { apiClient, apiClientJson } from "@/lib/apiClient";
import { CreateInvoiceRequest, UpdateInvoiceRequest } from "@/types/request/invoice";
import { InvoiceDetailResponse, InvoicePaginationResponse } from "@/types/response/invoice";

export const invoiceApi = {
  create: async (data: CreateInvoiceRequest): Promise<InvoiceDetailResponse> => {
    return apiClientJson<InvoiceDetailResponse>("/invoices", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  getById: async (id: number): Promise<InvoiceDetailResponse> => {
    return apiClientJson<InvoiceDetailResponse>(`/invoices/${id}`);
  },

  update: async (id: number, data: UpdateInvoiceRequest): Promise<InvoiceDetailResponse> => {
    return apiClientJson<InvoiceDetailResponse>(`/invoices/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  delete: async (id: number): Promise<void> => {
    await apiClient(`/invoices/${id}`, {
      method: "DELETE",
    });
  },

  getAll: async (
    page: number = 1,
    limit: number = 10,
    search: string = ""
  ): Promise<InvoicePaginationResponse> => {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    if (search) {
      params.append("search", search);
    }
    return apiClientJson<InvoicePaginationResponse>(`/invoices?${params.toString()}`);
  },
};
