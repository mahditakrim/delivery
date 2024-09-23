# Delivery Simulator

## API info

Server listens on port 8080  
core route: localhost:8080/core/order  
delivery webhook: localhost:8080/delivery/webhook  

## Screen logger

there is 5 stages for each delivery order  

1. pending (order has been accepted in system)
2. searching (order is sent to 3pl for searching driver)
3. found (3pl has been found a driver)
4. not found (3pl could not find a driver in this attempt)
5. delivered (driver sent delivery to customer)

### logger prints number of orders in each of above statges in interval mode  

## Simulation delay

There is some `time.Sleep` commands in 3pl API for delay simulation.  

## Seeder

You can find seeder app in ./seeder directory.  
It sends concurrent requests every 200 milliseconds.  
The timestamp range of delivery requests are within now till 2h from now.  

## How to run

You can `cd` to either seeder or delivery simulator project directory and simply run `go run .`  
`ctrl+c` will trigger graceful shutdown in delivery simulator.
