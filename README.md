# Hetty

Hetty is an HTTP toolkit for security research. It aims to become an open source
alternative to commercial software like Burp Suite Pro, with powerful features
tailored to the needs of the infosec and bug bounty community.

![Screenshot](docs/screenshot.png)

## Features

- Man-in-the-middle (MITM) HTTP/1.1 proxy with logs
- Project management (e.g. for organizing work across bug bounty programs)
- Intercept requests and responses for manual review/editing
- Search/filter proxy logs
- Sender module for sending manual HTTP requests (similar to Burp's Repeater)
- Scope management to help keep work within defined boundaries

## Requirements

- Go 1.21+
- Node.js 18+ (for building the frontend)

## Installation

### From source

```bash
git clone https://github.com/dstotijn/hetty.git
cd hetty
make build
```

### Docker

```bash
docker pull ghcr.io/dstotijn/hetty:latest
docker run -p 8080:8080 ghcr.io/dstotijn/hetty:latest
```

## Usage

Run the Hetty binary:

```bash
./hetty
```

By default, Hetty listens on `127.0.0.1:8080`. Open your browser and navigate
to `http://localhost:8080` to access the web interface.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--addr` | `127.0.0.1:8080` | Address to listen on |
| `--db` | `hetty.db` | Path to the database file |
| `--cert` | `` | Path to CA certificate |
| `--key` | `` | Path to CA private key |
| `--upstream` | `` | Upstream proxy URL |
| `--verbose` | `false` | Enable verbose logging |

### Setting up the CA certificate

To intercept HTTPS traffic, you need to install Hetty's CA certificate in your
browser or system trust store.

1. Start Hetty and navigate to the web interface
2. Go to **Settings** and download the CA certificate
3. Install the certificate in your browser or OS trust store
4. Configure your browser to use `127.0.0.1:8080` as an HTTP/HTTPS proxy

## Development

### Backend

```bash
go run ./cmd/hetty
```

### Frontend

```bash
cd admin
npm install
npm run dev
```

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) before
submitting a pull request.

## License

[Apache License 2.0](LICENSE)

## Acknowledgements

This project is a fork of [dstotijn/hetty](https://github.com/dstotijn/hetty).
All credit for the original work goes to the original author and contributors.
