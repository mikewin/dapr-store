#!/bin/bash

dapr stop --app-id products
dapr stop --app-id orders
dapr stop --app-id api-gateway
docker rm -f api-gateway
pkill frontend-host
