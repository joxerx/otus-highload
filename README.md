# Highload Architect

- [Начало работы](#начало-работы)
- [Makefile](#makefile)


## Начало работы
После клонирования необходимо проинициализировать проект.
```
make init
```
Затем нужно указать значения переменных в ```.env```.

Когда установлены подходящие значения, можно запустить проект
```
make build
```
После запуска и подключения приложения к БД, необходимо выполнить миграции
```
make migrate-up  
```
Когда миграции выполнены успешно, приложение доступно к использованию на порте указанном в `DOCKER_EXPOSE_PORT`



### Makefile
Список команд:

```
make init
make build
make down
make stop
make web-logs
make logs
make migrate-up
make migrate-down
make pg-shell
```