import { ItemResponse } from "./item";
import { TagResponse } from "./tag";

export interface InvoiceItemResponse {
  id: number;
  item_id: number;
  item: ItemResponse;
  quantity: number;
  unit_price: number;
  total_price: number;
}

export interface InvoiceDetailResponse {
  id: number;
  grand_price: number;
  items: InvoiceItemResponse[];
  tags: TagResponse[];
  created_at: string;
  updated_at: string;
}

export interface InvoiceListItem {
  id: number;
  grand_price: number;
  tags: TagResponse[];
  totalItem: number;
  created_at: string;
  updated_at: string;
}

export interface InvoicePaginationMeta {
  totalData: number;
  page: number;
  limit: number;
  totalPage: number;
}

export interface InvoicePaginationResponse {
  data: InvoiceListItem[];
  meta: InvoicePaginationMeta;
}
