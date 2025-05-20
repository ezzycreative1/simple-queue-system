# --- Build Stage ---
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY . .

RUN npm run build


# --- NGINX Stage ---
FROM nginx:stable-alpine

# Hapus default config nginx (optional tapi recommended)
RUN rm -rf /etc/nginx/conf.d

# Copy custom config (jika punya)
# COPY nginx.conf /etc/nginx/nginx.conf

# Copy build hasil React ke folder nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
