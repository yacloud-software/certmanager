package db

/*
 This file was created by mkdb-client.
 The intention is not to modify thils file, but you may extend the struct DBCertificate
 in a seperate file (so that you can regenerate this one from time to time)
*/

/*
 PRIMARY KEY: ID
*/

/*
 postgres:
 create sequence certificate_seq;

Main Table:

 CREATE TABLE certificate (id integer primary key default nextval('certificate_seq'),host text not null  ,pemcertificate text not null  ,pemprivatekey text not null  ,pemca text not null  ,created integer not null  ,expiry integer not null  ,creatoruser text not null  ,creatorservice text not null  ,lastattempt integer not null  ,lasterror text not null  );

Alter statements:
ALTER TABLE certificate ADD COLUMN host text not null default '';
ALTER TABLE certificate ADD COLUMN pemcertificate text not null default '';
ALTER TABLE certificate ADD COLUMN pemprivatekey text not null default '';
ALTER TABLE certificate ADD COLUMN pemca text not null default '';
ALTER TABLE certificate ADD COLUMN created integer not null default 0;
ALTER TABLE certificate ADD COLUMN expiry integer not null default 0;
ALTER TABLE certificate ADD COLUMN creatoruser text not null default '';
ALTER TABLE certificate ADD COLUMN creatorservice text not null default '';
ALTER TABLE certificate ADD COLUMN lastattempt integer not null default 0;
ALTER TABLE certificate ADD COLUMN lasterror text not null default '';


Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE certificate_archive (id integer unique not null,host text not null,pemcertificate text not null,pemprivatekey text not null,pemca text not null,created integer not null,expiry integer not null,creatoruser text not null,creatorservice text not null,lastattempt integer not null,lasterror text not null);
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
	default_def_DBCertificate *DBCertificate
)

type DBCertificate struct {
	DB                  *sql.DB
	SQLTablename        string
	SQLArchivetablename string
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

// archive. It is NOT transactionally save.
func (a *DBCertificate) Archive(ctx context.Context, id uint64) error {

	// load it
	p, err := a.ByID(ctx, id)
	if err != nil {
		return err
	}

	// now save it to archive:
	_, e := a.DB.ExecContext(ctx, "archive_DBCertificate", "insert into "+a.SQLArchivetablename+" (id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror) values ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ", p.ID, p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt, p.LastError)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBCertificate) Save(ctx context.Context, p *savepb.Certificate) (uint64, error) {
	qn := "DBCertificate_Save"
	rows, e := a.DB.QueryContext(ctx, qn, "insert into "+a.SQLTablename+" (host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id", p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt, p.LastError)
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
func (a *DBCertificate) SaveWithID(ctx context.Context, p *savepb.Certificate) error {
	qn := "insert_DBCertificate"
	_, e := a.DB.ExecContext(ctx, qn, "insert into "+a.SQLTablename+" (id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror) values ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ", p.ID, p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt, p.LastError)
	return a.Error(ctx, qn, e)
}

func (a *DBCertificate) Update(ctx context.Context, p *savepb.Certificate) error {
	qn := "DBCertificate_Update"
	_, e := a.DB.ExecContext(ctx, qn, "update "+a.SQLTablename+" set host=$1, pemcertificate=$2, pemprivatekey=$3, pemca=$4, created=$5, expiry=$6, creatoruser=$7, creatorservice=$8, lastattempt=$9, lasterror=$10 where id = $11", p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt, p.LastError, p.ID)

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
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where id = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByID: error scanning (%s)", e))
	}
	if len(l) == 0 {
		return nil, a.Error(ctx, qn, fmt.Errorf("No Certificate with id %v", p))
	}
	if len(l) != 1 {
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) Certificate with id %v", len(l), p))
	}
	return l[0], nil
}

// get it by primary id (nil if no such ID row, but no error either)
func (a *DBCertificate) TryByID(ctx context.Context, p uint64) (*savepb.Certificate, error) {
	qn := "DBCertificate_TryByID"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where id = $1", p)
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
		return nil, a.Error(ctx, qn, fmt.Errorf("Multiple (%d) Certificate with id %v", len(l), p))
	}
	return l[0], nil
}

// get all rows
func (a *DBCertificate) All(ctx context.Context) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_all"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" order by id")
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

// get all "DBCertificate" rows with matching Host
func (a *DBCertificate) ByHost(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByHost"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where host = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByHost: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByHost: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeHost(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeHost"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where host ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByHost: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByHost: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemCertificate
func (a *DBCertificate) ByPemCertificate(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemCertificate"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where pemcertificate = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCertificate: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCertificate: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemCertificate(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemCertificate"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where pemcertificate ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCertificate: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCertificate: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemPrivateKey
func (a *DBCertificate) ByPemPrivateKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemPrivateKey"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where pemprivatekey = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemPrivateKey: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemPrivateKey: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemPrivateKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemPrivateKey"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where pemprivatekey ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemPrivateKey: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemPrivateKey: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemCA
func (a *DBCertificate) ByPemCA(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByPemCA"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where pemca = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCA: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCA: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikePemCA(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikePemCA"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where pemca ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCA: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByPemCA: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching Created
func (a *DBCertificate) ByCreated(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreated"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where created = $1", p)
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
func (a *DBCertificate) ByLikeCreated(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeCreated"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where created ilike $1", p)
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

// get all "DBCertificate" rows with matching Expiry
func (a *DBCertificate) ByExpiry(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByExpiry"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where expiry = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExpiry: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExpiry: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeExpiry(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeExpiry"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where expiry ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExpiry: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByExpiry: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching CreatorUser
func (a *DBCertificate) ByCreatorUser(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreatorUser"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where creatoruser = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorUser: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorUser: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeCreatorUser(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeCreatorUser"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where creatoruser ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorUser: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorUser: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching CreatorService
func (a *DBCertificate) ByCreatorService(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByCreatorService"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where creatorservice = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorService: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorService: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeCreatorService(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeCreatorService"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where creatorservice ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorService: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByCreatorService: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching LastAttempt
func (a *DBCertificate) ByLastAttempt(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLastAttempt"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where lastattempt = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastAttempt: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastAttempt: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeLastAttempt(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeLastAttempt"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where lastattempt ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastAttempt: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastAttempt: error scanning (%s)", e))
	}
	return l, nil
}

// get all "DBCertificate" rows with matching LastError
func (a *DBCertificate) ByLastError(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLastError"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where lasterror = $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastError: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastError: error scanning (%s)", e))
	}
	return l, nil
}

// the 'like' lookup
func (a *DBCertificate) ByLikeLastError(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	qn := "DBCertificate_ByLikeLastError"
	rows, e := a.DB.QueryContext(ctx, qn, "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror from "+a.SQLTablename+" where lasterror ilike $1", p)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastError: error querying (%s)", e))
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, a.Error(ctx, qn, fmt.Errorf("ByLastError: error scanning (%s)", e))
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Query
**********************************************************************/

