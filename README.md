# panbot - A TODOlist Discord bot

## Getting started

Create an .env file with the following variables
```dotenv
BOT_TOKEN="YOUR_BOT_TOKEN"
```

After that is done, if you have docker feel free to run `make docker-up`, otherwise you can simply use `make run` to run
the bot locally. The Go version required is `1.22`.

Once you are done with the bot, you can simply terminate the process and if you decided to use docker you can use
`make docker-down`.