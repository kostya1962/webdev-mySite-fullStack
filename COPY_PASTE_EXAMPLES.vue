<!-- 🚀 КОПИРУЙ-ВСТАВЛЯЙ РЕШЕНИЕ -->
<!-- Используйте эти готовые примеры прямо в вашем проекте -->

<!-- =========================================================== -->
<!-- ФАЙЛ 1: Страница товара (/catalog/sup-[id].vue) -->
<!-- =========================================================== -->

<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface'

const route = useRoute()
const productId = Array.isArray(route.params.id) 
  ? route.params.id[0].replace('sup-', '') 
  : route.params.id.replace('sup-', '')

// Получить данные товара
const { data: product, pending: isLoading } = await useFetch<Product>(
  `/api/products/${productId}`,
  { key: `product-${productId}` }
)
</script>

<template>
  <div v-if="isLoading" class="loading">Загрузка товара...</div>

  <div v-else-if="product" class="product-container">
    <div class="product-gallery">
      <!-- Ваша галерея изображений -->
    </div>

    <div class="product-info">
      <h1 class="product-title">{{ product.name }}</h1>

      <div class="product-description">
        {{ product.short_description }}
      </div>

      <div class="product-price-section">
        <div class="price">{{ product.price }} ₽</div>

        <!-- 👇 ДОБАВЬТЕ ВОТ ЭТУ СТРОКУ 👇 -->
        <AddToCart :product="product" />
        <!-- 👆 ДОБАВЬТЕ ВОТ ЭТУ СТРОКУ 👆 -->
      </div>

      <div class="product-full-description">
        <h2>Описание товара</h2>
        <p>{{ product.long_description }}</p>
      </div>
    </div>
  </div>

  <div v-else class="error">Товар не найден</div>
</template>

<!-- =========================================================== -->
<!-- ФАЙЛ 2: Страница корзины (/cart.vue) -->
<!-- =========================================================== -->

<script setup lang="ts">
const cartStore = useCartStore()
const authStore = useAuthStore()
const router = useRouter()

// При загрузке восстанавливаем корзину с сервера
onMounted(async () => {
  if (authStore.email) {
    await cartStore.restore(authStore.email)
  }
})

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
    // Отправляем заказ
    const response = await $fetch('/api/orders/auth', {
      method: 'POST',
      headers: { Authorization: `Bearer ${authStore.token}` },
      body: {
        productIDs: cartStore.cartItems.map(item => item.product.id),
      },
    })

    if (response.success) {
      cartStore.clearCart()
      await router.push('/account')
      alert('Заказ успешно создан!')
    }
  } catch (error) {
    console.error('Ошибка:', error)
    alert('Ошибка при создании заказа')
  }
}
</script>

<template>
  <div class="cart-page">
    <h1>Корзина</h1>

    <div v-if="cartStore.cartItems.length === 0" class="empty-cart">
      <p>Корзина пуста</p>
      <NuxtLink to="/catalog" class="btn">Продолжить покупки</NuxtLink>
    </div>

    <div v-else class="cart-content">
      <div class="cart-items">
        <div v-for="item in cartStore.cartItems" :key="item.product.id" class="cart-item">
          <div class="item-info">
            <h3>{{ item.product.name }}</h3>
            <p class="price">{{ item.product.price }} ₽</p>
          </div>

          <div class="item-quantity">
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
            {{ item.product.price * item.quantity }} ₽
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
          <span class="label">Сумма:</span>
          <span class="total">{{ cartStore.getTotalPrice() }} ₽</span>
        </div>
        <button class="btn-checkout" @click="checkout">Оформить заказ</button>
      </div>
    </div>
  </div>
</template>

<!-- =========================================================== -->
<!-- ФАЙЛ 3: Header компонент (Layout/LayoutHeader.vue) -->
<!-- =========================================================== -->

<script setup lang="ts">
const cartStore = useCartStore()
</script>

<template>
  <header class="header">
    <div class="logo">Simple Shopper</div>
    
    <nav class="nav">
      <NuxtLink to="/">Главная</NuxtLink>
      <NuxtLink to="/catalog">Каталог</NuxtLink>
      <NuxtLink to="/favorites">Избранное</NuxtLink>
      <NuxtLink to="/about">О нас</NuxtLink>
    </nav>

    <div class="header-actions">
      <NuxtLink to="/cart" class="cart-link">
        🛒 Корзина
        <!-- 👇 Показывать количество товаров если есть -->
        <span v-if="cartStore.getItemsCount() > 0" class="badge">
          {{ cartStore.getItemsCount() }}
        </span>
        <!-- 👆 Показывать количество товаров если есть -->
      </NuxtLink>
      
      <NuxtLink to="/account" class="account-link">👤 Аккаунт</NuxtLink>
    </div>
  </header>
