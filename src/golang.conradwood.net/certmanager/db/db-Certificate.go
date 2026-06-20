package db

/*
 This file was created by mkdb-client.
 The intention is not to modify this file, but you may extend the struct DBCertificate
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence certificate_seq;

Main Table:

 CREATE TABLE certificate (id integer primary key default nextval('certificate_seq'),host text not null  ,pemcertificate text not null  ,pemprivatekey text not null  ,pemca text not null  ,created integer not null  ,expiry integer not null  ,creatoruser text not null  ,creatorservice text not null  ,lastattempt integer not null  ,lasterror text not null  ,islocalca boolean not null  ,islocalcert boolean not null  ,pempublickey text not null  );

Alter statements:
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS host text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS pemcertificate text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS pemprivatekey text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS pemca text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS created integer not null default 0;
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS expiry integer not null default 0;
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS creatoruser text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS creatorservice text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS lastattempt integer not null default 0;
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS lasterror text not null default '';
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS islocalca boolean not null default false;
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS islocalcert boolean not null default false;
ALTER TABLE certificate ADD COLUMN IF NOT EXISTS pempublickey text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE certificate_archive (id integer unique not null,host text not null,pemcertificate text not null,pemprivatekey text not null,pemca text not null,created integer not null,expiry integer not null,creatoruser text not null,creatorservice text not null,lastattempt integer not null,lasterror text not null,islocalca boolean not null,islocalcert boolean not null,pempublickey text not null);
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
	default_def_DBCertificate *DBCertificate
)

type DBCertificate struct {
	DB                   *sql.DB
	SQLTablename         string
	SQLArchivetablename  string
	customColumnHandlers []CustomColumnHandler
	lock                 sync.Mutex
}

func init() {
	RegisterDBHandlerFactory(func() Handler {
		return DefaultDBCertificate()
	})
}

func DefaultDBCertificate() *DBCertificate {
	if default_def_DBCertificate != nil {
		return default_def_DBCertificate
	}
	psql, err := sql.Open()
	if err != nil {
		fmt.Printf("Failed to open database: %s\n", err)
		os.Exit(10)
	}
	res := NewDBCertificate(psql)
	ctx := context.Background()
	err = res.CreateTable(ctx)
	if err != nil {
		fmt.Printf("Failed to create table: %s\n", err)
		os.Exit(10)
	}
	default_def_DBCertificate = res
	return res
}
func NewDBCertificate(db *sql.DB) *DBCertificate {
	foo := DBCertificate{DB: db}
	foo.SQLTablename = "certificate"
	foo.SQLArchivetablename = "certificate_archive"
	return &foo
}

func (a *DBCertificate) GetCustomColumnHandlers() []CustomColumnHandler {
	return a.customColumnHandlers
}
func (a *DBCertificate) AddCustomColumnHandler(w CustomColumnHandler) {
	a.lock.Lock()
	a.customColumnHandlers = append(a.customColumnHandlers, w)
	a.lock.Unlock()
}

func (a *DBCertificate) NewQuery() *Query {
	return newQuery(a)
}

// archive. It is NOT transactionally save.
func (a *DBCertificate) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBCertificate", "insert into "+a.SQLArchivetablename+" (id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror, islocalca, islocalcert, pempublickey) values ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) ", p.ID, p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt, p.LastError, p.IsLocalCA, p.IsLocalCert, p.PemPublicKey)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// return a map with columnname -> value_from_proto
func (a *DBCertificate) buildSaveMap(ctx context.Context, p *savepb.Certificate) (map[string]interface{}, error) {
	extra, err := extraFieldsToStore(ctx, a, p)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["id"] = a.get_col_from_proto(p, "id")
	res["host"] = a.get_col_from_proto(p, "host")
	res["pemcertificate"] = a.get_col_from_proto(p, "pemcertificate")
	res["pemprivatekey"] = a.get_col_from_proto(p, "pemprivatekey")
	res["pemca"] = a.get_col_from_proto(p, "pemca")
	res["created"] = a.get_col_from_proto(p, "created")
	res["expiry"] = a.get_col_from_proto(p, "expiry")
	res["creatoruser"] = a.get_col_from_proto(p, "creatoruser")
	res["creatorservice"] = a.get_col_from_proto(p, "creatorservice")
	res["lastattempt"] = a.get_col_from_proto(p, "lastattempt")
	res["lasterror"] = a.get_col_from_proto(p, "lasterror")
	res["islocalca"] = a.get_col_from_proto(p, "islocalca")
	res["islocalcert"] = a.get_col_from_proto(p, "islocalcert")
	res["pempublickey"] = a.get_col_from_proto(p, "pempublickey")
	if extra != nil {
		for k, v := range extra {
			res[k] = v
		}
	}
	return res, nil
}

func (a *DBCertificate) Save(ctx context.Context, p *savepb.Certificate) (uint64, error) {
	qn := "save_DBCertificate"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return 0, err
	}
	delete(smap, "id") // save without id
	return a.saveMap(ctx, qn, smap, p)
}

// Save using the ID specified
func (a *DBCertificate) SaveWithID(ctx context.Context, p *savepb.Certificate) error {
	qn := "insert_DBCertificate"
	smap, err := a.buildSaveMap(ctx, p)
	if err != nil {
		return err
	}
	_, err = a.saveMap(ctx, qn, smap, p)
	return err
}

// use a hashmap of columnname->values to store to database (see buildSaveMap())
func (a *DBCertificate) saveMap(ctx context.Context, queryname string, smap map[string]interface{}, p *savepb.Certificate) (uint64, error) {
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
func (a *DBCertificate) SaveOrUpdate(ctx context.Context, p *savepb.Certificate) error {
	if p.ID == 0 {
		_, err := a.Save(ctx, p)
		return err
	}
	return a.Update(ctx, p)
}
func (a *DBCertificate) Update(ctx context.Context, p *savepb.Certificate) error {
	qn := "DBCertificate_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set host=$1, pemcertificate=$2, pemprivatekey=$3, pemca=$4, created=$5, expiry=$6, creatoruser=$7, creatorservice=$8, lastattempt=$9, lasterror=$10, islocalca=$11, islocalcert=$12, pempublickey=$13 where id = $14", a.get_Host(p), a.get_PemCertificate(p), a.get_PemPrivateKey(p), a.get_PemCA(p), a.get_Created(p), a.get_Expiry(p), a.get_CreatorUser(p), a.get_CreatorService(p), a.get_LastAttempt(p), a.get_LastError(p), a.get_IsLocalCA(p), a.get_IsLocalCert(p), a.get_PemPublicKey(p), p.ID)

	return a.Error(ctx, qn, e)
}

// delete by id field
func (a *DBCertificate) DeleteByID(ctx context.Context, p uint64) error {
	qn := "deleteDBCertificate_ByID"
	_, e := a.DB.ExecContext(ctx, qn, "delete from "+a.SQLTablename+" where id = $1", p)
	return a.Error(ctx, qn, e)
}

// get it by primary id
func (a *DBCertificate) ByID(ctx context.Context, p uint64) (*savepb.Certificate, error) {
	qn := "DBCertificate_ByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, errors.Errorf("No Certificate with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) Certificate with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBCertificate) TryByID(ctx context.Context, p uint64) (*savepb.Certificate, error) {
	qn := "DBCertificate_TryByID"
	l, e := a.fromQuery(ctx, qn, "id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, nil
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, errors.Errorf("Multiple (%d) Certificate with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by multiple primary ids
func (a *DBCertificate) ByIDs(ctx context.Context, p []uint64) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByIDs"
	l, e := a.fromQuery(ctx, qn, "id in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("TryByID: error scanning (%s)", e))
	}
	return l, nil
}

// get all rows
func (a *DBCertificate) All(ctx context.Context) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_all"
	l, e := a.fromQuery(ctx, qn, "true")
	if e != nil {
		return nil, errors.Errorf("All: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* GetBy[FIELD] functions
**********************************************************************/

