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

    // Ensure cartItems is always an array
    const ensureArray = () => {
      if (!Array.isArray(cartItems.value)) {
        cartItems.value = [];
      }
    };

    // Добавить товар в корзину (или увеличить количество)
    function addToCart(product: Product, quantity: number) {
      // Инициализируем корзину если она пуста или null
      ensureArray();

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
      ensureArray();

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
      ensureArray();
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
      return (
        cartItems.value.find((item) => item.product.id === productId)
          ?.quantity ?? 0
      );
    }

    function getTotalPrice(): number {
      return cartItems.value.reduce(
        (total, item) => total + item.product.price * item.quantity,
        0
      );
    }

    function getItemsCount(): number {
      return cartItems.value.reduce((count, item) => count + item.quantity, 0);
    }

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

    async function restore(email: string) {
      try {
        const data = await $fetch<CartItem[]>(`${API_URL}/cart`, {
          query: {
            email: email,
          },
        });
        cartItems.value = Array.isArray(data) ? data : [];
      } catch (error) {
        console.error("Ошибка при загрузке корзины:", error);
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
