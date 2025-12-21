export default defineEventHandler(async (event) => {
  const { email } = getQuery<{ email: string }>(event);
  if (!email) {
    throw createError({
      statusCode: 400,
      statusMessage: "Email not defined in query",
    });
  }
  const response = await useStorage("db").getItem(email);
  return response;
});