// get all "DBCertificate" rows with matching Host
func (a *DBCertificate) ByHost(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByHost"
	l, e := a.fromQuery(ctx, qn, "host = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByHost: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching Host
func (a *DBCertificate) ByMultiHost(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByHost"
	l, e := a.fromQuery(ctx, qn, "host in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByHost: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeHost(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeHost"
	l, e := a.fromQuery(ctx, qn, "host ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByHost: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemCertificate
func (a *DBCertificate) ByPemCertificate(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemCertificate"
	l, e := a.fromQuery(ctx, qn, "pemcertificate = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemCertificate: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching PemCertificate
func (a *DBCertificate) ByMultiPemCertificate(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemCertificate"
	l, e := a.fromQuery(ctx, qn, "pemcertificate in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemCertificate: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemCertificate(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemCertificate"
	l, e := a.fromQuery(ctx, qn, "pemcertificate ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemCertificate: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemPrivateKey
func (a *DBCertificate) ByPemPrivateKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemPrivateKey"
	l, e := a.fromQuery(ctx, qn, "pemprivatekey = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemPrivateKey: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching PemPrivateKey
func (a *DBCertificate) ByMultiPemPrivateKey(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemPrivateKey"
	l, e := a.fromQuery(ctx, qn, "pemprivatekey in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemPrivateKey: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemPrivateKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemPrivateKey"
	l, e := a.fromQuery(ctx, qn, "pemprivatekey ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemPrivateKey: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemCA
func (a *DBCertificate) ByPemCA(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemCA"
	l, e := a.fromQuery(ctx, qn, "pemca = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemCA: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching PemCA
func (a *DBCertificate) ByMultiPemCA(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemCA"
	l, e := a.fromQuery(ctx, qn, "pemca in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemCA: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemCA(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemCA"
	l, e := a.fromQuery(ctx, qn, "pemca ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemCA: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching Created
func (a *DBCertificate) ByCreated(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreated"
	l, e := a.fromQuery(ctx, qn, "created = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching Created
func (a *DBCertificate) ByMultiCreated(ctx context.Context, p []uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreated"
	l, e := a.fromQuery(ctx, qn, "created in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeCreated(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeCreated"
	l, e := a.fromQuery(ctx, qn, "created ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreated: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching Expiry
func (a *DBCertificate) ByExpiry(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByExpiry"
	l, e := a.fromQuery(ctx, qn, "expiry = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByExpiry: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching Expiry
func (a *DBCertificate) ByMultiExpiry(ctx context.Context, p []uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByExpiry"
	l, e := a.fromQuery(ctx, qn, "expiry in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByExpiry: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeExpiry(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeExpiry"
	l, e := a.fromQuery(ctx, qn, "expiry ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByExpiry: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching CreatorUser
func (a *DBCertificate) ByCreatorUser(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreatorUser"
	l, e := a.fromQuery(ctx, qn, "creatoruser = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreatorUser: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching CreatorUser
func (a *DBCertificate) ByMultiCreatorUser(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreatorUser"
	l, e := a.fromQuery(ctx, qn, "creatoruser in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreatorUser: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeCreatorUser(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeCreatorUser"
	l, e := a.fromQuery(ctx, qn, "creatoruser ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreatorUser: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching CreatorService
func (a *DBCertificate) ByCreatorService(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreatorService"
	l, e := a.fromQuery(ctx, qn, "creatorservice = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreatorService: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching CreatorService
func (a *DBCertificate) ByMultiCreatorService(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreatorService"
	l, e := a.fromQuery(ctx, qn, "creatorservice in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreatorService: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeCreatorService(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeCreatorService"
	l, e := a.fromQuery(ctx, qn, "creatorservice ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByCreatorService: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching LastAttempt
func (a *DBCertificate) ByLastAttempt(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLastAttempt"
	l, e := a.fromQuery(ctx, qn, "lastattempt = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByLastAttempt: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching LastAttempt
func (a *DBCertificate) ByMultiLastAttempt(ctx context.Context, p []uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLastAttempt"
	l, e := a.fromQuery(ctx, qn, "lastattempt in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByLastAttempt: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeLastAttempt(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeLastAttempt"
	l, e := a.fromQuery(ctx, qn, "lastattempt ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByLastAttempt: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching LastError
func (a *DBCertificate) ByLastError(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLastError"
	l, e := a.fromQuery(ctx, qn, "lasterror = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByLastError: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching LastError
func (a *DBCertificate) ByMultiLastError(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLastError"
	l, e := a.fromQuery(ctx, qn, "lasterror in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByLastError: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeLastError(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeLastError"
	l, e := a.fromQuery(ctx, qn, "lasterror ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByLastError: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching IsLocalCA
func (a *DBCertificate) ByIsLocalCA(ctx context.Context, p bool) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByIsLocalCA"
	l, e := a.fromQuery(ctx, qn, "islocalca = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByIsLocalCA: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching IsLocalCA
func (a *DBCertificate) ByMultiIsLocalCA(ctx context.Context, p []bool) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByIsLocalCA"
	l, e := a.fromQuery(ctx, qn, "islocalca in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByIsLocalCA: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeIsLocalCA(ctx context.Context, p bool) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeIsLocalCA"
	l, e := a.fromQuery(ctx, qn, "islocalca ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByIsLocalCA: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching IsLocalCert
func (a *DBCertificate) ByIsLocalCert(ctx context.Context, p bool) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByIsLocalCert"
	l, e := a.fromQuery(ctx, qn, "islocalcert = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByIsLocalCert: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching IsLocalCert
func (a *DBCertificate) ByMultiIsLocalCert(ctx context.Context, p []bool) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByIsLocalCert"
	l, e := a.fromQuery(ctx, qn, "islocalcert in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByIsLocalCert: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeIsLocalCert(ctx context.Context, p bool) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeIsLocalCert"
	l, e := a.fromQuery(ctx, qn, "islocalcert ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByIsLocalCert: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemPublicKey
func (a *DBCertificate) ByPemPublicKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemPublicKey"
	l, e := a.fromQuery(ctx, qn, "pempublickey = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemPublicKey: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with multiple matching PemPublicKey
func (a *DBCertificate) ByMultiPemPublicKey(ctx context.Context, p []string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemPublicKey"
	l, e := a.fromQuery(ctx, qn, "pempublickey in $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemPublicKey: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemPublicKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemPublicKey"
	l, e := a.fromQuery(ctx, qn, "pempublickey ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, errors.Errorf("ByPemPublicKey: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* The field getters
**********************************************************************/

// getter for field "ID" (ID) [uint64]
func (a *DBCertificate) get_ID(p *savepb.Certificate) uint64 {
	return uint64(p.ID)
}

// getter for field "Host" (Host) [string]
func (a *DBCertificate) get_Host(p *savepb.Certificate) string {
	return string(p.Host)
}

// getter for field "PemCertificate" (PemCertificate) [string]
func (a *DBCertificate) get_PemCertificate(p *savepb.Certificate) string {
	return string(p.PemCertificate)
}

// getter for field "PemPrivateKey" (PemPrivateKey) [string]
func (a *DBCertificate) get_PemPrivateKey(p *savepb.Certificate) string {
	return string(p.PemPrivateKey)
}

// getter for field "PemCA" (PemCA) [string]
func (a *DBCertificate) get_PemCA(p *savepb.Certificate) string {
	return string(p.PemCA)
}

// getter for field "Created" (Created) [uint32]
func (a *DBCertificate) get_Created(p *savepb.Certificate) uint32 {
	return uint32(p.Created)
}

// getter for field "Expiry" (Expiry) [uint32]
func (a *DBCertificate) get_Expiry(p *savepb.Certificate) uint32 {
	return uint32(p.Expiry)
}

// getter for field "CreatorUser" (CreatorUser) [string]
func (a *DBCertificate) get_CreatorUser(p *savepb.Certificate) string {
	return string(p.CreatorUser)
}

// getter for field "CreatorService" (CreatorService) [string]
func (a *DBCertificate) get_CreatorService(p *savepb.Certificate) string {
	return string(p.CreatorService)
}

// getter for field "LastAttempt" (LastAttempt) [uint32]
func (a *DBCertificate) get_LastAttempt(p *savepb.Certificate) uint32 {
	return uint32(p.LastAttempt)
}

// getter for field "LastError" (LastError) [string]
func (a *DBCertificate) get_LastError(p *savepb.Certificate) string {
	return string(p.LastError)
}

// getter for field "IsLocalCA" (IsLocalCA) [bool]
func (a *DBCertificate) get_IsLocalCA(p *savepb.Certificate) bool {
	return bool(p.IsLocalCA)
}

// getter for field "IsLocalCert" (IsLocalCert) [bool]
func (a *DBCertificate) get_IsLocalCert(p *savepb.Certificate) bool {
	return bool(p.IsLocalCert)
}

// getter for field "PemPublicKey" (PemPublicKey) [string]
func (a *DBCertificate) get_PemPublicKey(p *savepb.Certificate) string {
	return string(p.PemPublicKey)
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBCertificate) ByDBQuery(ctx context.Context, query *Query) ([]*savepb.Certificate, error) {
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

func (a *DBCertificate) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.Certificate, error) {
	return a.fromQuery(ctx, "custom_query_"+a.Tablename(), query_where, args...)
}

// from a query snippet (the part after WHERE)
func (a *DBCertificate) fromQuery(ctx context.Context, queryname string, query_where string, args ...interface{}) ([]*savepb.Certificate, error) {
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
func (a *DBCertificate) get_col_from_proto(p *savepb.Certificate, colname string) interface{} {
	if colname == "id" {
		return a.get_ID(p)
	} else if colname == "host" {
		return a.get_Host(p)
	} else if colname == "pemcertificate" {
		return a.get_PemCertificate(p)
	} else if colname == "pemprivatekey" {
		return a.get_PemPrivateKey(p)
	} else if colname == "pemca" {
		return a.get_PemCA(p)
	} else if colname == "created" {
		return a.get_Created(p)
	} else if colname == "expiry" {
		return a.get_Expiry(p)
	} else if colname == "creatoruser" {
		return a.get_CreatorUser(p)
	} else if colname == "creatorservice" {
		return a.get_CreatorService(p)
	} else if colname == "lastattempt" {
		return a.get_LastAttempt(p)
	} else if colname == "lasterror" {
		return a.get_LastError(p)
	} else if colname == "islocalca" {
		return a.get_IsLocalCA(p)
	} else if colname == "islocalcert" {
		return a.get_IsLocalCert(p)
	} else if colname == "pempublickey" {
		return a.get_PemPublicKey(p)
	}
	panic(fmt.Sprintf("in table \"%s\", column \"%s\" cannot be resolved to proto field name", a.Tablename(), colname))
}

func (a *DBCertificate) Tablename() string {
	return a.SQLTablename
}

func (a *DBCertificate) SelectCols() string {
	return "id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror, islocalca, islocalcert, pempublickey"
}
func (a *DBCertificate) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".host, " + a.SQLTablename + ".pemcertificate, " + a.SQLTablename + ".pemprivatekey, " + a.SQLTablename + ".pemca, " + a.SQLTablename + ".created, " + a.SQLTablename + ".expiry, " + a.SQLTablename + ".creatoruser, " + a.SQLTablename + ".creatorservice, " + a.SQLTablename + ".lastattempt, " + a.SQLTablename + ".lasterror, " + a.SQLTablename + ".islocalca, " + a.SQLTablename + ".islocalcert, " + a.SQLTablename + ".pempublickey"
}

func (a *DBCertificate) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.Certificate, error) {
	var res []*savepb.Certificate
	for rows.Next() {
		// SCANNER:
		foo := &savepb.Certificate{}
		// create the non-nullable pointers
		// create variables for scan results
		scanTarget_0 := &foo.ID
		scanTarget_1 := &foo.Host
		scanTarget_2 := &foo.PemCertificate
		scanTarget_3 := &foo.PemPrivateKey
		scanTarget_4 := &foo.PemCA
		scanTarget_5 := &foo.Created
		scanTarget_6 := &foo.Expiry
		scanTarget_7 := &foo.CreatorUser
		scanTarget_8 := &foo.CreatorService
		scanTarget_9 := &foo.LastAttempt
		scanTarget_10 := &foo.LastError
		scanTarget_11 := &foo.IsLocalCA
		scanTarget_12 := &foo.IsLocalCert
		scanTarget_13 := &foo.PemPublicKey
		err := rows.Scan(scanTarget_0, scanTarget_1, scanTarget_2, scanTarget_3, scanTarget_4, scanTarget_5, scanTarget_6, scanTarget_7, scanTarget_8, scanTarget_9, scanTarget_10, scanTarget_11, scanTarget_12, scanTarget_13)
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
func (a *DBCertificate) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),host text not null ,pemcertificate text not null ,pemprivatekey text not null ,pemca text not null ,created integer not null ,expiry integer not null ,creatoruser text not null ,creatorservice text not null ,lastattempt integer not null ,lasterror text not null ,islocalca boolean not null ,islocalcert boolean not null ,pempublickey text not null );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),host text not null ,pemcertificate text not null ,pemprivatekey text not null ,pemca text not null ,created integer not null ,expiry integer not null ,creatoruser text not null ,creatorservice text not null ,lastattempt integer not null ,lasterror text not null ,islocalca boolean not null ,islocalcert boolean not null ,pempublickey text not null );`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS host text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS pemcertificate text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS pemprivatekey text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS pemca text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS created integer not null default 0;`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS expiry integer not null default 0;`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS creatoruser text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS creatorservice text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS lastattempt integer not null default 0;`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS lasterror text not null default '';`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS islocalca boolean not null default false;`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS islocalcert boolean not null default false;`,
		`ALTER TABLE ` + a.SQLTablename + ` ADD COLUMN IF NOT EXISTS pempublickey text not null default '';`,

		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS host text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS pemcertificate text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS pemprivatekey text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS pemca text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS created integer not null  default 0;`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS expiry integer not null  default 0;`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS creatoruser text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS creatorservice text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS lastattempt integer not null  default 0;`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS lasterror text not null  default '';`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS islocalca boolean not null  default false;`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS islocalcert boolean not null  default false;`,
		`ALTER TABLE ` + a.SQLTablename + `_archive  ADD COLUMN IF NOT EXISTS pempublickey text not null  default '';`,
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
func (a *DBCertificate) Error(ctx context.Context, q string, e error) error {
	if e == nil {
		return nil
	}
	return errors.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}

