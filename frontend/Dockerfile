# 前端多阶段构建：Vite 打包 -> Nginx 托管
FROM node:20-alpine AS builder
WORKDIR /app
ENV npm_config_registry=https://registry.npmmirror.com
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
