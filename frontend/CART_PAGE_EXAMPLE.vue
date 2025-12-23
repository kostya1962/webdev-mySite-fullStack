<!-- Пример обновления страницы корзины для использования нового store -->
<!-- Файл: frontend/pages/cart.vue -->

<script setup lang="ts">
const cartStore = useCartStore()
const authStore = useAuthStore()
const router = useRouter()

// При загрузке страницы восстанавливаем корзину с сервера
onMounted(async () => {
  if (authStore.email) {
    await cartStore.restore(authStore.email)
  }
})

// Функция для оформления заказа
async function checkout() {
  if (!authStore.token) {
    await router.push('/auth/login')
    return
  }

  if (cartStore.cartItems.length === 0) {
    alert('Корзина пуста!')
    return
  }

  try {
    const response = await $fetch('/api/orders/auth', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${authStore.token}`,
      },
      body: {
        productIDs: cartStore.cartItems.map(item => item.product.id),
        quantities: cartStore.cartItems.map(item => item.quantity),
      },
    })

    if (response.success) {
      cartStore.clearCart()
      await router.push('/account')
      alert('Заказ успешно создан!')
    }
  } catch (error) {
    console.error('Ошибка при создании заказа:', error)
    alert('Ошибка при создании заказа')
  }
}
</script>

<template>
  <div class="cart-page">
    <h1>Корзина</h1>

    <div v-if="cartStore.cartItems.length === 0" class="empty-cart">
      <p>Корзина пуста</p>
      <NuxtLink to="/catalog" class="btn-continue">
        Продолжить покупки
      </NuxtLink>
    </div>

    <div v-else class="cart-content">
      <div class="cart-items">
        <div v-for="item in cartStore.cartItems" :key="item.product.id" class="cart-item">
          <div class="item-image">
            <!-- Изображение товара -->
            <img :src="`/api/images/${item.product.images[0]}`" :alt="item.product.name" />
          </div>

          <div class="item-info">
            <h3>{{ item.product.name }}</h3>
            <p class="sku">{{ item.product.sku }}</p>
            <p class="price">{{ item.product.price }} ₽</p>
          </div>

          <div class="item-quantity">
            <label>Количество:</label>
            <input
              type="number"
              :value="item.quantity"
              min="1"
              @change="
                cartStore.updateQuantity(item.product.id, parseInt($event.target.value) || 1)
              "
            />
          </div>

          <div class="item-total">
            <p>{{ item.product.price * item.quantity }} ₽</p>
          </div>

          <button class="btn-remove" @click="cartStore.removeFromCart(item.product.id)">
            ✕
          </button>
        </div>
      </div>

      <div class="cart-summary">
        <h2>Итого</h2>

        <div class="summary-row">
          <span>Товаров:</span>
          <span>{{ cartStore.getItemsCount() }} шт.</span>
        </div>

        <div class="summary-row">
          <span>Сумма:</span>
          <span class="total-price">{{ cartStore.getTotalPrice() }} ₽</span>
        </div>

        <button class="btn-checkout" @click="checkout">
          Оформить заказ
        </button>

        <NuxtLink to="/catalog" class="btn-continue-shopping">
          Продолжить покупки
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cart-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

h1 {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 40px;
  color: var(--color-black);
}

.empty-cart {
  text-align: center;
  padding: 100px 20px;
}

.empty-cart p {
  font-size: 20px;
  color: var(--color-gray);
  margin-bottom: 30px;
}

.cart-content {
  display: grid;
  grid-template-columns: 1fr 350px;
  gap: 40px;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.cart-item {
  display: grid;
  grid-template-columns: 100px 1fr 120px 120px 40px;
  gap: 20px;
  align-items: center;
  padding: 20px;
  border: 1px solid var(--color-gray);
  border-radius: 8px;
}

.item-image {
  width: 100px;
  height: 100px;
  overflow: hidden;
  border-radius: 8px;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--color-black);
}

.sku {
  font-size: 12px;
  color: var(--color-gray);
  margin: 0;
}

.price {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-black);
  margin: 8px 0 0 0;
}

.item-quantity {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-quantity label {
  font-size: 12px;
  color: var(--color-gray);
}

.item-quantity input {
  padding: 8px;
  border: 1px solid var(--color-gray);
  border-radius: 4px;
  font-size: 14px;
  width: 100%;
}

.item-total {
  text-align: right;
  font-weight: 700;
  font-size: 16px;
  color: var(--color-black);
}

.btn-remove {
  width: 40px;
  height: 40px;
  padding: 0;
  border: 1px solid var(--color-gray);
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.3s;
}

.btn-remove:hover {
  background: var(--color-gray);
}

.cart-summary {
  position: sticky;
  top: 20px;
  border: 1px solid var(--color-gray);
  border-radius: 8px;
  padding: 20px;
  height: fit-content;
}

.cart-summary h2 {
  font-size: 20px;
  margin: 0 0 20px 0;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--color-gray);
  color: var(--color-black);
}

.summary-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 14px;
}

.total-price {
  font-weight: 700;
  font-size: 18px;
  color: var(--color-black);
}

.btn-checkout,
.btn-continue,
.btn-continue-shopping {
  width: 100%;
  padding: 12px 20px;
  margin-top: 20px;
  border: 2px solid var(--color-black);
  border-radius: 8px;
  background: var(--color-black);
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  text-align: center;
  transition: all 0.3s;
  display: block;
}

.btn-checkout:hover,
.btn-continue:hover,
.btn-continue-shopping:hover {
  background: white;
  color: var(--color-black);
}

.btn-continue-shopping {
  background: transparent;
  color: var(--color-black);
}

.btn-continue-shopping:hover {
  background: var(--color-gray);
}

@media (max-width: 768px) {
  .cart-content {
    grid-template-columns: 1fr;
  }

  .cart-item {
    grid-template-columns: 80px 1fr 30px;
  }

  .item-quantity,
  .item-total {
    display: none;
  }

  .cart-summary {
    position: static;
  }
}
</style>
