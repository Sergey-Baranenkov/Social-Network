package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGXConnection struct {
	Conn *pgxpool.Pool
}

func (rc *PGXConnection) InitDatabasesIfNotExist() (err error) {
	if _, err = rc.Conn.Exec(context.Background(), DropDB); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), "create extension if not exists ltree;"); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), UserTable); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), MusicTable); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), ObjectsTable); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), LikesTable); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), PostInfoTable); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), RelationsTable); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), ImagesTable); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), VideoTable); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), VideoTriggers); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), Triggers); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), SelectPostsCommentsFunctions); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), FriendsSubscribersFunctions); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), MessagesTables); err!=nil{
		return err
	}
	if _, err = rc.Conn.Exec(context.Background(), MessagesFunctions); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), InitTestSQL); err!=nil{
		return err
	}
	return nil
}

func (rc *PGXConnection) CreateConnection(path string) (err error) {
	Conn, err := pgxpool.Connect(context.Background(), path)
	if err != nil {
		return err
	}
	rc.Conn = Conn
	return nil
}