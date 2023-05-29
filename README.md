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

* Post API to add relay

POST http://{{serve}}/relays

EX:POST http://localhost:8800/relays

Example of API body when POST http://{{serve}}/relays

    {
        "url":"wss://relay.nekolicio.us/"
    }    

* Get API to get list of relay

GET http://{{serve}}/relays

## Stop database which build from docker
 ````
 docker compose down
 ````
