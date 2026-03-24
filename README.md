# iSCSI Volume Metrics Exporter

This Go application automatically discovers iSCSI volumes, collects their storage utilization metrics (total and available space), and exposes them via an HTTP endpoint for Prometheus to scrape.

## Motivation

In environments utilizing iSCSI for storage, it's crucial to monitor the capacity and usage of these volumes. This exporter provides a simple, automated way to integrate iSCSI storage metrics into a Prometheus monitoring stack. It addresses the need to dynamically identify iSCSI targets and report their disk space statistics without manual intervention.

## Features

*   **Dynamic iSCSI Volume Discovery:** Automatically finds connected iSCSI targets using the `iscsiadm` command.
*   **Detailed Storage Metrics:** Reports total size and available space for each discovered volume in bytes.
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
*   **`iscsiadm`** - Command-line utility for iSCSI management. Must be installed and configured on the host where the exporter will run.
    *   On Debian/Ubuntu: `sudo apt install open-iscsi`
    *   On RHEL/CentOS/Fedora: `sudo dnf install iscsi-initiator-utils`
*   Access to iSCSI targets.

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

### iSCSI Discovery with `iscsiadm`

The `iscsiadm` command-line utility is employed to discover available iSCSI targets.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues on the project's repository.

### 🧪 Tested Environments

| iSCSI Provider | CSI Plugin | Status |
|---|---|---|
| Synology NAS | [democratic-csi](https://github.com/democratic-csi/democratic-csi) | ✅ Tested |

If you successfully use this exporter with another iSCSI provider or CSI plugin, please open a PR or issue to update this table!
