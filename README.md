# Whois.bi
Tool for monitoring domains. Scan for common DNS records, or manually add them.
Periodically check for changes on the records. Periodically Check for Whois
changes. 

## Stack
Currently uses PostgreSQL to store state and RabbitMQ to queue and distribute
jobs. It is split in to 4 services:

- **frontend** with svelte
- **frontend api** with go
- **worker manager** manage the creation of jobs (record and whois lookups)
- **worker** process jobs

For development there is a `docker-compose.yml` file that creates the entire
stack and does hot reloading of all components on file change (see
`services/Dockerfile.dev` for more details on how this works).

Production is currently targeting kubernetes, see `manifests` for more details.
There is also a toolbox image that uses `pkg/cmd` to provide some sysadmin
functionality.

## Developing

There are utilities provided my `make` located in `scripts/make`, notable ones
being:

- `make build` build local dev images
- `make push` build and push production images
- `make logs` *dev mode*
- `make stop` *dev mode*
- `make clean` *dev mode*

## Testing
None as of yet as this is more a playground for playing with kubernetes.

## Similar Tools
Collection of services or tools.

- https://dnsspy.io/
