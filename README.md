# Nostr Static Site Generator

A static web site generator written in Go that creates HTML pages from long-form Nostr content. This tool downloads events by their IDs and generates static HTML pages that can be hosted anywhere.

## Demo

https://blog.nostrize.me

## Version 0.1

## Features

- Downloads Nostr events by their IDs
- Generates static HTML pages from long-form content
- Supports multiple relays
- Customizable layout and styling
- Easy deployment to GitHub Pages

## Todo

- [ ] support naddr
- [ ] profiles
- [ ] comments
- [ ] dynamic theme option

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

article_ids:
  - e053bd17c3e36ea81d69ae382aefedf0f53188fc4f95e67b21ce35a38dc9a1e2
  - e9afc5076f13880a18b57c06600c49aa432a7b3ffa92895620f0aeb55ca3218c
```

## Getting Started

1. Fork this repository
2. Modify the `config.yaml` file with your desired configuration
3. Add your Nostr article IDs to the `article_ids` section
4. Place a `logo.png` file (or another file name, but need to update logo) into the project folder
5. Build and run the command: `go build -o nostr-static ./src && ./nostr-static`
5. Commit and push
6. Deploy using GitHub Pages

## Deployment

### GitHub Pages Deployment

1. Go to your repository's Settings
2. Navigate to "Pages" in the menu
3. Under "Build and deployment" > "Source", select "GitHub Actions"

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