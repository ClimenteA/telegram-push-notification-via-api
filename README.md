# Push Notifications API (Telegram)

Send from your server to your smartphone notifications using [Telegram Bot API](https://telegram.org/).
Messages are also saved in a sqlite database for later retrival. 
Using this API you can keep your Bot Api secrets a bit safer than just making them public in a static site.  

<p align="center">
  <img src="./pics/telegram_messages.jpeg" width="200">
</p>


## Quickstart

Search `how to make a telegram bot` and follow the steps to get the following information needed in the `.env` file:

```shell
TELEGRAM_API_TOKEN=BOTTOKEN:FROM:BotFather
CHAT_ID=number-chat-id-after-first-exchange-message
API_KEY=openssl rand -base64 33 | tr '+/' '-_' or something else
PORT=4500 - port on which the GoFiber server will start
```


- download the zip from Releases or build go binary (see Go code); 
- install docker on your VPS or laptop; 
- add your static site(s) or maybe you just use this as an another service; 
- modify Caddyfile domains and links to static sites (see Releases zip file);
- run `docker-compose build --no-cache` to build images;
- run `docker-compose up -d` to serve static sites in production;


Checkout `static-site\index.html` from this repo to see an example on how you can send a message from a contact form to telegram.

From any static site you just make a fetch post request to `/push-notification-to-telegram` with body:
```json
{
  "messageType": "contact",
  "email": "alin@gmail.com",
  "message": "some message here",
  "apikey": "YjqkpUhZX9MFxhelTTyzg6cbzN4KYu4pbROsyYP5yc"
}
```


To copy sqlite db from inside the container you can do the following:
- `docker-compose up -d --force-recreate`;
- `docker ps` - to get the container id of the push-api service;
```shell
CONTAINER ID   IMAGE              COMMAND                  CREATED         STATUS         PORTS                                                                                         NAMES
51a176bf8048 <<THIS   release-push-api   "/home/server"           2 minutes ago   Up 5 seconds   0.0.0.0:4500->4500/tcp, :::4500->4500/tcp                                                     telegram-push-notification-bot-api
f151059a2be1   caddy:2-alpine     "caddy run --config …"   17 hours ago    Up 5 seconds   0.0.0.0:80->80/tcp, :::80->80/tcp, 0.0.0.0:443->443/tcp, :::443->443/tcp, 443/udp, 2019/tcp   telegram-static-website-caddy-proxy
```

- `docker cp 51a176bf8048:/home/database.sqlite ./`;


## Why?

I needed a way to be notified if someone filled the contact form on the website and hit send.
