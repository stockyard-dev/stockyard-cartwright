package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Cartwright</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.rpt{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.rpt:hover{border-color:var(--leather)}
.rpt-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.rpt-title{font-size:.85rem;font-weight:700}
.rpt-desc{font-size:.7rem;color:var(--cd);margin-top:.1rem}
.rpt-query{font-size:.6rem;color:var(--cm);margin-top:.3rem;background:var(--bg);padding:.3rem .5rem;border:1px solid var(--bg3);font-family:var(--mono);word-break:break-all}
.rpt-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.rpt-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-run{border-color:var(--green);color:var(--green)}.btn-run:hover{background:var(--green);color:#fff}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:500px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> CARTWRIGHT</h1><button class="btn btn-p" onclick="openForm()">+ New Report</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search reports..." oninput="render()"></div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/reports').then(function(r){return r.json()});items=r.reports||[];renderStats();render();}
function renderStats(){var total=items.length;var sources={};items.forEach(function(r){if(r.data_source)sources[r.data_source]=true});
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Reports</div></div><div class="st"><div class="st-v">'+Object.keys(sources).length+'</div><div class="st-l">Sources</div></div><div class="st"><div class="st-v">-</div><div class="st-l">&nbsp;</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(r){return(r.title||'').toLowerCase().includes(q)||(r.data_source||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No reports configured.</div>';return;}
var h='';f.forEach(function(r){
h+='<div class="rpt"><div class="rpt-top"><div style="flex:1">';
h+='<div class="rpt-title">'+esc(r.title)+'</div>';
if(r.description)h+='<div class="rpt-desc">'+esc(r.description)+'</div>';
h+='</div><div class="rpt-actions">';
h+='<button class="btn btn-sm btn-run" onclick="runReport(''+r.id+'')">&#9654; Run</button>';
h+='<button class="btn btn-sm" onclick="openEdit(''+r.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+r.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
if(r.query)h+='<div class="rpt-query">'+esc(r.query)+'</div>';
h+='<div class="rpt-meta">';
if(r.chart_type)h+='<span class="badge">'+esc(r.chart_type)+'</span>';
if(r.data_source)h+='<span>'+esc(r.data_source)+'</span>';
if(r.schedule)h+='<span>'+esc(r.schedule)+'</span>';
if(r.last_run_at)h+='<span>Last: '+ft(r.last_run_at)+'</span>';
h+='</div></div>';});
document.getElementById('list').innerHTML=h;}
async function runReport(id){await fetch(A+'/reports/'+id+'/run',{method:'POST'}).catch(function(){});load();}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/reports/'+id,{method:'DELETE'});load();}
function formHTML(rpt){var i=rpt||{title:'',description:'',query:'',chart_type:'bar',data_source:'',schedule:''};var isEdit=!!rpt;
var h='<h2>'+(isEdit?'EDIT':'NEW')+' REPORT</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(i.title)+'"></div>';
h+='<div class="fr"><label>Description</label><input id="f-desc" value="'+esc(i.description)+'"></div>';
h+='<div class="fr"><label>Query</label><textarea id="f-query" rows="4" placeholder="SELECT ...">'+esc(i.query)+'</textarea></div>';
h+='<div class="row2"><div class="fr"><label>Chart Type</label><input id="f-chart" value="'+esc(i.chart_type)+'" placeholder="bar, line, pie"></div>';
h+='<div class="fr"><label>Data Source</label><input id="f-src" value="'+esc(i.data_source)+'"></div></div>';
h+='<div class="fr"><label>Schedule</label><input id="f-sched" value="'+esc(i.schedule)+'" placeholder="cron expression"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var r=null;for(var j=0;j<items.length;j++){if(items[j].id===id){r=items[j];break;}}if(!r)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(r);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var title=document.getElementById('f-title').value.trim();if(!title){alert('Title required');return;}
var body={title:title,description:document.getElementById('f-desc').value.trim(),query:document.getElementById('f-query').value.trim(),chart_type:document.getElementById('f-chart').value.trim(),data_source:document.getElementById('f-src').value.trim(),schedule:document.getElementById('f-sched').value.trim()};
if(editId){await fetch(A+'/reports/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/reports',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();}
function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
