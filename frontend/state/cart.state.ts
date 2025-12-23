import { useAuthStore } from "./auth.state";
import type { CartItem } from "~/interfaces/cart.interface";
import type { Product } from "~/interfaces/product.interface";
import { useAPI } from "~/composables/useAPI";

export const useCartStore = defineStore(
  "cart",
  () => {
    const authStore = useAuthStore();
    const cartItems = ref<CartItem[]>([]);
    const API_URL = useAPI();

    // Добавить товар в корзину (или увеличить количество)
    function addToCart(product: Product, quantity: number) {
      // Инициализируем корзину если она пуста или null
      if (!cartItems.value) {
        cartItems.value = [];
      }

      const existingItem = cartItems.value.find(
        (item) => item.product.id === product.id
      );

      if (existingItem) {
        existingItem.quantity += quantity;
      } else {
        cartItems.value.push({
          product,
          quantity,
        });
      }

      if (authStore.email) {
        save(product.id, quantity);
      }
    }

    // Обновить количество товара
    function updateQuantity(productId: number, quantity: number) {
      if (!cartItems.value) {
        cartItems.value = [];
        return;
      }

      const item = cartItems.value.find(
        (item) => item.product.id === productId
      );
      if (item) {
        item.quantity = quantity;
        if (authStore.email) {
          save(productId, quantity);
        }
      }
    }

    // Удалить товар из корзины
    function removeFromCart(productId: number) {
      if (!cartItems.value) {
        return;
      }

      cartItems.value = cartItems.value.filter(
        (item) => item.product.id !== productId
      );
      if (authStore.email) {
        removeFromServer(productId);
      }
    }

    // Очистить корзину
    function clearCart() {
      cartItems.value = [];
    }

    // Получить количество товара в корзине
    function getQuantity(productId: number): number {
      if (!cartItems.value) {
        return 0;
      }

      return (
        cartItems.value.find((item) => item.product.id === productId)
          ?.quantity ?? 0
      );
    }

    // Получить общую стоимость
    function getTotalPrice(): number {
      if (!cartItems.value) {
        return 0;
      }

      return cartItems.value.reduce(
        (total, item) => total + item.product.price * item.quantity,
        0
      );
    }

    // Получить количество товаров
    function getItemsCount(): number {
      if (!cartItems.value) {
        return 0;
      }

      return cartItems.value.reduce((count, item) => count + item.quantity, 0);
    }

    // Сохранить товар на сервере
    async function save(productId: number, quantity: number) {
      try {
        await $fetch<{ success: boolean }>(`${API_URL}/cart`, {
          method: "POST",
          body: {
            email: authStore.email,
            productID: productId,
            quantity: quantity,
          },
        });
      } catch (error) {
        console.error("Ошибка при добавлении в корзину:", error);
      }
    }

    // Удалить товар с сервера
    async function removeFromServer(productId: number) {
      try {
        await $fetch<{ success: boolean }>(`${API_URL}/cart`, {
          method: "DELETE",
          body: {
            email: authStore.email,
            productID: productId,
          },
        });
      } catch (error) {
        console.error("Ошибка при удалении из корзины:", error);
      }
    }

    // Загрузить корзину с сервера
    async function restore(email: string) {
      try {
        const data = await $fetch<CartItem[]>(`${API_URL}/cart`, {
          query: {
            email: email,
          },
        });
        cartItems.value = data && Array.isArray(data) ? data : [];
      } catch (error) {
        console.error("Ошибка при загрузке корзины:", error);
        // При ошибке инициализируем пустой массив
        if (!cartItems.value) {
          cartItems.value = [];
        }
      }
    }

    return {
      cartItems,
      addToCart,
      updateQuantity,
      removeFromCart,
      clearCart,
      getQuantity,
      getTotalPrice,
      getItemsCount,
      restore,
    };
  },
  {
    persist: true,
  }
);
