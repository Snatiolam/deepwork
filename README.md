# Deepwork

A distraction-blocking focus timer with a built-in accountability system.

## Overview

Deepwork is a command-line tool that helps you stay focused by temporarily blocking access to distracting websites. When you start a focus session, it modifies your system's hosts file to redirect blocked domains to `0.0.0.0`, effectively making them unreachable.

### Features

- **Website Blocking**: Blocks common time-wasting sites (Reddit, Twitter/X, YouTube, etc.)
- **Focus Timer**: Configurable session duration with a beautiful TUI countdown
- **Shame Mode**: If you try to quit early, you must type "I surrender to my distractions" — a psychological barrier to giving up
- **DNS Flushing**: Automatically clears DNS cache when blocking/unblocking
- **Cross-platform**: Works on Linux systems with systemd-resolved, resolvectl, or nscd

## Installation

```bash
git clone https://github.com/Snatiolam/deepwork.git
cd deepwork
go build -o deepwork
sudo cp deepwork /usr/local/bin/
```

## Usage

```bash
sudo deepwork -min 25
```

### Options

- `-min`: Session duration in minutes (default: 1)
- `-dir`: Custom blocklist directory (default: `~/.config/deepwork/blocklists`)

### Example

```bash
# Start a 50-minute focus session
sudo deepwork -min 50
```

## Blocklist

The default blocklist is stored at `~/.config/deepwork/blocklists/blocklist.txt`.

### Default blocked domains

- reddit.com / www.reddit.com
- twitter.com / x.com
- youtube.com / www.youtube.com
- Various DNS providers to prevent DNS bypass

### Adding custom domains

Edit `~/.config/deepwork/blocklists/blocklist.txt` and add one domain per line:

```
example.com
www.example.com
```

## How It Works

1. **Startup**: Requires root privileges (sudo) to modify `/etc/hosts`
2. **Blocking**: Appends blocked domains mapped to `0.0.0.0` in `/etc/hosts`
3. **DNS Flush**: Clears system DNS cache to apply changes immediately
4. **Timer**: Displays TUI with countdown; pressing Ctrl+C activates "Shame Mode"
5. **Completion**: Restores `/etc/hosts` to original state and flushes DNS again

## Shame Mode

When you attempt to quit mid-session (Ctrl+C or 'q'), you enter Shame Mode. To exit early, you must type exactly:

```
I surrender to my distractions
```

This creates a psychological hurdle that makes it easier to stay focused than to give in.

## Requirements

- Go 1.25+
- Root/sudo access
- Linux with one of: `resolvectl`, `systemd-resolve`, or `nscd`

## License

MIT License - see [LICENSE](LICENSE) file.
