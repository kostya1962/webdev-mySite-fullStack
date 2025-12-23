<script setup lang="ts">
import GallerayProd from '~/components/GallerayProd.vue';
import AddToCart from '~/components/AddToCart.vue';
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
            <div class="up__info__name">{{ productData?.product.name }}</div>
            <div class="up__info__price">{{ formattedPrice}}</div>
            <div class="up__info__description">{{ productData?.product.short_description }}</div>
            <RatingStars  :rating="averageRating" :reviews-count="productData?.reviews.length ?? 0"/>
            <div>
                <AddToCart v-if="productData?.product" :product="productData.product" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.up{
    display: flex;
    gap: 3%;
}    

.up__gallery{
    width: 50%;
}

.up__info{
    display: flex;
    flex-direction: column;
    gap: 23px;
    margin-top: 1%;
}

.up__info__name{
    font-size: 30px;
    line-height: 30px;
    font-weight: 100;
    color: var(--color-black);
}

.up__info__price{
    font-size: 20px;
    font-weight: 700;
    line-height: 26px;
    text-transform: capitalize;
    color: var(--color-accent);
    margin-bottom: 15%;
}

.up__info__description{
    color: var(--color-dark-gray);
    font-size: 18px;
    line-height: 20px;
}
</style>