-- CREATE TYPE chat_types AS ENUM ('PRIVATE', 'GROUP');
-- CREATE TYPE message_types AS ENUM ('TEXT', 'IMAGE', 'FILE', 'SYSTEM');
-- CREATE TYPE user_roles AS ENUM ('USER', 'ADMIN');
-- CREATE TYPE workspace_user_roles AS ENUM ('CREATOR', 'ADMIN', 'MEMBER');
-- CREATE TYPE chat_user_roles AS ENUM ('USER', 'ADMIN');
-- CREATE TYPE project_user_roles AS ENUM ('CREATOR', 'EDITOR', 'VIEWER', 'ADMIN);


CREATE TABLE IF NOT EXISTS users (
                                     id uuid PRIMARY KEY,
                                     username VARCHAR(50) NOT NULL,
                                     email VARCHAR(100) NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     user_icon_url TEXT,
                                     role user_roles,
                                     created_at timestamptz DEFAULT NOW(),
                                     last_online_at timestamptz,
                                     is_online BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS workspaces (
                                          id uuid PRIMARY KEY,
                                          creator_id uuid REFERENCES users(id) NOT NULL,
                                          name VARCHAR(20) NOT NULL,
                                          created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_workspaces (
                                               user_id uuid REFERENCES users(id),
                                               workspace_id uuid REFERENCES workspaces(id),
                                               role workspace_user_roles NOT NULL,
                                               position VARCHAR(50),
                                               joined_at timestamptz DEFAULT NOW(),
                                               PRIMARY KEY (user_id, workspace_id)
);

CREATE TABLE IF NOT EXISTS task_columns (
                                            id UUID PRIMARY KEY,
                                            workspace_id UUID REFERENCES workspaces(id),
                                            project_id UUID REFERENCES projects(id),
                                            name TEXT NOT NULL,
                                            position INTEGER DEFAULT 0,
                                            color VARCHAR(7)
);

CREATE TABLE IF NOT EXISTS tasks (
                                     id UUID PRIMARY KEY,
                                     column_id UUID REFERENCES task_columns(id),
                                     workspace_id UUID REFERENCES workspaces(id),
                                     project_id UUID REFERENCES projects(id),
                                     sprint_id UUID REFERENCES sprints(id),
                                     author_id UUID REFERENCES users(id),
                                     created_at TIMESTAMPTZ DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ DEFAULT NOW(),
                                     title VARCHAR(30) NOT NULL,
                                     description TEXT,
                                     is_done BOOLEAN DEFAULT FALSE NOT NULL,
                                     due_date TIMESTAMPTZ,
                                     priority INTEGER DEFAULT 0,
                                     position INT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks_assignment (
                                                task_id UUID REFERENCES tasks(id),
                                                user_id UUID REFERENCES users(id),
                                                PRIMARY KEY (task_id, user_id)
);

CREATE TABLE IF NOT EXISTS chats (
                                     id TEXT PRIMARY KEY,
                                     workspace_id UUID REFERENCES workspaces(id),
                                     creator_id UUID REFERENCES users(id),
                                     type chat_types NOT NULL,
                                     created_at timestamptz DEFAULT now() NOT NULL,
                                     updated_at timestamptz DEFAULT now() NOT NULL,
                                     last_message_time timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS group_chats (
                                           chat_id TEXT PRIMARY KEY REFERENCES chats(id),
                                           name VARCHAR(50) NOT NULL,
                                           icon_url TEXT
);

CREATE TABLE IF NOT EXISTS user_chats (
                                          chat_id TEXT REFERENCES chats(id),
                                          user_id UUID REFERENCES users(id),
                                          workspace_id UUID REFERENCES workspaces(id),
                                          muted BOOLEAN DEFAULT FALSE NOT NULL,
                                          pinned BOOLEAN DEFAULT FALSE NOT NULL,
                                          notification BOOLEAN DEFAULT TRUE NOT NULL,
                                          role chat_user_roles NOT NULL,
                                          joined_at timestamptz NOT NULL DEFAULT NOW(),
                                          PRIMARY KEY (chat_id, user_id, workspace_id)
);

CREATE TABLE IF NOT EXISTS pinned_chat_position (
                                                    chat_id TEXT REFERENCES chats(id),
                                                    user_id UUID REFERENCES users(id),
                                                    position INTEGER NOT NULL,
                                                    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
                                        id TEXT PRIMARY KEY,
                                        chat_id TEXT REFERENCES chats(id),
                                        sender_id UUID REFERENCES users(id),
                                        content TEXT,
                                        message_type message_types NOT NULL,
                                        created_at timestamptz DEFAULT now() NOT NULL,
                                        updated_at timestamptz NOT NULL,
                                        is_edited BOOLEAN DEFAULT FALSE NOT NULL,
                                        is_deleted BOOLEAN DEFAULT FALSE NOT NULL,
                                        reply_to TEXT REFERENCES messages(id),
                                        is_read bool DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS message_attachments (
                                                   id UUID PRIMARY KEY,
                                                   message_id TEXT REFERENCES messages(id),
                                                   file_url TEXT NOT NULL,
                                                   file_type VARCHAR(20) NOT NULL,
                                                   file_name TEXT NOT NULL,
                                                   file_size INTEGER NOT NULL,
                                                   uploaded_at timestamptz DEFAULT now() NOT NULL,
                                                   chat_id text REFERENCES chats(id)
);

CREATE TABLE IF NOT EXISTS message_reads (
                                             user_id UUID REFERENCES users(id),
                                             message_id TEXT REFERENCES messages(id),
                                             chat_id TEXT REFERENCES chats(id),
                                             read_at timestamptz,
                                             PRIMARY KEY (user_id, message_id)
);

CREATE TABLE IF NOT EXISTS projects (
                                        id UUID PRIMARY KEY,
                                        workspace_id UUID REFERENCES workspaces(id) NOT NULL,
                                        creator_id UUID REFERENCES users(id) NOT NULL,
                                        name VARCHAR(100) NOT NULL,
                                        description TEXT,
                                        is_done BOOLEAN DEFAULT FALSE NOT NULL,
                                        created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                        updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                        prefix VARCHAR(10) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sprints (
                                       id UUID PRIMARY KEY,
                                       project_id UUID REFERENCES projects(id) NOT NULL,
                                       name VARCHAR(100) NOT NULL,
                                       start_date DATE NOT NULL,
                                       end_date DATE NOT NULL,
                                       is_done BOOLEAN DEFAULT FALSE NOT NULL,
                                       created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                       updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS project_members (
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               project_id UUID REFERENCES projects(id) ON DELETE CASCADE NOT NULL,
                                               user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                                               role project_user_roles NOT NULL DEFAULT 'VIEWER',
                                               joined_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
                                               UNIQUE (project_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_sprint_id ON tasks(sprint_id);

-- TODO ALTER TABLE tasks ALTER COLUMN position SET NOT NULL;