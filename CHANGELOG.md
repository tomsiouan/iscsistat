# v0.0.3

Released 2026-03-25

- refactor(iscsi): replace `iscsiadm` + `df` + regex parsing with modular parser system using `/sys/class/iscsi_session` + `syscall.Statfs` for volume discovery and metrics collection
- add support for multiple parsers with configurable naming schemes
- add logging of the parser being used in the configuration

# v0.0.2

Released 2026-03-24

- fix(iscsi): generalize target name and df header parsing
- Fix unhandled error on `logger.Sync()` (errcheck)
- Fix unhandled error on `http2.ConfigureServer` (errcheck)
- Fix non-constant format string in `fmt.Errorf` (govet)

# v0.0.1

Released 2026-03-24

- initial release