// from a query snippet (the part after WHERE)
func (a *DBCertificate) FromQuery(ctx context.Context, query_where string, args ...interface{}) ([]*savepb.Certificate, error) {
	rows, err := a.DB.QueryContext(ctx, "custom_query_"+a.Tablename(), "select "+a.SelectCols()+" from "+a.Tablename()+" where "+query_where, args...)
	if err != nil {
		return nil, err
	}
	return a.FromRows(ctx, rows)
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBCertificate) Tablename() string {
	return a.SQLTablename
}

func (a *DBCertificate) SelectCols() string {
	return "id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt, lasterror"
}
func (a *DBCertificate) SelectColsQualified() string {
	return "" + a.SQLTablename + ".id," + a.SQLTablename + ".host, " + a.SQLTablename + ".pemcertificate, " + a.SQLTablename + ".pemprivatekey, " + a.SQLTablename + ".pemca, " + a.SQLTablename + ".created, " + a.SQLTablename + ".expiry, " + a.SQLTablename + ".creatoruser, " + a.SQLTablename + ".creatorservice, " + a.SQLTablename + ".lastattempt, " + a.SQLTablename + ".lasterror"
}

func (a *DBCertificate) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.Certificate, error) {
	var res []*savepb.Certificate
	for rows.Next() {
		foo := savepb.Certificate{}
		err := rows.Scan(&foo.ID, &foo.Host, &foo.PemCertificate, &foo.PemPrivateKey, &foo.PemCA, &foo.Created, &foo.Expiry, &foo.CreatorUser, &foo.CreatorService, &foo.LastAttempt, &foo.LastError)
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
func (a *DBCertificate) CreateTable(ctx context.Context) error {
	csql := []string{
		`create sequence if not exists ` + a.SQLTablename + `_seq;`,
		`CREATE TABLE if not exists ` + a.SQLTablename + ` (id integer primary key default nextval('` + a.SQLTablename + `_seq'),host text not null  ,pemcertificate text not null  ,pemprivatekey text not null  ,pemca text not null  ,created integer not null  ,expiry integer not null  ,creatoruser text not null  ,creatorservice text not null  ,lastattempt integer not null  ,lasterror text not null  );`,
		`CREATE TABLE if not exists ` + a.SQLTablename + `_archive (id integer primary key default nextval('` + a.SQLTablename + `_seq'),host text not null  ,pemcertificate text not null  ,pemprivatekey text not null  ,pemca text not null  ,created integer not null  ,expiry integer not null  ,creatoruser text not null  ,creatorservice text not null  ,lastattempt integer not null  ,lasterror text not null  );`,
	}
	for i, c := range csql {
		_, e := a.DB.ExecContext(ctx, fmt.Sprintf("create_"+a.SQLTablename+"_%d", i), c)
		if e != nil {
			return e
		}
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
	return fmt.Errorf("[table="+a.SQLTablename+", query=%s] Error: %s", q, e)
}
