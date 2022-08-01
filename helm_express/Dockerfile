# FROM mongo-express:1.0.0-alpha.4
FROM node:16-alpine3.16
# Create app directory
WORKDIR /usr/node/app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 8080
CMD ["npx", "nodemon", "app.js"]