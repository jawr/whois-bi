FROM node 

WORKDIR /build

COPY package.json ./
COPY public ./public
COPY src ./src

RUN npm install

ENTRYPOINT npm start
