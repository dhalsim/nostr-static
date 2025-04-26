# Nostr Static Site Generator

A static web site generator written in Go that creates HTML pages from long-form Nostr content. This tool downloads events by their IDs and generates static HTML pages that can be hosted anywhere.

## Demo

https://blog.nostrize.me

## Version 0.5

## Features

- Downloads Nostr events by their naddr1 addresses
- Generates static HTML pages from long-form content
- Supports multiple relays
- Customizable layout and styling
- Easy deployment to any static hosting service

## Todo

- [x] support naddr
- [x] profiles
- [x] comments, zapthreads
- [ ] atom, rss feeds
- [ ] dynamic theme option

## Prequisites

### Terminal Knowledge
- Basic understanding of command-line interface (CLI) operations
- Familiarity with navigating directories using `cd` command
- Understanding of basic file operations (copy, move, delete)

### Git Installation

#### Windows
1. Download Git from [git-scm.com](https://git-scm.com/download/win)
2. Run the installer and follow the installation wizard
3. Verify installation by opening Command Prompt or PowerShell and running:
   ```bash
   git --version
   ```

#### macOS
1. Install via Homebrew (recommended):
   ```bash
   brew install git
   ```
2. Or download from [git-scm.com](https://git-scm.com/download/mac)
3. Verify installation:
   ```bash
   git --version
   ```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install git
git --version
```

### Go Installation

#### Windows
1. Download Go from [golang.org/dl](https://golang.org/dl/)
2. Run the installer and follow the installation wizard
3. Verify installation by opening Command Prompt or PowerShell:
   ```bash
   go version
   ```

#### macOS
1. Install via Homebrew (recommended):
   ```bash
   brew install go
   ```
2. Or download from [golang.org/dl](https://golang.org/dl/)
3. Verify installation:
   ```bash
   go version
   ```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang-go
go version
```

## Configuration

The tool is configured using a `config.yaml` file. Here's an example configuration:

```yaml
layout:
  color: dark  # Options: light, dark
  logo: logo.png  # Logo image file name
  title: "Nostr Articles"

relays:
  - wss://relay.damus.io
  - wss://nostr.wine
  - wss://relay.nostr.band
  - wss://nos.lol
  - wss://relay.primal.net

features:
  - comments: true

articles:
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq2j6ezsgu69j7n92cmxxmfsgyeyyvjtxfuk7lwjq6s
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq24wmjfwp6rv6t8v935ujfhv4yr2wzzdfz5gl5quve

profiles:
  - nprofile1qqsxue9c8s0kwnaspf03nqtv99akg99lvlcptz2wqnw5cet7jsgza6qpp4mhxue69uhkummn9ekx7mq8k7c9l
```

## Scheduled Deployment

Since this is a static site generator, content is not dynamic. When an article is updated on Nostr, the changes won't be reflected automatically on your site. To keep your site up-to-date, you can automate the update process using your system's task scheduler or cron jobs. The scripts below will periodically check for updates and deploy them to your site.

### Windows Setup

1. Set up Windows Task Scheduler:
   - Open Task Scheduler
   - Create a new Basic Task
   - Set the trigger (e.g., daily at a specific time)
   - Action: Start a program
   - Program/script: `powershell.exe`
   - Arguments: `-ExecutionPolicy Bypass -File "path\to\scripts\schedule-deploy.ps1"`
   - Complete the wizard

The script will run `nostr-static generate` and commit any changes to your repository.

**Tip:** You can monitor the `cron.log` file to view job execution logs and track updates.

### Unix/Linux Setup

1. Run the setup script:
   ```bash
   ./scripts/setup-cron.sh
   ```

This will set up a daily cron job that runs at midnight. To run hourly instead, edit the script and uncomment the hourly line.

The script will run `nostr-static` and commit any changes to your repository.

**Tips:** 
- Use `crontab -l` to view your scheduled jobs
- Use `crontab -r` to remove all scheduled jobs
- Check `cron.log` for detailed execution logs

## Getting Started

1. Fork this repository
2. Modify the `config.yaml` file with your desired configuration
3. Add Nostr naddr to the `articles` list that you want to serve
4. Place a `logo.png` file (or another file name, but don't forget to update logo in the `config.yaml` file) into the project folder
5. Build and run the command: `go build -o nostr-static ./src && ./nostr-static`
6. Commit and push your changes:
   ```bash
   git add .
   git commit -m "Added my events, changed logo, title, and light theme"
   git push origin main
   ```

## Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/nostr-static.git

# Build the project
go build -o nostr-static ./src

# Run the generator
./nostr-static
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 