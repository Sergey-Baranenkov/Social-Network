package postgres

var DropDB = `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`

var UserTable = `
create table if not exists users (
    /*secret_info*/
    user_id bigserial primary key,
	email text not null unique,
	token bytea not null,
	first_name text not null check ( length(first_name) > 0 ),
	last_name text not null check ( length(last_name) > 0),
    sex char(1) not null check ( sex in ('М', 'Ж')),
    /*basic_info*/
    avatar_ref text default 'hash_path/def_avatar.jpg', /*исправить путь*/
    bg_ref text default 'hash_path/def_bg',
	tel decimal(20),
	city text default '',
	birthday date, /*обработать*/
	status int2 check (0 <= status and status <= 5) default 0, /*0-не указано 1 1-женат 2-не женат 3-влюблен 4-все сложно 5-в акт.поиске */

    /*hobbies*/
	hobby text default '',
	fav_music text default '',
	fav_films text default '',
	fav_books text default '',
	fav_games text default '',
	other_interests text default '',

	/*privacy*/
	who_can_message text not null default 'all' check(who_can_message in ('all','fo')), /*fo - friends only*/
	who_can_see_info text not null default 'all' check(who_can_see_info in ('all','fo')),

    /*edu_and_emp*/
    edu_and_emp_info jsonb
); create index if not exists users_user_id_idx on users (user_id);
create index if not exists users_full_name_id_idx on users (first_name, last_name);
`

var ObjectsTable = `
create table if not exists objects(
    path ltree primary key,
    auth_id bigint references users (user_id) on delete cascade,
    text text not null,
    num_likes integer default 0 not null,
    creation_time      timestamptz default current_timestamp not null ,
    modification_time   timestamptz
); create index if not exists object_path_idx on objects using gist(path);
create index if not exists object_auth_id_idx on objects (auth_id);
`

var PostInfo = `
create table if not exists post_info(
    path ltree primary key references objects (path) on delete cascade ,
    num_comments integer default 0 not null,
    ref_id bigint
); create index if not exists post_obj_id_idx on post_info using gist(path);
create sequence if not exists num_posts;
`

var LikesTable = `
create table if not exists likes(
    path ltree references objects (path) on delete cascade ,
    auth_id bigint references users (user_id) on delete cascade,
    primary key (path, auth_id)
); create index if not exists likes_path_idx on likes using gist(path);
`


var Triggers = `CREATE OR REPLACE FUNCTION  insert_object_before_process() RETURNS trigger AS $insert_object_before_process$
    BEGIN
        if new.path = '' then
            new.path = text2ltree(nextval('num_posts')::text);
            else
                new.path = new.path || (select coalesce(max(right(path::text, 1)::bigint), 0) + 1 from objects o where o.path ~ (new.path::text || '.*{1}')::lquery)::text;
        end if;
       	return new;
    END;
$insert_object_before_process$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS insert_object_before_trigger ON public.objects;
CREATE TRIGGER insert_object_before_trigger
	BEFORE INSERT ON objects
FOR EACH ROW EXECUTE PROCEDURE insert_object_before_process ();


CREATE OR REPLACE FUNCTION insert_object_after_process() RETURNS trigger AS $insert_object_after_process$
    BEGIN
        if nlevel(new.path) = 1 then /*post*/
            insert into post_info (path) values (new.path);
        else /*comment*/
            update post_info
                set num_comments = num_comments + 1
            where post_info.path = subpath(new.path,0, 1);
        end if;

       	Return new;
    END;
$insert_object_after_process$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS insert_object_after_trigger ON public.objects;
CREATE TRIGGER insert_object_after_trigger
	AFTER INSERT ON objects
FOR EACH ROW EXECUTE PROCEDURE insert_object_after_process ();


CREATE OR REPLACE FUNCTION add_like_process() RETURNS trigger AS $add_like_process$
    begin
       	update objects
            set num_likes = num_likes + 1
        where objects.path = new.path;
       	return new;
    end;
$add_like_process$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS add_like_to_post_trigger ON public.likes;
CREATE TRIGGER add_like_to_post_trigger
	AFTER INSERT ON likes
FOR EACH ROW EXECUTE PROCEDURE add_like_process ();

CREATE OR REPLACE FUNCTION revoke_like_process() RETURNS trigger AS $revoke_like_process$
    begin
       	update objects
            set num_likes = num_likes - 1
        where objects.path = old.path;
       	return new;
    end;
$revoke_like_process$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS revoke_like_to_post_trigger ON public.likes;
CREATE TRIGGER revoke_like_to_post_trigger
	AFTER DELETE ON likes
FOR EACH ROW EXECUTE PROCEDURE revoke_like_process ();

CREATE OR REPLACE FUNCTION delete_object_process() RETURNS trigger AS $delete_object_process$
    begin
       	delete from objects where path <@ old.path;
       	return new;
    end;
$delete_object_process$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS delete_object_trigger ON public.objects;
CREATE TRIGGER delete_object_trigger
	AFTER DELETE ON objects
FOR EACH ROW EXECUTE PROCEDURE delete_object_process ();

CREATE OR REPLACE FUNCTION update_post_text_process() RETURNS trigger AS $update_post_text_process$
    BEGIN
       IF new.text <> old.text
        then
            new.modification_time = now();
       end if;
       return new;
    END;
$update_post_text_process$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_post_text_trigger ON public.objects;
CREATE TRIGGER update_post_text_trigger
	BEFORE UPDATE ON objects
FOR EACH ROW EXECUTE PROCEDURE update_post_text_process ();`

