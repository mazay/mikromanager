# MikroManager

**The project is WIP and is not stable at this point.**

## Appreciations

- [Bootstrap](https://getbootstrap.com/)
- [gocron](https://github.com/go-co-op/gocron)
- [GORM](https://gorm.io/)
- [zap](https://github.com/uber-go/zap)
- [go-routeros](https://github.com/go-routeros/routeros)
- [highlight.js](https://highlightjs.org/)

## Getting started

Create a config file using the example - https://github.com/mazay/mikromanager/blob/main/config.yml

Pay attention to S3 settings, at a bare minimum you will need to set `s3Bucket`. S3 credentials are also mandatory but can be set via the environment variables.

At the moment there's only prebuilt Docker image, so that's the easiest way to use it:

```bash
docker run -d \
    -p 8000:8000 \ # port mapping
    -v /mikromanager:/app \ # the persistence directory, should contain the config and "database" subdirectory
    -e "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE" \
    -e "AWS_SECRET_ACCESS_KEY=je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY" \
    ghcr.io/mazay/mikromanager:main
```

Alternatively, you can build the binary, it's recommended to use go `1.24` or later versions:

```bash
go mod download && go build .
```

**Default username and password is `admin` make sure to change it.**

Would appreciate any [feedback](https://github.com/mazay/mikromanager/issues/new).

**Notes**

The `mikromanager` will try to find a management IP using comment filter `MGMT`, if found device IP will be updated. This should help with subnet migrations, just make sure you have only one address with that comment, `mikromanager` will use the first found.
