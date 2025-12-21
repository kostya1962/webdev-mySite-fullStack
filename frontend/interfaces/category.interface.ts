export interface Category {
  id: number;
  name: string;
  alias: string;
}

export interface GetCategoryResponse {
  categories: Category[];
}
