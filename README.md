# Task-Test-devBack
This is a repository for a test Task for the MEDODS company.

# Сборка

1) Добавить .env файл в репозиторий, изменить конфиг (для работы smtp сервера)
в .env файл нужно добавить
POSTGRES_PASSWORD=
FROM_EMAIL=
SMTP_USER=
SMTP_PASSWORD=

в конфиге поменять сервер или можно оставить smtp.mailersend.net

Затем запускаем автосборку

```
docker-compose up
```

# Документация 

по localhost:8000 будет доступна swagger документация, можно посмотреть работу сервиса.

# Тесты

тесты не выносли в отдельную папку, поэтому нужно перейте в /internal/service/

``` go

go test

```
