<script setup lang="ts">
    import { useAuthStore } from '~/state/auth.state';
    import { useAPI } from '~/composables/useAPI';
    import type { Order, GetOrdersResponse } from '~/interfaces/order.interface';
    import type { User } from '~/interfaces/user.interface';

    const API_URL = useAPI();

    definePageMeta({
        middleware: 'auth',
    });

    const authStore = useAuthStore();
    const orders = ref<Order[]>([]);
    const user = ref<User | null>(null);
    const isLoading = ref(true);
    const errorMessage = ref('');

    useSeoMeta({
        title: 'Мой аккаунт',
        description: 'Профиль пользователя и история заказов',
    });

    onMounted(async () => {
        if (!authStore.token) {
            await navigateTo('/auth/login');
            return;
        }

        try {
            isLoading.value = true;
            const res = await $fetch<GetOrdersResponse>(`${API_URL}/orders`, {
                headers: {
                    Authorization: `Bearer ${authStore.token}`
                }
            });
            orders.value = res.orders ?? [];
            user.value = res.user ?? null;
        } catch (e) {
            console.error('Failed to load orders', e);
            errorMessage.value = 'Не удалось загрузить заказы. Пожалуйста, попробуйте позже.';
        } finally {
            isLoading.value = false;
        }
    });

    const handleLogout = () => {
        authStore.clearToken();
        navigateTo('/auth/login');
    };

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return new Intl.DateTimeFormat('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }).format(date);
    };
</script>

<template>
    <section class="account-page">
        <h1>Мой аккаунт</h1>

        <div v-if="isLoading" class="loading">Загрузка...</div>
        <div v-else-if="errorMessage" class="error">{{ errorMessage }}</div>
        <div v-else class="content">
            <!-- Профиль -->
            <div class="profile-section">
                <h2>Профиль</h2>
                <div v-if="user" class="profile-info">
                    <div class="info-row">
                        <span class="label">Email:</span>
                        <span class="value">{{ user.email }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">Фамилия и имя:</span>
                        <span class="value">{{ user.name || '—' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">Телефон:</span>
                        <span class="value">{{ user.phone || '—' }}</span>
                    </div>
                    <div class="info-row">
                        <span class="label">Адрес доставки:</span>
                        <span class="value">{{ user.delivery_address || '—' }}</span>
                    </div>
                </div>
            </div>

            <!-- Заказы -->
            <div class="orders-section">
                <h2>История заказов</h2>
                <div v-if="orders.length === 0" class="no-orders">
                    У вас пока нет заказов
                </div>
                <div v-else class="orders-list">
                    <div v-for="order in orders" :key="order.id" class="order-card">
                        <div class="order-header">
                            <div class="order-info">
                                <h3>Заказ #{{ order.id }}</h3>
                                <span class="order-date">{{ formatDate(order.created_at) }}</span>
                            </div>
                            <div class="order-status" :class="order.status">
                                {{ order.status }}
                            </div>
                        </div>

                        <div class="order-products">
                            <div v-if="order.products && order.products.length > 0">
                                <div v-for="product in order.products" :key="product.id" class="product-item">
                                    <div class="product-name">{{ product.name }}</div>
                                    <div class="product-details">
                                        <span class="price">{{ product.price }} ₽</span>
                                        <span class="sku">{{ product.sku }}</span>
                                    </div>
                                </div>
                            </div>
                            <div v-else class="no-products">Информация о товарах недоступна</div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Кнопка выхода -->
            <div class="logout-section">
                <button class="logout-btn" @click="handleLogout">Выход</button>
            </div>
        </div>
    </section>
</template>

<style scoped>
.account-page {
    max-width: 1000px;
    margin: 0 auto;
    padding: 24px;
}

.account-page h1 {
    font-size: 32px;
    margin-bottom: 32px;
    color: var(--color-black);
}

.loading,
.error {
    padding: 40px;
    text-align: center;
    font-size: 16px;
}

.error {
    color: #d32f2f;
    background-color: #ffebee;
    border-radius: 8px;
}

.content {
    display: flex;
    flex-direction: column;
    gap: 32px;
}

/* Profile Section */
.profile-section {
    border: 1px solid var(--color-gray);
    border-radius: 12px;
    padding: 24px;
    background: #fafafa;
}

.profile-section h2 {
    font-size: 20px;
    margin-bottom: 20px;
    color: var(--color-black);
}

.profile-info {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
}

.info-row {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.info-row .label {
    font-weight: 600;
    font-size: 14px;
    color: var(--color-dark-gray);
}

.info-row .value {
    font-size: 16px;
    color: var(--color-black);
}

/* Orders Section */
.orders-section h2 {
    font-size: 20px;
    margin-bottom: 20px;
    color: var(--color-black);
}

.no-orders {
    text-align: center;
    padding: 40px;
    color: var(--color-gray);
    background: #fafafa;
    border-radius: 8px;
}

.orders-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.order-card {
    border: 1px solid var(--color-gray);
    border-radius: 12px;
    padding: 20px;
    background: #fff;
    transition: box-shadow 0.3s ease;
}

.order-card:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.order-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--color-gray);
}

.order-info h3 {
    margin: 0 0 8px 0;
    font-size: 18px;
    color: var(--color-black);
}

.order-date {
    font-size: 14px;
    color: var(--color-dark-gray);
}

.order-status {
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    background: #e0e0e0;
    color: #424242;
}

.order-status.новый {
    background: #e3f2fd;
    color: #1976d2;
}

.order-status.обработка {
    background: #fff3e0;
    color: #f57c00;
}

.order-status.отправлен {
    background: #f3e5f5;
    color: #7b1fa2;
}

.order-status.доставлен {
    background: #e8f5e9;
    color: #388e3c;
}

.order-products {
    margin-top: 12px;
}

.product-item {
    padding: 12px 0;
    border-bottom: 1px solid #f0f0f0;
}

.product-item:last-child {
    border-bottom: none;
}

.product-name {
    font-weight: 600;
    color: var(--color-black);
    margin-bottom: 8px;
}

.product-details {
    display: flex;
    gap: 16px;
    font-size: 14px;
    color: var(--color-dark-gray);
}

.price {
    font-weight: 600;
    color: var(--color-accent);
}

.no-products {
    text-align: center;
    padding: 16px;
    color: var(--color-gray);
    font-size: 14px;
}

/* Logout Section */
.logout-section {
    display: flex;
    gap: 12px;
}

.logout-btn {
    padding: 12px 24px;
    background: var(--color-black);
    color: #fff;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 16px;
    font-weight: 600;
    transition: background 0.3s ease;
}

.logout-btn:hover {
    background: #333;
}

@media (max-width: 768px) {
    .profile-info {
        grid-template-columns: 1fr;
    }

    .order-header {
        flex-direction: column;
        gap: 12px;
    }

    .product-details {
        flex-direction: column;
    }
}
</style>