package messages_repository

import "database/sql"

type MessagesRepositorySqlite struct {
	db *sql.DB
}

func NewMessagesRepositorySqlite(db *sql.DB) *MessagesRepositorySqlite {
	return &MessagesRepositorySqlite{db: db}
}

func (r *MessagesRepositorySqlite) Get(f Filter) ([]Message, error) {
	stmt, err := r.db.Prepare(`
		with inp(
			i_merger_msg_id,i_reply_merger_msg_id, i_find_null_reply,
			i_chat_id, i_msg_id, i_sender_id,
			i_sender_first_name, i_kind, i_has_media,i_unix_sec 
		) as ( select ?,?,?,?,?,?,?,?,?,? )
		select merger_msg_id, reply_merger_msg_id,
		  chat_id, msg_id, sender_id,
		  sender_first_name, kind, has_media, unix_sec 
		from inp, messages
		where 
			(i_merger_msg_id is null or inp.i_merger_msg_id = merger_msg_id) and 
			(
			    (i_find_null_reply = 0 and i_reply_merger_msg_id is null) or 
			    (i_reply_merger_msg_id = reply_merger_msg_id)
			) and 
			(i_chat_id is null or i_chat_id = chat_id) and
			(i_msg_id is null or i_msg_id = msg_id) and
			(i_sender_id is null or i_sender_id = sender_id) and
			(i_sender_first_name is null or i_sender_first_name = sender_first_name) and
			(i_kind is null or i_kind = kind) and
			(i_has_media is null or i_has_media = has_media) and
			(i_unix_sec is null or i_unix_sec = unix_sec)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(
		f.MergerMsgId, f.ReplyMergerMsgId, f.FindNullReply,
		f.ChatId, f.MsgId, f.SenderId,
		f.SenderFirstName, f.Kind, f.HasMedia, f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]Message, 0)
	for rows.Next() {
		m := Message{}
		err = rows.Scan(
			&m.MergerMsgId, &m.ReplyMergerMsgId, &m.ChatId,
			&m.MsgId, &m.SenderId, &m.SenderFirstName,
			&m.Kind, &m.HasMedia, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessagesRepositorySqlite) Add(m Message) error {
	stmt, err := r.db.Prepare(`
		insert into messages (
		  merger_msg_id, reply_merger_msg_id,
		  chat_id, msg_id, sender_id,
		  sender_first_name, kind, has_media, unix_sec
	  	)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		m.MergerMsgId, m.ReplyMergerMsgId,
		m.ChatId, m.MsgId, m.SenderId,
		m.SenderFirstName, m.Kind, m.HasMedia, m.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
