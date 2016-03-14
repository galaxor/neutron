# neutron

Self-hosted server for [Protonmail client](https://github.com/ProtonMail/WebClient)

## Usage

```bash
# Build client
git submodule init
git submodule update
make build-client

# Start server
make start
```

Credentials:
* Username: `neutron`
* Password: `neutron`
* Mailbox password: `neutron` (public/private PGP keys are stored in `data/`)
