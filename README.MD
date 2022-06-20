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
### Compilar
    go build -o ubuntu-iot-hub-api .

### Configurar como servicio
    sudo touch /lib/systemd/system/iot_hub_api.service

    sudo vim /lib/systemd/system/iot_hub_api.service

description=IOT HUB API BKND
After=network.target
[Service]
Type=simple
User=root
LimitNOFILE=1024
Restart=on-failure
RestartSec=10
startLimitIntervalSec=60
Environment="IOTENV=prod"
WorkingDirectory=/repository/iot-hub-api
ExecStart=/repository/iot-hub-api/ubuntu-iot-hub-api
PermissionsStartOnly=true
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=IOT-HUB-API
[Install]
WantedBy=multi-user.target

Habilitar el servicio
    sudo chmod 644 /lib/systemd/system/iot_hub_api.service
    sudo systemctl enable iot_hub_api

Si cambia el file de config del servicio:
    systemctl daemon-reload
    systemctl restart iot_hub_api

Ver el estado:
systemctl status iot_hub_api

Ver los logs:
tail -f -n 20 /var/log/syslog



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




