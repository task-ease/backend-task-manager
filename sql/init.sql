CREATE TABLE IF NOT EXISTS users (
                                     id uuid PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL,
                                     email VARCHAR(100) NOT NULL UNIQUE,
                                     password_hash text NOT NULL,
                                     user_icon_url TEXT,
                                     role VARCHAR(20) DEFAULT 'user',
                                     created_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workspaces (
                                          id uuid PRIMARY KEY,
                                          creator_id uuid REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                                          name VARCHAR(20) NOT NULL,
                                          created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_workspaces (
                                               user_id uuid REFERENCES users(id) ON DELETE CASCADE,
                                               workspace_id uuid REFERENCES workspaces(id) ON DELETE CASCADE,
                                               role VARCHAR(20) DEFAULT 'member',
                                               joined_at timestamptz DEFAULT NOW(),
                                               PRIMARY KEY (user_id, workspace_id)
);

CREATE TABLE IF NOT EXISTS tasks (
                                     id uuid PRIMARY KEY,
                                     author_id uuid REFERENCES users(id) ON DELETE CASCADE,
                                     created_at TIMESTAMPTZ DEFAULT NOW(),
                                     title VARCHAR(20) NOT NULL,
                                     description TEXT,
                                     is_finished BOOLEAN DEFAULT FALSE,
                                     due_date TIMESTAMPTZ,
                                     priority INTEGER,
                                     updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks_assignment (
                                                task_id uuid REFERENCES tasks(id) ON DELETE CASCADE,
                                                user_id uuid REFERENCES users(id) ON DELETE CASCADE,
                                                PRIMARY KEY (task_id, user_id)
);