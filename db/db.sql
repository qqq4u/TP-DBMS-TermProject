CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE "user"
(
    Nickname CITEXT COLLATE "C" PRIMARY KEY,
    FullName TEXT NOT NULL,
    About    TEXT NOT NULL DEFAULT '',
    Email    CITEXT COLLATE "C" UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Title   TEXT NOT NULL,
    "user"  CITEXT COLLATE "C",
    Slug    CITEXT COLLATE "C" PRIMARY KEY,
    Posts   INT DEFAULT 0,
    Threads INT DEFAULT 0
);

CREATE UNLOGGED TABLE thread
(
    Id      SERIAL PRIMARY KEY,
    Title   TEXT NOT NULL,
    Author  CITEXT COLLATE "C" REFERENCES "user" (Nickname),
    Forum   CITEXT COLLATE "C" REFERENCES "forum" (Slug),
    Message TEXT NOT NULL,
    Votes   INT                      DEFAULT 0,
    Slug    CITEXT COLLATE "C",
    Created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNLOGGED TABLE post
(
    Id       SERIAL PRIMARY KEY,
    Author   CITEXT COLLATE "C",
    Created  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    Forum    CITEXT COLLATE "C",
    IsEdited BOOLEAN                  DEFAULT FALSE,
    Message  CITEXT COLLATE "C" NOT NULL,
    Parent   INT                      DEFAULT 0,
    Thread   INT,
    Path     INTEGER[],
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    FOREIGN KEY (author) REFERENCES "user" (nickname)
);

CREATE UNLOGGED TABLE vote
(
    ID     SERIAL PRIMARY KEY,
    Author CITEXT COLLATE "C" REFERENCES "user" (Nickname),
    Voice  INT NOT NULL,
    Thread INT,
    FOREIGN KEY (thread) REFERENCES "thread" (id),
    UNIQUE (Author, Thread)
);


CREATE UNLOGGED TABLE user_forum
(
    Nickname CITEXT COLLATE "C" NOT NULL,
    FullName TEXT               NOT NULL,
    About    TEXT,
    Email    CITEXT COLLATE "C",
    Slug     CITEXT COLLATE "C" NOT NULL,
    FOREIGN KEY (Nickname) REFERENCES "user" (Nickname),
    FOREIGN KEY (Slug) REFERENCES "forum" (Slug),
    UNIQUE (Nickname, Slug)
);
