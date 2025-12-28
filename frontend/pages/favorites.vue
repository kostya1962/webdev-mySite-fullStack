<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface';
import { useFavoriteStore } from '~/state/favorite.state';

useSeoMeta({
    title: 'Избранное',
    description: 'Избранные товары интернет магазина Shopper',
    ogDescription: 'Избранные товары интернет магазина Shopper',
});

const favoriteState = useFavoriteStore();
const API_URL = useAPI();
const products = ref<Product[]>([]);


watchEffect(async () => {
    if (!Array.isArray(favoriteState.favoriteIDs)) {
        favoriteState.favoriteIDs = [];
    }
    
    const results = await Promise.allSettled(
        favoriteState.favoriteIDs.map(id => {
            return $fetch<{product: Product}>(API_URL + '/products/' + id);
        })
    );
    
    const notFoundIds: number[] = [];
    const loadedProducts = results
        .map((result, index) => {
            if (result.status === 'fulfilled') {
                return result.value.product;
            } else {
                // Товар не найден (404), запомнить ID для удаления
                const id = favoriteState.favoriteIDs[index];
                if (id !== undefined) {
                    notFoundIds.push(id);
                }
                return null;
            }
        })
        .filter((product): product is Product => product !== null);
    
    products.value = loadedProducts;
    
    // Удалить несуществующие товары из состояния
    if (notFoundIds.length > 0) {
        notFoundIds.forEach(id => {
            favoriteState.removeFavoriteId(id);
        });
    }
});
</script>


<template>
    <div>
        <h1>Избранное</h1>

        <div v-show="!products.length" class="empty">Вы не выбрали ни одного товара</div>

        <div class="catalog__list">
            <CatalogCard 
                        v-for="product in products"
                        :key="product.id"
                        v-bind="product"/>
        </div>
    </div>
</template>

<style scoped>
    .empty{
        padding:40px;
        text-align:center;
        color:var(--color-gray);
    }

    .catalog__list{
        display: grid;
        width: 100%;
        grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
        gap: 64px 12px;
    }
</style>