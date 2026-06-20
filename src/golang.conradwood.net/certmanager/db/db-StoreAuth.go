package db

/*
 This file was created by mkdb-client.
 The intention is not to modify this file, but you may extend the struct DBStoreAuth
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
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/sql"
	"os"
	"sync"
)

var (
	default_def_DBStoreAuth *DBStoreAuth
)

type DBStoreAuth struct {
	DB                   *sql.DB
	SQLTablename         string
	SQLArchivetablename  string
	customColumnHandlers []CustomColumnHandler
	lock                 sync.Mutex
}

func init() {
	RegisterDBHandlerFactory(func() Handler {
		return DefaultDBStoreAuth()
	})
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

func (a *DBStoreAuth) GetCustomColumnHandlers() []CustomColumnHandler {
	return a.customColumnHandlers
}
func (a *DBStoreAuth) AddCustomColumnHandler(w CustomColumnHandler) {
	a.lock.Lock()
	a.customColumnHandlers = append(a.customColumnHandlers, w)
	a.lock.Unlock()
}

func (a *DBStoreAuth) NewQuery() *Query {
	return newQuery(a)
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

// return a map with columnname -> value_from_proto
func (a *DBStoreAuth) buildSaveMap(ctx context.Context, p *savepb.StoreAuth) (map[string]interface{}, error) {
	extra, err := extraFieldsToStore(ctx, a, p)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["id"] = a.get_col_from_proto(p, "id")
	res["domain"] = a.get_col_from_proto(p, "domain")
	res["token"] = a.get_col_from_proto(p, "token")
	res["keyauth"] = a.get_col_from_proto(p, "keyauth")
	res["created"] = a.get_col_from_proto(p, "created")
	if extra != nil {
		for k, v := range extra {
			res[k] = v
		}
	}
	return res, nil
}

func (a *DBStoreAuth) Save(ctx context.Context, p *savepb.StoreAuth) (uint64, error) {
	qn := "save_DBStoreAuth"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return 0, err
	}
	delete(smap, "id") // save without id
	return a.saveMap(ctx, qn, smap, p)
}

// Save using the ID specified
func (a *DBStoreAuth) SaveWithID(ctx context.Context, p *savepb.StoreAuth) error {
	qn := "insert_DBStoreAuth"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return err
	}
	_, err = a.saveMap(ctx, qn, smap, p)
	return err
}

// use a hashmap of columnname->values to store to database (see buildSaveMap())
func (a *DBStoreAuth) saveMap(ctx context.Context, queryname string, smap map[string]interface{}, p *savepb.StoreAuth) (uint64, error) {
	// Save (and use database default ID generation)

	var rows *gosql.Rows
	var e error

	q_cols := ""
	q_valnames := ""
	q_vals := make([]interface{}, 0)
	deli := ""
	i := 0
	// build the 2 parts of the query (column names and value names) as well as the values themselves
	for colname, val := range smap {
		q_cols = q_cols + deli + colname
		i++
		q_valnames = q_valnames + deli + fmt.Sprintf("$%d", i)
		q_vals = append(q_vals, val)
		deli = ","
	}
	rows, e = a.DB.QueryContext(ctx, queryname, "insert into "+a.SQLTablename+" ("+q_cols+") values ("+q_valnames+") returning id", q_vals...)
	if e != nil {
		return 0, a.Error(ctx, queryname, e)
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, a.Error(ctx, queryname, errors.Errorf("No rows after insert"))
	}
	var id uint64
	e = rows.Scan(&id)
	if e != nil {
		return 0, a.Error(ctx, queryname, errors.Errorf("failed to scan id after insert: %s", e))
	}
	p.ID = id
	return id, nil
}

// if ID==0 save, otherwise update
func (a *DBStoreAuth) SaveOrUpdate(ctx context.Context, p *savepb.StoreAuth) error {
	if p.ID == 0 {
		_, err := a.Save(ctx, p)
		return err
	}
	return a.Update(ctx, p)
}
func (a *DBStoreAuth) Update(ctx context.Context, p *savepb.StoreAuth) error {
	qn := "DBStoreAuth_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set domain=$1, token=$2, keyauth=$3, created=$4 where id = $5", a.get_Domain(p), a.get_Token(p), a.get_KeyAuth(p), a.get_Created(p), p.ID)

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
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, errors.Errorf("No StoreAuth with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) StoreAuth with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBStoreAuth) TryByID(ctx context.Context, p uint64) (*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_TryByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) StoreAuth with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by multiple primary ids
func (a *DBStoreAuth) ByIDs(ctx context.Context, p []uint64) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByIDs"
	l, e := a.fromQuery(ctx, qn, "id in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	return l, nil
}

// get all rows
func (a *DBStoreAuth) All(ctx context.Context) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_all"
	l, e := a.fromQuery(ctx, qn, "true")
	if e != nil {
		return nil, errors.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBStoreAuth" rows with matching Domain
func (a *DBStoreAuth) ByDomain(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByDomain"
	l, e := a.fromQuery(ctx, qn, "domain = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with multiple matching Domain
func (a *DBStoreAuth) ByMultiDomain(ctx context.Context, p []string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByDomain"
	l, e := a.fromQuery(ctx, qn, "domain in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeDomain(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeDomain"
	l, e := a.fromQuery(ctx, qn, "domain ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByDomain: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching Token
func (a *DBStoreAuth) ByToken(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByToken"
	l, e := a.fromQuery(ctx, qn, "token = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByToken: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with multiple matching Token
func (a *DBStoreAuth) ByMultiToken(ctx context.Context, p []string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByToken"
	l, e := a.fromQuery(ctx, qn, "token in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByToken: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeToken(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeToken"
	l, e := a.fromQuery(ctx, qn, "token ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByToken: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching KeyAuth
func (a *DBStoreAuth) ByKeyAuth(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByKeyAuth"
	l, e := a.fromQuery(ctx, qn, "keyauth = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByKeyAuth: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with multiple matching KeyAuth
func (a *DBStoreAuth) ByMultiKeyAuth(ctx context.Context, p []string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByKeyAuth"
	l, e := a.fromQuery(ctx, qn, "keyauth in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByKeyAuth: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeKeyAuth(ctx context.Context, p string) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeKeyAuth"
	l, e := a.fromQuery(ctx, qn, "keyauth ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByKeyAuth: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with matching Created
func (a *DBStoreAuth) ByCreated(ctx context.Context, p uint32) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByCreated"
	l, e := a.fromQuery(ctx, qn, "created = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBStoreAuth" rows with multiple matching Created
func (a *DBStoreAuth) ByMultiCreated(ctx context.Context, p []uint32) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByCreated"
	l, e := a.fromQuery(ctx, qn, "created in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBStoreAuth) ByLikeCreated(ctx context.Context, p uint32) ([]*savepb.StoreAuth, error) {
	qn := "DBStoreAuth_ByLikeCreated"
	l, e := a.fromQuery(ctx, qn, "created ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* The field getters
**********************************************************************/

