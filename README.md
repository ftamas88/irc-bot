# IRC AutoDL Bot

## Table of Contents
* [About](#about)
* [Environment Variables](#environment-variables)
* [Configuration file](#configuration-file)
* [Quickstart](#quickstart)
* [Tour](#Tour-/-Screenshots)

## About
This bot connects to a list of IRC servers and downloads the .torrent files based on the configuration given

- Multiple IRC server support (TODO)
- Apply filter (TODO)


## Environment Variables

> Before you can start the application, you have to make a copy of the included `.env.dist` file and save it as `.env` in the root folder.

Please update the values according the description

| Key               | Type   | Mandatory | Description                                 |
|:------------------|:-------|:----------|:--------------------------------------------|
| `TRACKERS_FOLDER` | string | Yes       | Link to the trackers folder (`./trackers/`) |

## Configuration file

## Quickstart
### Run the service
```
make run
```
**PS:** This command only available in *WSL* Linux terminal

or

```
go run cmd/app/app.go
```

> The easier option would be to include a compiled binary (.exe) and place it on Gitlab which you can run alongside the `.env` as a standalone file

### Formatting, linting and tests
```
make fmt lint test
```

## Tour / Screenshots
![downloader](https://gitlab.com/ftamas88/irc-bot/-/raw/main/docs/screenshots/irc_downloader.png)
![files](https://gitlab.com/ftamas88/irc-bot/-/raw/main/docs/screenshots/downloaded.png)