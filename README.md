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
MAX_REQUESTS_PER_HOUR=1 - a naive way of not getting super spammed by bots or malicious requests
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

Here is how you can download locally the database.sqlite file:

```shell
curl -X POST \
  'localhost:4500/download-db' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "telegramApiToken": "your bot father api token"
}' --output backup_database.sqlite
```

Once you have a backup of your dababase locally you can delete all rows:

```shell
curl  -X DELETE \
  'localhost:4500/clear-db' \
  --header 'Accept: */*' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "telegramApiToken": "your bot father api token"
}'
```

**NOTE:** By default only one request per hour is allowed. That's because a contact form is usually completed once by one person. You can remove/modify the rate limiter from server.go and rebuild. 
You can also increase the number of requests in `MAX_REQUESTS_PER_HOUR` variable from `.env` file. 

## Why?

I needed a way to be notified if someone filled the contact form on the website and hit send.
