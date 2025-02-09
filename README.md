# MikroManager

**The project is WIP and is not stable at this point.**

## Appreciations

- [Bootstrap](https://getbootstrap.com/)
- [gocron](https://github.com/go-co-op/gocron)
- [cloverDB](https://github.com/ostafen/clover)
- [zap](https://github.com/uber-go/zap)
- [go-routeros](https://github.com/go-routeros/routeros/tree/v2)
- [highlight.js](https://highlightjs.org/)

## Getting started

At the moment there's only prebuilt Docker image, so that's the easiest way to use it:

```bash
docker run -d \
    -p 8000:8000 \ # port mapping
    -v /database:/app/database.clover \ # DB persistence, optional
    -v /backups:/app/backups \ # MikroTik exports/backups persistence, optional
    ghcr.io/mikromanager:main
```

Alternatively, you can build the binary, it's recommended to use go `1.19`:

```bash
go mod download && go build .
```

**Default username and password is `admin` make sure to change it.**

Would appreciate any [feedback](https://github.com/mazay/mikromanager/issues/new).

**Notes**

The `mikromanager` will try to find a management IP using comment filter `MGMT`, if found device IP will be updated. This should help with subnet migrations, just make sure you have only one address with that comment, `mikromanager` will use the first found.
