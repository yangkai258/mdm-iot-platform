INSERT INTO sys_users (username, password, nickname, tenant_id, status)
VALUES ('admin', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/X4aMqNw6FKP6PZQ7i', 'Administrator', 'e6cbcb82-9bd6-4803-8bf7-b4b1af8eaec2', 1)
RETURNING id;
