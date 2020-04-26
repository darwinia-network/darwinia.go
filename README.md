<h1 align="center">
𝓭𝓪𝓻𝔀𝓲𝓷𝓲𝓪.𝓰𝓸
</h1>

[![Golang CI][workflow-badge]][github]

## Install

```
go get install "github.com/darwinia-network/darwinia.go/dargo"
```

## Config

`dargo` use the same config file with `darwinia.js`, if you don't know what 
`darwinia.js` is, run the scripts below before you start

```
mkdir ~/.darwinia
echo '{"eth": { api: "infura api with your key"}}' > ~/.darwinia/config.json
```

## Usage

```sh
$ dargo
The way to Go

Usage:
  dargo [command]

Available Commands:
  epoch       Calculate epoch cache
  header      Get eth block by number
  help        Help about any command
  proof       Proof the block by number
  shadow      Start shadow service
  version     Print the version number of dargo

Flags:
  -h, --help   help for dargo

Use "dargo [command] --help" for more information about a command.

```

## Trouble Shooting

Everytime you run `proof` in error, please delete `~/.ethashproof` and `~/.ethash` 
and retry.

## LICENSE

GPL-3.0


[github]: https://github.com/darwinia-network/darwinia.go
[workflow-badge]: https://github.com/darwinia-network/darwinia.go/workflows/Golang%20CI/badge.svg
