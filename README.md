Squishy 🧽
=======

Squishy is a lightweight link proxy/shortener.

Installation
-------------

### Standalone

You can obtain the standalone binary on the [releases page](https://gitlab.com/smncd/squishy/-/releases).

Simply move this to the location of your choice and make sure it's executable:
```bash
curl -LO https://gitlab.com/smncd/squishy/-/releases/v0.4.0-dev.4/downloads/squishy-linux-arm64-0.4.0-dev.4
chmod +x ./squishy-linux-arm64-0.4.0-dev.4
```

Binaries are available for:
- [`linux-amd64`](https://gitlab.com/smncd/squishy/-/releases/v0.4.0-dev.4/downloads/squishy-linux-amd64-0.4.0-dev.4)
- [`linux-arm64`](https://gitlab.com/smncd/squishy/-/releases/v0.4.0-dev.4/downloads/squishy-linux-arm64-0.4.0-dev.4)

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

License and Ownership
---------------------

Copyright © 2025 Simon Lagerlöf <contact@smn.codes>

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.
