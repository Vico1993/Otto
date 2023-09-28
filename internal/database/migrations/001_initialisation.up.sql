CREATE TABLE IF NOT EXISTS feeds (
    "id" uuid NOT NULL PRIMARY KEY,
    "url" character(125) NOT NULL,
    "created_at" time without time zone NOT NULL,
    "updated_at" time without time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
    "id" uuid NOT NULL PRIMARY KEY,
    "feed_id" uuid NOT NULL,
    "title" text NOT NULL,
    "published" text NOT NULL,
    "source" character(256) NOT NULL,
    "link" character(256) NOT NULL,
    "tags" character(256)[] NOT NULL,
    "created_at" time without time zone NOT NULL,
    "updated_at" time without time zone DEFAULT now() NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

CREATE TABLE IF NOT EXISTS chats (
    "id" uuid NOT NULL PRIMARY KEY,
    "telegram_chat_id" character(50) NOT NULL,
    "telegram_user_id" character(50),
    "tags" character(256)[] NOT NULL,
    "created_at" time without time zone NOT NULL,
    "updated_at" time without time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_feed (
    "id" uuid NOT NULL PRIMARY KEY,
    "chat_id" uuid NOT NULL,
    "feed_id" uuid NOT NULL,
    "created_at" time without time zone NOT NULL,
    "updated_at" time without time zone DEFAULT now() NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id)
);

CREATE TABLE IF NOT EXISTS chat_article (
    "id" uuid NOT NULL PRIMARY KEY,
    "chat_id" uuid NOT NULL,
    "article_id" uuid NOT NULL,
    "created_at" time without time zone NOT NULL,
    "updated_at" time without time zone DEFAULT now() NOT NULL,
    FOREIGN KEY (article_id) REFERENCES articles(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id)
);