# Task-Test-devBack
This is a repository for a test Task for the MEDODS company.

# Сборка

1) В .env файл нужно добавить
```
POSTGRES_PASSWORD=
FROM_EMAIL=
SMTP_USER=
SMTP_PASSWORD=
```

2) В конфиге по пути configs/config.yml меняем настройки сервера

3) Запускаем автосборку

```
docker-compose up
```

# Документация 

В браузере будет доступна swagger документация, можно посмотреть работу сервиса.
```
localhost:8000 
```

# Тесты

Тесты не выносил в отдельную папку, поэтому нужно перейти в /internal/service/

``` go

go test

```
