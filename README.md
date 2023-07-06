# Otto

[![Coverage](https://coveralls.io/repos/github/Vico1993/Otto/badge.svg?branch=main)](https://coveralls.io/github/Vico1993/Otto?branch=main)

Otto is a bot designed to help you stay up-to-date with the latest news by monitoring RSS feeds and sending notifications via Telegram.

This part of the code only contains the API with the gestion of the Mongo DB

## Table of Contents

-   [Getting Started](#getting-started)
-   [Usage](#usage)
-   [Contributing](#contributing)
-   [License](#license)

## Getting Started

To get started with Otto, clone the repository to your local machine:

```sh
git clone https://github.com/Vico1993/Otto.git
cd Otto
```

## Prerequisites

Make sure you have the following tools installed on your machine:

-   Go (at least version 1.20)
-   A valid Telegram bot:
    -   look at the great [bot father](https://core.telegram.org/bots)
-   Setup an `.env` file

```sh
# Bot token given by the @BotFather
TELEGRAM_BOT_TOKEN=<TELEGRAM_BOT>
# Used by the bot to know which chat to interact with
TELEGRAM_USER_CHAT_ID=<CHAT_ID>

#MONGO srv url you want to use with username + password for the connection
MONGO_URI=<URI>
```

## Installing

To install Otto, run the following command:

```sh
make ensure_deps
```

## Running Tests

To run tests, use the following command:

```sh
make test
```

## TIPS

To make sure it's easy to build, I use: gow. Once install:

```sh
make watch
```

## Usage

To use Otto, run the following command:

```sh
make build && ./bin/bot
```

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) file for more information.

## License

This project is licensed under the [LICENSE](./LICENSE) file in the root directory of this repository.
