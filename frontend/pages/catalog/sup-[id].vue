<script setup lang="ts">
import GallerayProd from '~/components/GallerayProd.vue';
import type { ProductIDRsponse } from '~/interfaces/productID.interface';




const route = useRoute(); //извлекаю текущий маршрут (путь и параметры)
const API_URL = useAPI();


const {data: productData } = await useFetch<ProductIDRsponse>(
    API_URL + '/products/' + route.params.id
);

useSeoMeta({
    title: `Купить ${productData.value?.product.name}`,
    description: productData.value?.product.short_description,
});

const formattedPrice = computed(() => {
        const value = Number( productData.value?.product.price || 0);
        return new Intl.NumberFormat('ru-RU', { style: 'currency', currency: 'RUB' }).format(value);
    });

    const averageRating = computed(() => {
    return productData.value?.reviews.length
        ? productData.value.reviews.reduce((sum, r) => sum + r.rating, 0) / productData.value.reviews.length
        : 0
    })

    console.log(averageRating.value);
</script>

<template>
    <div class="up">
        <div class="up__gallery">
            <GallerayProd :images="productData?.product.images ?? []" />
        </div>
        <div class="up__info">
            <h1>{{ productData?.product.name }}</h1>
            <p>{{ formattedPrice}}</p>
            <p>{{ productData?.product.short_description }}</p>
            <RatingStars  :rating="averageRating" :reviews-count="productData?.reviews.length ?? 0"/>
        </div>
    </div>
</template>

<style scoped>
.up{
    display: flex;
    gap: 10%;
}    
</style>