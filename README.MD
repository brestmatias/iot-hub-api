# RUNING TESTS
    go test -v ./...

# CURLS
curl --location --request POST 'localhost:6980/station/discover'
curl --location --request POST '192.168.1.117:8080/station/discover'
curl --location --request POST '192.168.1.117:8080/station/handshake'
curl --location --request POST 'localhost:6980/dispatcher/reloadtasks'

# SERVERS
## VAIO SERVER
    Conexion por eth
    ssh mati@192.168.1.117


# STATION 
    Salidas:
   [TIPO][ID]
   Tipos de salida
        R = relee
        T = transistor
        L = led
   values:
        00 = off
        01 = on

        10 = manual off
        11 = manual on


## CURLS - MQTT

### API
curl --location --request GET 'http://192.168.1.116/beacon
curl --location --request GET 'http://192.168.1.116/ping

curl --location -g --request PUT '{{sta}}/station' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "STA02",
    "broker":"192.168.1.101"
}'

### MQTT

Topic commands:

    command/%s
    body:
    {"interface":"L01", "value":0}




