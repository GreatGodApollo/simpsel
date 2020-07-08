<h1 align="center">simpsel</h1>
<p align="center"><i>Made with :heart: by <a href="https://github.com/GreatGodApollo">@GreatGodApollo</a></i></p>

An assembly language & accompanying VM implemented in Go

**TODO: Write Tests** 

## Built With
Plain old Go & the standard library. Now also uses [gliderlabs/ssh](https://github.com/gliderlabs/ssh) to provide an SSH
access point to the REPL.

## Usage
To launch the repl: `./simpsel`

To use the SSH server make sure to generate a key first `ssh-keygen -t ed25519 -f ./host.key`. Can than be started with
`./simpsel -ssh`

To run a file directly: `./simpsel -file test.sasm`

## Licensing

This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/)

## Authors

* [Brett Bender](https://github.com/GreatGodApollo)