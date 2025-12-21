<script setup lang=ts>
    import type { NuxtError } from '#app';

    const { error } = defineProps<{
        error: NuxtError;
    }>();

    function goHome() {
        navigateTo('/');
    } 
</script>

<template>
    <div class="error-wrapper">
        <h1><span class="error-bold">{{ error.statusCode }}</span> ОШИБКА</h1>
        <div  v-if="error.statusCode == 404" class="error-message">
            Страница не найдена. Попробуйте перейти на главную страницу.
        </div>
        <div  v-else class="error-message">{{ error.statusMessage }}</div>
        <ActionButton color="ghost" @click="goHome">На главную</ActionButton>
    </div>
</template>

<style scoped>
.error-message{
    color: var(--color-dark-gray);
    text-align: center;
    font-weight: 400;
    font-size: 20px;
    line-height: 26px;
    margin-bottom: 40px;
}

.error-wrapper{
    display:flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 24px;
}

.error-bold{
    font-weight: 800;
}
</style>