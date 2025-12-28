import { useAuthStore } from "./auth.state";

export const useFavoriteStore = defineStore(
  "favorites",
  () => {
    const authStore = useAuthStore();
    const favoriteIDs = ref<number[]>([]);

    // Ensure favoriteIDs is always an array
    const ensureArray = () => {
      if (!Array.isArray(favoriteIDs.value)) {
        favoriteIDs.value = [];
      }
    };

    function toggleFavorite(id: number) {
      ensureArray();
      if (!favoriteIDs.value.includes(id)) {
        favoriteIDs.value.push(id);
        if (authStore.email) {
          save();
        }
        return;
      }
      favoriteIDs.value = favoriteIDs.value.filter((favID) => favID !== id);
      if (authStore.email) {
        save();
      }
    }

    function isFavorite(id: number) {
      ensureArray();
      return favoriteIDs.value.find((f) => f == id);
    }

    function removeFavoriteId(id: number) {
      ensureArray();
      favoriteIDs.value = favoriteIDs.value.filter((favID) => favID !== id);
      if (authStore.email) {
        save();
      }
    }

    async function save() {
      ensureArray();
      await $fetch<{ success: boolean }>("/api/favorites", {
        method: "POST",
        body: {
          email: authStore.email,
          productIDs: favoriteIDs.value,
        },
      });
    }

    async function restore(email: string) {
      try {
        const data = await $fetch<number[]>("/api/favorites", {
          query: {
            email: email,
          },
        });
        favoriteIDs.value = Array.isArray(data) ? data : [];
      } catch (e) {
        // If backend favorites endpoint is missing or request fails, keep local favorites
        favoriteIDs.value = favoriteIDs.value || [];
      }
    }

    return {
      favoriteIDs,
      toggleFavorite,
      isFavorite,
      removeFavoriteId,
      restore,
    };
  },
  {
    persist: true,
  }
);
