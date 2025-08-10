# Практика 03: CRUD REST API (In-memory)

## Задание

Реализуй http-сервер на Go, который:

- Работает с сущностями типа Note (заметка):
    - id (целое число, auto increment)
    - text (строка)

- Поддерживает маршруты:
    - POST /notes — создать заметку (json: {"text": "..."})   — возвращает note с id
    - GET /notes — получить все заметки (json: массив)
    - GET /notes/{id} — получить одну заметку по id (json)
    - DELETE /notes/{id} — удалить заметку по id

- Все данные хранятся в памяти (map[int]Note или slice)
- Ответы — только application/json
- Для всех ошибок возвращай http 400 и json вида {"error": "сообщение ошибки"}

## Пример:

POST /notes   (body {"text":"abc"})
Ответ: {"id":1, "text":"abc"}

GET /notes
Ответ: [{...}, {...}]

GET /notes/1
Ответ: {"id":1, "text":"abc"}

DELETE /notes/1
Ответ: {"result": "deleted"}

---

main.go — шаблон уже в папке, начни писать туда.