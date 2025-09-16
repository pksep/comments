# Работа с миграциями

В этом руководстве описано, как работать с миграциями базы данных в проекте.

## Установка инструментов миграции

### Установка golang-migrate

#### macOS (с помощью Homebrew)
```bash
brew install golang-migrate
```

#### Linux (смотреть актуальный релиз на GitHub)
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.21.3/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate.linux-amd64 /usr/local/bin/migrate
```

## Создание новой миграции

1. Создайте новую миграцию с помощью команды:
   ```bash
   make migration-create NAME=description_of_changes
   ```
   Это создаст два файла в директории `migrations/`:
   - `{timestamp}_description_of_changes.up.sql` - для применения миграции
   - `{timestamp}_description_of_changes.down.sql` - для отката миграции

2. Отредактируйте созданные файлы, добавив SQL-запросы:
   - В `.up.sql` - запросы для применения изменений
   - В `.down.sql` - запросы для отката изменений

## Применение миграций
- Миграции применяются автоматически при запуске сервера

## Рекомендации по созданию миграций

1. Каждая миграция должна быть атомарной и идемпотентной
2. Всегда создавайте откатываемые миграции (down)
3. Тестируйте миграции на тестовой базе данных перед применением в production
4. Храните миграции в системе контроля версий
5. Не изменяйте уже примененные миграции - создавайте новую миграцию для изменений

## Пример миграции

`20230911120000_create_users_table.up.sql`:
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    initials TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
```

`20230911120000_create_users_table.down.sql`:
```sql
DROP TABLE IF EXISTS users;
```
