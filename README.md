# subscription-aggregator-go

REST API сервис для управления подписками пользователей.  

### Cтарт
1. Клонирование<br>
`git clone https://github.com/yourusername/subscription-service.git`<br>
`cd subscription-service`

2. Сборка и запуск<br>
`make up`

### API Endpoints
```
|   POST  | `/api/subscriptions`            |  Создать подписку   |
|   GET   | `/api/subscriptions/:id`        |  Получить подписку  |
|   PUT   | `/api/subscriptions/:id`        |  Обновить подписку  |
|  DELETE | `/api/subscriptions/:id`        |  Удалить подписку   |
|   GET   | `/api/subscriptions`            |  Список (пагинация) |
|   GET   | `/api/subscriptions/total-cost` | Стоимость за период |
```
### Технологии
- Go 1.26 + Gin
- PostgreSQL + sqlx
- Docker + docker-compose
- Swagger
- Zap (логи)

