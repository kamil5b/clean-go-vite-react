export interface InvoiceItemInput {
  item_id: number;
  quantity: number;
  unit_price: number;
}

export interface CreateInvoiceRequest {
  grand_price: number;
  items: InvoiceItemInput[];
  tags: number[];
}

export interface UpdateInvoiceRequest {
  grand_price: number;
  items: InvoiceItemInput[];
  tags: number[];
}
