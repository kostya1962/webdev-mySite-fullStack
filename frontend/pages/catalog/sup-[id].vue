<!-- eslint-disable vue/no-multiple-template-root -->
<script setup lang="ts">
import GallerayProd from '~/components/GallerayProd.vue';
import AddToCart from '~/components/AddToCart.vue';
import type { ProductIDRsponse } from '~/interfaces/productID.interface';
import { useFavoriteStore } from '~/state/favorite.state';

const route = useRoute(); //извлекаю текущий маршрут (путь и параметры)
const API_URL = useAPI();
const favoriteState = useFavoriteStore();
const activeFlag = ref<number>(0);


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

  const countReviews = computed(() => {
        return productData.value?.reviews.length ?? 0;
    });

    function setActiveFlag(val: number) {
        activeFlag.value = val;
    }
</script>

<template>
<div>
    <div class="up">
        <div class="up__gallery">
            <GallerayProd :images="productData?.product.images ?? []" />
        </div>
        <div class="up__info">
            <div class="up__info__name">{{ productData?.product.name }}</div>
            <div class="up__info__price">{{ formattedPrice}}</div>
            <div class="up__info__description">{{ productData?.product.short_description }}</div>
            <div class="up__info__rating">
                <RatingStars  :rating="averageRating" />
                <span class="reviews">
                    {{ countReviews?? 0 }} 
                    <span v-show="countReviews > 4">отзывов</span>
                    <span v-show="countReviews == 1">отзыв</span>
                    <span v-show="countReviews > 1 && countReviews < 5">отзыва</span>
                </span>
            </div>
            <div class="push-card">
                <AddToCart v-if="productData?.product" :product="productData.product" />
            </div>
            <div class="up__info__additional">
                <div>
                    <AddFavorite v-show="favoriteState.isFavorite(Number(route.params.id))" :id="productData?.product.id ?? 0" :is-shown="true" />
                    <DelFavorite v-show="!favoriteState.isFavorite(Number(route.params.id))" :id="productData?.product.id ?? 0" :is-shown="true" />
                </div>
                <div class="up__info__hr"></div>
                <span class="up__info__social">
                    <NuxtLink to="#">
                        <Icon name="icons:vk" size="20px"/>
                    </NuxtLink>
                    <NuxtLink to="#">
                        <Icon name="icons:tgc" size="20px"/>
                    </NuxtLink>
                </span>

            </div>
            <div class="info-row">
                <span class="label">SKU:</span>
                <span class="value">{{ productData?.product.sku || '—' }}</span>
            </div>
            <div class="info-row">
                <span class="label">Категория:</span>
                <span class="value">{{ productData?.product.category.name || '—' }}</span>
            </div>
        </div>
    </div>

    
    <div class="dawn">
        <div class="dawn__header">
            <div class="dawn__header__btn" :class="{ blacklight: activeFlag === 0 }" @click="setActiveFlag(0)">Описание</div>
            <button class="dawn__header__btn" :class="{ blacklight: activeFlag === 1 }" @click="setActiveFlag(1)">Отзывы({{ countReviews }})</button>
        </div>

        <div v-show="activeFlag === 0" class="dawn__description">
            <p>{{ productData?.product.long_description }}</p>
        </div>

        <div v-show="activeFlag === 1" class="dawn_panel">
            <div class="review-list">
                <ReviewOne 
                    v-for="review in productData?.reviews" 
                    :key="review.id" 
                    v-bind="review"/>
            </div>
            <div class="review-form">
                <ReviewForm />
            </div>
        </div>
    </div>
</div></template>

<style scoped>
.up{
    display: flex;
    gap: 3%;
    margin-bottom: 101px;
}    

.up__gallery{
    width: 50%;
}

.up__info{
    display: flex;
    flex-direction: column;
    gap: 5px;
    margin-top: 1%;
    width: 50%;
}

.up__info__name{
    font-size: 30px;
    line-height: 30px;
    font-weight: 100;
    margin-bottom: 3%;
    color: var(--color-black);
}

.up__info__price{
    font-size: 20px;
    font-weight: 700;
    line-height: 26px;
    text-transform: capitalize;
    color: var(--color-accent);
    margin-bottom: 10%;
}

.up__info__description{
    color: var(--color-dark-gray);
    font-size: 18px;
    line-height: 20px;
    margin-bottom: 3%;
}   

.up__info__rating{
    display: flex;
    align-items: center;
    gap: 30px;
    margin-bottom: 5%;
}

.reviews {
    font-size: 14px;
        color: var(--color-dark-gray); 
    }

.push-card{
    margin-top: 10%;
    margin-bottom: 5%;
    width: 100%;
}



.up__info__hr{
    border-left: 1px solid var(--color-dark-gray);
    height: 17px;
    align-self: flex-start;
    margin-top: 4px;
}

.up__info__additional   {
    display: flex;
    gap: 20px;
    align-items: center;
    width: 100%;
    margin-bottom: 5%;
}

.up__info__social   {
    display: flex;
    color: var(--color-dark-gray);
    gap: 20px;
    
}

.up__info__social a{
    color: var(--color-dark-gray);
}

.up__info__social a:hover{
    color: var(--color-black);
}

.info-row {
    display: flex;
    gap: 10px;
    font-size: 16px;
}

.info-row .label {
    font-weight: 600;
    color: var(--color-black); 
}

.info-row .value {
    color: var(--color-dark-gray);
}

.dawn__header {
    border-bottom: 1px solid var(--color-gray);
    display: flex;
    gap: 85px;
    margin-bottom: 43px;
}

.dawn__header__btn{
    font-size: 20px;
    line-height: 27px;
    color: var(--color-black);
    text-decoration: none;
    padding-bottom: 16px;
    margin-right: 32px;
    background: none;
    cursor: pointer;
    border: none
}

.blacklight{
    border-bottom: 1px solid var(--color-black);
}

.dawn_panel{
    display: flex;
    gap: 10%;
}

.review-list{
    width: 45%;
}
.review-form{
    width: 45%;
}
</style>