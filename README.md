# NVDA AT Automation
Implementation of the AT Automation Driver specification for the NVDA screen reader.

This implementation contains two parts: a Python addon for the NVDA screen reader, and a Golang websocket command-and-control server.

## Usage

### NVDA Add-on

* Move the folder `NVDAPlugin` into your NVDA addons directory — `userConfig\addons` in case of a portable NVDA, `%appdata%\nvda\addons` in case of a system install — and (re)start NVDA.
* A http server should be available at `localhost:8765`.

### Golang Server
* Change to the `Server` folder.
* Compile the application by running `go build main\main.go`.
* Start the server by running `.\main`.
* A websocket server should be available at `ws://localhost:3031` (configurable in `Server\.env`).

## Useful External Links
* [W3C ARIA-AT Automation Repo](https://github.com/w3c/aria-at-automation)
* [AT Automation API Roadmap](https://github.com/w3c/aria-at-automation/issues/15)
* [Protocol Design Issue](https://github.com/w3c/aria-at-automation/issues/20)
* [AT Driver Specification](https://w3c.github.io/at-driver/)
* [NVDA System Tests](https://github.com/nvaccess/nvda/tree/master/tests/system)
