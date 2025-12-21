export function useAPI() {
  const config = useRuntimeConfig();
  return config.public.apiurl;
}
