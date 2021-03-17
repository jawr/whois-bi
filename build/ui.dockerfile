FROM node AS builder

WORKDIR /build

COPY package.json ./
COPY rollup.config.js ./
COPY public/ ./public
COPY src/ ./src

RUN npm install
RUN npm run build

FROM nginx:1.19.0
COPY ../build/ui.nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /build/public/ /usr/share/nginx/html
