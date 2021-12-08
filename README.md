# IRC AutoDL Bot

## Table of Contents
* [About](#about)
* [Quickstart](#quickstart)
* [Environment Variables](#environment-variables)
* [Tour](#tour)

## About
This bot connects to a list of IRC servers and downloads the .torrent files based on the configuration given

- Multiple IRC server support (TODO)
- Apply filter (TODO)

## Quickstart
### Run the service
```
make run
```

or

```
go run cmd/app/app.go
```

### Formatting, linting and tests
```
make fmt lint test
```

## Environment Variables

| Key               | Type   | Mandatory | Description               |
|:------------------|:-------|:----------|:--------------------------|
| `PASSKEY`         | string | Yes       | Your tracker passkey      |
| `INVITE_CODE`     | string | Yes       | iRC invite code           |
| `SERVER`          | string | Yes       | iRC server IP/Hostname    |
| `PORT`            | integer| Yes       | iRC server port number    |
| `NICK`            | string | Yes       | iRC bot nickname          |
| `DOWNLOAD_DIR`    | string | Yes       | Download folder relative to the root  |
| `DOWNLOAD_LINK`   | string | Yes       | Don't change it, template to download |

## Tour / Screenshots
![downloader](https://gitlab.com/ftamas88/irc-bot/-/raw/main/docs/screenshots/irc_downloader.png)
![files](https://gitlab.com/ftamas88/irc-bot/-/raw/main/docs/screenshots/downloaded.png)