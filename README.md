
# Gopher slack bot

Slack bot Go app, which be deployed to Heroku.

## Prepare Locally

```sh
git clone https://github.com/masumomo/gopher-slackbot.git
cd gopher-slackbot
go build -o bin/gopher-slackbot -v . # or `go build -o bin/gopher-slackbot.exe -v .` in git bash
```

## Deploying to Heroku

### Set up environmental valuable

You have to set up these env valuable below
- SLACK_BOT_TOKEN
- SLACK_VERIFY_TOKEN
- CHANNEL_ID

### Deploy

```sh
heroku create
git push heroku master
heroku open
```
or

```sh
heroku logs -tail
```


## Endpoints

- `/events`
- `/interactions`

