<script setup lang="ts">
import type { Banner, GetBannersResponse } from '~/interfaces/banner.interface'

const API_URLimage = useAPIimage();
const API_URL = useAPI();
const current = ref(0);

function goTo(index: number) {
    current.value = index - 1;
}
const { data: bannersData, pending: isLoading } = await useFetch<GetBannersResponse>(API_URL + '/banners', {
    key: 'get-banners-home',
});
const slides = computed<Banner[]>(() => {
    return bannersData.value?.banners ?? []
});

console.log(`url(${API_URLimage}${slides.value[0]?.image})`);

function next() { 
    current.value = (current.value + 1) % slides.value.length 
    }
</script>

<template>
    <section v-if="slides.length || isLoading" class="home-banner">
        <div v-if="slides.length" class="banner-rect"  :style="{ backgroundImage: slides[current] && slides[current]?.image ? `url(${API_URLimage}${slides[current]?.image})` : '' }" @click="next">

            <!-- slides -->
            <div v-for="(s, idx) in slides" :key="s.id" class="banner-slide"  :aria-hidden="idx !== current" :style="{ display: idx === current ? 'block' : 'none' }">
                <div class="banner-content">
                    <h2 class="banner-title">{{ s.product.name }}</h2>
                    <div class="banner-price">{{ s.product.price }} ₽</div>
                    <NuxtLink :to="`/catalog/sup-${s.product.id}`" class="banner-cta">Перейти к товару</NuxtLink>
                </div>

                <div class="banner-product-img">
                </div>
            </div>

            <div class="banner-dots" aria-hidden="true">
                <button v-for="(s, idx) in slides" :key="s.id" :class="['dot', { active: idx === current }]" aria-hidden="true" @click="goTo(idx)"  />
            </div>
        </div>

        <div v-else class="banner-placeholder">Загрузка баннера...</div>
    </section>
</template>

<style scoped>
    .home-banner {
        margin-bottom: 28px;
    }
    .banner-rect {
        position: relative;
        width: 1200px;
        height: 600px;
        border-radius: 20px;
        overflow: hidden;
        margin: 0 auto 28px auto;
        box-shadow: 0 6px 24px rgba(0,0,0,0.08);
        background-size: cover;
        background-position: center center;
        background-repeat: no-repeat;
    }
    .banner-dots {
        position: absolute;
        left: 50%;
        transform: translateX(-50%);
        bottom: 22px;
        display: flex;
        gap: 12px;
        z-index: 3;
        align-items: center;
    }
    .dot {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        background: rgba(255,255,255,0.2);
        border: 2px solid transparent;
    }
    .dot.active {
        width: 16px;
        height: 16px;
        background: transparent;
        border: 2px solid var(--color-black);
    }
    .banner-content {
        position: absolute;
        left: 46px;
        top: 230px;
        z-index: 4;
        color: var(--color-black);
    }
    .banner-title {
        margin: 0;
        font-size: 33px;
        line-height: 1.1;
        color: var(--color-black);
        max-width: 900px;
    }
    .banner-price {
        margin-top: 35px;
        font-size: 28px;
        color: var(--color-black);
        font-weight: 700;
    }
    .banner-cta {
        display: inline-block;
        margin-top: 24px;
        padding: 12px 22px;
        font-size: 20px;
        color: var(--color-black);
        border: 2px solid var(--color-black);
        background: rgba(255,255,255,0.06);
        border-radius: 8px;
        text-decoration: none;
    }
    .banner-cta:hover {
        background-color: var(--color-gray);
        
    }
    .banner-product-img {
        position: absolute;
        right: 60px;
        top: 140px;
        z-index: 4;
        }
    .banner-product-img img {
        max-width: 620px;
        max-height: 900px;
        object-fit: contain;
    }
</style>
