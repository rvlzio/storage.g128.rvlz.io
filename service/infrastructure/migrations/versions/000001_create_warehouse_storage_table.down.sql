BEGIN;

REVOKE ALL PRIVILEGES ON warehouse_storage FROM storage_service_api;

DROP TABLE IF EXISTS warehouse_storage;

END;
