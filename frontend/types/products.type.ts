export interface Product {
  id: number;
  reference: string;
  title: string;
  price: number;
  image_url: string;
  product_url: string;
  category: string;
  event_type: string;
  event_date: string;
  in_stock: boolean;
}