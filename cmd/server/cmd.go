package main

import (
	"flag"
	"fmt"
	"os"
)

// handleCLIWithFlagSet parses command-line arguments and handles help/env flags.
func handleCLIWithFlagSet(fs *flag.FlagSet) error {
	helpFlag := fs.Bool("help", false, "Show help message")
	envFlag := fs.Bool("env", false, "List environment variables and their default values")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}
	if *envFlag {
		printEnv()
		os.Exit(0)
	}

	return nil
}

// printHelp displays the usage instructions for the binary.
func printHelp() {
	fmt.Println("Usage: iscsistat-exporter [options]")
	fmt.Println()
	fmt.Println("A Prometheus exporter for iSCSI volume statistics.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help     Show this help message")
	fmt.Println("  --env      List all supported environment variables")
	fmt.Println()
	fmt.Println("If no options are provided, the server will start using environment variables or defaults.")
}

// printEnv lists all environment variables used by the application and their defaults.
func printEnv() {
	fmt.Println("Environment Variables Configuration:")
	fmt.Println()
	fmt.Println("  Server settings:")
	fmt.Println("    HTTP_HOST                Host to bind the server to (default: 0.0.0.0)")
	fmt.Println("    HTTP_PORT                Port for the HTTP server (default: 9101)")
	fmt.Println()
	fmt.Println("  Security (TLS/mTLS):")
	fmt.Println("    HTTP_TLS_ENABLED         Enable HTTPS (true/false, default: false)")
	fmt.Println("    HTTP_TLS_CERT_FILE       Path to the server certificate file")
	fmt.Println("    HTTP_TLS_KEY_FILE        Path to the server private key file")
	fmt.Println("    HTTP_TLS_CLIENT_CA_FILE  Path to CA file to enable mTLS client verification")
	fmt.Println()
	fmt.Println("  Metrics collection:")
	fmt.Println("    METRICS_COLLECT_INTERVAL Frequency of disk stats collection (default: 15s)")
	fmt.Println("                             Supports units: s, m, h (e.g., 30s, 1m)")
	fmt.Println()
}
