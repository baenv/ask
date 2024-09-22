-- Seed data for users table
INSERT INTO users (id, user_id, platform, username, created_at, last_active)
VALUES
    (1, '123456789', 'discord', 'JohnDoe', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, '987654321', 'discord', 'JaneSmith', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 'alice_wonder', 'telegram', 'Alice', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 'bob_builder', 'telegram', 'Bob', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Seed data for user_agent_configs table
INSERT INTO user_agent_configs (id, user_id, api_key, endpoint_url, created_at, updated_at)
VALUES
    (1, 1, 'api_key_john', 'https://api.example.com/john', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 3, 'api_key_alice', 'https://api.example.com/alice', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Seed data for servers table
INSERT INTO servers (id, server_id, platform, server_name, created_at)
VALUES
    (1, '111222333', 'discord', 'Gaming Club', CURRENT_TIMESTAMP),
    (2, '444555666', 'discord', 'Book Club', CURRENT_TIMESTAMP),
    (3, '1001234567890', 'telegram', 'Tech Talk', CURRENT_TIMESTAMP);

-- Seed data for server_admin_configs table
INSERT INTO server_admin_configs (id, server_id, api_key, endpoint_url, created_at, updated_at)
VALUES
    (1, 1, 'api_key_gaming', 'https://api.example.com/gaming', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 3, 'api_key_tech', 'https://api.example.com/tech', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
