# nostr-ex

## Usage

* Rename .env.example to .env
* Change value of MQ_URL in .env for QabbitMQ
* Run docker-compose to start db
 ````
 docker compose up
 ````
* Run project at  ws://127.0.0.1:8800/ for default
* Client will connect to wss://relay.nekolicio.us/ for default relay


## Stop
 ````
 docker compose down
 ````
