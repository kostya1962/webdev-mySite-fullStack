<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface';

    const config = useRuntimeConfig();
    const product = defineProps<Product>();
    const image = computed(() => `url(${config.public.imageurl}${product.images?.[0] ?? ''})`);

    const formattedPrice = computed(() => {
        const value = Number(product.price || 0);
        return new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB' }).format(value);
    });

    const isHovered = ref(false);
</script>

<template>
    <NuxtLink 
        class="card" 
        :to="`/catalog/sup-${product.id}`" 
        @mouseenter="isHovered=true"
        @mouseleave="isHovered=false"
        >
        <div class="card__image">
            <span v-if="product.discount > 0" class="card__discount">
                - {{ product.discount }}%
            </span>
            <span v-else></span>
            <AddFavorite :id="product.id" :is-shown="isHovered" />
        </div>
        <div class="card__info">
            <div class="card__name">
                {{ product.name }}
            </div>
            <div class="card__price">
                {{ formattedPrice }}
            </div>
        </div>
    </NuxtLink>
</template>

<style scoped>
    .card{
        display: flex;
        flex-direction: column;
        gap: 24px;
        width: 100%;
        text-decoration: none;
    }

    .card__image{
        aspect-ratio: 1/1;
        border-radius: 8px;
        min-width: 320px;
        width: 100%;
        background-size: cover;
        background-position: center;
        background-repeat: no-repeat;
        background-color: lightgray;
        padding: 16px;
        background-image: v-bind(image);
        display: flex;
        justify-content: space-between;
        align-items: start;
    }

    .card__discount{
        background: var(--color-accent);
        border-radius: 4px;
        padding: 2px 8px;
        font-size: 12px;
        color: var(--color-white);
    }

    .card__name{
        font-size: 20px;
        line-height: 26px;
        color: var(--color-black);
    }

    .card__price {
        font-size: 20px;
        font-weight: 500;
        line-height: 26px;
        text-transform: capitalize;
        color: var(--color-accent);
    }
    
    .card__info{
        display:flex;
        gap:16px;
        flex-direction: column;
    }
</style>