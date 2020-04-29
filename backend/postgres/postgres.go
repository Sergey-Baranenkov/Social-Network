package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type RegistrationConn struct {
	Conn *pgx.Conn
}

func (rc *RegistrationConn) InitDatabasesIfNotExist() (err error) {

	if _, err = rc.Conn.Exec(context.Background(), DropDB); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), "create extension if not exists ltree;"); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), UserTable); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), ObjectsTable); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), LikesTable); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), PostInfo); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), Triggers); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), SelectFunctions); err!=nil{
		return err
	}

	if _, err = rc.Conn.Exec(context.Background(), InitTestSQL); err!=nil{
		return err
	}
	return nil
}

func (rc *RegistrationConn) CreateConnection(path string) (err error) {
	Conn, err := pgx.Connect(context.Background(), path)
	if err != nil {
		return err
	}
	rc.Conn = Conn
	return nil
}