# Project 6 — использование

## Что делает проект
Список задач с хранением в базе SQLite (Lesson 7). HTTP-сервер на `net/http` позволяет получать, создавать, отмечать выполненными и удалять задачи. Данные лежат в файле `tasks.db` и сохраняются между перезапусками.

## Запуск
```bash
cd project6-databases
go run .
```
При первом запуске создаётся `tasks.db` с двумя задачами («Learn Go», «Build a web server»). Чтобы начать с чистой базы, остановите сервер и удалите `tasks.db`.

## Эндпоинты
Данные в POST и PUT передаются как **форма** (`-d key=value`), а не как JSON. Ответы приходят в JSON.

| Метод | Путь | Данные | Ответ |
|-------|------|--------|-------|
| GET | `/tasks` | — | `200` массив задач |
| POST | `/tasks` | `title` (обязательно), `description` | `201` созданная задача / `400 Title is required` |
| PUT | `/tasks/{id}` | `done=true` или `done=false` | `200` `{"id":1,"done":true}` / `404 Task not found` |
| DELETE | `/tasks/{id}` | — | `200` `{"deleted":true}` / `404 Task not found` |
| любой | `/tasks/abc` | — | `400 Invalid task ID` |
| другой метод | `/tasks` или `/tasks/{id}` | — | `405 Method not allowed` |

Формат задачи:
```json
{"id": 1, "title": "Learn Go", "description": "Complete the Go web programming course", "done": false}
```

### Примеры
```bash
# Все задачи
curl http://localhost:8080/tasks

# Создать задачу
curl -X POST http://localhost:8080/tasks -d "title=Read a book" -d "description=Chapter 1"

# Отметить выполненной / вернуть обратно
curl -X PUT http://localhost:8080/tasks/1 -d "done=true"
curl -X PUT http://localhost:8080/tasks/1 -d "done=false"

# Удалить
curl -X DELETE http://localhost:8080/tasks/1

# Ошибки
curl -i -X POST http://localhost:8080/tasks -d "title="     # 400
curl -i -X DELETE http://localhost:8080/tasks/999           # 404

# Кавычки в названии — JSON остаётся корректным
curl -X POST http://localhost:8080/tasks --data-urlencode 'title=say "hi"'
```

Содержимое базы можно посмотреть любым SQLite-клиентом, например DB Browser for SQLite: откройте `tasks.db`, таблица `tasks`.

## Как это работает под капотом
1. **Драйвер.** `database/sql` — общий интерфейс для любых СУБД, а конкретную базу подключает драйвер. Строка `_ "modernc.org/sqlite"` импортирует пакет только ради его `init()`, который регистрирует драйвер под именем `"sqlite"`. Затем `sql.Open("sqlite", "./tasks.db")` его использует.
   - Используется `modernc.org/sqlite` — SQLite, переписанный на чистый Go. Популярный `github.com/mattn/go-sqlite3` — обёртка над C-кодом, ему нужны gcc и `CGO_ENABLED=1`.
2. **Пул соединений.** `sql.DB` — это не одно соединение, а пул: соединения открываются лениво и переиспользуются. `SetMaxOpenConns(1)` ограничивает пул одним соединением, потому что SQLite допускает только одного писателя одновременно.
3. **Схема.** `CREATE TABLE IF NOT EXISTS` создаёт таблицу при первом запуске. `INTEGER PRIMARY KEY AUTOINCREMENT` — автоинкрементный ID.
4. **Три вида запросов:**
   - `db.Query` — много строк (`SELECT ... FROM tasks`). Результат читается в цикле `rows.Next()` + `rows.Scan(...)`, затем нужен `rows.Close()` (через `defer`) и проверка `rows.Err()`;
   - `db.QueryRow` — одна строка (`SELECT COUNT(*)`);
   - `db.Exec` — без результата (`INSERT`, `UPDATE`, `DELETE`). Возвращает `LastInsertId()` и `RowsAffected()`.
5. **Подготовленные выражения.** При заполнении стартовых данных `db.Prepare` один раз компилирует SQL, после чего `stmt.Exec(...)` вызывается с разными значениями.
6. **Защита от SQL-инъекций.** Значения передаются через плейсхолдеры `?`, а не склеиваются со строкой SQL. Драйвер передаёт их отдельно от текста запроса, поэтому ввод вида `'; DROP TABLE tasks; --` останется просто строкой.
7. **404.** Если `UPDATE` или `DELETE` затронул 0 строк (`RowsAffected() == 0`), значит задачи нет, и функция возвращает `sql.ErrNoRows`, который превращается в `404`.
8. **JSON.** `writeJSON` использует `encoding/json`, который корректно экранирует кавычки и спецсимволы. Теги `json:"title"` задают имена полей. `tasks := []Task{}` вместо `var tasks []Task` гарантирует, что пустой список сериализуется как `[]`, а не `null`.
9. **Маршрутизация без фреймворка.** `/tasks` обрабатывает коллекцию, `/tasks/` — отдельные задачи. ID берётся из пути через `strings.TrimPrefix` и `strconv.ParseInt`, а метод выбирается в `switch r.Method`. Сравните с Gin в project7, где это делает роутер.
