FROM node AS builder

WORKDIR /build

COPY package.json ./
COPY public/ ./public
COPY src/ ./src

RUN npm install
RUN npm run build

FROM nginx:1.19.0
COPY --from=builder /build/build/ /usr/share/nginx/html
