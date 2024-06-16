# API Wrapper for Telegram Bot Push Notifications

Useful for server notifications. Send a `POST` request to `/send-push-notification-to-telegram` from your server and be notified on your telegram app.


# Quickstart

Get the binary from releases or use the Dockerfile from this repo to add it to your docker-compose.yml file. 

Next to your binary create a file named `telegram.config.json` which will contain the following:

```json
{
    "TELEGRAM_API_TOKEN": "TOKEN:FROM-BOTFATHER",
    "CHAT_ID": 123123123, # a bit tricky to get your hands on
    "API_KEY": "generate an apikey with: openssl rand -base64 32",
    "PORT": 3000, # server will start at this port
    "BOT_NAME": "@my_bot", # optional 
    "BOT_SHORT_NAME": "@my_bot", # optional
    "BOT_URL": "https://t.me/my_bot" # optional,
}
```

Start server with ./server (if you are using the binary).
Now from your app send a post request like:

```shell
curl  -X POST \
  'localhost:3000/send-push-notification-to-telegram' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "message": "some message here, even some string json",
  "apikey": "YjqkpUhZX9MFxhelTTyzg6cbzN4KYu4pb/ROsyYP5yc="
}'
```

Adapt the curl POST request for the programing language your are using.


One way to get the chat_id you can find it [here](https://dev.to/climentea/push-notifications-from-server-with-telegram-bot-api-32b3).
