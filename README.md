# gostore

A fiber app used to manage files in a server using HTTP requests.

Next steps:

* Add tests. Starting with object handler and object service
* Implement monitoring of disk usage
* Implement backups and the posibility to upload them to a cloud solution like S3

setup:

```bash
cd gostore
go install
```

Install air https://github.com/air-verse/air

use air with the config file added in the repo

```bash
air -c .air.linux.conf
```

Create .env file

```
cp .env-template .env
```

