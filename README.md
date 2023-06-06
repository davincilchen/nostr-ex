# nostr-ex

## Usage

* Rename .env.example to .env
* Set value of AUTOMIGRATE=1 in .env to migrate Database table at first time
* Change value of MQ_URL in .env for RabbitMQ
* Run docker-compose to start db
 ````
 docker compose up
 ````
* Run project at  http://localhost:8800/ for default
* Aggregator will connect to wss://relay.nekolicio.us/ for default relay

* Post API to add relay

    * POST http://{{serve}}/relays
    * EX:POST http://localhost:8800/relays

* Example of API body when POST http://{{serve}}/relays
```
    {
        "url":"wss://relay.nekolicio.us/"
    }    
```
* Get API to get list of relays
    
    * GET http://{{serve}}/relays


* view:  http://localhost:8800/watcher for default
* view:  http://localhost:2223/metrics for phase 5 to check metrics

## Stop database which is builded from docker compose
 ````
 docker compose down
 ````
