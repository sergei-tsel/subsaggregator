CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,              -- название сервиса, предоставляющего подписку
    price INTEGER NOT NULL,                  -- стоимость месячной подписки в рублях
    user_id TEXT NOT NULL,                 	 -- уникальный идентификатор пользователя
    start_date DATE NOT NULL,                -- дата начала подписки
    end_date DATE DEFAULT NULL               -- дата окончания подписки
)
