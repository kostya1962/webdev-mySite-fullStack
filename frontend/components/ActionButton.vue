<script setup lang="ts">
const { color = 'primary' } = defineProps<{
    color?: 'primary' | 'ghost';
}>()
</script>

<template>
    <button :class="color">
        <slot />
    </button>
</template>

<style scoped>
    .primary{
        border: 1px solid var(--color-black);
        background-color: var(--color-black);
        border-radius: 4px;
        font-family: var(--font);
        cursor: pointer;
        padding: 16px;
        font-size: 16px;
        font-weight: 700;
        color: var(--color-white);
    }

    .primary:hover{
        background-color: var(--color-black-hover);
    }

    .ghost{
        border: 1px solid var(--color-black);
        border-radius: 4px;
        font-family: var(--font);
        cursor: pointer;
        padding: 16px;
        font-size: 16px;
        font-weight: 700;
        color: (--color-black);
        background-color: var(--color-white);
    }

    .ghost:hover{
        background-color: var(--color-black-hover);
        color: var(--color-white);
    }
</style>