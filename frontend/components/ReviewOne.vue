<script setup lang="ts">
    import { computed } from 'vue'
    import RatingStars from '@/components/RatingStars.vue'
    import type { Review } from '~/interfaces/review.interface'


    const props = defineProps<Review>()


    const formattedDate = computed(() => {
    const date = typeof props.created_at === 'string'
        ? new Date(props.created_at)
        : props.created_at

    return date.toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
    })
    })
</script>

<template>
    <div class="review">
        <div class="review__header">
            <span class="review__author">
                {{ props.name }}
            </span>

            <span class="review__date">
                {{ formattedDate }}
            </span>
        </div>

        <div class="review__rating">
            <RatingStars :rating="props.rating" />
        </div>

        <div class="review__text">
        {{ props.text }}
        </div>
    </div>
</template>

<style scoped>
    .review {
        max-width: 600px;
        font-family: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
        color: #1f1f1f;
        border-bottom: 1px solid var(--color-gray);
        padding-bottom: 52px;
        margin-bottom: 30px;
    }

    .review__header {
        display: flex;
        gap: 36px;
        align-items: baseline;
        margin-bottom: 20px;
    }

    .review__author {
    font-weight: 400;
    font-size: 20px;
    }

    .review__date {
    font-size: 14px;
    color: #8a8a8a;
    }

    .review__rating {
        margin: 4px 0 23px 0;
    }

    .review__text {
    font-size: 14px;
    line-height: 1.4;
    }
</style>