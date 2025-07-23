# Сервис для учета онлайн-подписок

RESTful-сервис для агрегации данных об онлайн-подписках пользователей.

## Стек технологий

*   **Язык:** Go 1.24
*   **Веб-фреймворк:** Gin
*   **База данных:** PostgreSQL 16
*   **Драйвер БД:** pgx/v5
*   **Миграции:** golang-migrate
*   **Конфигурация:** Viper
*   **Логирование:** slog (стандартная библиотека)
*   **Документация API:** Swag (Swagger)
*   **Контейнеризация:** Docker & Docker Compose

## Архитектура

Проект построен с использованием принципов чистой архитектуры (Clean Architecture) с четким разделением на слои:

-   **Handler (Обработчик):** Слой, отвечающий за прием HTTP-запросов, валидацию данных и отправку ответов.
-   **Service (Сервис):** Слой, содержащий основную бизнес-логику приложения.
-   **Repository (Репозиторий):** Слой, отвечающий за взаимодействие с базой данных.

## Запуск проекта

### Предварительные требования
-   Установленный [Docker Desktop](https://www.docker.com/get-started/)
-   Установленный [Git](https://git-scm.com/downloads/)

### Инструкция по запуску

1.  **Клонируйте репозиторий:**
    ```bash
    git clone <URL_РЕПОЗИТОРИЯ>
    cd <ПАПКА_ПРОЕКТА>
    ```

2.  **Создайте конфигурационный файл `.env`:**
    В корне проекта создайте файл `.env` и скопируйте в него содержимое ниже. Эти значения будут использованы для инициализации базы данных.

    ```dotenv
    # PostgreSQL Credentials
    POSTGRES_USER=user
    POSTGRES_PASSWORD=strongpassword
    POSTGRES_DB=subscriptions_db
    POSTGRES_PORT=5432
    ```

3.  **Запустите проект с помощью Docker Compose:**
    Выполните команду в корневой директории проекта:

    ```bash
    docker compose up --build -d
    ```
    *   `--build`: Эта опция соберет Docker-образ для Go-приложения при первом запуске.
    *   `-d`: Запустит контейнеры в фоновом режиме.

    При первом запуске команда автоматически соберет Go-приложение, запустит контейнеры с приложением и базой данных. Миграции базы данных необходимо применить вручную (см. раздел "Работа с миграциями").

## API Документация (Swagger)

После успешного запуска сервиса интерактивная документация API будет доступна в вашем браузере по адресу:

**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

В интерфейсе Swagger можно не только изучить все доступные эндпоинты, но и выполнять тестовые запросы.

## Работа с миграциями

Для управления схемой базы данных используется `golang-migrate`.

**Установка (если требуется):**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Применить все миграции:**
```bash
migrate -path migrations -database 'postgres://user:strongpassword@localhost:5432/subscriptions_db?sslmode=disable' up
```

**Откатить последнюю миграцию:**
```bash
migrate -path migrations -database 'postgres://user:strongpassword@localhost:5432/subscriptions_db?sslmode=disable' down 1
```