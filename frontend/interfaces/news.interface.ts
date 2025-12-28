export interface NewsItem {
  id: number;
  title: string;
  description: string;
  created_at: string;
  image: string;
}
export interface GetNewsResponse {
  length: number;
  news: NewsItem[];
}
