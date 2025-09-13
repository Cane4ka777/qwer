# QWER Band API

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Status](https://img.shields.io/badge/Status-Active-success)](#)
[![License](https://img.shields.io/badge/License-Copyright%20Kiy0w0%202025-blue)](#license)


![Qwer Band](public/qwer/group.webp)
Simple REST API for QWER band data (albums, songs, members, awards). Built with Go

## Quick start

```bash
go mod tidy
go run main.go
# visit http://localhost:8080/api
```

## Endpoints

- GET /api
- GET /api/band
- GET /api/members[?id=|name=]
- GET /api/albums[?id=|title=]
- GET /api/songs[?id=|title=|album=]

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=kiy0w0/qwer&type=Date)](https://star-history.com/#kiy0w0/qwer)

## Todo

- All Members fancam
- ~~Website~~
- All qwer images

## License

Copyright Kiy0w0 2025

Made with ♥️ By Pemuja QWER
