CREATE TABLE IF NOT EXISTS users (
                                     id uuid PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL,
                                     email VARCHAR(100) NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     user_icon_url TEXT,
                                     role VARCHAR(20) DEFAULT 'user',
                                     created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workspaces (
                                          id SERIAL PRIMARY KEY,
                                          name VARCHAR(20) NOT NULL,
                                          created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_workspaces (
                                               user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                                               workspace_id INTEGER REFERENCES workspaces(id) on delete CASCADE,
                                               role VARCHAR(20) DEFAULT 'member',
                                               joined_at TIMESTAMP DEFAULT NOW(),
                                               PRIMARY KEY (user_id, workspace_id)

);