// getter for field "ID" (ID) [uint64]
func (a *DBStoreAuth) get_ID(p *savepb.StoreAuth) uint64 {
	return uint64(p.ID)
}

// getter for field "Domain" (Domain) [string]
func (a *DBStoreAuth) get_Domain(p *savepb.StoreAuth) string {
	return string(p.Domain)
}

// getter for field "Token" (Token) [string]
func (a *DBStoreAuth) get_Token(p *savepb.StoreAuth) string {
	return string(p.Token)
}

// getter for field "KeyAuth" (KeyAuth) [string]
func (a *DBStoreAuth) get_KeyAuth(p *savepb.StoreAuth) string {
	return string(p.KeyAuth)
}

// getter for field "Created" (Created) [uint32]
func (a *DBStoreAuth) get_Created(p *savepb.StoreAuth) uint32 {
	return uint32(p.Created)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBStoreAuth) ByDBQuery(ctx context.Context, query *Query) ([]*savepb.StoreAuth, error) {
	extra_fields, err := extraFieldsToQuery(ctx, a)
	if err != nil {
		return nil, err
	}
	i := 0
	for col_name, value := range extra_fields {
		i++
		/*
		   efname:=fmt.Sprintf("EXTRA_FIELD_%d",i)
		   query.Add(col_name+" = "+efname,QP{efname:value})
		*/
		query.AddEqual(col_name, value)
	}

	gw, paras := query.ToPostgres()
	queryname := "custom_dbquery"
	rows, err := a.DB.QueryContext(ctx, queryname, "select "+a.SelectCols()+" from "+a.Tablename()+" where "+gw, paras...)
	if err != nil {
		return nil, err
	}
	res, err := a.FromRows(ctx, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *DBStoreAuth) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.StoreAuth, error) {
	return a.fromQuery(ctx, "custom_query_"+a.Tablename(), query_where, args...)
}

// from a query snippet (the part after WHERE)
func (a *DBStoreAuth) fromQuery(ctx context.Context, queryname string, query_where string, args ...interface{}) ([]*savepb.StoreAuth, error) {
	extra_fields, err := extraFieldsToQuery(ctx, a)
	if err != nil {
		return nil, err
	}
	eq := ""
	if extra_fields != nil && len(extra_fields) > 0 {
		eq = " AND ("
		// build the extraquery "eq"
		i := len(args)
		deli := ""
		for col_name, value := range extra_fields {
			i++
			eq = eq + deli + col_name + fmt.Sprintf(" = $%d", i)
			deli = " AND "
			args = append(args, value)
		}
		eq = eq + ")"
	}
	rows, err := a.DB.QueryContext(ctx, queryname, "select "+a.SelectCols()+" from "+a.Tablename()+" where ( "+query_where+") "+eq, args...)
	if err != nil {
		return nil, err
	}
	res, err := a.FromRows(ctx, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return res, nil
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBStoreAuth) get_col_from_proto(p *savepb.StoreAuth, colname string) interface{} {
	if colname == "id" {
		return a.get_ID(p)
	} else if colname == "domain" {
		return a.get_Domain(p)
	} else if colname == "token" {
		return a.get_Token(p)
	} else if colname == "keyauth" {
		return a.get_KeyAuth(p)
	} else if colname == "created" {
		return a.get_Created(p)
	}
	panic(fmt.Sprintf("in table \"%s\", column \"%s\" cannot be resolved to proto field name", a.Tablename(), colname))
}

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
		// SCANNER:
		foo := &savepb.StoreAuth{}
		// create the non-nullable pointers
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.Domain
		scanTarget_2 := &foo.Token
		scanTarget_3 := &foo.KeyAuth
		scanTarget_4 := &foo.Created
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2, scanTarget_3, scanTarget_4)
		// END SCANNER

		if err != nil {
			return nil, a.Error(ctx, "fromrow-scan", err)
		}
		res = append(res, foo)
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
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS domain text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS token text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS keyauth text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS created integer not null default 0;`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS domain text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS token text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS keyauth text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS created integer not null  default 0;`,
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
	return errors.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}

