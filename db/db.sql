 CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE UNLOGGED TABLE "user"
(
    Nickname CITEXT COLLATE "C" PRIMARY KEY,
    FullName TEXT NOT NULL,
    About    TEXT NOT NULL DEFAULT '',
    Email    CITEXT COLLATE "C"
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
    parent_thread int;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT thread FROM post WHERE id = new.parent INTO parent_thread;
        IF NOT FOUND OR parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'NOT FOUND OR parent_thread != NEW.thread' USING ERRCODE = '22409';
        end if;

        SELECT path FROM post WHERE id = new.parent INTO parent_path;
        NEW.path := parent_path || new.id;
    END IF;
    UPDATE forum SET Posts=Posts + 1 WHERE forum.slug = new.forum;
    RETURN new;
END
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER u_p_trigger
    BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE updatePath();

CREATE INDEX IF NOT EXISTS thread_slug_idx ON "thread" (slug);
CREATE INDEX IF NOT EXISTS thread_id_idx ON "thread" (id);

CREATE INDEX IF NOT EXISTS post_thread_idx ON "post" (thread);
CREATE INDEX IF NOT EXISTS post_thread_id_idx ON "post" (thread, id);
CREATE INDEX IF NOT EXISTS post_id_path1_idx ON "post" (id, (path[1]));
CREATE INDEX IF NOT EXISTS post_path1_idx ON "post" ((path[1]));
CREATE INDEX IF NOT EXISTS post_thread_path_idx ON "post" (thread, path, id);
CREATE INDEX IF NOT EXISTS post_thread_path1_parent_idx ON "post" (thread, id, (path[1]), parent);

CREATE INDEX IF NOT EXISTS user_nickname_hash_idx ON "user" USING hash (nickname);
CREATE INDEX IF NOT EXISTS user_email_hash_idx ON "user" USING hash (email);
CREATE INDEX IF NOT EXISTS user_full_idx ON "user" (email, nickname);

CREATE INDEX IF NOT EXISTS forum_slug_hash_idx ON "forum" USING hash (slug);

CREATE INDEX IF NOT EXISTS thread_date_idx ON thread (created);
CREATE INDEX IF NOT EXISTS thread_forum_date_idx ON thread (forum, created);
CREATE INDEX IF NOT EXISTS thread_forum_hash_idx ON thread USING hash (forum);
CREATE INDEX IF NOT EXISTS thr_slug_hash_idx ON thread USING hash (slug);

CREATE INDEX IF NOT EXISTS vote_full_idx ON vote (author, voice, thread);

CREATE INDEX IF NOT EXISTS fu_forum ON user_forum USING hash(Slug);

CREATE UNIQUE INDEX IF NOT EXISTS forum_users_unique ON user_forum (slug, nickname);
CREATE UNIQUE INDEX IF NOT EXISTS vote_unique ON vote (author, thread);

VACUUM;
VACUUM ANALYSE;
