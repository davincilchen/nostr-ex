# Distributed Systems Project 

## Prerequisites

- Go: v1.20.4

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
    * 目前還在研究如何整合到gafana，暫時只能用以上頁面顯示

## Stop database which is builded from docker compose
 ````
 docker compose down
 ````

 ## Questions

 #### **Why did you choose these 5 metrics?**
- queue size: (搜尋queue_size)因為存在生產者消費者問題，觀察queue size的大小可以知道此一關聯。但在rabbit queue已有dashboard可以觀看這些指標。但如果可以集中式控管應該更為方便，只是可能可以用直接接入rabbit queue來取得這些數據
- database: 如果接入大量relay會有大量的寫入，如果有很多連線連接到watcher會有大量的讀取，因此觀察資料庫的運作狀況。目前用otelsql接入，但還沒有研究如何導出到UI顯示。
- duration: (搜尋duration_in_milliseconds)觀察函數的運作時間是否在合理範圍

- success counter: (搜尋success_counter) 觀察函數運作成功的數量，可以藉由失敗和成功的比例評估是否合理


- fail counter: (搜尋fail_counter) 觀察函數運作失敗的數量，可以藉由失敗和成功的比例評估是否合理
