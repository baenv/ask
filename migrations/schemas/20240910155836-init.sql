
-- +migrate Up
-- Create an enum for platform types
DO $$ BEGIN
    CREATE TYPE platform_type AS ENUM ('discord', 'telegram');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,  -- Numeric primary key
    user_id VARCHAR(255) NOT NULL,  -- Platform-specific user identifier
    platform platform_type NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_active TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, platform)
);

-- User agent configurations table
CREATE TABLE IF NOT EXISTS user_agent_configs (
    id BIGINT PRIMARY KEY,  -- Numeric primary key
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key VARCHAR(255) NOT NULL,
    endpoint_url VARCHAR(255) NOT NULL,
    command VARCHAR(50) NOT NULL,  -- New field for custom command
    description TEXT,  -- New column for configuration description
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (user_id, command)  -- Ensure unique command per user
);

-- Servers table
CREATE TABLE IF NOT EXISTS servers (
    id BIGINT PRIMARY KEY,  -- Numeric primary key
    server_id VARCHAR(255) NOT NULL,  -- Platform-specific server identifier
    platform platform_type NOT NULL,
    server_name VARCHAR(255) NOT NULL,
    owner_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  -- Add reference to users table
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (server_id, platform)
);

-- Server admin configurations table
CREATE TABLE IF NOT EXISTS server_admin_configs (
    id BIGINT PRIMARY KEY,  -- Numeric primary key
    server_id BIGINT NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    api_key VARCHAR(255) NOT NULL,
    endpoint_url VARCHAR(255) NOT NULL,
    command VARCHAR(50) NOT NULL,  -- New field for custom command
    description TEXT,  -- New column for configuration description
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (server_id, command)  -- Ensure unique command per server
);

-- Create indexes for improved query performance
CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);
CREATE INDEX IF NOT EXISTS idx_servers_server_id ON servers(server_id);
CREATE INDEX IF NOT EXISTS idx_servers_owner_id ON servers(owner_id);  -- Add index for owner_id
CREATE INDEX IF NOT EXISTS idx_user_agent_configs_command ON user_agent_configs(command);
CREATE INDEX IF NOT EXISTS idx_server_admin_configs_command ON server_admin_configs(command);

-- +migrate Down
-- Drop tables
DROP TABLE IF EXISTS server_admin_configs;
DROP TABLE IF EXISTS user_agent_configs;
DROP TABLE IF EXISTS servers;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS platform_type;

-- Drop indexes (optional, as they will be automatically dropped when the tables are dropped)
DROP INDEX IF EXISTS idx_servers_server_id;
DROP INDEX IF EXISTS idx_users_user_id;
DROP INDEX IF EXISTS idx_servers_owner_id;
DROP INDEX IF EXISTS idx_user_agent_configs_command;
DROP INDEX IF EXISTS idx_server_admin_configs_command;
