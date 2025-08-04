Squishy 🧽
=======

Squishy is a lightweight link proxy/shortener, configured through a single yaml file. 

It's intended for simple scenarios as it can (for now) only be configured by editing the config file directly on the server, and is not suited for multi-user setups.

Installation
-------------

### Standalone

You can obtain the standalone binary on the [releases page](https://gitlab.com/smncd/squishy/-/releases).

Simply move this to the location of your choice and make sure it's executable:
```bash
curl -LO https://gitlab.com/smncd/squishy/-/releases/v0.4.0-dev.6/downloads/squishy-linux-arm64-0.4.0-rc.1
chmod +x ./squishy-linux-arm64-0.4.0-rc.1
```

Binaries are available for:
- [`linux-amd64`](https://gitlab.com/smncd/squishy/-/releases/v0.4.0-dev.6/downloads/squishy-linux-amd64-0.4.0-rc.1)
- [`linux-arm64`](https://gitlab.com/smncd/squishy/-/releases/v0.4.0-dev.6/downloads/squishy-linux-arm64-0.4.0-rc.1)

### Docker

Docker images are available on the [Gitlab Container Registry](https://gitlab.com/smncd/squishy/container_registry).

You can get started with a simple docker compose file:
```yaml
services:
  squishy:
    image: registry.gitlab.com/smncd/squishy:latest
    ports:
      - 1394:1394
    volumes:
      - ./squishy.yaml:/squishy.yaml
```

**Note**: You need to set the `config.host` option to `0.0.0.0` when running Squishy in Docker.

Then, simply start the service:
```bash
docker compose up -d
```

### From source

You can build Squishy from this repo. You will need `git` and `go` 1.24 or higher.
```bash
git clone https://gitlab.com/smncd/squishy.git
cd squishy
go mod download
make build-all
```

You'll find the produced binaries in `./bin`.

Configuration
-------------

Squishy is configured with the `squishy.yaml` file. This file must be located in the same directory as the binary and readable by it. Squishy will not run without this file.

The file is split into two main sections `config` and `routes`.

### Config section

```yaml
config:
  host: localhost # The server will listen on localhost
  port: 1394      # The server will use port 1394
  debug: false    # Debugging is disabled
```

### Routes section

```yaml
routes:
  _index: https://example.com # Redirects the root path to https://example.com
  hello:
    _index: https://example2.com # Redirects /hello/ to https://example2.com
    there: https://example3.com # Redirects /hello/there to https://example3.com
```

### Example Configuration File

```yaml
# The config object contains your settings, such as host and port
config:
  host: localhost # Host server listens on
  port: 1394      # Port server listens on
  debug: false    # Enable or disable debugging

# The routes object contains the paths used to access the redirect URLs
routes:
  _index: https://example.com # The root path redirects to https://example.com
  hello:
    _index: https://example2.com # /hello/ redirects to https://example2.com
    there: https://example3.com # /hello/there redirects to https://example3.com
```

### Wildcard routes

Route targets can contain a wildcard indicator: `/*`. By configuring routes like this:
```yaml
# ...
routes:
  example: https://example.com/* # <-- wildcard

```

Squishy will redirect like:
| Squishy path 				| Target   			                   |
| ---		   				| ---    			                   |
| `/example`   				| `https://example.com`                |
| `/example/sub-path` 		| `https://example.com/sub-path`	   |
| `/example/nested/sub-path`| `https://example.com/nested/sub-path`|

**Note**: Wildcard routes are currently not supported for `_index` routes.

License and Ownership
---------------------

Copyright © 2025 Simon Lagerlöf <contact@smn.codes>

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.
