# 🎨 Визуальный гайд компонента AddToCart

## Внешний вид компонента

### Состояние 1️⃣: Кнопка "Добавить в корзину" (неавторизованный или начальное состояние)

```
┌─────────────────────────────────────┐
│                                     │
│   🖱️  Добавить в корзину            │
│                                     │
│   (Черный фон, белый текст)         │
│   (Размер: полная ширина)           │
│                                     │
└─────────────────────────────────────┘
```

При клике:
- Если **не авторизован** → Редирект на `/auth/login`
- Если **авторизован** → Переход в состояние 2️⃣

---

### Состояние 2️⃣: Счётчик + кнопки (авторизованный пользователь)

```
┌──────────┬──────────────┬──────────┐
│          │              │          │
│ Удалить  │  Счётчик     │ Добавить │
│          │              │          │
│ (Серая)  │  [1][2][3]   │ (Чёрная) │
│ (Граница)│  Количество  │(Граница) │
│          │              │          │
└──────────┴──────────────┴──────────┘
```

Компоненты:
- **Кнопка "Удалить"** (слева)
  - Фон: прозрачный
  - Граница: черная
  - При наведении: серый фон
  - При клике: возврат в состояние 1️⃣

- **CounterFiled** (в центре)
  - Управление количеством товара
  - Передает значение обратно в родительский компонент

- **Кнопка "Добавить"** (справа)
  - Фон: черный
  - Текст: белый
  - При наведении: инверсия (белый фон, черный текст)
  - При клике: добавляет товар в корзину и возвращается в состояние 1️⃣

---

## Поток взаимодействия пользователя

```
┌──────────────────────────────────────────────────────────────┐
│  Пользователь смотрит товар на странице /catalog/sup-[id]    │
└──────────┬───────────────────────────────────────────────────┘
           │
           ├─→ Нажимает "Добавить в корзину"
           │
           ├─→ AddToCart.vue: Проверка авторизации
           │   │
           │   ├─ Если token НЕ существует
           │   │  └─→ Редирект на /auth/login 🔐
           │   │
           │   └─ Если token СУЩЕСТВУЕТ
           │      └─→ isAdded.value = true
           │         └─→ Показ CounterFiled + кнопки
           │
           ├─→ Пользователь вводит количество в CounterFiled
           │   └─→ quantity.value = X
           │
           ├─→ Нажимает "Удалить" (опционально)
           │   └─→ isAdded.value = false
           │      └─→ Возврат к состоянию 1️⃣
           │
           └─→ Нажимает "Добавить"
               ├─→ cartStore.addToCart(product, quantity)
               │  │
               │  ├─→ Добавляет в локальный store
               │  ├─→ Сохраняет в localStorage
               │  └─→ POST /api/cart (на сервер)
               │
               └─→ isAdded.value = false
                  └─→ quantity.value = 1
                     └─→ Возврат к состоянию 1️⃣ ✅
```

---

## Проверка в браузере (DevTools)

### Вкладка Console (F12)
```javascript
// Проверить состояние корзины:
cartStore = useCartStore()
cartStore.cartItems          // Все товары в корзине
cartStore.getItemsCount()    // Количество товаров
cartStore.getTotalPrice()    // Сумма корзины

// Проверить авторизацию:
authStore = useAuthStore()
authStore.token              // Токен (должен существовать)
authStore.email              // Email пользователя
```

### Вкладка Network
```
Должны быть запросы:
- POST /api/cart          (при добавлении товара)
- GET  /api/cart          (при загрузке корзины)
- DELETE /api/cart        (при удалении товара)

Статус: 200 OK ✅
```

### Вкладка Application → LocalStorage
```
Ключ: cart
Содержимое: {"cartItems": [{"product": {...}, "quantity": 2}]}
```

---

## Интеграция в существующие компоненты

### На странице товара

```vue
<script setup lang="ts">
import type { Product } from '~/interfaces/product.interface'

// Получить данные товара
const { data: product } = await useFetch<Product>(`/api/products/${id}`)
</script>

<template>
  <div class="product">
    <h1>{{ product.name }}</h1>
    <p>{{ product.price }} ₽</p>
    
    <!-- 👇 ДОБАВЬТЕ ВОТ ЭТО 👇 -->
    <AddToCart :product="product" />
    <!-- 👆 ДОБАВЬТЕ ВОТ ЭТО 👆 -->
  </div>
</template>
```

### На странице корзины

```vue
<script setup lang="ts">
const cartStore = useCartStore()

// Восстановить корзину при загрузке
onMounted(async () => {
  if (authStore.email) {
    await cartStore.restore(authStore.email)
  }
})
</script>

<template>
  <div class="cart">
    <div v-for="item in cartStore.cartItems" :key="item.product.id">
      {{ item.product.name }} × {{ item.quantity }}
    </div>
    
    <p>Итого: {{ cartStore.getTotalPrice() }} ₽</p>
  </div>
</template>
```

---

## CSS Переменные (Customization)

Компонент использует CSS переменные, которые вы можете переопределить:

```css
:root {
  --color-black: #000000;
  --color-gray: #f5f5f5;
}
```

Переопределение в вашем CSS:
```css
:root {
  --color-black: #1a1a1a;
  --color-gray: #e0e0e0;
}
```

---

## Примеры данных

### Product (товар)
```json
{
  "id": 1,
  "name": "Кольцо из золота",
  "price": 12999,
  "short_description": "Красивое кольцо",
  "long_description": "Полное описание...",
  "sku": "RING-001",
  "discount": 10,
  "images": ["ring1.jpg", "ring2.jpg"],
  "category_id": 1,
  "category": {"id": 1, "name": "Украшения"},
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### CartItem (элемент корзины)
```json
{
  "product": { /* полные данные товара */ },
  "quantity": 2
}
```

### API Request (добавление в корзину)
```json
{
  "email": "user@example.com",
  "productID": 1,
  "quantity": 2
}
```

### API Response (успешно)
```json
{
  "success": true,
  "message": "Item added to cart"
}
```

---

## Проверочный список перед продакшеном

- [ ] Компонент AddToCart.vue существует в components/
- [ ] Store cart.state.ts существует в state/
- [ ] Интерфейсы cart.interface.ts созданы
- [ ] Backend обработчики cart.go добавлены
- [ ] Маршруты в main.go прописаны
- [ ] Таблица cart_items создаётся при инициализации БД
- [ ] Проверено добавление товара (не авторизованный)
- [ ] Проверено добавление товара (авторизованный)
- [ ] Проверено удаление товара из корзины
- [ ] Проверена синхронизация с сервером
- [ ] Проверена загрузка корзины после авторизации
- [ ] Проверены все стили на разных экранах

---

## Быстрые ссылки на файлы

- 📄 [AddToCart компонент](../frontend/components/AddToCart.vue)
- 📦 [Cart Store](../frontend/state/cart.state.ts)
- 🔗 [Cart Интерфейсы](../frontend/interfaces/cart.interface.ts)
- 🌐 [Backend обработчики](../backend/handlers/cart.go)
- 📚 [Полная документация](./CART_SYSTEM_DOCUMENTATION.md)
- ✅ [Чеклист внедрения](./IMPLEMENTATION_CHECKLIST.md)

