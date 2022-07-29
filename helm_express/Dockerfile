FROM node:8
# Create app directory
WORKDIR /usr/node/app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 8080
CMD ["npm", "start"]