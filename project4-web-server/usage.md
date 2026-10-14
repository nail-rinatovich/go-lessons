# Project 4 — использование

## Что делает проект
Веб-сервер на стандартном пакете `net/http` (Lesson 4). Показывает страницу с текущим временем сервера и форму. После отправки формы сервер здоровается с пользователем по имени. Также раздаёт статические файлы (CSS).

## Запуск
```bash
cd project4-web-server
go run .
```
Запускайте из папки проекта: путь к `./static` считается от текущей директории.

## Эндпоинты

| Метод | Путь | Данные | Ответ |
|-------|------|--------|-------|
| GET | `/` | — | `200` HTML-страница со временем и формой |
| GET | любой другой путь, кроме перечисленных | — | `404 page not found` |
| POST | `/greet` | форма `name=...` (`application/x-www-form-urlencoded`) | `200` HTML `Hello, <name>!` |
| POST | `/greet` | пустое `name` | `400 Name is required` |
| GET, PUT, … | `/greet` | — | `405 Method not allowed` |
| GET | `/static/style.css` | — | `200` файл из папки `static/` |

### Браузер
1. Откройте http://localhost:8080. Время на странице обновляется при каждой перезагрузке.
2. Введите имя и нажмите **Greet Me**. Браузер отправит `POST /greet`.

### curl
```bash
# Главная страница
curl http://localhost:8080/

# Отправить форму (так же, как это делает браузер)
curl -X POST http://localhost:8080/greet -d "name=Alice"

# Пустое имя → 400
curl -i -X POST http://localhost:8080/greet -d "name="

# GET на /greet → 405
curl -i http://localhost:8080/greet

# Статический файл
curl http://localhost:8080/static/style.css

# Неизвестный путь → 404
curl -i http://localhost:8080/nope

# Защита от XSS: тег будет выведен как текст, а не выполнен
curl -X POST http://localhost:8080/greet -d "name=<script>alert(1)</script>"
# → <h1>Hello, &lt;script&gt;alert(1)&lt;/script&gt;! ...
```

## Как это работает под капотом
1. **Маршрутизация.** Три обработчика регистрируются в `http.DefaultServeMux`:
   - `/` — шаблон-«ловушка» для всех путей, поэтому внутри `helloHandler` стоит проверка `r.URL.Path != "/"`, которая отдаёт 404;
   - `/greet` — точное совпадение;
   - `/static/` — префикс: всё, что начинается с `/static/`.
2. **Динамический HTML.** `helloHandler` подставляет время в HTML-строку через `fmt.Sprintf`. Формат `"2006-01-02 15:04:05"` — особенность Go: даты форматируются по образцу эталонной даты 2 января 2006, 15:04:05.
3. **Заголовки.** `w.Header().Set("Content-Type", "text/html; charset=utf-8")` нужно вызвать **до** первой записи в тело ответа, потом заголовки уже отправлены.
4. **Проверка метода.** `net/http` не фильтрует методы сам, поэтому `greetHandler` проверяет `r.Method` и отвечает `405`.
5. **Разбор формы.** `r.ParseForm()` читает тело запроса и строку запроса. `r.FormValue("name")` возвращает значение поля или пустую строку.
6. **Ошибки.** `http.Error(w, msg, code)` выставляет статус и пишет текст ошибки.
7. **XSS.** Пользовательский ввод экранируется через `html.EscapeString`: `<` превращается в `&lt;` и т. д. Без этого злоумышленник мог бы внедрить на страницу JavaScript.
8. **Статика.** `http.FileServer(http.Dir("./static"))` отдаёт файлы с диска. `http.StripPrefix("/static/", ...)` отрезает префикс, чтобы запрос `/static/style.css` искал файл `./static/style.css`, а не `./static/static/style.css`. `FileServer` сам ставит `Content-Type` по расширению и поддерживает кэширование (`Last-Modified`, `304 Not Modified`).
