export interface Review {
  id: number;
  product_id: number;
  name: string;
  text: string;
  rating: number;
  created_at: string;
}

export interface ListReviewResponse {
  reviews: Review[];
}
