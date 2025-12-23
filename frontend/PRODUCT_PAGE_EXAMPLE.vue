<!-- Пример интеграции AddToCart компонента на странице товара -->
<!-- Файл: frontend/pages/catalog/sup-[id].vue -->

<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface'

const route = useRoute()
const productId = Array.isArray(route.params.id) 
  ? route.params.id[0].replace('sup-', '') 
  : route.params.id.replace('sup-', '')

const { data: product, pending: isLoading } = await useFetch<Product>(
  `/api/products/${productId}`,
  { key: `product-${productId}` }
)
</script>

<template>
  <div v-if="isLoading" class="loading">
    Загрузка товара...
  </div>

  <div v-else-if="product" class="product-container">
    <div class="product-gallery">
      <!-- Галерея изображений товара -->
      <GallerayProd :images="product.images" />
    </div>

    <div class="product-info">
      <h1 class="product-title">{{ product.name }}</h1>

      <div class="product-rating">
        <!-- Рейтинг товара -->
        <RatingStars :rating="4.5" :reviews-count="128" />
      </div>

      <div class="product-description">
        {{ product.short_description }}
      </div>

      <div class="product-meta">
        <div class="sku">SKU: {{ product.sku }}</div>
        <div v-if="product.discount > 0" class="discount">
          Скидка: {{ product.discount }}%
        </div>
      </div>

      <div class="product-price-section">
        <div class="price-wrapper">
          <span class="price">{{ product.price }} ₽</span>
          <span v-if="product.discount > 0" class="old-price">
            {{ Math.round(product.price / (1 - product.discount / 100)) }} ₽
          </span>
        </div>

        <!-- ЗДЕСЬ ИСПОЛЬЗУЕМ КОМПОНЕНТ AddToCart -->
        <AddToCart :product="product" />
      </div>

      <div class="product-full-description">
        <h2>Описание товара</h2>
        <p>{{ product.long_description }}</p>
      </div>
    </div>
  </div>

  <div v-else class="error">
    Товар не найден
  </div>
</template>

<style scoped>
.product-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 60px;
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.product-gallery {
  position: sticky;
  top: 20px;
  height: fit-content;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.product-title {
  font-size: 36px;
  font-weight: 700;
  margin: 0;
  line-height: 1.2;
  color: var(--color-black);
}

.product-rating {
  display: flex;
  align-items: center;
  gap: 12px;
}

.product-description {
  font-size: 16px;
  color: var(--color-gray);
  line-height: 1.6;
}

.product-meta {
  display: flex;
  gap: 24px;
  padding: 20px 0;
  border-top: 1px solid var(--color-gray);
  border-bottom: 1px solid var(--color-gray);
}

.sku,
.discount {
  font-size: 14px;
  color: var(--color-gray);
}

.product-price-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 20px 0;
}

.price-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
}

.price {
  font-size: 40px;
  font-weight: 700;
  color: var(--color-black);
}

.old-price {
  font-size: 18px;
  color: var(--color-gray);
  text-decoration: line-through;
}

.product-full-description {
  margin-top: 40px;
  padding-top: 40px;
  border-top: 1px solid var(--color-gray);
}

.product-full-description h2 {
  font-size: 24px;
  margin-bottom: 16px;
}

.product-full-description p {
  font-size: 16px;
  line-height: 1.8;
  color: var(--color-gray);
}

.loading,
.error {
  text-align: center;
  padding: 100px 20px;
  font-size: 18px;
  color: var(--color-gray);
}

@media (max-width: 768px) {
  .product-container {
    grid-template-columns: 1fr;
    gap: 30px;
  }

  .product-title {
    font-size: 24px;
  }

  .price {
    font-size: 28px;
  }
}
</style>
