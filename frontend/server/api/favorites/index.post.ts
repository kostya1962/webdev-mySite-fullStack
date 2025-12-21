import type { CreateFavorite } from "~/interfaces/favorites.interface";

export default defineEventHandler(async (event) => {
  const { email, productIDs } = await readBody<CreateFavorite>(event);
  await useStorage("db").setItem(email, productIDs);
  setResponseStatus(event, 201);
  return {
    success: true,
  };
});
