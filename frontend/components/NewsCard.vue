<script setup lang="ts">
    import type { NewsItem } from '~/interfaces/news.interface';


    const props = defineProps<{
        news: NewsItem
    }>()

    const imagePrefix = useAPIimage();

    const image = computed(() => {
        const path = props.news.image ?? '';
        if (!path) return '';
        const prefix = (imagePrefix ?? '').replace(/\/+$/, '');
        const cleanPath = path.replace(/^\/+/, '');
        const full = prefix ? `${prefix}/${cleanPath}` : `/${cleanPath}`;
        return `url(${full})`;
    })

    const formattedDateTime = computed(() => {
        const date = new Date(props.news.created_at)

        const formattedDate = date.toLocaleDateString('ru-RU')
        const formattedTime = date.toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit'
        })

        return `${formattedDate} ${formattedTime}`
    })
</script>

<template>
    <div class="news-card">
        
        <div
            class="news-card__image"
            :style="{ backgroundImage: image }"
        />

            
        <div class="news-card__content">
            <div class="news-card__header">
                <h3 class="news-card__title">
                    {{ news.title }}
                </h3>
                <span class="news-card__date">
                    {{ formattedDateTime }}
                </span>
            </div>

            <p class="news-card__description">
                {{ news.description }}
            </p>
        </div>
    </div>
</template>

<style scoped>
    .news-card {
        display: flex;
        width: 100%;
        height: 300px;
        background: #ffffff;
        border-radius: 16px;
        overflow: hidden;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
        color: #000;
    }


    .news-card__image {
    width: 40%;
    background-size: cover;
    background-position: center;
    }

    .news-card__content {
    width: 60%;
    padding: 16px 20px;
    display: flex;
    flex-direction: column;
    gap: 20%
    }


    .news-card__header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 10px;
    }

    .news-card__title {
    font-size: 20px;
    font-weight: 600;
    margin: 0;
    line-height: 1.2;
    }

    .news-card__date {
    font-size: 13px;
    color: #555;
    white-space: nowrap;
    }


    .news-card__description {
    font-size: 14px;
    color: #222;
    line-height: 1.4;
    }
</style>