# MikroManager

**The project is WIP and is not stable at this point.**

## Appreciations

- [Bootstrap](https://getbootstrap.com/)
- [gocron](https://github.com/go-co-op/gocron)
- [cloverDB](https://github.com/ostafen/clover)
- [Logrus](https://github.com/sirupsen/logrus)
- [go-routeros](https://github.com/go-routeros/routeros/tree/v2)

## Getting started

At the moment there's only prebuilt Docker image, so that's the easiest way to use it:

```bash
docker run -d \
    -p 8000:8000 \ # port mapping
    -v /database:/app/database.clover \ # DB persistence, optional
    -v /backups:/app/backups \ # MikroTik exports/backups persistence, optional
    zmazay/mikromanager:main
```

Alternatively, you can build the binary, it's recommended to use go `1.19`:

```bash
go mod download && go build .
```

**The app has no authentication at the moment so be careful with exposing it both internally and externally.**

Would appreciate any [feedback](https://github.com/mazay/mikromanager/issues/new).
