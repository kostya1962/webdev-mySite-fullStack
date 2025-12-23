<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface';
import { useAuthStore } from '~/state/auth.state';
import { useCartStore } from '~/state/cart.state';

const { product } = defineProps<{
product: Product;
}>();

const authStore = useAuthStore();
const cartStore = useCartStore();
const router = useRouter();

const isAdded = ref(false);
const quantity = ref(1);


async function handleAddClick() {
if (!authStore.token) {
    await router.push('/auth/login');
    return;
}

// Переводим в режим выбора количества
isAdded.value = true;
quantity.value = 1;
}


function handleConfirmAdd() {
cartStore.addToCart(product, quantity.value);

isAdded.value = false;
quantity.value = 1;
}


function handleQuantityChange(newQuantity: number) {
quantity.value = newQuantity;
}
</script>

<template>
    <div class="add-to-cart">
        <div v-if="!isAdded" class="before-add">
            <ActionButton @click="handleAddClick">
                Добавить в корзину
            </ActionButton>
        </div>


        <div v-else>
            <div class="cart-controls">
                <CounterFiled @quantity-change="handleQuantityChange" />

                <button class="confirm-btn" @click="handleConfirmAdd">
                    Добавить
                </button>

                <button class="remove-btn" @click="isAdded = false">
                    Отменить
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.add-to-cart {
width: 100%;
}

.cart-controls {
display: flex;
gap: 12px;
align-items: center;
width: 100%;
}

.remove-btn,
.confirm-btn {
padding: 12px 22px;
font-size: 16px;
border: 2px solid var(--color-black);
border-radius: 8px;
cursor: pointer;
transition: all 0.3s ease;
font-weight: 600;
min-height: 44px;
display: flex;
align-items: center;
justify-content: center;
width: 40%;
}

.before-add button {
    width: 100%;
}

.remove-btn {
background: transparent;
color: var(--color-black);
flex: 0 1 auto;
}

.remove-btn:hover {
background: var(--color-gray);
}

.confirm-btn {
background: var(--color-black);
color: white;
flex: 0 1 auto;
}

.confirm-btn:hover {
background: var(--color-gray);
color: var(--color-black);
}
</style>
