<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import type { GetProductsResponse, Product } from '~/interfaces/product.interface'
import HomeBanner from '~/components/HomeBanner.vue'

useSeoMeta({
    title: 'Главная страница',
    description: 'Главная страница магазина Shopper с указанием последних поступлений и новостей',
    ogDescription: 'Главная страница магазина Shopper с указанием последних поступлений и новостей',
});

const API_URL = useAPI()
const productsData = ref<GetProductsResponse | null>(null)
const isLoading = ref(true)

onMounted(async () => {
    try {
        const { data } = await useFetch<GetProductsResponse>(API_URL + '/products', {
            query: { limit: 1000, offset: 0 },
        })
        if (data.value) {
            productsData.value = data.value
        }
    } catch (e) {
        console.error('Ошибка загрузки товаров:', e)
    } finally {
        isLoading.value = false
    }
})

const recentProducts = computed(() => {
    const items: Product[] = productsData.value?.products ?? []
    return items
        .slice()
        .filter((p) => !!p.created_at)
        .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
        .slice(0, 6)
})
</script>

<template>
    <div class="home-page">
        <ClientOnly>
            <HomeBanner />
        
        <section class="latest-section">
            <div class="latest-header">
                <h1>Последние поступления</h1>
                <span class="latest-all"> 
                    <NuxtLink to="/catalog">Все</NuxtLink>
                </span>
            </div>

            <div v-if="isLoading" class="loading-message">Загрузка товаров...</div>
            <div v-else class="latest-grid">
                <CatalogCard
                    v-for="product in recentProducts.slice(0, 8)"
                    :key="product.id"
                    v-bind="product"
                />
            </div>
        </section>
        </ClientOnly>
    </div>
</template>

<style scoped>

.latest-section {
    margin-bottom: 32px;
}
.latest-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 50px;
}
.latest-all {
    color: var(--color-accent);
    font-size: 18px;
    font-weight: 600;
}
.latest-all a {
    color: inherit;
    text-decoration: none;
    font-size: inherit;
    font-weight: inherit;
}
.latest-grid {
    display: grid;
    width: 100%;
    gap: 30px 60px;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

.home-page {
    margin-bottom: 250px;
}

.loading-message {
    font-size: 16px;
    color: #707070;
    text-align: center;
    padding: 40px;
}
</style>