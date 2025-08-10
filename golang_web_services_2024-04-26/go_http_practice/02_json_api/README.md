# Практика 02: HTTP+JSON API

## Задание

Реализуй http-сервер с маршрутом POST /echojson, который:

1. Принимает JSON-объект вида {"message": "текст"} в теле POST-запроса
2. Возвращает тот же объект в ответе (Content-Type: application/json), только текст должен быть в верхнем регистре.

## Пример curl-запроса:

```sh
curl -X POST -d '{"message": "привет"}' \
-H "Content-Type: application/json" \
http://localhost:8080/echojson
```

## Ожидаемый ответ:

```
{"message": "ПРИВЕТ"}
```

Если запрос пустой или без поля message — возвращай http 
400 и json-ошибку {"error": "bad request"}.

---

main.go — шаблон уже в папке, начни писать туда.