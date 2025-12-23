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

// Обработчик первого клика на кнопку
async function handleAddClick() {
  if (!authStore.token) {
    // Если не авторизован, перенаправляем на страницу входа
    await router.push('/auth/login');
    return;
  }

  // Переводим в режим выбора количества
  isAdded.value = true;
  quantity.value = 1;
}

// Обработчик окончательного добавления в корзину
function handleConfirmAdd() {
  cartStore.addToCart(product, quantity.value);
  // Возвращаемся в исходное состояние
  isAdded.value = false;
  quantity.value = 1;
}

// Обработчик изменения количества из CounterFiled
function handleQuantityChange(newQuantity: number) {
  quantity.value = newQuantity;
}
</script>

<template>
  <div class="add-to-cart">
    <!-- Начальное состояние - просто кнопка -->
    <template v-if="!isAdded">
      <ActionButton @click="handleAddClick">
        Добавить в корзину
      </ActionButton>
    </template>

    <!-- Состояние с количеством -->
    <template v-else>
      <div class="cart-controls">
        <!-- Кнопка удаления (без заливки) -->
        <button class="remove-btn" @click="isAdded = false">
          Удалить
        </button>

        <!-- Счётчик -->
        <CounterFiled @quantity-change="handleQuantityChange" />

        <!-- Кнопка подтверждения добавления -->
        <button class="confirm-btn" @click="handleConfirmAdd">
          Добавить
        </button>
      </div>
    </template>
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
