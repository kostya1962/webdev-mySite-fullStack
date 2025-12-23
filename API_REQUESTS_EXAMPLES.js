// Примеры HTTP запросов для тестирования API корзины
// Используйте Postman, curl или любой другой инструмент

// ============================================================
// 1. ДОБАВИТЬ ТОВАР В КОРЗИНУ
// ============================================================

// cURL
curl -X POST http://localhost:3000/api/cart \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "productID": 1,
    "quantity": 2
  }'

// JavaScript Fetch
fetch('http://localhost:3000/api/cart', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    productID: 1,
    quantity: 2
  })
})
  .then(res => res.json())
  .then(data => console.log('Добавлено в корзину:', data))
  .catch(err => console.error('Ошибка:', err))

// JavaScript Axios
const axios = require('axios');

axios.post('http://localhost:3000/api/cart', {
  email: 'user@example.com',
  productID: 1,
  quantity: 2
})
  .then(response => console.log('Ответ:', response.data))
  .catch(error => console.error('Ошибка:', error))

// ============================================================
// 2. ПОЛУЧИТЬ КОРЗИНУ ПОЛЬЗОВАТЕЛЯ
// ============================================================

// cURL
curl 'http://localhost:3000/api/cart?email=user@example.com'

// JavaScript Fetch
fetch('http://localhost:3000/api/cart?email=user@example.com')
  .then(res => res.json())
  .then(data => console.log('Корзина:', data))
  .catch(err => console.error('Ошибка:', err))

// Response Example:
/*
[
  {
    "product": {
      "id": 1,
      "name": "Кольцо золотое",
      "price": 12999,
      "short_description": "Красивое кольцо",
      "long_description": "Полное описание товара...",
      "sku": "RING-001",
      "discount": 10,
      "images": ["ring1.jpg"],
      "category_id": 1,
      "category": {"id": 1, "name": "Украшения"},
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "quantity": 2
  }
]
*/

// ============================================================
// 3. УДАЛИТЬ ТОВАР ИЗ КОРЗИНЫ
// ============================================================

// cURL
curl -X DELETE http://localhost:3000/api/cart \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "productID": 1
  }'

// JavaScript Fetch
fetch('http://localhost:3000/api/cart', {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    productID: 1
  })
})
  .then(res => res.json())
  .then(data => console.log('Удалено из корзины:', data))
  .catch(err => console.error('Ошибка:', err))

// ============================================================
// 4. ПОЭТАПНОЕ ТЕСТИРОВАНИЕ
// ============================================================

// Шаг 1: Создать пользователя (если нет)
const userEmail = 'testuser@example.com'

// Шаг 2: Добавить товар в корзину
async function addToCart(productID, quantity) {
  const response = await fetch('http://localhost:3000/api/cart', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: userEmail,
      productID,
      quantity
    })
  })
  const data = await response.json()
  console.log('✅ Товар добавлен:', data)
  return data
}

// Шаг 3: Получить корзину
async function getCart() {
  const response = await fetch(`http://localhost:3000/api/cart?email=${userEmail}`)
  const data = await response.json()
  console.log('📦 Корзина:', data)
  return data
}

// Шаг 4: Увеличить количество товара
async function updateCartItem(productID, newQuantity) {
  const response = await fetch('http://localhost:3000/api/cart', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: userEmail,
      productID,
      quantity: newQuantity
    })
  })
  const data = await response.json()
  console.log('✅ Количество обновлено:', data)
  return data
}

// Шаг 5: Удалить товар
async function removeFromCart(productID) {
  const response = await fetch('http://localhost:3000/api/cart', {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: userEmail,
      productID
    })
  })
  const data = await response.json()
  console.log('🗑️  Товар удалён:', data)
  return data
}

// Полный тестовый сценарий
async function testCartAPI() {
  console.log('=== Начинаем тестирование API корзины ===')
  
  try {
    // 1. Добавить товар
    await addToCart(1, 2)
    
    // 2. Добавить ещё один товар
    await addToCart(2, 1)
    
    // 3. Получить полную корзину
    let cart = await getCart()
    
    // 4. Обновить количество первого товара
    await updateCartItem(1, 3)
    
    // 5. Получить обновлённую корзину
    cart = await getCart()
    
    // 6. Удалить товар
    await removeFromCart(2)
    
    // 7. Получить финальную корзину
    cart = await getCart()
    
    console.log('=== Тестирование завершено ✅ ===')
  } catch (error) {
    console.error('❌ Ошибка при тестировании:', error)
  }
}

