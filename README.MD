# RUNING TESTS
    go test -v ./...

# CURLS
curl --location --request POST 'localhost:6980/station/discover'
curl --location --request POST '192.168.1.117:8080/station/discover'
curl --location --request POST '192.168.1.117:8080/station/handshake'

# SERVERS
## VAIO SERVER
    Conexion por eth
    ssh mati@192.168.1.117
