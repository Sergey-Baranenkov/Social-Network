package postgres

var DropDB = `DROP SCHEMA public CASCADE; CREATE SCHEMA public;`

var UserTable = `
create table if not exists users (
    user_id bigserial primary key,
    
    /*secret info*/
	email text not null unique,
	token bytea not null, 
	
	first_name text not null check ( length(first_name) > 0 ),
	last_name text not null check ( length(last_name) > 0),
	full_name tsvector,
    sex char(1) not null check ( sex in ('М', 'Ж')),
    /*basic_info*/
	tel decimal(20) default 0,
	city text default '',
	country text default '',
	birthday text default '',
	status int2 check (0 <= status and status <= 5) default 0, /*0-не указано 1 1-женат 2-не женат 3-влюблен 4-все сложно 5-в акт.поиске */

    /*hobbies*/
	hobby text default '',
	fav_music text default '',
	fav_films text default '',
	fav_books text default '',
	fav_games text default '',
	other_interests text default '',

    /*edu_and_emp*/
    edu_and_emp_info jsonb,

	/*music*/
	music_list bigint[],
	
	/*video*/
	video_list bigint[],
	
	/*images*/
	images_list bigint[]
	
); create index if not exists users_user_id_idx on users (user_id);
create index if not exists users_full_name_id_idx on users using gin(full_name);
`
var MusicTable = `
create table if not exists music (
    music_id bigserial primary key,
    adder_id bigint references users(user_id),
    name text not null default 'undefined',
    author text not null default 'undefined',
    document tsvector
); create index if not exists music_doc_idx on music using gin(document);
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

var PostInfoTable = `
create table if not exists post_info(
    path ltree primary key references objects (path) on delete cascade ,
    num_comments integer default 0 not null
); create index if not exists post_obj_id_idx on post_info using gist(path);
create sequence if not exists num_posts;
`

var RelationsTable = `
create table if not exists relations__friends (
    user_id1 bigint references users(user_id) on delete cascade on update cascade,
    user_id2 bigint references users(user_id) on delete cascade on update cascade,
    primary key (user_id1, user_id2),
	check (user_id1 != user_id2)
);
create index if not exists rfr_uid1 on relations__friends(user_id1);
create index if not exists rfr_uid2 on relations__friends(user_id2);

create table if not exists relations__subscribers (
    subscriber_id bigint references users(user_id) on delete cascade on update cascade,
    subscribed_id bigint references users(user_id) on delete cascade on update cascade,
    primary key (subscriber_id, subscribed_id),
    check ( subscribed_id != subscriber_id )
);
create index if not exists rsub_uid1 on relations__subscribers(subscriber_id);
create index if not exists rsub_uid2 on relations__subscribers(subscribed_id);
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

CREATE OR REPLACE FUNCTION delete_object_after_process() RETURNS trigger AS $$
    BEGIN
        if nlevel(old.path) > 1 then /*comment*/
            update post_info
                set num_comments = num_comments - 1
            where post_info.path = subpath(old.path,0, 1);
        end if;
       	return new;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS delete_object_after_trigger ON public.objects;
CREATE TRIGGER delete_object_after_trigger
	AFTER DELETE ON objects
FOR EACH ROW EXECUTE PROCEDURE delete_object_after_process();


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
FOR EACH ROW EXECUTE PROCEDURE update_post_text_process ();

CREATE OR REPLACE FUNCTION  add_music() RETURNS trigger AS $$
    BEGIN
        new.document = to_tsvector(new.name) || to_tsvector(new.author);
       	return new;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS add_music_trigger ON public.music;
CREATE TRIGGER add_music_trigger
	BEFORE INSERT ON music
FOR EACH ROW EXECUTE PROCEDURE add_music ();

CREATE OR REPLACE FUNCTION  update_full_name() RETURNS trigger AS $$
    BEGIN
        new.full_name = to_tsvector(new.first_name) || to_tsvector(new.last_name);
       	return new;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_full_name_trigger ON public.users;
CREATE TRIGGER update_full_name_trigger
	BEFORE INSERT OR UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE update_full_name ();
`
var VideoTable = `
create table if not exists video (
    video_id bigserial primary key,
    adder_id bigserial references users(user_id),
    name text not null default 'undefined',
    document tsvector
); create index if not exists video_doc_idx on video using gin(document);
`

var VideoTriggers = `
CREATE OR REPLACE FUNCTION  add_video() RETURNS trigger AS $$
    BEGIN
        new.document = to_tsvector(new.name);
       	return new;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS add_video_trigger ON public.video;
CREATE TRIGGER add_video_trigger
	BEFORE INSERT ON video
