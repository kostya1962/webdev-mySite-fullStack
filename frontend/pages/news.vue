<script setup lang="ts">
    import { computed } from 'vue';
    import type { GetNewsResponse, NewsItem } from '~/interfaces/news.interface';

    useSeoMeta({
        title: 'Новости',
        description: 'Последние новости интернет магазина Shopper',
        ogDescription: 'Последние новости интернет магазина Shopper',
    });

    const API_URL = useAPI();

    const { data: newsData, pending, error } = await useFetch<GetNewsResponse>(`${API_URL}/news`);

    const newsList = computed<NewsItem[]>(() => {
        const val = newsData?.value;
        return Array.isArray(val) ? (val as NewsItem[]) : [];
    });

</script>

<template>
<div class="news-page">

        <h1 class="news-page__title">Новости</h1>


        <div class="news-page__content">
            <div v-if="pending" class="loading">Загрузка новостей...</div>
            <div v-else-if="error" class="error">Ошибка загрузки новостей</div>
            <div v-else-if="!newsList || newsList.length === 0" class="no-news">
                Новостей пока нет
            </div>
            <div v-else>
                <div class="news-page__list">
                    <NewsCard
                        v-for="news in newsList"
                        :key="news.id"
                        :news="news"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.news-page {
max-width: 900px;
margin: 0 auto;
padding: 40px 20px;
display: flex;
flex-direction: column;
align-items: center;
}

.news-page__title {
font-size: 36px;
font-weight: 700;
margin-bottom: 40px;
text-align: center;
}

.news-page__content {
width: 100%;
display: flex;
flex-direction: column;
gap: 20px;
}

.news-page__item {
display: flex;
justify-content: center;
}

.loading,
.error,
.no-news {
font-size: 18px;
color: #555;
text-align: center;
margin-top: 20px;
}

.news-page__list{
    display: flex;
    flex-direction: column;
    gap: 24px;
}
</style>
