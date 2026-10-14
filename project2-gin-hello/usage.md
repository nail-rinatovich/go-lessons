# Project 2 — использование

## Что делает проект
JSON-API для блога на фреймворке Gin (Lesson 6). Посты хранятся в памяти: можно получить список и добавить новый. После перезапуска сервера список возвращается к двум исходным постам.

## Запуск
```bash
cd project2-gin-hello
go run .
```
Сервер слушает http://localhost:8080. При старте Gin печатает таблицу зарегистрированных маршрутов, а затем лог каждого запроса.

## Эндпоинты

| Метод | Путь | Тело запроса | Ответ |
|-------|------|--------------|-------|
| GET | `/` | — | `200` `{"message": "Hello, World! ..."}` |
| GET | `/posts` | — | `200` `{"posts": [...]}` |
| POST | `/posts` | JSON `{"title": "...", "body": "..."}` | `201` `{"message": "Post created", "post": {...}}` или `400` `{"error": "..."}` |

Оба поля `title` и `body` обязательны.

### Примеры
```bash
# Приветствие
curl http://localhost:8080/

# Список постов
curl http://localhost:8080/posts

# Создать пост
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"New Post","body":"Post content here"}'

# Ошибка валидации: нет body → 400
curl -i -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Only title"}'

# Некорректный JSON → 400
curl -i -X POST http://localhost:8080/posts -d 'not json'
```

Пример ответа `GET /posts`:
```json
{"posts":[{"id":1,"title":"First Post","body":"This is the first post."},{"id":2,"title":"Second Post","body":"This is the second post."}]}
```

> В Windows PowerShell одинарные кавычки и `\` для переноса строки ведут себя иначе. Запускайте эти команды в Git Bash или пишите всю команду в одну строку.

## Как это работает под капотом
1. `gin.Default()` создаёт роутер с двумя middleware:
   - **Logger** — печатает каждую строку `[GIN] ... | 200 | GET "/posts"`;
   - **Recovery** — если обработчик упадёт с `panic`, сервер не завершится, а ответит `500`.
2. `r.GET(...)` и `r.POST(...)` регистрируют маршруты. В отличие от `net/http`, Gin сам проверяет HTTP-метод: например, `PUT /posts` получит `404`.
3. Обработчик получает `*gin.Context`. В нём собраны запрос, ответ и вспомогательные методы.
4. `c.JSON(status, obj)` сериализует Go-значение в JSON и выставляет `Content-Type: application/json`. `gin.H` — это просто `map[string]any`.
5. `c.ShouldBindJSON(&newPost)` читает тело запроса, разбирает JSON в структуру `Post` и запускает валидацию по тегам `binding:"required"` (библиотека `go-playground/validator`). При ошибке возвращает `err`, и мы отвечаем `400`.
6. Теги `json:"title"` задают имена полей в JSON.
7. Новый ID — это ID последнего поста + 1.
8. **Конкурентность**: Gin обрабатывает запросы параллельно, в разных горутинах. Глобальный срез `posts` защищён `sync.Mutex`, иначе два одновременных POST могли бы получить одинаковый ID или повредить срез.
9. `r.Run(":8080")` внутри вызывает тот же `http.ListenAndServe`, что и в project1.
