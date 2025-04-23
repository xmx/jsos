package jsmod

import (
	"database/sql"

	"github.com/xmx/jsos/jsvm"
)

func NewSQL() jsvm.ModuleRegister {
	return &stdSQL{}
}

type stdSQL struct {
	eng jsvm.Engineer
}

func (sq *stdSQL) RegisterModule(eng jsvm.Engineer) error {
	sq.eng = eng
	vals := map[string]any{
		"open":     sq.open,
		"register": sql.Register,
		"drivers":  sql.Drivers,
	}
	eng.RegisterModule("database/sql", vals, true)

	return nil
}

func (sq *stdSQL) open(driverName, dataSourceName string) (map[string]any, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	dri := &stdDBConn{eng: sq.eng, db: db}
	sq.eng.AddFinalizer(dri.close)
	vals := map[string]any{
		"close":           dri.close,
		"stats":           dri.stats,
		"exec":            dri.exec,
		"query":           dri.query,
		"setMaxIdleConns": dri.setMaxIdleConns,
	}

	return vals, nil
}

type stdDBConn struct {
	eng jsvm.Engineer
	db  *sql.DB
}

func (sc *stdDBConn) close() error {
	return sc.db.Close()
}

func (sc *stdDBConn) stats() sql.DBStats {
	return sc.db.Stats()
}

func (sc *stdDBConn) setMaxIdleConns(n int) {
	sc.db.SetMaxIdleConns(n)
}

func (sc *stdDBConn) query(strSQL string, args ...any) (any, error) {
	rows, err := sc.db.Query(strSQL, args...)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	size := len(cols)

	ret := make([]map[string]any, 0, 10)
	for rows.Next() {
		lines := make([]any, size)
		if err = rows.Scan(lines...); err != nil {
			return nil, err
		}
		row := make(map[string]any)
		for i, col := range cols {
			row[col] = lines[i]
		}
		ret = append(ret, row)
	}

	return ret, nil
}

func (sc *stdDBConn) exec(strSQL string, args ...any) (sql.Result, error) {
	return sc.db.Exec(strSQL, args...)
}
