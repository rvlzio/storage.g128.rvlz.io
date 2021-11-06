BEGIN;

CREATE TABLE IF NOT EXISTS warehouse_storage (
    id VARCHAR UNIQUE NOT NULL,
    warehouse_id VARCHAR UNIQUE NOT NULL,
    capacity BIGINT NOT NULL,
    claimed_storage BIGINT NOT NULL,
    file_reservations JSONB NOT NULL
);

GRANT SELECT, INSERT, UPDATE ON warehouse_storage TO storage_service_api;

END;
