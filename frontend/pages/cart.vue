<script setup lang="ts">
import { useAuthStore } from '~/state/auth.state';
import { useCartStore } from '~/state/cart.state';

useSeoMeta({
    title: 'Корзина товаров',
    description: 'Корзина товаров интернет магазина Shopper',
    ogDescription: 'Корзина товаров интернет магазина Shopper',
});

const cartStore = useCartStore();
const authStore = useAuthStore();
const router = useRouter();
const API_URL = useRuntimeConfig().public.apiurl || '';

const fullName = ref('');
const phone = ref('');
const deliveryAddress = ref('');

onMounted(async () => {
    if (authStore.email) {
        await cartStore.restore(authStore.email);
    }
});

async function submitOrder() {
    if (!authStore.token) {
        await router.push('/auth/login');
        return;
    }

    if (!cartStore.cartItems.length) return;

    const product_ids = cartStore.cartItems.map((i) => i.product.id);
    
    console.log('[submitOrder] API_URL:', API_URL);
    console.log('[submitOrder] Token:', authStore.token.substring(0, 20) + '...');
    console.log('[submitOrder] Request payload:', {
        product_ids,
        name: fullName.value,
        phone: phone.value,
        delivery_address: deliveryAddress.value,
    });

    try {
        const response = await $fetch(`${API_URL}/orders/auth`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${authStore.token}`,
            },
            body: {
                product_ids,
                name: fullName.value,
                phone: phone.value,
                delivery_address: deliveryAddress.value,
            },
        });
        
        console.log('[submitOrder] Success response:', response);

        // Если успешно — очистим корзину и перейдём в профиль
        cartStore.clearCart();
        await router.push('/account');
    } catch (e: any) {
        console.error('Order submit error', e);
        console.error('[submitOrder] Error status:', e?.status);
        console.error('[submitOrder] Error data:', e?.data);
        console.error('[submitOrder] Error message:', e?.message);
    }
}

// computed-safe bindings to avoid null during hydration
const items = computed(() => cartStore.cartItems ?? []);
const itemsCount = computed(() => items.value.reduce((c, it) => c + (it.quantity || 0), 0));
const totalPrice = computed(() => items.value.reduce((t, it) => t + ((it.product?.price * (1 - (it.product?.discount || 0) * 0.01) || 0) * (it.quantity || 0)), 0));
</script>

<template>
    <section class="cart-page">
        <h1>Корзина</h1>

        <div v-if="!items.length" class="empty">Ваша корзина пуста</div>

        <div v-else class="content">
            <div class="items">
                <div v-for="item in items" :key="item.product.id" class="item">
                    <div class="item-left">
                        <div class="name">{{ item.product.name }}</div>
                        <div class="sku">{{ item.product.sku }}</div>
                    </div>
                    <div class="item-center">
                        <input
                            type="number"
                            :value="item.quantity"
                            min="1"
                            @change="cartStore.updateQuantity(item.product.id, parseInt(($event.target as HTMLInputElement)?.value) || 1)"
                        />
                    </div>
                    <div class="item-right">
                        <div class="price">{{ item.product.price }} ₽</div>
                        <button class="remove" @click="cartStore.removeFromCart(item.product.id)">Удалить</button>
                    </div>
                </div>
            </div>

            <aside class="summary">
                <div class="row"><span>Товаров:</span><span>{{ itemsCount }} шт.</span></div>
                <div class="row"><span>Сумма:</span><span>{{ totalPrice }} ₽</span></div>

                <div class="personal-form">
                    <label>Полное имя и фамилия</label>
                    <input v-model="fullName" type="text" placeholder="Иван Иванов" />

                    <label>Телефон</label>
                    <input v-model="phone" type="tel" placeholder="+7 900 000 00 00" />

                    <label>Адрес доставки</label>
                    <input v-model="deliveryAddress" type="text" placeholder="Улица, дом, квартира" />
                </div>

                <button class="checkout" @click="submitOrder">Оформить заказ</button>
            </aside>
        </div>
    </section>
</template>

<style scoped>
.cart-page{
    max-width:1200px;
    margin:0 auto;
    padding:24px
}
.empty{
    padding:40px;
    text-align:center;
    color:var(--color-gray)
}
.content{
    display:flex;gap:24px
}
.items{
    flex:1;display:flex;
    flex-direction:column;
    gap:12px
}
.item{
    display:flex;
    justify-content:space-between;
    align-items:center;
    padding:12px;
    border:1px solid var(--color-gray);
    border-radius:8px
}
.item-left .name{
    font-weight:700
}
.item-center input{
    width:80px;padding:6px
}
.summary{
    width:320px;
    border:1px solid var(--color-gray);
    padding:16px;border-radius:8px
}
.summary .row{
    display:flex;
    justify-content:space-between;
    margin-bottom:12px
}
.checkout{width:100%;
    padding:12px;
    background:var(--color-black);
    color:#fff;
    border:none;
    border-radius:8px;
    cursor:pointer
}
.remove{
    background:transparent;
    border:1px solid var(--color-gray);
    padding:6px 10px;
    border-radius:6px;
    cursor:pointer
    }
.personal-form{display:flex;flex-direction:column;gap:8px;margin:12px 0}
.personal-form label{font-size:13px;color:var(--color-dark-gray)}
.personal-form input{padding:8px;border:1px solid var(--color-gray);border-radius:6px}
</style>