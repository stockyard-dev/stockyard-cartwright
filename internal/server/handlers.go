package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-cartwright/internal/store")
func(s *Server)handleListDashboards(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListDashboards();if list==nil{list=[]store.Dashboard{}};writeJSON(w,200,list)}
func(s *Server)handleCreateDashboard(w http.ResponseWriter,r *http.Request){var d store.Dashboard;json.NewDecoder(r.Body).Decode(&d);if d.Name==""{writeError(w,400,"name required");return};s.db.CreateDashboard(&d);writeJSON(w,201,d)}
func(s *Server)handleDeleteDashboard(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteDashboard(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListCharts(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListCharts(id);if list==nil{list=[]store.Chart{}};writeJSON(w,200,list)}
func(s *Server)handleAddChart(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var c store.Chart;json.NewDecoder(r.Body).Decode(&c);c.DashboardID=id;if c.Title==""{writeError(w,400,"title required");return};if c.ChartType==""{c.ChartType="bar"};if c.Config==""{c.Config="{}"};s.db.AddChart(&c);writeJSON(w,201,c)}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
