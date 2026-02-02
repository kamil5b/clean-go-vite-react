export interface ItemResponse {
  id: number;
  name: string;
  desc: string;
  created_at: string;
  updated_at: string;
}

export interface ItemPaginationMeta {
  totalData: number;
  page: number;
  limit: number;
  totalPage: number;
}

export interface ItemPaginationResponse {
  data: ItemResponse[];
  meta: ItemPaginationMeta;
}
