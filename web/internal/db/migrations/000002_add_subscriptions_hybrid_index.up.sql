CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE INDEX subscriptions_hybrid_index ON subscriptions USING GIST (
    user_id,
    service_name,
    daterange(start_date, end_date, '[]')
);
