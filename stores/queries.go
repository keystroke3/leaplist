package stores

import (
	"context"
	"database/sql"
)

var stmtsStrings = map[string]string{
	"createUser": `INSERT INTO user (id, username, display_name, passphrase) VALUES (?, ?, ?, ?)`,
	"updateUser": `UPDATE user SET  (display_name) VALUES (?) WHERE id = ?`,
	"deleteUser": `DELETE FROM user WHERE id = ?`,
	"createRelay": `INSERT INTO relay (title, alias, destination, note, station_id)
					VALUES (?, ?) RETURNING id`,
	"updateRelay": `UPDATE relay SET (title, alias, destination, note)
					VALUES (?, ?, ?, ?) WHERE id = ?`,
	"deleteRelay": `DELETE FROM tag where relay_id = ?`,

	"createTag":  `INSERT INTO tag (title, station_id) VALUES (?, ?) RETURNING id`,
	"tagRelay":   `INSERT INTO relay_tag (relay_id, tag_id, station_id) VALUES (?, ?, ?)`,
	"untagRelay": `DELETE FROM relay_tag where relay_id = ? and tag_id = ?`,
	"deleteTag":  `DELETE FROM tag where relay_id = ?`,

	"getRelaysByTag": `SELECT r.id FROM relays r
						JOIN relay_tags rt ON r.id = rt.relay_id
						JOIN tags t ON rt.tag_id = t.id
						WHERE t.label = ? AND station_id = ?`,
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Close() error
}

type Database struct {
	Db    DBTX
	stmts map[string]*sql.Stmt
}

func (s *Database) PrepareStatements(ctx context.Context) error {
	for name, value := range stmtsStrings {
		stmt, err := s.Db.PrepareContext(ctx, value)
		if err != nil {
			return err
		}
		s.stmts[name] = stmt
	}
	return nil
}

func (s *Database) Close() {
	for _, stmt := range s.stmts {
		if stmt != nil {
			stmt.Close()
		}
	}
	if s.Db != nil {
		s.Db.Close()
	}
}

func (s *Database) CreateRelay(ctx context.Context, title, alias, destination, note, station_id string) (int64, error) {
	var relayID int64
	stmt := s.stmts["createRelay"]
	err := stmt.QueryRowContext(ctx, title, alias, destination, note, station_id).Scan(&relayID)
	return relayID, err
}

func (s *Database) CreateTag(ctx context.Context, label string, stationId string) (int64, error) {
	var tagID int64
	stmt := s.stmts["createTag"]
	err := stmt.QueryRowContext(ctx, label, stationId).Scan(&tagID)
	return tagID, err
}

func (s *Database) TagRelay(ctx context.Context, relayID, tagID int64) error {
	query := `INSERT INTO relay_tags (relay_id, tag_id) VALUES (?, ?)`
	_, err := s.Db.ExecContext(ctx, query, relayID, tagID)
	return err
}

func (s *Database) UntagRelay(ctx context.Context, relayID, tagID int64) error {
	query := `DELETE FROM relay_tas where ralay_id = ? AND tag_id = ?`
	_, err := s.Db.ExecContext(ctx, query, relayID, tagID)
	return err
}

func (s *Database) GetRelaysByTag(ctx context.Context, label string, station_id string) ([]int64, error) {
	stmt := s.stmts["getRelaysByTag"]
	rows, err := stmt.QueryContext(ctx, label, station_id)
	if err != nil {
		return nil, err
	}
	var relayIDs []int64
	for rows.Next() {
		var relayID int64
		if err := rows.Scan(&relayID); err != nil {
			return nil, err
		}
		relayIDs = append(relayIDs, relayID)
	}
	return relayIDs, nil
}
