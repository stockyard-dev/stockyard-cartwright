package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Report struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Query string `json:"query"`
	ChartType string `json:"chart_type"`
	DataSource string `json:"data_source"`
	Schedule string `json:"schedule"`
	LastRunAt string `json:"last_run_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"cartwright.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS reports(id TEXT PRIMARY KEY,title TEXT NOT NULL,description TEXT DEFAULT '',query TEXT DEFAULT '',chart_type TEXT DEFAULT 'table',data_source TEXT DEFAULT '',schedule TEXT DEFAULT '',last_run_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Report)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO reports(id,title,description,query,chart_type,data_source,schedule,last_run_at,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Description,e.Query,e.ChartType,e.DataSource,e.Schedule,e.LastRunAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Report{var e Report;if d.db.QueryRow(`SELECT id,title,description,query,chart_type,data_source,schedule,last_run_at,created_at FROM reports WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Description,&e.Query,&e.ChartType,&e.DataSource,&e.Schedule,&e.LastRunAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Report{rows,_:=d.db.Query(`SELECT id,title,description,query,chart_type,data_source,schedule,last_run_at,created_at FROM reports ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Report;for rows.Next(){var e Report;rows.Scan(&e.ID,&e.Title,&e.Description,&e.Query,&e.ChartType,&e.DataSource,&e.Schedule,&e.LastRunAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Report)error{_,err:=d.db.Exec(`UPDATE reports SET title=?,description=?,query=?,chart_type=?,data_source=?,schedule=?,last_run_at=? WHERE id=?`,e.Title,e.Description,e.Query,e.ChartType,e.DataSource,e.Schedule,e.LastRunAt,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM reports WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM reports`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Report{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (title LIKE ? OR description LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    rows,_:=d.db.Query(`SELECT id,title,description,query,chart_type,data_source,schedule,last_run_at,created_at FROM reports WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Report;for rows.Next(){var e Report;rows.Scan(&e.ID,&e.Title,&e.Description,&e.Query,&e.ChartType,&e.DataSource,&e.Schedule,&e.LastRunAt,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
