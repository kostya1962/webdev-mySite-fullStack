import type { Product } from "./product.interface";
import type { User } from "./user.interface";

export interface Order {
  id: number;
  user_id: number;
  product_ids: number[];
  products: Product[];
  status: string;
  created_at: string;
  user?: User;
  price: number;
}

export interface GetOrdersResponse {
  orders: Order[];
  user: User;
}
