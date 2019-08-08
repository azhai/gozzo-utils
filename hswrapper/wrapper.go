package hswrapper

import (
	"errors"
	"fmt"
)

const DefaultIndexName = "PRIMARY"

type HandlerSocketIndex struct {
	Socket    *HandlerSocket
	indexNo   int //1-base
	indexName string
	dbName    string
	table     string
	columns   []string
}

type HandlerSocketWrapper struct {
	Socket *HandlerSocket
	lastNo int
}

func NewWrapper(host string, rport, wport int) *HandlerSocketWrapper {
	auth := &HandlerSocketAuth{}
	auth.host = host
	auth.readPort = DefaultReadPort
	auth.writePort = DefaultWritePort
	if rport > 0 {
		auth.readPort = rport
	}
	if wport > 0 {
		auth.writePort = wport
	}
	obj := &HandlerSocketWrapper{lastNo: 0}
	obj.Socket = New()
	obj.Socket.auth = auth
	return obj
}

func (w *HandlerSocketWrapper) Close() error {
	if w.Socket.connected {
		return w.Socket.Close()
	}
	return nil
}

func (w *HandlerSocketWrapper) WrapIndex(dbName, table, indexName string, columns ...string) *HandlerSocketIndex {
	w.lastNo++
	if indexName == "" {
		indexName = DefaultIndexName
	}
	index := &HandlerSocketIndex{
		dbName: dbName, table: table, columns: columns, indexName: indexName,
	}
	index.Socket = w.Socket
	index.indexNo = w.lastNo
	index.Socket.OpenIndex(index.indexNo, dbName, table, indexName, columns...)
	return index
}

func (w *HandlerSocketIndex) FindAll(limit int, offset int, oper string, where ...string) ([]HandlerSocketRow, error) {
	rows, err := w.Socket.Find(w.indexNo, oper, limit, offset, where...)
	if err != nil {
		panic(err)
	}
	return rows, err
}

func (w *HandlerSocketIndex) FindOne(oper string, where ...string) (HandlerSocketRow, error) {
	rows, err := w.FindAll(1, 0, oper, where...)
	if rows == nil || len(rows) == 0 {
		err = errors.New("Nothing found")
		return HandlerSocketRow{}, err
	}
	return rows[0], err
}

func (w *HandlerSocketIndex) DeleteString(limit int, oper string, where []string) (int, error) {
	return w.Socket.Modify(w.indexNo, oper, limit, 0, "D", where, nil)
}

func (w *HandlerSocketIndex) Delete(limit int, oper string, where []interface{}) (int, error) {
	var conds []string
	for _, wh := range where {
		conds = append(conds, ToString(wh))
	}
	return w.DeleteString(limit, oper, conds)
}

func (w *HandlerSocketIndex) InsertString(vals ...string) error {
	return w.Socket.Insert(w.indexNo, vals...)
}

func (w *HandlerSocketIndex) Insert(vals ...interface{}) error {
	var row []string
	for _, val := range vals {
		row = append(row, ToString(val))
	}
	return w.InsertString(row...)
}

func (w *HandlerSocketIndex) UpdateString(limit int, oper string, where []string, vals ...string) (int, error) {
	return w.Socket.Modify(w.indexNo, oper, limit, 0, "U", where, vals)
}

func (w *HandlerSocketIndex) Update(limit int, oper string, where []interface{}, vals ...interface{}) (int, error) {
	var row, conds []string
	for _, val := range vals {
		row = append(row, ToString(val))
	}
	for _, wh := range where {
		conds = append(conds, ToString(wh))
	}
	return w.UpdateString(limit, oper, conds, row...)
}

func ToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
