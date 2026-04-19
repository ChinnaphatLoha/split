FROM node:20-alpine

WORKDIR /app

COPY gateway/package.json gateway/package-lock.json* ./
RUN npm install

COPY gateway/ .

# Copy proto files for gRPC client loading
COPY proto /app/proto

RUN npm run build

EXPOSE 4000

CMD ["npm", "start"]
