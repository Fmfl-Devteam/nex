<p align="center">
  <img src="assets/nex-gopher.png" alt="nex Gopher with system and network diagnostic tools" width="420">
</p>

# nex

`nex` is a small, cross-platform command-line toolkit for system administration,
development, networking, and diagnostics. It currently provides process,
network, and system inspection commands and is structured so that new command
groups can be added without coupling them to platform-specific code.

## Supported platforms

- Linux
- Windows

## Build and install

Go 1.27 or newer is required by the module.

```bash
go install github.com/fmfl-devteam/nex@latest
```

To build a local checkout:

```bash
go build -o nex .
```

On Windows, use `go build -o nex.exe .`.

## Commands

```text
nex proc list                List visible processes
nex proc info <pid>          Show details for a process
nex proc kill <pid>          Terminate a process
nex net ping <host>          Send one ICMP reachability probe
nex net dns <host>           Resolve IPv4 and IPv6 addresses
nex net ports <host>         Check a bounded list of TCP ports
nex sys info                 Show host, CPU, and memory information
```

Run `nex <command> --help` for command-specific flags and usage. Every result
command supports the global `--json` flag for stable, machine-readable output.

## Examples

```bash
nex proc list
nex proc info 1234
nex net ping example.com
nex net dns example.com
nex net ports example.com --ports 22,80,443
nex net ports example.com --timeout 2s
nex sys info
nex --json sys info
```

The port command checks a conservative default set when `--ports` is omitted.
Checks run concurrently with a fixed limit, and a single invocation accepts at
most 128 explicitly listed ports.

## Development

```bash
gofmt -w .
go vet ./...
go test ./...
GOOS=linux GOARCH=amd64 go build ./...
GOOS=windows GOARCH=amd64 go build ./...
```

Command wiring lives under `cmd/`; reusable operations and stable data types
live under `internal/`. Linux and Windows implementations use build-tagged
files, so command packages contain no runtime platform branching.

On Linux, process inspection reads procfs. On Windows, process and system detail
collection uses the built-in `tasklist.exe` and Windows PowerShell/CIM. Optional
process fields may be unavailable because of permissions. Ping uses the
platform's built-in `ping` program so it does not require raw-socket privileges.

## License

MIT. See [LICENSE](LICENSE).
