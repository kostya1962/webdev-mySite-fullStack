export const useAuthStore = defineStore(
  "auth",
  () => {
    const token = ref<string | undefined>();
    const email = ref<string | undefined>();

    function setToken(val: string) {
      token.value = val;
    }

    function setEmail(val: string) {
      email.value = val;
    }

    function clearToken() {
      token.value = undefined;
    }

    return { token, email, setToken, setEmail, clearToken };
  },
  {
    persist: true,
  }
);
