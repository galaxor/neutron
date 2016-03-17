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

Default credentials:
* Username: `neutron`
* Password: `neutron`
* Mailbox password: `neutron`

## Roadmap

- [ ] Implement dummy server (see [#1](https://github.com/emersion/neutron/issues/1))
- [x] Define backend interface (see https://github.com/emersion/neutron/blob/master/backend/backend.go)
- [ ] Implement IMAP + SMTP interface

## Backends

All backends must implement the [backend interface](https://github.com/emersion/neutron/blob/master/backend/backend.go).

Currently, only a simple memory backend is available. Nothing is saved on disk, everything is destroyed when the server is shut down.

## License

MIT
