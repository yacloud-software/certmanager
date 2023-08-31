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

 CREATE TABLE storeauth (id integer primary key default nextval('storeauth_seq'),domain text not null  ,token text not null  ,keyauth text not null  ,created integer not null  );

Alter statements:
ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS domain text not null default '';
ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS token text not null default '';
ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS keyauth text not null default '';
ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS created integer not null default 0;


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE storeauth_archive (id integer unique not null,domain text not null,token text not null,keyauth text not null,created integer not null);
*/

import (
	"context"
	gosql "database/sql"
	"fmt"
	savepb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/sql"
	"os"
)

var (
	default_def_DBStoreAuth *DBStoreAuth
)

type DBStoreAuth struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
}

func DefaultDBStoreAuth() *DBStoreAuth {
	if default_def_DBStoreAuth != nil {
		return default_def_DBStoreAuth
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBStoreAuth(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBStoreAuth = res
	return res
}
func NewDBStoreAuth(db *sql.DB) *DBStoreAuth {
	foo := DBStoreAuth{DB: db}
	foo.SQLTablename = "storeauth"
	foo.SQLArchivetablename = "storeauth_archive"
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
	_, e := a.DB.ExecContext(ctx, "archive_DBStoreAuth", "insert into "+a.SQLArchivetablename+" (id,domain, token, keyauth, created) values ($1,$2, $3, $4, $5) ", p.ID, p.Domain, p.Token, p.KeyAuth, p.Created)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBStoreAuth) Save(ctx context.Context, p *savepb.StoreAuth) (uint64, error) {
	qn := "DBStoreAuth_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (domain, token, keyauth, created) values ($1, $2, $3, $4) returning id", p.Domain, p.Token, p.KeyAuth, p.Created)
	if e != nil {
		return 0, a.Error(ctx, qn, e)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, a.Error(ctx, qn, fmt.Errorf("No rows after insert"))
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, a.Error(ctx, qn, fmt.Errorf("failed to scan id after insert: %s", e))
	}
	p.ID = id
	return id, nil
}

// Save using the ID specified
func (a *DBStoreAuth) SaveWithID(ctx context.Context, p *savepb.StoreAuth) error {
	qn := "insert_DBStoreAuth"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,domain, token, keyauth, created) values ($1,$2, $3, $4, $5) ", p.ID, p.Domain, p.Token, p.KeyAuth, p.Created)
	return a.Error(ctx, qn, e)
}

func (a *DBStoreAuth) Update(ctx context.Context, p *savepb.StoreAuth) error {
	qn := "DBStoreAuth_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set domain=$1, token=$2, keyauth=$3, created=$4 where id = $5", p.Domain, p.Token, p.KeyAuth, p.Created, p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBStoreAuth) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBStoreAuth_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBStoreAuth) ByID(ctx context.Context, p uint64) (*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No StoreAuth with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) StoreAuth with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBStoreAuth) TryByID(ctx context.Context, p uint64) (*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) StoreAuth with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBStoreAuth) All(ctx context.Context) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" order by id")
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("All: error querying (%s)", e))
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
	qn := "DBStoreAuth_ByDomain"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where domain = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeDomain(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeDomain"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where domain ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching Token
func (a *DBStoreAuth) ByToken(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByToken"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where token = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByToken: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByToken: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeToken(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeToken"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where token ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByToken: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByToken: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching KeyAuth
func (a *DBStoreAuth) ByKeyAuth(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByKeyAuth"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where keyauth = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKeyAuth: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKeyAuth: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeKeyAuth(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeKeyAuth"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where keyauth ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKeyAuth: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByKeyAuth: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching Created
func (a *DBStoreAuth) ByCreated(ctx context.Context, p uint32) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByCreated"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where created = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreated: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeCreated(ctx context.Context, p uint32) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeCreated"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,domain, token, keyauth, created from "+a.SQLTablename+" where created ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreated: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBStoreAuth) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.StoreAuth, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBStoreAuth) Tablename() string {
	return a.SQLTablename
}

func (a *DBStoreAuth) SelectCols() string {
	return "id,domain, token, keyauth, created"
}
func (a *DBStoreAuth) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".domain, " + a.SQLTablename + ".token, " + a.SQLTablename + ".keyauth, " + a.SQLTablename + ".created"
}

func (a *DBStoreAuth) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.StoreAuth, error) {
	var res []*savepb.StoreAuth
	for rows.Next() {
		foo := savepb.StoreAuth{}
		err := rows.Scan(&foo.ID, &foo.Domain, &foo.Token, &foo.KeyAuth, &foo.Created)
		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, &foo)
	}
	return res, nil
}

/**********************************************************************
* Helper to create table and columns
**********************************************************************/
func (a *DBStoreAuth) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),domain text not null ,token text not null ,keyauth text not null ,created integer not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),domain text not null ,token text not null ,keyauth text not null ,created integer not null );`,
		`ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS domain text not null default '';`,
		`ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS token text not null default '';`,
		`ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS keyauth text not null default '';`,
		`ALTER TABLE storeauth ADD COLUMN IF NOT EXISTS created integer not null default 0;`,

		`ALTER TABLE storeauth_archive ADD COLUMN IF NOT EXISTS domain text not null default '';`,
		`ALTER TABLE storeauth_archive ADD COLUMN IF NOT EXISTS token text not null default '';`,
		`ALTER TABLE storeauth_archive ADD COLUMN IF NOT EXISTS keyauth text not null default '';`,
		`ALTER TABLE storeauth_archive ADD COLUMN IF NOT EXISTS created integer not null default 0;`,
	}
	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
	}

	// these are optional, expected to fail
	csql = []string{
		// Indices:

		// Foreign keys:

	}
	for i, c := range csql {
		a.DB.ExecContextQuiet(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
	}
	return nil
}

/**********************************************************************
* Helper to meaningful errors
**********************************************************************/
func (a *DBStoreAuth) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}
