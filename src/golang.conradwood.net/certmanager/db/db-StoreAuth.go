package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBStoreAuth
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence storeauth_seq;

Main Table:

 CREATE TABLE storeauth (id integer primary key default nextval('storeauth_seq'),domain varchar(2000) not null,token varchar(2000) not null,keyauth varchar(2000) not null,created integer not null);

Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE storeauth_archive (id integer unique not null,domain varchar(2000) not null,token varchar(2000) not null,keyauth varchar(2000) not null,created integer not null);
*/

import (
	gosql "database/sql"
	"fmt"
	savepb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/sql"
	"golang.org/x/net/context"
)

type DBStoreAuth struct {
	DB *sql.DB
}

func NewDBStoreAuth(db *sql.DB) *DBStoreAuth {
	foo := DBStoreAuth{DB: db}
	return &foo
}

// archive. It is NOT transactionally save.
func (a *DBStoreAuth) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "insert_DBStoreAuth", "insert into storeauth_archive (id,domain, token, keyauth, created) values ($1,$2, $3, $4, $5) ", p.ID, p.Domain, p.Token, p.KeyAuth, p.Created)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBStoreAuth) Save(ctx context.Context, p *savepb.StoreAuth) (uint64, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_Save", "insert into storeauth (domain, token, keyauth, created) values ($1, $2, $3, $4) returning id", p.Domain, p.Token, p.KeyAuth, p.Created)
	if e != nil {
		return 0, e
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, fmt.Errorf("No rows after insert")
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, fmt.Errorf("failed to scan id after insert: %s", e)
	}
	p.ID = id
	return id, nil
}

// Save using the ID specified
func (a *DBStoreAuth) SaveWithID(ctx context.Context, p *savepb.StoreAuth) error {
	_, e := a.DB.ExecContext(ctx, "insert_DBStoreAuth", "insert into storeauth (id,domain, token, keyauth, created) values ($1,$2, $3, $4, $5) ", p.ID, p.Domain, p.Token, p.KeyAuth, p.Created)
	return e
}

func (a *DBStoreAuth) Update(ctx context.Context, p *savepb.StoreAuth) error {
	_, e := a.DB.ExecContext(ctx, "DBStoreAuth_Update", "update storeauth set domain=$1, token=$2, keyauth=$3, created=$4 where id = $5", p.Domain, p.Token, p.KeyAuth, p.Created, p.ID)

	return e
}

// delete by id field
func (a *DBStoreAuth) DeleteByID(ctx context.Context, p uint64) error {
	_, e := a.DB.ExecContext(ctx, "deleteDBStoreAuth_ByID", "delete from storeauth where id = $1", p)
	return e
}

// get it by primary id
func (a *DBStoreAuth) ByID(ctx context.Context, p uint64) (*savepb.StoreAuth, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_ByID", "select id,domain, token, keyauth, created from storeauth where id = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByID: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByID: error scanning (%s)", e)
	}
	if len(l) == 0 {
		return nil, fmt.Errorf("No StoreAuth with id %d", p)
	}
	if len(l) != 1 {
		return nil, fmt.Errorf("Multiple (%d) StoreAuth with id %d", len(l), p)
	}
	return l[0], nil
}

// get all rows
func (a *DBStoreAuth) All(ctx context.Context) ([]*savepb.StoreAuth, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_all", "select id,domain, token, keyauth, created from storeauth order by id")
	if e != nil {
		return nil, fmt.Errorf("All: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBStoreAuth" rows with matching Domain
func (a *DBStoreAuth) ByDomain(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_ByDomain", "select id,domain, token, keyauth, created from storeauth where domain = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByDomain: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByDomain: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching Token
func (a *DBStoreAuth) ByToken(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_ByToken", "select id,domain, token, keyauth, created from storeauth where token = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByToken: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByToken: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching KeyAuth
func (a *DBStoreAuth) ByKeyAuth(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_ByKeyAuth", "select id,domain, token, keyauth, created from storeauth where keyauth = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByKeyAuth: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByKeyAuth: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching Created
func (a *DBStoreAuth) ByCreated(ctx context.Context, p uint32) ([]*savepb.StoreAuth, error) {
	rows, e := a.DB.QueryContext(ctx, "DBStoreAuth_ByCreated", "select id,domain, token, keyauth, created from storeauth where created = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByCreated: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByCreated: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBStoreAuth) Tablename() string {
	return "storeauth"
}

func (a *DBStoreAuth) SelectCols() string {
	return "id,domain, token, keyauth, created"
}

func (a *DBStoreAuth) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.StoreAuth, error) {
	var res []*savepb.StoreAuth
	for rows.Next() {
		foo := savepb.StoreAuth{}
		err := rows.Scan(&foo.ID, &foo.Domain, &foo.Token, &foo.KeyAuth, &foo.Created)
		if err != nil {
			return nil, err
		}
		res = append(res, &foo)
	}
	return res, nil
}
