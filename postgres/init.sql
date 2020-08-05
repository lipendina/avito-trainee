CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, username TEXT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, UNIQUE(username));
CREATE TABLE IF NOT EXISTS chats (id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, name TEXT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE TABLE IF NOT EXISTS messages (id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, chat UUID REFERENCES chats(id), author UUID REFERENCES users(id), "text" TEXT NOT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL);
CREATE TABLE IF NOT EXISTS chats_users (id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, chat_id UUID REFERENCES chats(id), user_id UUID REFERENCES users(id));
CREATE INDEX IF NOT EXISTS chats_users_chat_id_idx ON chats_users (chat_id);
CREATE INDEX IF NOT EXISTS chats_users_user_id_idx ON chats_users (user_id);
CREATE INDEX IF NOT EXISTS messages_author_idx ON messages (author);
CREATE INDEX IF NOT EXISTS messages_chat_idx ON messages (chat);