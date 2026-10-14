# Project 9 — использование

## Что делает проект
Полноценный менеджер задач (Lesson 10). Фронтенд на React общается с REST API на Gin, а задачи хранятся в SQLite. Можно создавать, редактировать, отмечать выполненными и удалять задачи, а также фильтровать их по статусу.

## Запуск

### Режим разработки (два терминала)
```bash
cd project9-fullstack/backend
go run .
```
```bash
cd project9-fullstack/frontend
npm install
npm run dev
```
Откройте **http://localhost:5173**. Изменения во фронтенде видны сразу. Изменения в Go требуют перезапустить `go run .`.

### Production (один сервер)
```bash
cd project9-fullstack/frontend
npm run build
```
```bash
cd project9-fullstack/backend
go run .
```
Откройте **http://localhost:8080**: Gin отдаёт и API, и собранный React.

### Переменные окружения бэкенда

| Переменная | По умолчанию | Назначение |
|------------|--------------|------------|
| `PORT` | `8080` | порт сервера |
| `DB_PATH` | `./tasks.db` | файл базы |
| `CORS_ORIGIN` | `http://localhost:5173` | origin, которому разрешён CORS |
| `FRONTEND_DIST` | `../frontend/dist` | папка со сборкой React |

## Как взаимодействовать через интерфейс
- **Добавить задачу.** Заполните название (обязательно, до 100 символов) и описание (до 500), затем нажмите **Add task** или Enter. Ошибки валидации показываются под полями.
- **Выполнено / не выполнено.** Чекбокс слева от задачи.
- **Редактировать.** **Edit** → изменить поля → **Save** (или Enter); **Cancel** отменяет изменения.
- **Удалить.** **Delete**.
- **Фильтр.** **All** / **Active** / **Done**.
- Если бэкенд недоступен, появится красное сообщение об ошибке.

## REST API

| Метод | Путь | Тело (JSON) | Успех | Ошибки |
|-------|------|-------------|-------|--------|
| GET | `/api/health` | — | `200` `{"status":"ok"}` | — |
| GET | `/api/tasks?status=all\|active\|done` | — | `200` массив задач | `400` неверный status |
| GET | `/api/tasks/:id` | — | `200` задача | `400`, `404` |
| POST | `/api/tasks` | `{"title", "description"}` | `201` созданная задача | `400` |
| PUT | `/api/tasks/:id` | `{"title", "description", "done"}` | `200` обновлённая задача | `400`, `404` |
| DELETE | `/api/tasks/:id` | — | `204` без тела | `400`, `404` |

Задача:
```json
{"id": 1, "title": "Learn Gin", "description": "Lesson 10", "done": false, "createdAt": "2026-10-12T12:42:45Z"}
```
Список отсортирован от новых к старым. Ошибки имеют формат `{"error": "описание"}`.

### Примеры (бэкенд на порту 8080)
```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/tasks
curl "http://localhost:8080/api/tasks?status=active"

curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Gin","description":"Lesson 10"}'

# PUT заменяет все поля, поэтому передавайте и title, и description, и done
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Gin","description":"Lesson 10","done":true}'

curl -i -X DELETE http://localhost:8080/api/tasks/1

# Ошибки
curl -i -X POST http://localhost:8080/api/tasks -H "Content-Type: application/json" -d '{"title":"   "}'  # 400
curl -i http://localhost:8080/api/tasks/999                                                             # 404

# Проверка CORS: preflight-запрос, который браузер шлёт перед кросс-доменным PUT/DELETE
curl -i -X OPTIONS http://localhost:8080/api/tasks -H "Origin: http://localhost:5173"
```
Через dev-сервер Vite те же запросы можно делать на порт 5173: `curl http://localhost:5173/api/tasks`. Vite перенаправит их в Gin.

## Как это работает под капотом

