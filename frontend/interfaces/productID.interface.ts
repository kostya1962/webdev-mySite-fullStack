import type { Category } from "./category.interface";
import type { Review } from "./review.interface";

export interface Product {
  id: number;
  name: string;
  price: number;
  short_description: string;
  long_description: string;
  sku: string;
  discount: number;
  images: string[];
  category_id: number;
  category: Category;
  created_at: string;
  updated_at: string;
}

export interface ProductIDRsponse {
  product: Product;
  reviews: Review[];
}
