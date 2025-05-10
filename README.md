# 🚀 Распределённый вычислитель арифметических выражений

<p align="center">
  <img src="/assets/GOPHER_ACADEMY.png" alt="Gopher Logo" width="200"/>
</p>
<p align="center">
  <img alt="GoLang" src="https://img.shields.io/badge/Go-v1.24-blue?style=flat-square&logo=go"/>
  <img alt="Fiber" src="https://img.shields.io/badge/Fiber-v2.52.6-orange?style=flat-square"/>
</p>
<p align="center">
  <a href="https://t.me/SOKRAT_00">
    <img src="https://img.shields.io/badge/Telegram-sokrat_00-blue?style=for-the-badge&logo=telegram" alt="Telegram"/>
  </a>
</p>

Этот проект представляет распределённый вычислитель арифметических выражений. Система разделена на две основные части:

- **Оркестратор (Server):** Принимает арифметическое выражение от пользователя, разбивает его на последовательные задачи и следит за их выполнением.
- **Агент (Agent):** Получает задачи от оркестратора, выполняет вычисления (каждая операция обрабатывается отдельно с имитацией значительных вычислительных затрат) и отправляет результаты обратно.

> **Примечание:** Для имитации «дорогих» вычислений время выполнения каждой операции задается через переменные среды:
>
> - `TIME_ADDITION_MS` – время выполнения сложения (мс)
> - `TIME_SUBTRACTION_MS` – время выполнения вычитания (мс)
> - `TIME_MULTIPLICATIONS_MS` – время выполнения умножения (мс)
> - `TIME_DIVISIONS_MS` – время выполнения деления (мс)
>
> Количество параллельных вычислений регулируется переменной среды `COMPUTING_POWER`.

---

## 📖 Содержание

