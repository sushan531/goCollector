# Collector

A dockerised collector service written in go. 

### Third party tools 
1. Redis
2. ZeroMQ

### System requirements
1. A locally running redis server in the host machine at default redis port.
2. Docker and docker compose installed 
3. Access to internet to install the third party libraries.

### Config file:
All the configuration for the service is in this file in three sections.
1. RedisConfig 
   1. host `Redis host`
   2. port `Redis port to connect`
2. RedisStream 
   1. stream  `Stream topic "Sth similar to kafka topic"`
   2. dummy `boolean for generating dummy value's and storing to redis`
3. ZmqConfig 
   1. connType `Conn type(udp|tcp)`
   2. host `zmq host machine`
   3. port `zmq port to bind`


### USAGE

1. Clone this repo in local system
2. Change the config.yaml file with desired configuration.
3. Run `docker compose up --build`