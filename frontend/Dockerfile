FROM node:20-bullseye AS builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci --legacy-peer-deps
COPY frontend/ .
ENV NUXT_PUBLIC_APIURL=/api
ENV NUXT_PUBLIC_IMAGEURL=/
RUN npm run generate

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx/frontend.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]