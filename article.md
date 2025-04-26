# How to Create a Blog Out of Nostr Long-Form Articles

## Summary

nostr-static is a static site generator (SSG) tool that transforms Nostr long-form content into a beautiful, static blog. It fetches articles from the Nostr network and generates HTML pages that can be hosted anywhere. This tool bridges the gap between Nostr's decentralized content and traditional web hosting, making it easy to create and maintain a blog powered by Nostr content.

**Demo:** [https://blog.nostrize.me](https://blog.nostrize.me)

## Features

### Core Functionality
- **Index Page**: A homepage featuring your blog's title, logo, article summaries, and tags
- **Article Pages**: Individual pages for each article, including:
  - Title and logo
  - Article summary
  - Full content
  - Tags
  - Comments (via ZapThreads integration)

### Social Features
- **Comments**: Integrated with [ZapThreads](https://github.com/franzaps/zapthreads) for decentralized commenting
- **Nostr Connect**: Seamless integration with [window.nostr.js (wnj)](https://git.njump.me/wnj), supporting NIP-46 bunker connect

### Content Organization
- **Tag Pages**: Browse articles filtered by specific tags
- **Profile Pages**: View articles from specific authors
- **Manual Curation**: Select and order articles by adding their naddr strings (see [NIP-19](https://nips.nostr.com/19))

### Customization Options
- **Themes**: Choose between dark and light mode
- **Branding**: 
  - Custom logo
  - Custom blog title
- **Network**: Specify your preferred Nostr relays

### Technical Requirements
- **Profile Format**: Authors must be added in nprofile format (see [NIP-19](https://nips.nostr.com/19)) for consistency
- **Automatic Updates**: Built-in scripts for:
  - Windows Task Scheduler
  - Unix/Linux cron jobs

## Getting Started

1. **Fork and Clone**: 
   - Fork this repository to your GitHub account
   - Clone it to your local machine or use [GitHub Codespaces](https://docs.github.com/en/codespaces/quickstart) for a cloud-based development environment
   - Watch this [quick tutorial](https://youtu.be/_W9B7qc9lVc) to learn more about GitHub Codespaces

2. **Configuration**: Set up your `config.yaml` file with:
   - Blog title and logo
   - Theme preference
   - Relay list
   - Article naddr strings
   - Author nprofile strings

3. **Content Selection**: Add your desired articles by including their naddr strings in the configuration

4. **Author Selection**: You have to add the nprofile strings of the articles. This is needed for URL consistancy.

5. **Build & Run**: Follow the instruction in the README at https://github.com/dhalsim/nostr-static

6. **Deployment**: Choose your preferred static hosting service and deploy the generated HTML files

7. **Updates**: Set up automatic updates using the provided scripts for your operating system (For github pages)

## Deployment Options

### GitHub Pages (Recommended)

GitHub Pages provides free hosting for static websites. Here's how to set it up:

1. **Enable GitHub Pages**:
   - Go to your repository's Settings
   - Navigate to "Pages" in the menu
   - Under "Build and deployment" > "Source", select "GitHub Actions"
   - Enable Actions by following the [GitHub Actions settings guide](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#allowing-select-actions-and-reusable-workflows-to-run)
   - Go to the "Actions" tab in the top menu. If you see the message "Workflows aren't being run on this forked repository", click the "I understand my workflows, go ahead and enable them" button

2. **Custom Domain Setup**:
   - Purchase a domain from your preferred domain registrar
   - Create a CNAME record in your domain's DNS settings:
     - Type: CNAME
     - Name: @ or www or a subdomain you prefer (depending on your preference)
     - Value: YOUR_GITHUB_USERNAME.github.io
   - In your repository's GitHub Pages settings:
     - Enter your custom domain in the "Custom domain" field
     - Check "Enforce HTTPS" for secure connections
   - Wait for DNS propagation (can take up to 24 hours)
   - Your site will be available at your custom domain

### Other Hosting Options

You can also deploy your static site to any hosting service that supports static websites, such as:
- Netlify
- Vercel
- Cloudflare Pages
- Amazon S3
- Any traditional web hosting service

## Why nostr-static?

nostr-static offers a unique solution for bloggers who want to leverage Nostr's decentralized content while maintaining a traditional web presence. It combines the best of both worlds:

- **Decentralized Content**: Your articles live on the Nostr network
- **Traditional Web Presence**: A familiar blog interface for your readers
- **Easy Maintenance**: Simple configuration and automatic updates
- **Flexible Hosting**: Deploy anywhere that supports static websites
- **Social interactions**: Leverage nostr for comments

## Conclusion

nostr-static makes it easy to create a professional blog from your Nostr long-form content. Whether you're a seasoned Nostr user or new to the ecosystem, this tool provides a straightforward way to share your content with both the Nostr community and traditional web users.

Start your Nostr-powered blog today by visiting the [demo](https://blog.nostrize.me) and exploring the possibilities!