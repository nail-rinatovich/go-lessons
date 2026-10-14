# Project 7 — использование

## Что делает проект
REST API для каталога книг на Gin + SQLite (Lesson 8). Поддерживает полный CRUD, валидацию, пагинацию, фильтрацию, сортировку и простую аутентификацию по API-ключу. Данные хранятся в `books.db`.

## Запуск
```bash
cd project7-rest-api
go run .
```
При первом запуске в базу добавляются три книги. Ключ API по умолчанию — `secret-key`. Свой ключ можно задать переменной окружения:
```bash
API_KEY=my-key go run .
```
В PowerShell: `$env:API_KEY="my-key"; go run .`

## Ресурс «книга»
```json
{"id": 1, "title": "The Go Programming Language", "author": "Alan Donovan", "year": 2015, "isbn": "978-0134190440"}
```
Правила валидации для POST и PUT: `title` — обязателен, до 200 символов; `author` — обязателен, до 100 символов; `year` — обязателен, от 1450 до 2100; `isbn` — необязателен, до 20 символов.

## Эндпоинты

| Метод | Путь | Ключ | Успех | Ошибки |
|-------|------|------|-------|--------|
| GET | `/books` | нет | `200` страница книг | `400` неверные параметры |
| GET | `/books/:id` | нет | `200` книга | `400`, `404` |
| POST | `/books` | да | `201` книга + заголовок `Location` | `400`, `401` |
| PUT | `/books/:id` | да | `200` обновлённая книга | `400`, `401`, `404` |
| DELETE | `/books/:id` | да | `204` без тела | `400`, `401`, `404` |

Все ошибки имеют формат `{"error": "описание"}`.

### Параметры `GET /books`

| Параметр | По умолчанию | Описание |
|----------|--------------|----------|
| `page` | `1` | номер страницы, ≥ 1 |
| `limit` | `10` | книг на странице, 1–100 |
| `author` | — | точное совпадение автора |
| `q` | — | подстрока в названии (без учёта регистра для латиницы) |
| `sort` | `id` | `id`, `title`, `author` или `year` |
| `order` | `asc` | `asc` или `desc` |

Ответ:
```json
{"data": [ ...книги... ], "page": 1, "limit": 10, "total": 3}
```
`total` — общее число книг, подходящих под фильтр. Число страниц = `ceil(total / limit)`.

### Примеры
```bash
# Чтение (ключ не нужен)
curl http://localhost:8080/books
curl http://localhost:8080/books/1
curl "http://localhost:8080/books?sort=year&order=desc"
curl "http://localhost:8080/books?page=2&limit=1"
curl "http://localhost:8080/books?q=go&author=Jon%20Bodner"

# Создание
curl -i -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -H "X-API-Key: secret-key" \
  -d '{"title":"100 Go Mistakes","author":"Teiva Harsanyi","year":2022}'

# Полное обновление (PUT заменяет все поля)
curl -X PUT http://localhost:8080/books/4 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: secret-key" \
  -d '{"title":"100 Go Mistakes and How to Avoid Them","author":"Teiva Harsanyi","year":2022}'

# Удаление
curl -i -X DELETE http://localhost:8080/books/4 -H "X-API-Key: secret-key"

# Ошибки
curl -i -X POST http://localhost:8080/books -H "Content-Type: application/json" \
  -d '{"title":"X","author":"Y","year":2020}'                 # 401: нет ключа
curl -i -X POST http://localhost:8080/books -H "Content-Type: application/json" \
  -H "X-API-Key: secret-key" -d '{"title":"X","author":"Y"}'  # 400: нет year
curl -i "http://localhost:8080/books?sort=price"              # 400: недопустимая сортировка
curl -i http://localhost:8080/books/999                       # 404
```

Для Postman: создайте запросы с теми же URL, во вкладке **Headers** добавьте `X-API-Key: secret-key`, во вкладке **Body** выберите `raw` → `JSON`.

## Как это работает под капотом
1. **Структура кода.**
   - `store.go` — слой доступа к данным (`BookStore`): SQL, никакого HTTP.
   - `main.go` — HTTP-слой (`BookHandler`): разбор запроса, валидация, коды ответа.

   Такое разделение позволяет поменять базу, не трогая обработчики, и наоборот.
2. **REST.** Ресурс — книга. Коллекция находится по адресу `/books`, элемент — по `/books/:id`. Действие определяется HTTP-методом: GET читает, POST создаёт, PUT заменяет, DELETE удаляет. Коды ответов несут смысл: `201 Created` + `Location` указывает адрес новой книги, `204 No Content` означает «удалено, тела нет», `404` — «нет такой книги».
3. **Параметры маршрута.** `c.Param("id")` достаёт `:id` из пути, `strconv.ParseInt` проверяет, что это число.
4. **Валидация.** `c.ShouldBindJSON(&in)` разбирает JSON в `bookInput` и проверяет теги `binding:"required,min=...,max=..."`. Для валидации используется отдельная входная структура, а не `Book`, поэтому клиент не может подменить `id`.
5. **Аутентификация через middleware.**
   ```go
   protected := r.Group("/books", apiKeyAuth(apiKey))
   ```
   Группа маршрутов с middleware: `apiKeyAuth` выполняется до обработчика. Если заголовок `X-API-Key` неверный, `c.AbortWithStatusJSON(401, ...)` прерывает цепочку, и обработчик не вызывается. GET-маршруты зарегистрированы вне группы, поэтому они публичные.
6. **Динамический SQL без инъекций.** В `List` условия `WHERE` собираются из кусков `"author = ?"`, а значения кладутся в `args` и передаются плейсхолдерами. Имена колонок в `ORDER BY` плейсхолдерами передать нельзя, поэтому `sort` проверяется по белому списку `sortColumns`, а `order` — только `asc` или `desc`. Всё остальное отклоняется с `400`.
7. **Пагинация.** `LIMIT ? OFFSET ?`, где `OFFSET = (page - 1) * limit`. Отдельный запрос `SELECT COUNT(*)` с тем же `WHERE` даёт `total`.
8. **Обработка ошибок.** Слой данных возвращает `ErrNotFound`. `storeError` проверяет его через `errors.Is` и отвечает `404`. Прочие ошибки логируются на сервере, а клиент получает общее `500 internal server error`, без внутренних деталей.
