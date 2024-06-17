# Push Notifications API

Useful for server notifications. Send a `POST` request to `/push-notification-to-telegram` from your server and be notified on your telegram app.


# Quickstart


Go server binary can be created with this command:
```shell
GOOS=linux GOARCH=amd64 go build -o dist/server server.go
```
You can use the binary from `Releases` if you don't want to build it yourself.

Contents of `.env` file:

```shell
TELEGRAM_API_TOKEN=BOTTOKEN:FROM:BotFather
CHAT_ID=number-chat-id-after-first-exchange-message
API_KEY=openssl rand -base64 33 | tr '+/' '-_' or something else
PORT=4500 - port on which the GoFiber server will start
```

One way to get the chat_id you can find it [here](https://dev.to/climentea/push-notifications-from-server-with-telegram-bot-api-32b3).


Checkout `static-site\index.html` from this repo to see an example on how you can send a message from a contact form straight to telegram.

From any static site you just make a fetch post request to `/push-notification-to-telegram` with body:
```json
{
  "message": "some message here, even some string json",
  "apikey": "YjqkpUhZX9MFxhelTTyzg6cbzN4KYu4pbROsyYP5yc"
}
```

You can also use any server side language to make a post request (ex: another service in docker-compose.yml file).


# Why?

I need to be notified (for free) somehow that someone (or somebot) sent me a message thru my static sites contact forms. 