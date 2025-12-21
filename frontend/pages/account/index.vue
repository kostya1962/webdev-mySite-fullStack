<script setup>
    import { useAuthStore } from '~/state/auth.state';

    definePageMeta({
        middleware: 'auth',
    });

    const authStore = useAuthStore();
</script>

<template>
    <div>
        <NuxtLink to="/" @click="authStore.clearToken">
            Выход
        </NuxtLink>
    </div>
</template>