var SelectFunctions = `CREATE OR REPLACE FUNCTION get_comments(post_path text) RETURNS json AS $$
        BEGIN
            return (
select array_to_json(array_agg(res))
    from (
            WITH RECURSIVE
                c AS (
                    SELECT text,
                           num_likes,
                           path,
                           exists (select 1 from likes where likes.path = objects.path) as me_liked,
                           first_name,
                           last_name,
                           to_char(creation_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as creation_time,
                           to_char(modification_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as modification_time,
                           nlevel(path) AS lvl
                    FROM objects join users on user_id = auth_id where path <@ text2ltree(post_path)
                ),
                maxlvl AS (
                    SELECT max(lvl) maxlvl
                    FROM c
                ),
                j AS (
                    SELECT c.*,
                           json '[]' AS children
                    FROM c,
                         maxlvl
                    WHERE lvl = maxlvl
                    UNION ALL
                    SELECT (c).*,
                           CASE
                               WHEN COUNT(j) > 0
                                   THEN json_agg(j)
                               ELSE json '[]'
                               END AS children
                    FROM (
                             SELECT c,
                                    CASE
                                        WHEN c.path = subpath(j.path, 0, nlevel(j.path) - 1)
                                            THEN j
                                        END AS j
                             FROM j
                                      JOIN c ON c.lvl = j.lvl - 1
                         ) AS v
                    GROUP BY v.c
                )
            SELECT *
            FROM j
            WHERE lvl = 2
        ) res
                );
        end;
    $$ PARALLEL SAFE LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_posts(u_id bigint, my_id bigint) returns json as $$
        begin
            return
                (
select array_to_json(array_agg(row_to_json(t)))
    from (
      select o.text,
             o.path,
             o.num_likes,
             exists (select 1 from likes where o.path = likes.path and likes.auth_id = my_id) as me_liked,
             to_char(o.creation_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as creation_time,
             to_char(o.modification_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as modification_time,
             first_name,
             last_name,
             num_comments

      from objects o join users u on u.user_id = o.auth_id join post_info pi on o.path = pi.path where o.auth_id = u_id
    ) t
        );
        end;
    $$ PARALLEL SAFE LANGUAGE plpgsql;`

var InitTestSQL = `
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs@mail.ru', E'\a','Vladimir','Putin', 'М');
insert into users (email, token, first_name, last_name, sex) values ('lol@mail.ru', E'\a','Sergey','Baranenkov', 'М');

insert into objects (auth_id, path, text) values (1, '','lol');
insert into objects (auth_id, path, text) values (1, '','lol');
insert into objects (auth_id, path, text) values (1, '', 'lol');
insert into objects (auth_id, path, text) values (1, '1','lol');

insert into likes (path, auth_id) values ('1', 1);
insert into likes (path, auth_id) values ('1.1', 1);
insert into likes (path, auth_id) values ('2', 1);
`