// Запустить тест:
// testCartAPI()

// ============================================================
// 5. POSTMAN ПРИМЕРЫ
// ============================================================

/*
Импортируйте в Postman:

### POST /api/cart - Добавить товар
URL: http://localhost:3000/api/cart
Method: POST
Body (JSON):
{
  "email": "user@example.com",
  "productID": 1,
  "quantity": 2
}

### GET /api/cart - Получить корзину
URL: http://localhost:3000/api/cart?email=user@example.com
Method: GET

### DELETE /api/cart - Удалить из корзины
URL: http://localhost:3000/api/cart
Method: DELETE
Body (JSON):
{
  "email": "user@example.com",
  "productID": 1
}
*/

// ============================================================
// 6. ИНТЕГРАЦИЯ С FRONTEND (Nuxt)
// ============================================================

// В компоненте Vue/Nuxt используйте store:

/*
<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface'

const props = defineProps<{ product: Product }>()
const cartStore = useCartStore()

function handleAddToCart(quantity: number) {
  // Это автоматически делает POST запрос
  cartStore.addToCart(props.product, quantity)
}
</script>
*/

// ============================================================
// 7. ОБРАБОТКА ОШИБОК
// ============================================================

// Возможные ошибки API:

/*
❌ 400 Bad Request
{
  "error": "Invalid request body"
}

❌ 404 Not Found
{
  "error": "User not found"
}
или
{
  "error": "Product not found"
}

❌ 500 Internal Server Error
{
  "error": "Failed to add to cart"
}
*/

// Обработка в JavaScript:

async function addToCartWithErrorHandling(email, productID, quantity) {
  try {
    const response = await fetch('http://localhost:3000/api/cart', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, productID, quantity })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error)
    }

    const data = await response.json()
    console.log('✅ Успешно:', data)
    return data
  } catch (error) {
    console.error('❌ Ошибка:', error.message)
    // Показать пользователю уведомление об ошибке
  }
}

// ============================================================
// 8. ПРОВЕРКА В БРАУЗЕРЕ (DevTools)
// ============================================================

/*
1. Откройте F12 → Console
2. Выполните:

const cartStore = useCartStore()
const authStore = useAuthStore()

// Проверить токен авторизации
console.log('Токен:', authStore.token)
console.log('Email:', authStore.email)

// Проверить корзину
console.log('Товары в корзине:', cartStore.cartItems)
console.log('Количество товаров:', cartStore.getItemsCount())
console.log('Сумма корзины:', cartStore.getTotalPrice())

// Добавить товар вручную
cartStore.addToCart({
  id: 1,
  name: 'Тестовый товар',
  price: 1000,
  short_description: 'Описание',
  long_description: 'Полное описание',
  sku: 'TEST-001',
  discount: 0,
  images: [],
  category_id: 1,
  category: { id: 1, name: 'Категория' },
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString()
}, 2)

// Проверить в Network вкладке
// Должен быть запрос: POST /api/cart
// Статус: 200 OK
*/

// ============================================================
// 9. СЦЕНАРИЙ ИСПОЛЬЗОВАНИЯ В PRODUCTION
// ============================================================

/*
// На странице товара:
<template>
  <div class="product">
    <h1>{{ product.name }}</h1>
    <p>{{ product.price }} ₽</p>
    <AddToCart :product="product" />
  </div>
</template>

// На странице корзины:
<template>
  <div v-for="item in cartStore.cartItems" :key="item.product.id" class="cart-item">
    <h3>{{ item.product.name }}</h3>
    <p>{{ item.quantity }} × {{ item.product.price }} ₽</p>
    <button @click="cartStore.removeFromCart(item.product.id)">Удалить</button>
  </div>
  <p>Итого: {{ cartStore.getTotalPrice() }} ₽</p>
</template>

// В header (для показа количества товаров):
<div class="cart-badge">{{ cartStore.getItemsCount() }}</div>
*/
