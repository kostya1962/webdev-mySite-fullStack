// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: {
    enabled: true,

    timeline: {
      enabled: true,
    },
  },

  components: [
    {
      path: "~/components",
      pathPrefix: false,
    },
  ],

  runtimeConfig: {
    token: "",
    public: {
      apiurl: process.env.NUXT_PUBLIC_APIURL,
      imageurl: process.env.NUXT_PUBLIC_IMAGEURL,
    },
  },

  icon: {
    customCollections: [
      {
        prefix: "icons",
        dir: "./assets/icons",
      },
    ],
  },

  app: {
    head: {
      title: "Shopper - магазин ювелирных украшений",
      titleTemplate: "%s | Shopper",
      htmlAttrs: {
        lang: "ru",
      },
      link: [
        {
          rel: "icon",
          type: "image/png",
          href: "/favicon-32x32.png",
        },
        {
          rel: "mainifest",
          href: "/site.webmanifest",
          crossorigin: "anonymous",
        },
        {
          rel: "stylesheet",
          href: "https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.min.css",
        },
      ],
    },
  },

  modules: [
    "@nuxt/eslint",
    "@nuxt/fonts",
    "@nuxt/image",
    "@nuxt/scripts",
    "@nuxt/icon",
    "@pinia/nuxt",
    "pinia-plugin-persistedstate/nuxt",
    "@nuxtjs/sitemap",
    "@nuxtjs/robots",
  ],

  robots: {
    disallow: ["/account", "/auth/login", "/auth/register"],
  },

  sitemap: {
    sources: ["/api/sitemap/urls"],
    defaults: {
      lastmod: new Date().toISOString(),
      priority: 0.5,
      changefreq: "daily",
    },
  },

  nitro: {
    storage: {
      db: {
        driver: "fs-lite",
        base: "./db",
      },
    },
  },
});
