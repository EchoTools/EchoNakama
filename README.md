# EchoNakama


## Development Deployment

Move `local.yml.example` -> `local.yml`:
* Change all secrets
* change console username and password

Access via http://nakamaserverhost:7351/

## Production Deployment

### Encryption

Edit `nginx_conf.d/nakama.conf`
* Change `server_name` values to appropriate hostname


`docker-compose up`

Access via https://nakamaserverhost:7351/
