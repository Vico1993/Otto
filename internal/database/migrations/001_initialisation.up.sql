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
    "published" text NOT NULL,
    "source" text NOT NULL,
    "link" text NOT NULL,
    "tags" text[] NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
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
    FOREIGN KEY (feed_id) REFERENCES feeds(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id)
);

CREATE TABLE IF NOT EXISTS chat_article (
    "id" uuid NOT NULL PRIMARY KEY,
    "chat_id" uuid NOT NULL,
    "article_id" uuid NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY (article_id) REFERENCES articles(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id)
);