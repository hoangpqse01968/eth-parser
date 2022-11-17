# Ethereum Parser

## Solution

- A cronjob to fetch the latest block and add it to database, and notify to all subscribes (which are goroutines,
  responsible for extract transactions of an address)

- Create a goroutine to handle an address. The reason because when we subscribe new address, it may take so long to
  fetch the history transactions -> Therefore, we can use a separate goroutine to avoid block other addresses update its
  transactions

## Exposed HTTP API

- `/get-current-block-number` return latest block number in database
- `/subscribe?address={{address}}` register an address as subscriber
- `/unsubscribe?address={{address}}` unregister an address from subscriber list
- `/get-transactions?address={{address}}` get parsed transactions of an address

Those endpoints are already exported to `etth-parser.postman_collection.json`, so can import it for testing purpose.

## How to run

```
$ make build
$ ./bin/parser
```

## Improvements

- Change database to real DB so can store the parsed result
- Retrieving data for address with too many history transactions can be problem -> should get block from database by
  batch and process
