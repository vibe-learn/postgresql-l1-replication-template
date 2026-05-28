        # postgresql — Streaming replication, WAL и failover

        Homework-шаблон для урока **l1_replication_basics** (Streaming replication, WAL и failover) на платформе Vibe Learn.

        ## Что делать

        Дано: docker-compose с primary + 2 standby + Patroni-like wrapper. Реализуй на Go:
1) Health-check primary/standby; парсинг pg_stat_replication.
2) Программный failover: detect primary down (через сетевую заглушку),
   выбор лучшего standby по replay_lsn, promote.
3) Тест: пишет на primary, отключаем primary, проверяет, что новый primary
   не теряет committed transactions (если sync) или теряет минимально (если async).
Тесты в template проверят корректность detection, выбора кандидата, и сохранение
инвариантов после failover.

## Контекст (из transfer-задачи урока)

Ты архитектор. Тебе нужно спроектировать HA для PG в трёх сценариях:

(S1) **Финтех-платформа.** Финансовые транзакции. SLA: ZERO data loss. Один регион,
     2 AZ, RTT между AZ ~2мс. p99 COMMIT — допустим 30мс.
(S2) **Аналитическая платформа.** OLAP-запросы на свежих данных. SLA: 30 минут lag допустим,
     но writes должны быть быстрыми. Read-нагрузка в 10× больше write-нагрузки.
(S3) **Кросс-региональный SaaS.** Primary в Европе, read-replicas в US/Asia. SLA на
     запись: 50мс, на чтение: 100мс (regional). Допустима secondary потеря ~5 секунд
     при failover.

## Recap из урока

- **Streaming replication = physical replication через WAL-stream.** Standby получает WAL и реплеит его в свою копию данных.
- **Async (дефолт) vs sync** — фундаментальный трейд-оф: latency vs zero data loss. Финтех → sync, остальное обычно → async.
- **Hot standby отвечает на SELECT, но с replication lag.** Read-your-writes требует чтения с primary или sync replication с remote_apply.
- **Split-brain — главная опасность HA.** Лечится consensus-based failover (Patroni + etcd/Consul/ZK), а не таймерами и пингами.
- **Failover в проде делает Patroni или managed**, не руки DBA. Готовая логика: detection, candidate selection по LSN, promote, fencing старого primary.

        ## Как работать

        1. Платформа Vibe Learn создаёт копию этого репо в твоём GitHub-аккаунте по клику «Начать домашку» на странице урока (через GitHub `/generate`, codecrafters-pattern).
        2. Склонируй копию локально, реализуй TODO в `main.go`, прогони тесты, запушь.
        3. CI (`.github/workflows/ci.yml`) запускает `go vet` + `go test ./...` на каждый push. Платформа слушает результат через webhook от GitHub Actions и обновляет статус домашки на странице урока.

        ## Локальное окружение

        - Go 1.22+
        - Docker + docker-compose — `docker compose up -d` поднимает single-node PostgreSQL 16 на `localhost:5432` с healthcheck. DSN: `postgres://postgres:postgres@localhost:5432/postgres`. Переопределяется через env `DATABASE_URL`.

        ## Запуск

        ```bash
        # Поднять локальный PostgreSQL
        docker compose up -d

        # Прогнать тесты (интеграционный включается через PG_INTEGRATION=1)
        go test ./...
        PG_INTEGRATION=1 go test ./...

        # Запустить main (печатает marker; замени stub на реализацию)
        go run .
        ```

        ## Заметка автора

        Это baseline-шаблон, сгенерированный платформой. Бизнес-сущность задачи (что конкретно реализовать в `main.go`, какие тесты сделать строгими) расширяется по ходу итераций — параллельно с углублением теории урока.
