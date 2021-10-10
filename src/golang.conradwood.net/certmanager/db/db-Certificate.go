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

 CREATE TABLE certificate (id integer primary key default nextval('certificate_seq'),host varchar(2000) not null,pemcertificate varchar(2000) not null,pemprivatekey varchar(2000) not null,pemca varchar(2000) not null,created integer not null,expiry integer not null,creatoruser varchar(2000) not null,creatorservice varchar(2000) not null,lastattempt integer not null);

Archive Table: (structs can be moved from main to archive using Archive() function)

 CREATE TABLE certificate_archive (id integer unique not null,host varchar(2000) not null,pemcertificate varchar(2000) not null,pemprivatekey varchar(2000) not null,pemca varchar(2000) not null,created integer not null,expiry integer not null,creatoruser varchar(2000) not null,creatorservice varchar(2000) not null,lastattempt integer not null);
*/

import (
	gosql "database/sql"
	"fmt"
	savepb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/sql"
	"golang.org/x/net/context"
)

type DBCertificate struct {
	DB *sql.DB
}

func NewDBCertificate(db *sql.DB) *DBCertificate {
	foo := DBCertificate{DB: db}
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
	_, e := a.DB.ExecContext(ctx, "insert_DBCertificate", "insert into certificate_archive (id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt) values ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10) ", p.ID, p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt)
	if e != nil {
		return e
	}

	// now delete it.
	a.DeleteByID(ctx, id)
	return nil
}

// Save (and use database default ID generation)
func (a *DBCertificate) Save(ctx context.Context, p *savepb.Certificate) (uint64, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_Save", "insert into certificate (host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id", p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt)
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
func (a *DBCertificate) SaveWithID(ctx context.Context, p *savepb.Certificate) error {
	_, e := a.DB.ExecContext(ctx, "insert_DBCertificate", "insert into certificate (id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt) values ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10) ", p.ID, p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt)
	return e
}

func (a *DBCertificate) Update(ctx context.Context, p *savepb.Certificate) error {
	_, e := a.DB.ExecContext(ctx, "DBCertificate_Update", "update certificate set host=$1, pemcertificate=$2, pemprivatekey=$3, pemca=$4, created=$5, expiry=$6, creatoruser=$7, creatorservice=$8, lastattempt=$9 where id = $10", p.Host, p.PemCertificate, p.PemPrivateKey, p.PemCA, p.Created, p.Expiry, p.CreatorUser, p.CreatorService, p.LastAttempt, p.ID)

	return e
}

// delete by id field
func (a *DBCertificate) DeleteByID(ctx context.Context, p uint64) error {
	_, e := a.DB.ExecContext(ctx, "deleteDBCertificate_ByID", "delete from certificate where id = $1", p)
	return e
}

// get it by primary id
func (a *DBCertificate) ByID(ctx context.Context, p uint64) (*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByID", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where id = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByID: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByID: error scanning (%s)", e)
	}
	if len(l) == 0 {
		return nil, fmt.Errorf("No Certificate with id %d", p)
	}
	if len(l) != 1 {
		return nil, fmt.Errorf("Multiple (%d) Certificate with id %d", len(l), p)
	}
	return l[0], nil
}

// get all rows
func (a *DBCertificate) All(ctx context.Context) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_all", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate order by id")
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

// get all "DBCertificate" rows with matching Host
func (a *DBCertificate) ByHost(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByHost", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where host = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByHost: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByHost: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemCertificate
func (a *DBCertificate) ByPemCertificate(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByPemCertificate", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where pemcertificate = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByPemCertificate: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByPemCertificate: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemPrivateKey
func (a *DBCertificate) ByPemPrivateKey(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByPemPrivateKey", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where pemprivatekey = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByPemPrivateKey: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByPemPrivateKey: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching PemCA
func (a *DBCertificate) ByPemCA(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByPemCA", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where pemca = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByPemCA: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByPemCA: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching Created
func (a *DBCertificate) ByCreated(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByCreated", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where created = $1", p)
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

// get all "DBCertificate" rows with matching Expiry
func (a *DBCertificate) ByExpiry(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByExpiry", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where expiry = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByExpiry: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByExpiry: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching CreatorUser
func (a *DBCertificate) ByCreatorUser(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByCreatorUser", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where creatoruser = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByCreatorUser: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByCreatorUser: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching CreatorService
func (a *DBCertificate) ByCreatorService(ctx context.Context, p string) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByCreatorService", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where creatorservice = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByCreatorService: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByCreatorService: error scanning (%s)", e)
	}
	return l, nil
}

// get all "DBCertificate" rows with matching LastAttempt
func (a *DBCertificate) ByLastAttempt(ctx context.Context, p uint32) ([]*savepb.Certificate, error) {
	rows, e := a.DB.QueryContext(ctx, "DBCertificate_ByLastAttempt", "select id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt from certificate where lastattempt = $1", p)
	if e != nil {
		return nil, fmt.Errorf("ByLastAttempt: error querying (%s)", e)
	}
	defer rows.Close()
	l, e := a.FromRows(ctx, rows)
	if e != nil {
		return nil, fmt.Errorf("ByLastAttempt: error scanning (%s)", e)
	}
	return l, nil
}

/**********************************************************************
* Helper to convert from an SQL Row to struct
**********************************************************************/
func (a *DBCertificate) Tablename() string {
	return "certificate"
}

func (a *DBCertificate) SelectCols() string {
	return "id,host, pemcertificate, pemprivatekey, pemca, created, expiry, creatoruser, creatorservice, lastattempt"
}

func (a *DBCertificate) FromRows(ctx context.Context, rows *gosql.Rows) ([]*savepb.Certificate, error) {
	var res []*savepb.Certificate
	for rows.Next() {
		foo := savepb.Certificate{}
		err := rows.Scan(&foo.ID, &foo.Host, &foo.PemCertificate, &foo.PemPrivateKey, &foo.PemCA, &foo.Created, &foo.Expiry, &foo.CreatorUser, &foo.CreatorService, &foo.LastAttempt)
		if err != nil {
			return nil, err
		}
		res = append(res, &foo)
	}
	return res, nil
}
