import { useAuthStore } from "~/state/auth.state";

export default defineNuxtRouteMiddleware(() => {
  const auth = useAuthStore();

  if (!auth.token) {
    return navigateTo("/auth/login");
  }

  if (auth.role !== "admin") {
    return navigateTo("/");
  }
});
