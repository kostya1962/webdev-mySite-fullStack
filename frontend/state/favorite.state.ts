import { useAuthStore } from "./auth.state";

export const useFavoriteStore = defineStore(
  "favorites",
  () => {
    const authStore = useAuthStore();
    const favoriteIDs = ref<number[]>([]);

    function toggleFavorite(id: number) {
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
      return favoriteIDs.value.find((f) => f == id);
    }

    async function save() {
      await $fetch<{ success: boolean }>("/api/favorites", {
        method: "POST",
        body: {
          email: authStore.email,
          productIDs: favoriteIDs.value,
        },
      });
    }

    async function restore(email: string) {
      const data = await $fetch<number[]>("/api/favorites", {
        query: {
          email: email,
        },
      });
      favoriteIDs.value = data;
    }

    return {
      favoriteIDs,
      toggleFavorite,
      isFavorite,
      restore,
    };
  },
  {
    persist: true,
  }
);
