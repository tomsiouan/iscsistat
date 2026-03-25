# iSCSI Volume Metrics Exporter

This Go application automatically discovers iSCSI volumes, collects their storage utilization metrics (total and available space), and exposes them via an HTTP endpoint for Prometheus to scrape.

## Motivation

In environments utilizing iSCSI for storage, it's crucial to monitor the capacity and usage of these volumes. This exporter provides a simple, automated way to integrate iSCSI storage metrics into a Prometheus monitoring stack. It addresses the need to dynamically identify iSCSI targets and report their disk space statistics without manual intervention.

## Architecture — Client-Side Exporter

This exporter runs **on the iSCSI initiator side** (the client that mounts the volumes), not on the storage server.

**Advantages:**
- Reports the actual disk usage as seen by the host (filesystem level)
- No access to the storage backend required
- Works with any iSCSI target (NAS, SAN, CSI provisioner...)

**Limitations:**
- Must be deployed on every node that mounts iSCSI volumes
- Requires read access to `/sys/class/iscsi_session` and the mounted filesystems
- Only sees volumes that are currently connected and mounted on the node

## Features

*   **Dynamic iSCSI Volume Discovery:** Automatically finds connected iSCSI sessions by reading `/sys/class/iscsi_session`.
*   **Detailed Storage Metrics:** Reports total size and used space for each discovered volume in bytes, using `syscall.Statfs`.
*   **Prometheus-Compatible Metrics:** Exposes metrics on `/metrics` endpoint, ready for Prometheus scraping.
*   **Flexible Configuration:** Supports loading settings from a `config.toml` file or environment variables.

## Getting Started

### Prerequisites

*   **Go** (version 1.25 or higher recommended)
*   **Task** - A task runner / build tool used as a modern alternative to `make`.
    *   Installation: [taskfile.dev](https://taskfile.dev/installation/)
    *   All available tasks can be listed by running:
    ```bash
    task --list
    ```
*   Read access to `/sys/class/iscsi_session` on the host.
*   iSCSI volumes must be **connected and mounted** on the host.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/tomsiouan/iscsistat.git
    cd iscsistat
    ```

2.  **Build the application:**
    ```bash
    task build
    ```

### Configuration

See the `config.example.toml` file.

## Available Parsers

The exporter supports different parsers to handle various iSCSI target naming schemes. The following parsers are currently available:

| Parser Name          | Description                                                                 |
|----------------------|-----------------------------------------------------------------------------|
| `generic`            | Default parser that works with any iSCSI target (no special processing)    |
| `democratic-csi`     | Parser for volumes managed by the [Democratic-CSI](https://github.com/democratic-csi/democratic-csi) plugin |

Each parser can be configured with optional prefix/suffix trimming in the configuration file:

```toml
[iscsi]
  parser      = "democratic-csi"
  trim_prefix = "csi-"
  trim_suffix = "-data"
```


## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues on the project's repository.

### 🧪 Tested Environments

| iSCSI Provider | CSI Plugin | Status |
|---|---|---|
| Synology NAS | [democratic-csi](https://github.com/democratic-csi/democratic-csi) | ✅ Tested |

If you successfully use this exporter with another iSCSI provider or CSI plugin, please open a PR or issue to update this table!
