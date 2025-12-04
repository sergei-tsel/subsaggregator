CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,                 -- название сервиса, предоставляющего подписку
    price INTEGER NOT NULL,                     -- стоимость месячной подписки в рублях
    user_id UUID NOT NULL,                 	    -- уникальный идентификатор пользователя
    start_date DATE DEFAULT now(),
    end_date DATE DEFAULT null
)
