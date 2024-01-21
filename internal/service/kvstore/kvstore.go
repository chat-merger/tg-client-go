package kvstore

import (
	"database/sql"
	"errors"
	"merger-adapter/internal/service/merger"
)

type Scope string
type ContMsgID string

type MessagesMap interface {
	ByMergedID(scope Scope, message merger.ID) (*ContMsgID, bool, error)
	ByContID(scope Scope, message ContMsgID) (*merger.ID, bool, error)

	Save(scope Scope, mid merger.ID, cid ContMsgID) error
}

var _ MessagesMap = (*SqliteMessagesMap)(nil)

type SqliteMessagesMap struct {
	db *sql.DB
}

func NewSqliteMessagesMap(db *sql.DB) *SqliteMessagesMap {
	return &SqliteMessagesMap{db: db}
}

func (s *SqliteMessagesMap) ByMergedID(scope Scope, message merger.ID) (*ContMsgID, bool, error) {
	stmt, err := s.db.Prepare(`
		select controller_id
		from message_map 
		where scope = ? and merger_id = ? 
	`)
	if err != nil {
		return nil, false, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(scope, message)
	var msgid ContMsgID
	err = row.Scan(&msgid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return &msgid, true, nil
}

func (s *SqliteMessagesMap) ByContID(scope Scope, message ContMsgID) (*merger.ID, bool, error) {
	stmt, err := s.db.Prepare(`
		select merger_id
		from message_map 
		where scope = ? and controller_id = ? 
	`)
	if err != nil {
		return nil, false, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(scope, message)
	var msgid merger.ID
	err = row.Scan(&msgid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return &msgid, true, nil
}

func (s *SqliteMessagesMap) Save(scope Scope, mid merger.ID, cid ContMsgID) error {
	stmt, err := s.db.Prepare(`
		insert into message_map (scope, merger_id, controller_id)
		values (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(scope, mid, cid)
	if err != nil {
		return err
	}
	return nil
}
