# Nostr Static Site Generator

A static web site generator written in Go that creates HTML pages from long-form Nostr content. This tool downloads events by their naddr and generates static HTML pages that can be hosted anywhere.

## Demo

https://blog.nostrize.me

## Version 0.7

## Features

- Downloads Nostr events by their naddr1 addresses
- Generates static HTML pages from long-form content
- Choose between dark and light themes
- Lists articles by profiles and tags
- Has RSS and Atom feeds for index, profiles and tags
- Smart content discovery with article recommendations
- One-click deployment to GitHub Pages with automated workflows

## Todo

- [x] support naddr
- [x] profiles
- [x] comments, zapthreads
- [x] atom, rss feeds
- [x] discoveribility
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
  title: "Nostr Articles"  # Site title

blog_url: https://blog.nostrize.me

relays:
  - wss://relay.damus.io
  - wss://nostr.wine
  - wss://relay.nostr.band
  - wss://nos.lol
  - wss://relay.primal.net

features:
  comments: true
  nostr_links: njump.me
  tag_discovery: true

settings:
  tag_discovery:
    fetch_count_per_tag: 100      # fetch this many events per tag
    popular_articles_count: 10    # max number of articles to display in tag discovery
    weights:
      reaction_weight: 1.0       # weight for reaction count
      repost_weight: 2.0         # weight for repost count
      reply_weight: 1.5          # weight for reply count
      report_penalty: 5.0        # penalty for each report
      zap_amount_weight: 0.1     # weight for total zap amount
      zapper_count_weight: 0.5   # weight for number of unique zappers
      zap_avg_weight: 2.0        # weight for average zap amount
      zap_total_weight: 3.0      # multiplier for total zap score
      max_sats: 10000000.0       # maximum sats for normalization (0.1 BTC)

articles:
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq25x6r2d90hzaj529f9gur2x9h4yt2gfy69qm05rzu
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq2j6ezsgu69j7n92cmxxmfsgyeyyvjtxfuk7lwjq6s
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq24wmjfwp6rv6t8v935ujfhv4yr2wzzdfz5gl5quve
  - naddr1qvzqqqr4gupzqwlsccluhy6xxsr6l9a9uhhxf75g85g8a709tprjcn4e42h053vaqqyrzwrxvc6ngvfkxty9t8

profiles:
  - nprofile1qqsrhuxx8l9ex335q7he0f09aej04zpazpl0ne2cgukyawd24mayt8gpp4mhxue69uhkummn9ekx7mqprhgw8
  - nprofile1qqsxue9c8s0kwnaspf03nqtv99akg99lvlcptz2wqnw5cet7jsgza6qpp4mhxue69uhkummn9ekx7mq8k7c9l
```

## Getting Started

1. Fork this repository
2. Clone it into your local machine or use [codespaces](https://docs.github.com/en/codespaces/quickstart)
3. Modify the `config.yaml` file with your desired configuration
4. Add Nostr naddr to the `articles` list that you want to serve
5. Place a `logo.png` file (or another file name, but don't forget to update logo in the `config.yaml` file) into the project folder
6. Build and run the command: `go build -o nostr-static ./src && ./nostr-static`
7. Commit and push your changes:
   ```bash
   git add .
   git commit -m "Added my events, changed logo, title, and light theme"
   git push origin main
   ```

## Tag Discovery

The tag discovery feature helps users find popular and relevant content by analyzing engagement metrics across different tags. To enable this feature:

1. Set `features.tag_discovery: true` in your `config.yaml`
2. Run `./nostr-static` once to generate the initial `tags.txt` file
3. Create the content index by running:
   ```bash
   ./nostr-static index-tag-discovery
   ```
   This will fetch events from the Nostr network and their engagement statistics from [nostr.band API](https://api.nostr.band/).

4. Calculate popularity scores for articles:
   ```bash
   ./nostr-static calculate-tag-discovery
   ```
   The scoring system considers:
   - Engagement metrics (reposts, replies, reactions)
   - Zap statistics (total amount, unique zappers, average zap size)
   - Content quality signals (report penalties)

5. Finally, run `./nostr-static` again to generate the recommended articles section

You can customize the scoring weights in your `config.yaml` under `settings.tag_discovery.weights`. If not specified, the system will use default weights optimized for balanced content discovery. Setting a weight to 0 will disable that particular scoring factor.

Example:
```yaml
settings:
  tag_discovery:
    weights:
      report_penalty: 0      # disable report penalties
      reaction_weight: 2.0   # increase reaction weight
      # other weights will use defaults
```

## Deployment

### GitHub Pages Deployment

1. Go to your repository's Settings
2. Navigate to "Pages" in the menu
3. Under "Build and deployment" > "Source", select "GitHub Actions"
4. Enable Actions by following the [GitHub Actions settings guide](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#allowing-select-actions-and-reusable-workflows-to-run)
5. Go to the "Actions" tab in the top menu. If you see the message "Workflows aren't being run on this forked repository", click the "I understand my workflows, go ahead and enable them" button


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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 