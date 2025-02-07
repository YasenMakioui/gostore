# gostore

Project used to learn GO.

A fiber app used to manage files in a server using HTTP requests.

With gostore you have a file store that you can manage using filesystem paths inserted in the url similar to S3

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

GOSTOREPATH references the api base url
BASEDIR is the directory used to save and manage files