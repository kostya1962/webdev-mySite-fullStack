import type { Product } from "./product.interface";

export interface CartItem {
  product: Product;
  quantity: number;
}

export interface Cart {
  items: CartItem[];
  totalPrice: number;
}

export interface AddToCartRequest {
  email: string;
  productID: number;
  quantity: number;
}

export interface CartResponse {
  success: boolean;
  message?: string;
}
