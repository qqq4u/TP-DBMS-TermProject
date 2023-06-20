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
    Nickname CITEXT NOT NULL,
    FullName TEXT   NOT NULL,
    About    TEXT,
    Email    CITEXT,
    Slug     CITEXT NOT NULL,
    FOREIGN KEY (Nickname) REFERENCES "user" (Nickname),
    FOREIGN KEY (Slug) REFERENCES "forum" (Slug),
    UNIQUE (Nickname, Slug)
);

CREATE OR REPLACE FUNCTION updatePostUserForum() RETURNS TRIGGER AS
$update_forum_posts$
DECLARE
    t_fullname CITEXT;
    t_about    CITEXT;
    t_email    CITEXT;
BEGIN
    SELECT fullname, about, email FROM "user" WHERE nickname = NEW.author INTO t_fullname, t_about, t_email;
    INSERT INTO user_forum (nickname, fullname, about, email, Slug)
    VALUES (New.Author, t_fullname, t_about, t_email, NEW.forum)
    on conflict do nothing;
    return NEW;
end
$update_forum_posts$ LANGUAGE plpgsql;

CREATE TRIGGER p_i_user_forum
    AFTER INSERT
    ON "post"
    FOR EACH ROW
EXECUTE PROCEDURE updatePostUserForum();


CREATE OR REPLACE FUNCTION updateThreadUserForum() RETURNS TRIGGER AS
$update_forum_threads$
DECLARE
    a_nick     CITEXT;
    t_fullname CITEXT;
    t_about    CITEXT;
    t_email    CITEXT;
BEGIN
    SELECT Nickname, fullname, about, email
    FROM "user"
    WHERE Nickname = new.Author
    INTO a_nick, t_fullname, t_about, t_email;
    INSERT INTO "user_forum" (nickname, fullname, about, email, slug)
    VALUES (a_nick, t_fullname, t_about, t_email, NEW.forum)
    on conflict do nothing;
    return NEW;
end
$update_forum_threads$ LANGUAGE plpgsql;

CREATE TRIGGER t_i_forum_users
    AFTER INSERT
    ON "thread"
    FOR EACH ROW
EXECUTE PROCEDURE updateThreadUserForum();

CREATE OR REPLACE FUNCTION ThreadsCountInc() RETURNS TRIGGER AS
$update_forums$
BEGIN
    UPDATE forum SET Threads=(Threads + 1) WHERE slug = NEW.forum;
    return NEW;
end
$update_forums$ LANGUAGE plpgsql;

CREATE TRIGGER a_t_i_forum
    BEFORE INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE ThreadsCountInc();

CREATE OR REPLACE FUNCTION insertVotes() RETURNS TRIGGER AS
$update_vote$
BEGIN
    UPDATE thread SET votes=(votes + NEW.voice) WHERE id = NEW.thread;
    return NEW;
end
$update_vote$ LANGUAGE plpgsql;

CREATE TRIGGER a_voice
    BEFORE INSERT
    ON vote
    FOR EACH ROW
EXECUTE PROCEDURE insertVotes();

CREATE OR REPLACE FUNCTION updateVotes() RETURNS TRIGGER AS
$update_votes$
BEGIN
    IF OLD.Voice <> NEW.Voice THEN
        UPDATE thread SET votes=(votes + NEW.Voice * 2) WHERE id = NEW.Thread;
    END IF;
    return NEW;
end
$update_votes$ LANGUAGE plpgsql;

CREATE TRIGGER e_voice
    BEFORE UPDATE
    ON vote
    FOR EACH ROW
EXECUTE PROCEDURE updateVotes();

CREATE OR REPLACE FUNCTION updatePath() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path   INTEGER[];
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT path FROM "post" WHERE id = NEW.parent INTO parent_path;
        NEW.path := parent_path || NEW.id;
    END IF;
    UPDATE forum SET Posts=(Posts+1) WHERE Slug = NEW.Forum;
    return NEW;
END
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER on_insert_post
    BEFORE INSERT
    ON "post"
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();

CREATE INDEX IF NOT EXISTS users_nickname_index ON "user" USING hash (nickname);
CREATE INDEX IF NOT EXISTS users_email_index ON "user" USING hash (email);

CREATE INDEX IF NOT EXISTS forum_slug_index ON forum USING hash (slug);

CREATE INDEX IF NOT EXISTS thread_slug_index ON thread USING hash (slug);
CREATE INDEX IF NOT EXISTS thread_forum_index ON thread USING hash (forum);
CREATE INDEX IF NOT EXISTS thread_forum_date_index ON thread (forum, created);

CREATE UNIQUE INDEX IF NOT EXISTS forum_users_index ON user_forum (slug, nickname);

CREATE UNIQUE INDEX IF NOT EXISTS vote_index ON vote (Author, Thread);
CREATE INDEX IF NOT EXISTS post_id_index ON post USING hash (id);

CREATE INDEX IF NOT EXISTS post_thread_index ON post (thread);
CREATE INDEX IF NOT EXISTS post_thread_path_id_index ON post (thread, path, id);
CREATE INDEX IF NOT EXISTS post_thread_id_path_parent_index ON post (thread, id, (path[1]), parent);
CREATE INDEX IF NOT EXISTS post_path_index ON post ((path[1]));

VACUUM;
VACUUM ANALYSE;