### Путь запроса
```
Режим разработки:
  браузер ──> Vite :5173 ──(/api/*, proxy)──> Gin :8080 ──> SQLite
                  └─ остальное: React с hot reload

Production:
  браузер ──> Gin :8080 ─┬─ /api/*     → обработчики → SQLite
                         ├─ /assets/*  → файлы из frontend/dist/assets
                         └─ всё прочее → frontend/dist/index.html
```

### Бэкенд (`backend/`)
1. **`db.go` — слой данных.** `TaskStore` с методами `List/Get/Create/Update/Delete`. SQL выполняется только с плейсхолдерами `?`. Если `UPDATE` или `DELETE` затронул 0 строк, метод возвращает `ErrNotFound`. Драйвер `modernc.org/sqlite` написан на чистом Go и не требует gcc. `created_at DEFAULT CURRENT_TIMESTAMP` проставляет время создания на стороне базы.
2. **`handlers.go` — HTTP-слой.** `ShouldBindJSON` + теги `binding:"required,max=100"` проверяют вход, затем пробелы обрезаются, и название из одних пробелов отклоняется. `ErrNotFound` превращается в `404`, прочие ошибки логируются, а клиент получает `500` без деталей.
3. **`main.go` — сборка приложения.**
   - Маршруты API объединены в группу `r.Group("/api")`.
   - **CORS-middleware.** Браузер разрешает JavaScript со страницы `localhost:5173` читать ответы от `localhost:8080` (другой порт = другой origin), только если сервер вернёт заголовок `Access-Control-Allow-Origin`. Для PUT/DELETE и JSON браузер сначала шлёт **preflight** `OPTIONS`, и middleware сразу отвечает на него `204`.
   - **Раздача фронтенда.** `r.Static("/assets", ...)` отдаёт JS и CSS. `r.NoRoute` возвращает `index.html` на любой неизвестный путь, кроме `/api/*`: так работает клиентская маршрутизация SPA. Если сборки нет, сервер пишет предупреждение в лог и работает только как API.

### Фронтенд (`frontend/`)
1. **Прокси Vite** (`vite.config.js`): `'/api': 'http://localhost:8080'`. Браузер обращается к тому же origin (5173), а Vite пересылает запрос в Gin. Поэтому в режиме разработки CORS фактически не участвует, а код использует относительный адрес `/api/...`, который работает и в разработке, и в production.
2. **`api.js`** — обёртка над `fetch`: ставит `Content-Type: application/json`, для `204` возвращает `null`, при статусе 4xx/5xx бросает `Error` с текстом из `{"error": ...}`.
3. **`App.jsx`** — состояние `tasks`, `filter`, `loading`, `error`.
   - `useEffect([filter])` перезагружает список при смене фильтра. Флаг `ignore` отбрасывает устаревший ответ, если пользователь успел переключить фильтр до прихода предыдущего ответа (защита от гонки запросов).
   - После создания, изменения или удаления список обновляется локально по ответу сервера, без повторной загрузки всего списка. Задача, которая перестала подходить под текущий фильтр (например, отмечена выполненной во вкладке Active), сразу убирается из списка.
4. **`TaskForm.jsx`** — управляемая форма. Функция `validate` повторяет серверные правила, поэтому ошибка показывается мгновенно, без запроса к серверу. Серверная проверка всё равно остаётся: клиентской доверять нельзя. Пока идёт запрос, кнопка заблокирована. После успеха очищается только отправленный текст: если пользователь успел набрать новую задачу, она не сотрётся.
5. **`TaskItem.jsx`** — у каждой задачи свой режим просмотра или редактирования и флаг `busy`, который блокирует кнопки на время запроса.

### Production-сборка
`npm run build` минифицирует JS и CSS и добавляет к именам файлов хэш содержимого (`index-D9iFJi6V.js`). Благодаря этому браузер может кэшировать их навсегда, а после новой сборки автоматически загрузит новые файлы. Для деплоя достаточно Go-бинарника (`go build`) и папки `frontend/dist`. Запускайте с `GIN_MODE=release`, чтобы отключить отладочный вывод Gin.
