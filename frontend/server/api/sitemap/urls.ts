import { useAPI } from "~/composables/useAPI";
import type { GetProductsResponse } from "~/interfaces/product.interface";
import type { SitemapUrlInput } from "@nuxtjs/sitemap";

export default defineSitemapEventHandler(async () => {
  const API_URL = useAPI();
  const products = await $fetch<GetProductsResponse>(API_URL + "/products", {
    query: {
      limit: 1000,
      offset: 0,
    },
  });
  const pages = products.products.map((p) => ({
    loc: `/catalog/${p.id}`,
    changefreq: "daily",
    priority: 0.7,
  })) satisfies SitemapUrlInput[];
  return pages;
});
