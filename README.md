# Push Notifications API (Telegram)

[Telegram](https://telegram.org/) is an alternative to WhatsApp Messenger which has a free api which we can use to send messages from the server.
This repo is useful for people who have one or more static sites with contact forms and want to be notified when someone sent a message thru that contact form. 


## Quickstart

This repo contains all you need to ship a ton of static sites on a small 5$ VPS. Checkout [bunjucks](https://github.com/ClimenteA/bunjucks) a simple static site generator (SSG) if you have to create custom multi-page static websites (any SSG will do). Of course, you could just use this as another rest api service in your app.

Google `how to make a telegram bot` and follow the steps there to get info needed for the `.env` file.

Fill the `.env` file:

```shell
TELEGRAM_API_TOKEN=BOTTOKEN:FROM:BotFather
CHAT_ID=number-chat-id-after-first-exchange-message
API_KEY=openssl rand -base64 33 | tr '+/' '-_' or something else
PORT=4500 - port on which the GoFiber server will start
```

One way to get the chat_id you can find it [here](https://dev.to/climentea/push-notifications-from-server-with-telegram-bot-api-32b3).


Download the zip from Releases, install docker on your VPS or laptop, add your static site, modify Caddyfile and run `docker-compose up -d`;


## Development

This app is a small Go/Fiber api which makes a post request to Telegram (that's it). 

- clone this repo;
- install [Go](https://go.dev/);
- run `make build` to build binary needed in Dockerfile;
- run `make dev` just to see if everything is working fine;
- run `make up` to serve the app in production;


Checkout `static-site\index.html` from this repo to see an example on how you can send a message from a contact form to telegram.

From any static site you just make a fetch post request to `/push-notification-to-telegram` with body:
```json
{
  "message": "some message here",
  "apikey": "YjqkpUhZX9MFxhelTTyzg6cbzN4KYu4pbROsyYP5yc"
}
```

Using this API you can keep your Bot Api secrets a bit safer than just making them public in a static site. 