</template>

<!-- =========================================================== -->
<!-- БЫСТРЫЕ КОМАНДЫ -->
<!-- =========================================================== -->

<!-- Способ 1: Получить информацию о корзине в любом компоненте -->
<script setup lang="ts">
const cartStore = useCartStore()

// Все товары в корзине
console.log(cartStore.cartItems)

// Количество товаров
const itemCount = cartStore.getItemsCount()

// Сумма корзины
const totalPrice = cartStore.getTotalPrice()

// Количество конкретного товара
const quantity = cartStore.getQuantity(productId)

// Удалить товар
cartStore.removeFromCart(productId)

// Обновить количество
cartStore.updateQuantity(productId, 5)

// Очистить корзину
cartStore.clearCart()

// Загрузить корзину с сервера
await cartStore.restore(email)
</script>

<!-- =========================================================== -->
<!-- СПОСОБ 2: Использовать в шаблоне -->

<template>
  <!-- Показать количество товаров -->
  <span>{{ cartStore.getItemsCount() }} товаров в корзине</span>

  <!-- Показать сумму -->
  <span>Итого: {{ cartStore.getTotalPrice() }} ₽</span>

  <!-- Перебрать все товары -->
  <div v-for="item in cartStore.cartItems" :key="item.product.id">
    <h3>{{ item.product.name }}</h3>
    <p>{{ item.quantity }} × {{ item.product.price }} ₽</p>
  </div>

  <!-- Показать кнопку корзины с количеством -->
  <button class="cart-btn">
    🛒 Корзина
    <span v-if="cartStore.getItemsCount() > 0" class="badge">
      {{ cartStore.getItemsCount() }}
    </span>
  </button>
</template>

<!-- =========================================================== -->
<!-- СПОСОБ 3: Минималистичная страница корзины -->

<template>
  <div class="cart-minimal">
    <h1>Ваша корзина</h1>

    <div v-if="!cartStore.cartItems.length" class="empty">
      Корзина пуста
    </div>

    <div v-else>
      <div class="items">
        <div v-for="item in cartStore.cartItems" :key="item.product.id" class="item">
          {{ item.product.name }} — {{ item.quantity }} шт. — {{ item.product.price * item.quantity }} ₽
          <button @click="cartStore.removeFromCart(item.product.id)">Удалить</button>
        </div>
      </div>

      <div class="total">
        Итого: {{ cartStore.getTotalPrice() }} ₽
      </div>

      <button class="checkout-btn">Оформить заказ</button>
    </div>
  </div>
</template>

<!-- =========================================================== -->
<!-- СПОСОБ 4: Ошибки и решения -->

<script setup lang="ts">
// ❌ ОШИБКА: забыли передать product пропс
// <AddToCart /> ← НЕПРАВИЛЬНО!

// ✅ ПРАВИЛЬНО:
// <AddToCart :product="product" />

// ❌ ОШИБКА: забыли вызвать restore при загрузке
// onMounted(() => { /* ничего */ })

// ✅ ПРАВИЛЬНО:
// onMounted(async () => {
//   if (authStore.email) {
//     await cartStore.restore(authStore.email)
//   }
// })

// ❌ ОШИБКА: не проверили авторизацию перед checkout
// async function checkout() {
//   // НЕПРАВИЛЬНО! Может быть null
//   const response = await $fetch('/api/orders/auth', ...)
// }

// ✅ ПРАВИЛЬНО:
// async function checkout() {
//   if (!authStore.token) {
//     await router.push('/auth/login')
//     return
//   }
//   // Теперь безопасно
//   const response = await $fetch('/api/orders/auth', ...)
// }
</script>

<!-- =========================================================== -->
<!-- СПОСОБ 5: Условное отображение в зависимости от статуса -->

<template>
  <!-- Показать AddToCart только если товар не в корзине -->
  <AddToCart v-if="!cartStore.getQuantity(product.id)" :product="product" />

  <!-- Или показать информацию если уже в корзине -->
  <div v-else class="in-cart">
    В корзине: {{ cartStore.getQuantity(product.id) }} шт.
    <button @click="cartStore.removeFromCart(product.id)">Удалить из корзины</button>
  </div>
</template>

<!-- =========================================================== -->
<!-- ГОТОВО К ИСПОЛЬЗОВАНИЮ! -->
<!-- =========================================================== -->

<!-- Просто скопируйте нужные части и вставьте в ваши файлы -->
<!-- Система полностью готова к работе! 🚀 -->
