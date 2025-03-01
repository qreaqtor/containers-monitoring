FROM node:22.14.0 AS builder
WORKDIR /app
COPY . .
RUN npm install
ENV VITE_BACKEND_URL=localhost:50055
RUN npm run build

FROM syntaqx/serve:latest
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package*.json ./

EXPOSE 8080
CMD ["serve", "./dist"]
