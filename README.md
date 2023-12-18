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
`


## Debugging deployments

1. Edit `EchoRelay/Dockerfile` to build "Debug" (instead of "Release")
2. Create an ssh tunnel to `localhost:2375` (docker daemon)
3. Attach to process..., select "*Docker (Linux Container)*, click *Find...*
4. enter `localhost:2735`, hit *Enter*
5. select `echonakama-e`
8. Click OK
