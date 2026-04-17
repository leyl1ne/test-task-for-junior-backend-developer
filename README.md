# Task Service

Сервис для управления задачами с HTTP API на Go.

## Требования

- Go `1.23+`
- Docker и Docker Compose

## Быстрый запуск через Docker Compose

```bash
docker compose up --build
```

После запуска сервис будет доступен по адресу `http://localhost:8080`.

Если `postgres` уже запускался ранее со старой схемой, пересоздай volume:

```bash
docker compose down -v
docker compose up --build
```

Причина в том, что SQL-файл из `migrations/0001_create_tasks.up.sql` монтируется в `docker-entrypoint-initdb.d` и применяется только при инициализации пустого data volume.

## Swagger

Swagger UI:

```text
http://localhost:8080/swagger/
```

OpenAPI JSON:

```text
http://localhost:8080/swagger/openapi.json
```

## API

Базовый префикс API:

```text
/api/v1
```

Основные маршруты:

- `POST /api/v1/tasks`
- `GET /api/v1/tasks`
- `GET /api/v1/tasks/{id}`
- `PUT /api/v1/tasks/{id}`
- `DELETE /api/v1/tasks/{id}`


## Реализованный функционал в рамках тестового задания

### 1. Поддержка расписаний (Schedule)

Добавлена новая сущность `Schedule`, описывающая правила генерации задач.

Поддерживаются следующие типы периодичности:

- **Daily** — каждые N дней
- **Monthly** — в определённые дни месяца (1–30)
- **Specific Dates** — строго заданные даты
- **Even/Odd** — чётные или нечётные дни месяца

---

### 2. Генерация задач по расписанию

Реализован background scheduler, который:

- периодически опрашивает базу данных
- определяет, нужно ли создавать задачу на текущую дату
- создаёт задачи на основе расписания

---

### 3. Идемпотентность

Для предотвращения дублирования задач реализовано:

- добавлено поле `schedule_id` в таблицу `tasks`
- добавлено поле `date` (дата выполнения)
- создан уникальный индекс:

```sql
UNIQUE(schedule_id, date)
```

Это гарантирует, что задача для одного расписания не будет создана более одного раза в день.

---

### 4. Обработка ошибок базы данных

Добавлена обработка ошибки уникальности PostgreSQL:

- используется код ошибки `23505`
- маппится в доменную ошибку `ErrTaskAlreadyCreated`

---

### 5. Расширение API

Метод создания задачи `(POST /api/v1/tasks)` теперь поддерживает:

- создание обычной задачи (с date)
- создание задачи с расписанием (schedule)

---

### 6. Scheduler

Реализован фоновый процесс:

- запускается вместе с приложением
- работает по `ticker`
- выполняет проверку расписаний
- поддерживает graceful shutdown через `context`

---

### 7. Архитектурные решения
- `Schedule` выделен как отдельная доменная сущность
- логика определения запуска (`ShouldRunToday`) находится в domain
- scheduler вынесен в отдельный пакет `internal/scheduler`
- repository слой расширен без нарушения обратной совместимости

---

### Пример создания задачи с расписанием
```json
{
  "title": "Call patients",
  "description": "Daily calls",
  "status": "new",
  "schedule": {
    "type": "daily",
    "interval": 1,
    "start_date": "2026-04-17"
  }
}
```

---

### Пример обычной задачи
```json
{
  "title": "One-time task",
  "description": "Do something",
  "status": "new",
  "date": "2026-04-17"
}
```