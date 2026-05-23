# logslice

Fast log file slicer that extracts time-range segments from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Pipe output to a file:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log > slice.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | stdin |
| `--from` | Start of time range (RFC3339) | required |
| `--to` | End of time range (RFC3339) | required |
| `--format` | Timestamp format in log entries | `RFC3339` |

### Example Output

```
2024-01-15T08:12:44Z INFO  server started on :8080
2024-01-15T08:15:02Z INFO  request received method=GET path=/health
2024-01-15T08:58:31Z ERROR connection timeout host=db.internal
```

## How It Works

logslice uses binary search to efficiently locate the start of a time range within a log file, avoiding the need to scan the entire file from the beginning. This makes it significantly faster than `grep` on large log files.

## License

MIT — see [LICENSE](LICENSE) for details.