# Bakery

'Bakery' is a toy, distributed system for playing with Jaeger for instrumentation

## Overview

![Diagram of the general structure of the app](./doc/img/overview.png)

The app consists of:  
* a REST API 'reception' that accepts new orders and provides the status of orders  
* a queue that new orders are placed on  
* a queue consumer 'baker' that accepts orders off of the queue and bakes them
* a database that tracks orders and their status
