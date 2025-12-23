<script setup lang="ts">
    const emit = defineEmits<{
        (e: "update:count", value: number): void;
    }>();


    const count = ref<number>(0);


    const increment = (): void => {
        count.value++;
        emit("update:count", count.value);
    };

    const decrement = (): void => {
        count.value--;
        emit("update:count", count.value);
    };
</script>

<template>
    <div class="counter">
        <button class="btn" @click="decrement">−</button>
        <span class="value">{{ count }}</span>
        <button class="btn" @click="increment">+</button>
    </div>
</template>

<style scoped>
    .counter {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .btn {
        width: 32px;
        height: 32px;
        border: 1px solid #ccc;
        background-color: #f5f5f5;
        font-size: 18px;
        cursor: pointer;
    }

    .value {
        min-width: 20px;
        text-align: center;
        font-size: 16px;
    }
</style>