CREATE TABLE IF NOT EXISTS feeds (
    "id" uuid NOT NULL PRIMARY KEY,
    "url" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS articles (
    "id" uuid NOT NULL PRIMARY KEY,
    "feed_id" uuid NOT NULL,
    "title" text NOT NULL,
    "source" text NOT NULL,
    "link" text NOT NULL,
    "tags" text[] NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chats (
    "id" uuid NOT NULL PRIMARY KEY,
    "telegram_chat_id" text NOT NULL,
    "telegram_user_id" text,
    "tags" text[] NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chat_feed (
    "id" uuid NOT NULL PRIMARY KEY,
    "chat_id" uuid NOT NULL,
    "feed_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_article (
    "id" uuid NOT NULL PRIMARY KEY,
    "chat_id" uuid NOT NULL,
    "article_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE
);