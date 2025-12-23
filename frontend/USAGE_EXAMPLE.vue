<!-- Пример использования компонента AddToCart на странице товара -->

<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface'

// Получить данные товара (например, на странице /catalog/sup-[id].vue)
const product = ref<Product>({
  id: 1,
  name: "Пример товара",
  price: 1999,
  short_description: "Короткое описание",
  long_description: "Полное описание товара",
  sku: "SKU123",
  discount: 0,
  images: ["image1.jpg"],
  category_id: 1,
  category: { id: 1, name: "Категория" },
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
});
</script>

<template>
  <div class="product-page">
    <div class="product-info">
      <h1>{{ product.name }}</h1>
      <p>{{ product.long_description }}</p>
      <div class="price">{{ product.price }} ₽</div>
      
      <!-- Компонент добавления в корзину -->
      <AddToCart :product="product" />
    </div>
  </div>
</template>

<style scoped>
.product-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.price {
  font-size: 28px;
  font-weight: 700;
  margin: 20px 0;
  color: var(--color-black);
}
</style>
