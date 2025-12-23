<script lang="ts" setup>
    import { computed } from 'vue'

    interface Props {
        rating: number
    }

    const props = defineProps<Props>()

    const maxStars = 5

    const filledStars = computed(() => {
        return Math.floor(props.rating)
    })
</script>

<template>
    <div class="rating">
        <div class="stars">
        <Icon
            v-for="index in maxStars"
            :key="index"
            name="icons:star"
            size="21px"
            :class="{
            'star--active': index <= filledStars,
            'star--inactive': index > filledStars
            }"
        />
        </div>


    </div>
</template>


<style scoped>
    .rating {
        display: flex;
        align-items: center;
        gap: 50px;
    }

    .stars {
        display: flex;
        gap: 5px;
    }

    .star--active {
        color: var(--color-black); 
    }

    .star--inactive {
        color: var(--color-gray); 
    }

    .reviews {
    font-size: 14px;
        color: var(--color-dark-gray); 
    }
</style>