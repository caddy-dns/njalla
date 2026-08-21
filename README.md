Njalla for [Caddy](https://github.com/caddyserver/caddy)
=======================

DNS provider module backed by [libdns/njalla](https://github.com/libdns/njalla).
Use it to solve the ACME DNS-01 challenge and to manage records in Njalla zones.

Module name: `dns.providers.njalla`

## Install

Caddy does not ship this module. Build it in with
[xcaddy](https://github.com/caddyserver/xcaddy):

```sh
xcaddy build --with github.com/caddy-dns/njalla
```

The included `Dockerfile` builds a Caddy image with the module baked in:

```sh
docker build -t caddy-njalla .
```

## Authentication

Create a token under [Settings → API](https://njal.la/settings/api/). For the
DNS-01 challenge, pick the **ACME** option: it grants `add-record`,
`edit-record`, `remove-record`, and `list-records`, scoped to TXT records under
`_acme-challenge`.

The [libdns provider README](https://github.com/libdns/njalla) lists the full
permissions and the provider's behaviour.

## Usage

Keep the token in an environment variable. Caddy resolves
[placeholders](https://caddyserver.com/docs/conventions#placeholders) when it
provisions the module, and a literal string also works.

### Caddyfile

One site, through the
[`tls` directive](https://caddyserver.com/docs/caddyfile/directives/tls):

```caddy
tls {
	dns njalla {env.NJALLA_API_TOKEN}
}
```

Every site, through
[global options](https://caddyserver.com/docs/caddyfile/options):

```caddy
{
	acme_dns njalla {env.NJALLA_API_TOKEN}
}
```

The token also takes a block:

```caddy
tls {
	dns njalla {
		api_token {env.NJALLA_API_TOKEN}
	}
}
```

### JSON

Configure the
[ACME issuer](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/):

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "njalla",
        "api_token": "{env.NJALLA_API_TOKEN}"
      }
    }
  }
}
```

## Troubleshooting

| Symptom                                            | Cause                                                                                                                                                                 |
| -------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `API token is empty after resolving placeholders`  | The placeholder resolved to nothing. Set the variable for the process running Caddy. A systemd unit ignores your shell environment, so declare it with `Environment=` or `EnvironmentFile=`. |
| The challenge record stays in the zone              | The token cannot call `remove-record` for `_acme-challenge`. Njalla's ACME token option covers it.                                                                      |
| A wildcard certificate fails                        | Issue for the Njalla domain. Another provider serves any subdomain you have delegated away, and this module cannot write there.                                         |
| `unrecognized subdirective`                         | The block accepts `api_token` and nothing else. Check the spelling.                                                                                                    |

## Testing

```sh
go test ./...
```

Confirm the module links into a real Caddy build:

```sh
xcaddy build --with github.com/caddy-dns/njalla=.
./caddy list-modules | grep dns.providers.njalla
```
