# VAPE - RIPE Atlas CLI Tool

A command-line interface for interacting with the RIPE Atlas API.

## Installation

```bash
go build -o vape .
```

## Configuration

### Authentication Priority

1. `--key` CLI flag (highest priority)
2. `RIPE_ATLAS_KEY` environment variable
3. `~/.vape.yaml` config file (lowest priority)

### Initialize Config

```bash
vape config init
```

### Config File Format (~/.vape.yaml)

```yaml
api_key: "your-api-key"
output_format: table
default_af: 4
default_probes: 10
```

### View Current Config

```bash
vape config show
```

### Set Config Values

```bash
vape config set api_key "your-api-key"
vape config set output_format json
```

## Commands

### Probes

List probes:
```bash
vape probes list
vape probes list --country US
vape probes list --asn 15169
vape probes list --status Connected --limit 50
```

Get probe details:
```bash
vape probes get 12345
```

### Measurements

Create a measurement:
```bash
vape measurements create --type ping --target google.com
vape measurements create --type traceroute --target 8.8.8.8 --probes 20 --af 4
vape measurements create --type dns --target 8.8.8.8 --query-argument example.com --query-type A
vape measurements create --type http --target example.com --port 443
```

List measurements:
```bash
vape measurements list
vape measurements list --mine
vape measurements list --status Ongoing --type ping
```

Get measurement details:
```bash
vape measurements get 12345678
```

Get measurement results:
```bash
vape measurements results 12345678
vape measurements results 12345678 --limit 10
vape measurements results 12345678 --probe-ids 1234,5678
```

Stop a measurement:
```bash
vape measurements stop 12345678
```

### Quick Measurements

Run quick one-off measurements with sensible defaults:

```bash
# Ping
vape quick ping google.com
vape quick ping 8.8.8.8 --probes 20 --wait

# Traceroute
vape quick traceroute google.com
vape quick trace 8.8.8.8 --protocol UDP

# DNS
vape quick dns 8.8.8.8 example.com
vape quick dns 1.1.1.1 google.com --query-type AAAA
vape quick dns probe example.com  # Use probe's resolver

# SSL Certificate
vape quick ssl example.com
vape quick ssl example.com --port 8443

# HTTP
vape quick http example.com
vape quick http example.com --path /api/health --method HEAD

# NTP
vape quick ntp pool.ntp.org
vape quick ntp time.google.com --wait
```

### Credits

Check your credit balance:
```bash
vape credits
```

### Version

```bash
vape version
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--key` | RIPE Atlas API key |
| `--config` | Config file path (default: ~/.vape.yaml) |
| `-o, --output` | Output format: `table` or `json` |

## Output Formats

Table output (default):
```bash
vape probes list --country US --limit 5
```

JSON output:
```bash
vape probes list --country US --limit 5 -o json
```

## Aliases

| Command | Aliases |
|---------|---------|
| `measurements` | `msm`, `m` |
| `probes` | `probe`, `p` |
| `traceroute` | `tr`, `trace` |
| `ssl` | `tls`, `cert` |

## Measurement Types

- `ping` - ICMP/UDP/TCP ping
- `traceroute` - Network path tracing
- `dns` - DNS queries
- `sslcert` - SSL/TLS certificate checks
- `http` - HTTP requests
- `ntp` - NTP server checks

## Probe Selection

When creating measurements, you can specify probe selection criteria:

| Type | Description | Example |
|------|-------------|---------|
| `area` | Geographic area | `WW` (worldwide), `West`, `North-Central` |
| `country` | Country code | `US`, `DE`, `JP` |
| `asn` | AS number | `15169` |
| `prefix` | IP prefix | `8.8.8.0/24` |
| `msm` | Probes from another measurement | `12345678` |

Example:
```bash
vape measurements create --type ping --target google.com --probe-type country --probe-value DE --probes 10
```
