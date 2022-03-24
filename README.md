## Image Morph

A simple image manipulation tool written in Go

## Running the app

```bash
# Development
$ cd image-morph && go run main.go
```

## Generating GIFs

```bash
# Innstall convert with Homebrew
$ convert -size 100x100 -delay 0 -loop 0 gen/*.png out.gif
```
