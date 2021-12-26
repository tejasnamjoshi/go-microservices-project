FROM node:14

RUN mkdir /app
WORKDIR /app

COPY ./next-todo-tailwind/package*.json .
COPY ./next-todo-tailwind/yarn.lock .
RUN yarn install

COPY ./next-todo-tailwind .

RUN yarn build
EXPOSE 3000

CMD ["yarn", "start"]