FOR EACH ROW EXECUTE PROCEDURE add_video ();
`

var ImagesTable = `
create table images (
    image_id bigserial primary key,
    adder_id bigint references users(user_id)
)
`
var SelectPostsCommentsFunctions = `CREATE OR REPLACE FUNCTION get_comments(post_path text) RETURNS json AS $$
        BEGIN
            return (
select json_agg(to_jsonb(res) - 'lvl')
    from (
            WITH RECURSIVE

                c AS (
                    SELECT text,
                           num_likes,
                           o.path,
                           l2 is not null as me_liked,
                           first_name,
					       o.auth_id,
                           last_name,
                           to_char(creation_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as creation_time,
                           to_char(modification_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as modification_time,
                           nlevel(o.path) AS lvl
                    FROM objects o inner join users on user_id = auth_id
                                   left join likes l2 on o.path = l2.path
                    where o.path <@ text2ltree(post_path) order by o.creation_time
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
			where lvl = 2
        ) res
                );
        end;
    $$ PARALLEL SAFE LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION get_posts(u_id bigint, my_id bigint, _limit bigint, _offset bigint) returns json as $$
        begin
            return
                (
select json_agg(t) from (
      select o.text,
             o.path,
             o.num_likes,
			 o.auth_id,
             to_char(o.creation_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as creation_time,
             to_char(o.modification_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI') as modification_time,
             first_name,
             last_name,
             num_comments,
             l is not null as me_liked
      from objects o
          inner join users u on u.user_id = o.auth_id
          inner join post_info pi on o.path = pi.path
          left join likes l on o.path = l.path and l.auth_id = my_id
      where o.auth_id = u_id order by o.creation_time desc limit _limit offset _offset
    ) t
        );
        end;
    $$ PARALLEL SAFE LANGUAGE plpgsql;

create or replace function push_object(_auth_id bigint, _path ltree, _text text) returns json as $$
    declare result json;
    begin
       insert into objects (auth_id, path, text) values (_auth_id, _path, _text) 
       returning json_build_object('creation_time', to_char(creation_time at time zone 'Europe/Moscow', 'DD.MM.YY HH24:MI'), 'path', path) into result;
       return result;
    end;
$$ language plpgsql;

`

var FriendsSubscribersFunctions = `
create or replace function add_subscriber_to_friend(param_subscribed_id bigint, param_subscriber_id bigint) returns void AS $$
    declare is_deleted bool;
    BEGIN
        with t as (
            delete from relations__subscribers rst where rst.subscriber_id = param_subscriber_id
                                                         and
                                                         rst.subscribed_id = param_subscribed_id
                returning *
        ) select count(*) > 0 from t into is_deleted;
        if is_deleted then
            insert into relations__friends (user_id1, user_id2) values (param_subscribed_id, param_subscriber_id);
        else
            raise exception 'Пользователь не находится в подписчиках!';
        end if;
    END;
$$ LANGUAGE plpgsql;

create or replace function add_friend_to_subscriber (param_subscribed_id bigint, param_subscriber_id bigint) returns void as $$
    declare is_deleted bool;
    BEGIN
        with t as (
            delete from relations__friends where (user_id1 = param_subscriber_id and user_id2 = param_subscribed_id)
                                                     or -- cringe or best practices?
                                                 (user_id1 = param_subscribed_id and user_id2 = param_subscriber_id)
                returning *
        ) select count(*) > 0 from t into is_deleted;

        if is_deleted then
            insert into relations__subscribers (subscriber_id, subscribed_id) values (param_subscriber_id, param_subscribed_id);
        else
            raise exception 'Пользователь не находится в друзьях!';
        end if;
    END;
$$ language plpgsql;

create or replace function get_relationship (requester_id bigint, _user_id2 bigint) returns int2 as $$
    declare relType int2 := 0;

    BEGIN
        if exists(select * from relations__friends
            where user_id1 = requester_id and user_id2 = _user_id2
                     or
                  user_id1 = _user_id2 and user_id2 = requester_id) then
        relType = 3;
        else
            select case when subscriber_id = requester_id then 1 else 2 end from relations__subscribers
                where subscriber_id = requester_id and subscribed_id = _user_id2
                        or
                      subscriber_id = _user_id2 and subscribed_id = requester_id into relType;
        end if;
        return coalesce(relType, 0);
    END;
$$ language plpgsql;

`


var InitTestSQL = `
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs1@mail.ru', '1 2 3 4','Vladimir','Putin', 'М');
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs2@mail.ru', '1 2 3 4','Vladimir','Putin', 'М');
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs3@mail.ru', '1 2 3 4','Vladimir','Putin', 'М');
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs4@mail.ru', '1 2 3 4','Vladimir','Putin', 'М');
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs5@mail.ru', '1 2 3 4','Vladimir','Putin', 'М');
insert into users (email, token, first_name, last_name, sex) values ('baranenkovs6@mail.ru', '1 2 3 4','Vladimir','Putin', 'М');

insert into users (email, token, first_name, last_name, sex) values ('lol@mail.ru', '1 2 3 4','Sergey','Baranenkov', 'М');
insert into users (email, token, first_name, last_name, sex) values ('lol2@mail.ru',  '1 2 3 4','Jury','Dud', 'М');

insert into objects (auth_id, path, text) values (1, '','lol');
insert into objects (auth_id, path, text) values (1, '','lol');
insert into objects (auth_id, path, text) values (1, '', 'lol');
insert into objects (auth_id, path, text) values (1, '1','lol');
insert into objects (auth_id, path, text) values (1, '1.1','lol');

insert into relations__subscribers (subscriber_id, subscribed_id) values (1,2);
select add_subscriber_to_friend(2, 1);

insert into likes (path, auth_id) values ('1', 1);
insert into likes (path, auth_id) values ('1.1', 1);
insert into likes (path, auth_id) values ('2', 1);
select push_message(1,2,'lol1');
`
var MessagesTables = `
create table if not exists conversation (
    conversation_id bigserial primary key,
    user_1 bigint references users(user_id),
    user_2 bigint references users(user_id),
    last_message_id bigint,
    last_message_time timestamptz
    check ( user_1 <= user_2 )
); create index user_both_idx on conversation(user_1, user_2);
   create index user_1_idx on conversation(user_1);
   create index user_2_idx on conversation(user_2);


create table if not exists messages (
    message_id bigserial primary key,
    conversation_id bigint references conversation(conversation_id),
    message_from bigint references users(user_id),
    message_text text,
    created_at timestamptz default now()
); create index conversation_id_idx on messages(conversation_id);
   create index messages_created_at_idx on messages(created_at desc);
`
var MessagesFunctions = `
-- поиск последних бесед и последнее сообщение оттуда
create or replace function select_conversations_list(_user_id bigint, _limit bigint, _offset bigint) returns json as $$
    declare json_res json;
    begin
        with t as (
            select conversation_id, last_message_id, last_message_time, u.first_name, u.last_name, u.user_id as partner_id
                from conversation inner join users u on (case when user_1 = _user_id then user_2 else user_1 end) = u.user_id
            where user_1 = _user_id or user_2 = _user_id
            order by last_message_time desc limit _limit offset _offset
        ),

        j as (
            select m.message_from, m.message_text, t.partner_id, t.first_name, t.last_name, t.conversation_id from t
                inner join messages m on m.conversation_id = t.conversation_id and m.message_id = t.last_message_id order by last_message_time desc
        )
        select json_agg(j) from j into json_res;

        return json_res;
    end;
$$ language plpgsql;



create or replace function on_insert_message_function() returns trigger as $$
    begin
        update conversation set last_message_time = new.created_at, last_message_id = new.message_id
            where conversation_id =  new.conversation_id;
        return new;
    end;
$$ language plpgsql;
create trigger on_insert_message_trigger after insert on messages
    for row
execute procedure on_insert_message_function();



-- выбор для определенной беседы (быстрый)
create or replace function select_conversation_messages(_user_1 bigint, _user_2 bigint, _limit bigint, _offset bigint, out json_res json, out _conversation_id bigint) returns record as $$
    declare min bigint := _user_1;
    declare max bigint := _user_2;

    begin
        if _user_1 > _user_2 then
            min = _user_2;
            max = _user_1;
        end if;
        _conversation_id =  check_conversation_exists(min, max);

        if _conversation_id is not null then
            json_res = (select json_agg(t) from (select m.message_id, m.message_from, m.message_text
                from messages m where m.conversation_id = _conversation_id
            order by m.created_at desc offset _offset limit _limit) t);
            return;
        else
            return;
        end if;
    end;
$$ language plpgsql;


create or replace function check_conversation_exists (_user_1 bigint, _user_2 bigint) returns bigint as $$
    declare res bigint;
    begin
         select conversation_id from conversation where user_1 = _user_1 and user_2 = _user_2 limit 1 into res;
         return res;
    end;
$$ language plpgsql;


create or replace function push_message(_message_from bigint, _message_to bigint, _message_text text) returns json as $$
    declare min bigint := _message_from;
    declare max bigint := _message_to;
    declare _conversation_id bigint;
    declare result json;
    begin
        if _message_from > _message_to then
            min = _message_to;
            max = _message_from;
        end if;
        _conversation_id = check_conversation_exists(min, max);
        if _conversation_id is null then
            insert into conversation (user_1, user_2) values (min, max) returning conversation_id into _conversation_id;
        end if;
        insert into messages (conversation_id, message_from, message_text)
        values (_conversation_id, _message_from, _message_text)
        returning json_build_object('message_id',message_id, 'message_from', message_from, 'message_text',message_text, 'conversation_id', conversation_id) into result;
        return result;
    end;
$$ language plpgsql;

create or replace function get_short_profile_info(_conversation_id bigint, fetcher_id bigint) returns json as $$
    declare result json;
    begin
       select json_build_object('partner_id', user_id, 'first_name',first_name,'last_name',last_name)
        from conversation inner join users u on (case when user_1 = fetcher_id then user_2 else user_1 end) = u.user_id
       where conversation_id = _conversation_id into result;
       return result;
    end;
$$ language plpgsql;
`

var AboutMeFunctions = `
CREATE OR REPLACE FUNCTION  get_extended_info(_user_id bigint) RETURNS json AS $$
    declare result json;
    BEGIN
        select to_json(k) from (select first_name, last_name, sex, tel, city, country, birthday, status,
               hobby, fav_music, fav_films, fav_books, fav_games, other_interests, edu_and_emp_info from users
        where users.user_id = _user_id) k into result;
        return result;
    END;
$$ LANGUAGE plpgsql;
`