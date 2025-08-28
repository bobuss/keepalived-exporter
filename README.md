# Keepalived Exporter

[![Continuous Integration](https://github.com/mehdy/keepalived-exporter/workflows/Continuous%20Integration/badge.svg)](https://github.com/mehdy/keepalived-exporter/actions)

Prometheus exporter for [Keepalived](https://keepalived.org) metrics.

## Installation

### Debian/Ubuntu

Download the latest release from [here](https://github.com/mehdy/keepalived-exporter/releases) and install it using dpkg

```bash
dpkg -i keepalived-exporter-downloaded-pkg.deb
```

### RedHat/CentOS

Download the latest release from [here](https://github.com/mehdy/keepalived-exporter/releases) and install it using rpm

```bash
rpm -i keepalived-exporter-downloaded-pkg.rpm
```

### Binary releases

```bash
export VERSION=1.3.2
wget https://github.com/mehdy/keepalived-exporter/releases/download/v${VERSION}/keepalived-exporter-${VERSION}.linux-amd64.tar.gz
tar xvzf keepalived-exporter-${VERSION}.linux-amd64.tar.gz keepalived-exporter-${VERSION}.linux-amd64/keepalived-exporter
sudo mv keepalived-exporter-${VERSION}.linux-amd64/keepalived-exporter /usr/local/bin/
```

### From source

```bash
git clone --depth 1 https://github.com/mehdy/keepalived-exporter.git
cd keepalived-exporter
make build
sudo mv keepalived-exporter /usr/local/bin/
```

## Usage

Run keepalived-exporter

```bash
sudo keepalived-exporter [flags]
```

Help on flags

```bash
./keepalived-exporter --help
```

Name               | Description
-------------------|------------
web.listen-address | Address to listen on for web interface and telemetry, defaults to `:9165`.
web.telemetry-path | A path under which to expose metrics, defaults to `/metrics`.
web.config.file    | Configuration file for TLS and Basic Authentication.
ka.json            | Send SIGJSON and decode JSON file instead of parsing text files, defaults to `false`.
ka.pid-path        | A path for Keepalived PID, defaults to `/var/run/keepalived.pid`.
cs                 | Health Check script path to be execute for each VIP.
container-name     | Keepalived container name to export metrics from Keepalived container.
container-tmp-dir  | Keepalived container tmp volume path, defaults to `/tmp`.

**Note:** For `ka.json` option requirement is to have Keepalived compiled with `--enable-json` configure option.

### TLS and Basic Authentication

To enable TLS and/or Basic Authentication, create a configuration file and pass it via the `--web.config.file` flag:

```yaml
tls_server_config:
  cert_file: server.crt
  key_file: server.key
basic_auth_users:
  prometheus: $2y$10$X0jHn2QU5iiF/0wLGwSqEOCF1nA1M5CkQU8mhgvuQ4D4fXZJ0M3fO
```

For more details on the configuration format, see the [exporter-toolkit documentation](https://github.com/prometheus/exporter-toolkit).

### Keepalived on Docker and Keepalived Exporter on host

Set the `--container-name` to the Keepalived container name and set `--container-tmp-dir` to the Keepalived `/tmp` dir path that is volumed to the host

```bash
./keepalived-exporter --container-name keepalived --container-tmp-dir /tmp
```

### Keepalived and Keepalived Exporter on docker

We support two methods to signal keepalived running in a container.

#### Using Docker Signal

This is when the keepalived is running with PID 1 in the container so we can use the standard docker API to send signal to the keepalived process.

Volume docker socket (`/var/run/docker.sock`) to Keepalived Exporter cotnainer in the same path and pass the args like as using Keepalived on container

```bash
docker pull ghcr.io/mehdy/keepalived-exporter
docker run -v keepalived-data:/tmp/ ... $KEEPALIVED_IMAGE
docker run -v /var/run/docker.sock:/var/run/docker.sock -v keepalived-data:/tmp/keepalived-data:ro -p 9165:9165 ghcr.io/mehdy/keepalived-exporter --container-name keepalived --container-tmp-dir "/tmp/keepalived-data"
```

#### Exec to container with PID path

In case the keepalived process is not running with PID 1, this method will exec to the container and use the provided PID path to send the signal.

```bash
docker pull ghcr.io/mehdy/keepalived-exporter
docker run -v keepalived-data:/tmp/ -v keepalived-pid:/var/run/ ... $KEEPALIVED_IMAGE
docker run -v /var/run/docker.sock:/var/run/docker.sock -v keepalived-data:/tmp/keepalived-data:ro -v keepalived-pid:/var/run/:ro -p 9165:9165 ghcr.io/mehdy/keepalived-exporter --container-name keepalived --container-tmp-dir "/tmp/keepalived-data" --ka.container.pid-path "/var/run/keepalived.pid"
```

## Metrics

| Metric                                          | Notes
|-------------------------------------------------|------------------------------------
| keepalived_exporter_build_info                  | Exporter build info
| keepalived_up                                   | Status of Keepalived service
| keepalived_vrrp_state                           | State of vrrp
| keepalived_vrrp_excluded_state                  | State of vrrp with excluded VIP
| keepalived_exporter_check_script_status         | Check Script status for each VIP
| keepalived_gratuitous_arp_delay_total           | Gratuitous ARP delay
| keepalived_advertisements_received_total        | Advertisements received
| keepalived_advertisements_sent_total            | Advertisements sent
| keepalived_become_master_total                  | Became master
| keepalived_release_master_total                 | Released master
| keepalived_packet_length_errors_total           | Packet length errors
| keepalived_advertisements_interval_errors_total | Advertisement interval errors
| keepalived_ip_ttl_errors_total                  | TTL errors
| keepalived_invalid_type_received_total          | Invalid type errors
| keepalived_address_list_errors_total            | Address list errors
| keepalived_authentication_invalid_total         | Authentication invalid
| keepalived_authentication_mismatch_total        | Authentication mismatch
| keepalived_authentication_failure_total         | Authentication failure
| keepalived_priority_zero_received_total         | Priority zero received
| keepalived_priority_zero_sent_total             | Priority zero sent
| keepalived_script_status                        | Tracker Script Status
| keepalived_script_state                         | Tracker Script State

## Check Script

You can specify a check script like Keepalived script check to check if all the things is okay or not.
The script will run for each VIP and gives an arg `$1` that contains VIP.

**Note:** The script should be executable.

```bash
chmod +x check_script.sh
```

### Sample Check Script

```bash
#!/usr/bin/env bash

ping $1 -c 1 -W 1
```
