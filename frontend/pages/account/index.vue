<script setup>
    import { useAuthStore } from '~/state/auth.state';

    const API_URL = useRuntimeConfig().public.apiurl || '';

    definePageMeta({
        middleware: 'auth',
    });

    const authStore = useAuthStore();
    const orders = ref([])
    const user = ref(null)

    onMounted(async () => {
        try {
            const res = await $fetch(`${API_URL}/orders`, {
                headers: {
                    Authorization: `Bearer ${authStore.token}`
                }
            })
            orders.value = res.orders ?? []
            user.value = res.user ?? null
        } catch (e) {
            console.error('Failed to load orders', e)
        }
    })
</script>

<template>
    <div>
        <div class="account-page">
            <div class="profile">
                <h2>Профиль</h2>
                <div v-if="user">
                    <div><strong>Email:</strong> {{ user.email }}</div>
                    <div><strong>Имя:</strong> {{ user.name }}</div>
                    <div><strong>Телефон:</strong> {{ user.phone }}</div>
                    <div><strong>Адрес:</strong> {{ user.delivery_address }}</div>
                </div>
            </div>

            <div class="orders">
                <h2>Заказы</h2>
                <div v-if="orders.length === 0">Нет заказов</div>
                <ul>
                    <li v-for="o in orders" :key="o.id">
                        Заказ #{{ o.id }} — {{ o.status }} — {{ o.created_at }}
                    </li>
                </ul>
            </div>

            <NuxtLink to="/" @click="authStore.clearToken">Выход</NuxtLink>
        </div>
    </div>
</template>