export const useAuthStore = defineStore(
  "auth",
  () => {
    const token = ref<string | undefined>();
    const email = ref<string | undefined>();
    const role = ref<string | undefined>();

    function setToken(val: string) {
      token.value = val;
    }

    function setEmail(val: string) {
      email.value = val;
    }

    function setRole(val: string) {
      role.value = val;
    }

    function clearToken() {
      token.value = undefined;
    }

    function clearRole() {
      role.value = undefined;
    }

    return {
      token,
      email,
      role,
      setToken,
      setEmail,
      setRole,
      clearToken,
      clearRole,
    };
  },
  {
    persist: true,
  }
);
