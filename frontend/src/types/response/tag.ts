export interface TagResponse {
  id: number;
  name: string;
  color_hex: string;
  created_at: string;
  updated_at: string;
}

export interface TagPaginationMeta {
  totalData: number;
  page: number;
  limit: number;
  totalPage: number;
}

export interface TagPaginationResponse {
  data: TagResponse[];
  meta: TagPaginationMeta;
}