- [🚀 Распределённый вычислитель арифметических выражений](#-распределённый-вычислитель-арифметических-выражений)
  - [📖 Содержание](#-содержание)
  - [🖥️ Функционал проекта](#️-функционал-проекта)
    - [1. Оркестратор (Server)](#1-оркестратор-server)
    - [2. Агент (Agent)](#2-агент-agent)
    - [3. Аuth](#3-аuth)
  - [🚀 Быстрый старт](#-быстрый-старт)
  - [🔧 Установка и запуск](#-установка-и-запуск)
  - [📚 Использование API](#-использование-api)
  - [🧪 Тестирование](#-тестирование)
    - [Автоматизированные тесты](#автоматизированные-тесты)
    - [Ручное тестирование](#ручное-тестирование)
  - [⚙️ Переменные окружения](#️-переменные-окружения)
  - [Диаграмма проекта](#диаграмма-проекта)
  - [📂 Структура проекта](#-структура-проекта)
  - [🤝 Контакты](#-контакты)

---

## 🖥️ Функционал проекта

### 1. Оркестратор (Server)

- **Создание вычисления выражения:**
  Пользователь отправляет арифметическое выражение (например, `2+2*2`) через HTTP‑POST запрос. Оркестратор:

  - Разбирает строку, строит AST и генерирует набор задач для каждой операции.
  - Сохраняет выражение с уникальным идентификатором (например, `"expression-1"`).
  - Возвращает код **201** при успешном принятии, **422** при невалидных данных и **500** при внутренней ошибке.

- **Получение списка выражений:**
  Клиент может получить список всех выражений с их статусами (`pending`, `success`, `error`) и итоговыми результатами (если вычисление завершено).
  Код ответа: **200** (успех) или **500** (ошибка).

- **Получение конкретного выражения:**
  Позволяет запросить информацию об отдельном выражении по его идентификатору.
  Коды ответа: **200**, **404** (не найдено) или **500**.

- **Работа с задачами:**
  Оркестратор предоставляет следующие endpoint-ы:

  - **GET** `/internal/task`: возвращает задачу, готовую к выполнению.
  - **POST** `/internal/task/:id/:result`: принимает результат выполнения задачи, обновляет её и, при необходимости, итоговое выражение.
    Возможные коды ответа: **200**, **404**, **422**, **500**.

- **Юнит тесты**

### 2. Агент (Agent) (gRPC)

- **Получение задач:**
  Агент, работающий в виде демона, постоянно опрашивает оркестратора (через `GET /internal/task`) для получения задач.

- **Параллельная обработка:**
  Агент запускает несколько горутин (количество определяется переменной `COMPUTING_POWER`), каждая из которых независимо обрабатывает полученные задачи с заданными задержками (заданными переменными `TIME_*_MS`).

- **Отправка результатов:**
  После выполнения задачи Агент отправляет результат обратно на сервер через `POST /internal/task/:id/:result`, где оркестратор обновляет статус и итоговое значение выражения.

- **Юнит тесты**

### 3. Аuth (gRPC)

- **Авторизация и выдача JWT**

- **Регистрация**

- **Юнит + Интеграционные тесты**

---

## 🚀 Быстрый старт

Оркестратор (Server)

```bash
cd orchestrator
go run ./cmd/orchestrator/main.go
```

Агент (Agent)

```bash
cd agent
go run ./cmd/agent/main.go
```

Сервис аутентификации (Auth)

```bash
cd auth
go run ./cmd/main.go
```

---

## 🔧 Установка и запуск

```bash
git clone https://github.com/0sokrat0/CalcServiceYA.git
```

## 📚 Использование API

Существует готовая [Postman-сборка](https://0sokrat0-3578990.postman.co/workspace/0sokrat0's-Workspace~3b95fdf7-12dd-4ba0-8b9d-45f20c701886/collection/44694863-b36889dd-6f55-483e-8e96-98e309ca248f?action=share&creator=44694863), в которой собраны тестовые запросы.

В этой Postman‑коллекции «CalcYA API» собраны шесть запросов для взаимодействия с сервисом:

    Auth
    — POST {{baseUrl}}/api/login
    Отправка JSON‑тела { login, password } для получения токена доступа.

    Register
    — POST {{baseUrl}}/api/register
    Создание нового пользователя аналогичным JSON‑телом.

    Calculate Expression
    — POST {{baseUrl}}/api/v1/calculate
    Отправка выражения вида { "expression": "(1 + 2) * 3" } с заголовком Authorization: Bearer {{accessToken}}. Порождает новое вычисление.

    List Expressions
    — GET {{baseUrl}}/api/v1/expressions
    Получение списка всех выражений с их статусами и результатами (требует авторизации).

    Get Expression by ID
    — GET {{baseUrl}}/api/v1/expressions/{{expressionId}}
    Подробная информация об одном выражении по его UUID (требует авторизации).

    Get All Tasks (Internal)
    — GET {{baseUrl}}/internal/tasks
    Внутренний endpoint оркестратора: выдаёт очередь готовых к выполнению задач (без авторизации).

Переменные коллекции:

    baseUrl (по умолчанию http://localhost:8080)

    accessToken — используется в заголовке Authorization

    expressionId — подставляется в путь для запроса конкретного выражения

---

## 🧪 Тестирование

### Автоматизированные тесты

Запустите тесты из корня репозитория каждого микро командой:

```bash
go test ./...
```

### Ручное тестирование

Используйте приведённые выше примеры `curl`-запросов или Postman для проверки работы API в различных сценариях (валидные данные, ошибки валидации, отсутствие задач, проверка параллельных вычислений и т.д.).

---

## ⚙️ Переменные окружения

Для настройки временных задержек операций и параллельных вычислений используйте следующие переменные (указываются в `orchestrator/config/config.yaml` и `agent/config/config.yaml` или напрямую через окружение):

- `TIME_ADDITION_MS` — время выполнения сложения (мс)
- `TIME_SUBTRACTION_MS` — время выполнения вычитания (мс)
- `TIME_MULTIPLICATIONS_MS` — время выполнения умножения (мс)
- `TIME_DIVISIONS_MS` — время выполнения деления (мс)
- `COMPUTING_POWER` — количество параллельных вычислителей (горутин)

Значения по умолчанию прописаны в соответствующих `config.yaml`. При необходимости их можно переопределить, например:

```bash
export TIME_ADDITION_MS=500
export COMPUTING_POWER=5
```

Затем запустить Оркестратор/Агент.

---

## Диаграмма проекта

<p align="center">
  <img src="/assets/Diagram.png" alt="Diagram" width="400"/>
</p>

Диаграмма наглядно показывает, как Оркестратор принимает выражения, разбивает их на задачи, а Агент обрабатывает каждую задачу отдельно. После завершения задачи результат возвращается Оркестратору, который обновляет статус выражения.

---

## 📂 Структура проекта

```
.
├── agent
│   ├── api
│   │   └── task
│   │       └── task.proto
│   ├── cmd
│   │   └── agent
│   │       └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── deamon
│   │   │   ├── deamon.go
│   │   │   └── deamon_test.go
│   │   └── grpc
│   │       └── conn.go
│   ├── Makefile
│   └── pkg
│       ├── gen
│       │   └── api
│       │       └── task
│       │           ├── task_grpc.pb.go
│       │           └── task.pb.go
│       └── logger
│           └── logger.go
├── assets
│   ├── Diagram.png
│   └── GOPHER_ACADEMY.png
├── auth
│   ├── api
│   │   └── auth.proto
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── app
│   │   │   ├── app.go
│   │   │   ├── dto
│   │   │   │   ├── auth.go
│   │   │   │   └── register.go
│   │   │   ├── login.go
│   │   │   └── register.go
│   │   ├── domain
│   │   │   ├── entity
│   │   │   │   └── user.go
│   │   │   └── interfaces
│   │   │       └── user.go
│   │   ├── infrastructure
│   │   │   ├── auth
│   │   │   │   ├── interface.go
│   │   │   │   └── jwt.go
│   │   │   ├── hashPass
│   │   │   │   ├── passHash.go
│   │   │   │   └── passHash_test.go
│   │   │   └── persistence
│   │   │       ├── postgres
│   │   │       │   └── user.go
│   │   │       └── sqlite
│   │   │           └── user.go
│   │   └── presentation
│   │       └── grpc
│   │           ├── handlers
│   │           │   └── auth.go
│   │           └── server.go
│   ├── Makefile
│   ├── migrations
│   │   ├── 000001_init.down.sql
│   │   └── 000001_init.up.sql
│   ├── pkg
│   │   ├── db
│   │   │   ├── postgres
│   │   │   │   └── conn.go
│   │   │   └── sqlite_conn
│   │   │       ├── conn.go
│   │   │       └── models
│   │   │           └── user.go
│   │   ├── gen
│   │   │   └── api
│   │   │       ├── auth_grpc.pb.go
│   │   │       └── auth.pb.go
│   │   └── logger
│   │       └── logger.go
│   └── tests
├── CalcYA API.postman_collection.json
├── orchestrator
│   ├── api
│   │   ├── auth
│   │   │   └── auth.proto
│   │   └── task
│   │       └── task.proto
│   ├── cmd
│   │   └── orchestrator
│   │       └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── config.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── app
│   │   │   └── expr
│   │   │       ├── createExp.go
│   │   │       ├── expr_test.go
│   │   │       ├── genTask.go
│   │   │       ├── parse.go
│   │   │       └── tokenize.go
│   │   ├── domain
│   │   │   ├── entity
│   │   │   │   ├── expression.go
│   │   │   │   └── task.go
│   │   │   └── repository
│   │   │       ├── expression.go
│   │   │       └── task.go
│   │   ├── infrastructure
│   │   │   └── persistence
│   │   │       ├── exprStore.go
│   │   │       └── taskStore.go
│   │   └── presentation
│   │       ├── grpc
│   │       │   ├── handlers
│   │       │   │   └── handlers.go
│   │       │   └── server.go
│   │       └── http
│   │           ├── dto
│   │           │   └── dto.go
│   │           ├── handlers
│   │           │   ├── auth.go
│   │           │   └── expr.go
│   │           ├── handlers.go
│   │           ├── middleware
│   │           │   └── jwt.go
│   │           ├── routes.go
│   │           └── server.go
│   ├── Makefile
│   ├── migrations
│   │   └── models
│   │       ├── expression.go
│   │       └── task.go
│   └── pkg
│       ├── db
│       │   └── SQLite
│       │       └── conn.go
│       ├── gen
│       │   └── api
│       │       ├── auth
│       │       │   ├── auth_grpc.pb.go
│       │       │   └── auth.pb.go
│       │       └── task
│       │           ├── task_grpc.pb.go
│       │           └── task.pb.go
│       └── logger
│           └── logger.go
└── README.md

77 directories, 83 files

```

---

## 🤝 Контакты

Автор: **Даня (sokrat_00)**

- Telegram: [@SOKRAT_00](https://t.me/SOKRAT_00)
- GitHub: [github.com/0sokrat0](https://github.com/0sokrat0)

---
