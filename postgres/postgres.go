package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type RegistrationConn struct {
	Conn *pgx.Conn
}

func (rc *RegistrationConn) CreateRegTable() (err error) {
	_, err = rc.Conn.Exec(context.Background(), "create table if not exists registration (user_id serial primary key,"+
		"email text not null,"+
		"first_name text not null,"+
		"last_name text not null,"+
		"token bytea not null)")
	return err
}

func (rc *RegistrationConn) CreateConnection(path string) (err error) {
	Conn, err := pgx.Connect(context.Background(), path)
	if err != nil {
		return err
	}
	rc.Conn = Conn
	return nil
}

func (rc *RegistrationConn) Close()(err error){
	return rc.Conn.Close(nil)
}