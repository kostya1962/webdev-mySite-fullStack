export function useAPIimage() {
  const config = useRuntimeConfig();
  return config.public.imageurl;
}
