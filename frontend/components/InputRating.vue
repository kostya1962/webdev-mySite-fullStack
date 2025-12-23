<script lang="ts" setup>
    import { computed } from 'vue'

    interface Props {
        modelValue: number
    }

    const props = defineProps<Props>()

    const emit = defineEmits<{
        (e: 'update:modelValue', value: number): void
    }>()

    const maxStars = 5

    const filledStars = computed(() => {
        return Math.floor(props.modelValue)
    })

    function setRating(value: number) {
        emit('update:modelValue', value)
    }
</script>

<template>
    <div class="rating">
        <div class="stars">
            <Icon
            v-for="index in maxStars"
            :key="index"
            name="icons:star"
            size="21px"
            class="star"
            :class="{
                'star--active': index <= filledStars,
                'star--inactive': index > filledStars
            }"
            @click="setRating(index)"
            />
        </div>
    </div>
</template>


<style scoped>
    .rating {
        display: flex;
        align-items: center;
    }

    .stars {
        display: flex;
        gap: 5px;
    }

    .star {
        cursor: pointer;
    }

    .star--active {
        color: var(--color-black);
    }

    .star--inactive {
        color: var(--color-gray);
    }
</style>
