# Nostr Static Site Generator

A static web site generator written in Go that creates HTML pages from long-form Nostr content. This tool downloads events by their IDs and generates static HTML pages that can be hosted anywhere.

## Demo

https://blog.nostrize.me

## Version 0.3

## Features

- Downloads Nostr events by their naddr1 addresses
- Generates static HTML pages from long-form content
- Supports multiple relays
- Customizable layout and styling
- Easy deployment to GitHub Pages

## Todo

- [x] support naddr
- [x] profiles
- [ ] comments, zapthreads
- [ ] atom, rss feeds
- [ ] dynamic theme option
- [ ] release binaries

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

articles:
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq2j6ezsgu69j7n92cmxxmfsgyeyyvjtxfuk7lwjq6s
  - naddr1qvzqqqr4gupzqmnyhq7p7e60kq997xvpds5hkeq5hanlq9vffczd6nr9062pqthgqq24wmjfwp6rv6t8v935ujfhv4yr2wzzdfz5gl5quve
```

## Deployment

### GitHub Pages Deployment

1. Go to your repository's Settings
2. Navigate to "Pages" in the menu
3. Under "Build and deployment" > "Source", select "GitHub Actions"
4. Enable Actions by following the [GitHub Actions settings guide](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#allowing-select-actions-and-reusable-workflows-to-run)
5. Go to the "Actions" tab in the top menu. If you see the message "Workflows aren't being run on this forked repository", click the "I understand my workflows, go ahead and enable them" button

## Getting Started

1. Fork this repository
2. Modify the `config.yaml` file with your desired configuration
3. Add your Nostr article IDs to the `article_ids` section
4. Place a `logo.png` file (or another file name, but don't forget to update logo in the `config.yaml` file) into the project folder
5. Build and run the command: `go build -o nostr-static ./src && ./nostr-static`
6. Commit and push your changes:
   ```bash
   git add .
   git commit -m "Added my events, changed logo, title, and light theme"
   git push origin main
   ```
7. The GitHub Pages workflow will deploy your site automatically
8. Your website will be available at `https://YOUR_GITHUB_USERNAME.github.io/nostr-static/`

### Custom Domain Setup

To use a custom domain:

1. Add a CNAME record in your domain (root or subdomain)
   - Set the value to `YOUR_GITHUB_USERNAME.github.io`
2. In your repository's Pages Settings:
   - Set the Custom Domain field to your domain
   - Wait for DNS propagation (can take up to 24 hours)

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