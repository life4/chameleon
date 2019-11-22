# Chameleon

Chameleon is web application that reflects content from markdown files from Github.

[Demo](https://articles.life4web.ru/)

## Run from release

1. [Download release](https://github.com/orsinium/chameleon/releases).
2. Extract: `tar -xzf chameleon.tar.gz`
3. Edit config: `nano config.toml`
4. Run binary release for your platform: `./linux-amd64.bin`

## Run from source

```bash
git clone https://github.com/orsinium/chameleon
cd chameleon
cp config{_example,}.toml
go get .
go run *.go
```

## Run from build

```bash
go build .
./chameleon
```
