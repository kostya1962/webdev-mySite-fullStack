import type { Product } from "~/interfaces/product.interface";

export interface Banner {
  id: number;
  prodyct_id: number;
  product: Product;
  image: string;
  position: number;
}

export interface GetBannersResponse {
  banners: Banner[];
}
