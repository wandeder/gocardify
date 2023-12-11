# GoCardify

gocardify is a simple Go program designed to facilitate the transfer of messages from a Telegram chat to a RabbitMQ queue. A secondary worker then processes these messages to generate cards on the AnkiWeb site.

## Features

- **Telegram Integration:** Connects to a Telegram chat and retrieves messages.
- **RabbitMQ Queue:** Sends messages to a RabbitMQ queue for further processing.
- **AnkiWeb Card Generation:** A secondary worker consumes messages from the queue and generates cards on AnkiWeb.

## Installation

1. Clone the repository:
   `git clone https://github.com/yourusername/GoCardify.git`
2. Install dependencies:
   `go get -u github.com/yourusername/dependency`
3. Build the program:
   `go build` 

## Usage

1. Set up your Telegram bot and RabbitMQ server, sign up AnkiWeb.
2. Configure the program with your Telegram bot token and AnkiWeb credentials in the .env file.
3. Run the program: `./gocardify`.
4. Add your Telegram bot into chat and start to sent messages.
5. Monitor the AnkiWeb site for the generated cards.
6. Study it!

## License
This project is licensed under the MIT License - see the LICENSE file for details.
    