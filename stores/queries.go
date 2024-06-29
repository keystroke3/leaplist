package stores

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

var stmtsStrings = map[string]string{
	"createUser":    `INSERT INTO user (id, username, display_name, passphrase) VALUES (?, ?, ?, ?)`,
	"getUser":       `SELECT id, username, display_name, passphrase from user WHERE id=? OR username=?`,
	"updateUser":    `UPDATE user SET  (display_name) VALUES (?) WHERE id = ?`,
	"deleteUser":    `DELETE FROM user WHERE id = ?`,
	"createStation": `INSERT INTO station (id, user_id) VALUES (?, ?)`,
	"deleteStation": `DELETE FROM station WHERE id = ?`,
	"createRelay": `INSERT INTO relay (title, alias, destination, note, station_id)
					VALUES (?, ?) RETURNING id`,
	"updateRelay": `UPDATE relay SET (title, alias, destination, note)
					VALUES (?, ?, ?, ?) WHERE id = ?`,
	"deleteRelay": `DELETE FROM tag where relay_id = ?`,

	"createTag":  `INSERT INTO tag (title, station_id) VALUES (?, ?) RETURNING id`,
	"tagRelay":   `INSERT INTO relay_tag (relay_id, tag_id, station_id) VALUES (?, ?, ?)`,
	"untagRelay": `DELETE FROM relay_tag where relay_id = ? and tag_id = ?`,
	"deleteTag":  `DELETE FROM tag where relay_id = ?`,

	"getRelaysByTag": `SELECT r.id, r.title, r.alias, r.description, r.note FROM relays r
						JOIN relay_tags rt ON r.id = rt.relay_id
						JOIN tags t ON rt.tag_id = t.id
						WHERE t.label = ? AND station_id = ?`,
	"getRelayById":     `SELECT r.id, r.title, r.alias, r.description, r.note FROM relays r WHERE id = ?`,
	"getStationRelays": `SELECT r.id, r.title, r.alias, r.description, r.note FROM relays r WHERE station_id = ?`,
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

func (s *Database) CreateUser(ctx context.Context, username, display_name, passphrase string) error {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	stmt := s.stmts["createRelay"]

	_, err := stmt.ExecContext(ctx, username, display_name, passwordHash)
	return err
}
func (s *Database) CreateStation(ctx context.Context, id, user_id string) error {
	stmt := s.stmts["createStation"]
	_, err := stmt.ExecContext(ctx, id, user_id)
	return err
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
	stmt := s.stmts["tagRelay"]
	_, err := stmt.ExecContext(ctx, relayID, tagID)
	return err
}

func (s *Database) UntagRelay(ctx context.Context, relayID, tagID int64) error {
	stmt := s.stmts["untagRelay"]
	_, err := stmt.ExecContext(ctx, relayID, tagID)
	return err
}

func (s *Database) GetRelaysByTag(ctx context.Context, label string, station_id string) ([]Relay, error) {
	stmt := s.stmts["getRelaysByTag"]
	rows, err := stmt.QueryContext(ctx, label, station_id)
	if err != nil {
		return nil, err
	}
	var relays []Relay
	for rows.Next() {
		r := Relay{}
		if err := rows.Scan(&r.Id, &r.Title, &r.Alias, &r.Note); err != nil {
			return nil, err
		}
		relays = append(relays, r)
	}
	return relays, nil
}

func (s *Database) GetStationRelays(ctx context.Context, station_id string) ([]Relay, error) {
	stmt := s.stmts["getStationRelays"]
	rows, err := stmt.QueryContext(ctx, station_id)
	if err != nil {
		return nil, err
	}
	var relays []Relay
	for rows.Next() {
		r := Relay{}
		if err := rows.Scan(&r.Id, &r.Title, &r.Alias, &r.Note); err != nil {
			return nil, err
		}
		relays = append(relays, r)
	}
	return relays, nil
}

func (s *Database) GetRelayByID(ctx context.Context, id string) (Relay, error) {
	stmt := s.stmts["getRelayById"]
	var r = Relay{}
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return r, err
	}
	if err := rows.Scan(&r.Id, &r.Title, &r.Alias, &r.Note); err != nil {
		return r, err
	}
	return r, nil
}
