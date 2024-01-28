package server

import (
	"encoding/json"
	"github.com/pflow-dev/go-metamodel/v2/image"
	"github.com/pflow-dev/go-metamodel/v2/metamodel"
	"github.com/pflow-dev/go-metamodel/v2/model"
	"html/template"
	"net/http"
)

type HandlerWithVars = func(vars map[string]string, w http.ResponseWriter, r *http.Request)

type VarsFactory = func(r *http.Request) map[string]string

func WithVars(handler HandlerWithVars, getVarsFunc VarsFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(getVarsFunc(r), w, r)
	}
}

type BlobAccessor interface {
	Get(id int64) *model.Zblob
	GetByCid(cid string) *model.Zblob
	GetMaxId() int64
	Create(ipfsCid, base64Zipped, title, description, keywords, referrer string) (int64, error)
}
type Storage struct {
	Model   BlobAccessor
	Snippet BlobAccessor
}

type Service interface {
	IndexPage() *template.Template
	SandboxPage() *template.Template
	Event(eventType string, params map[string]interface{})
	GetState(r *http.Request) (state metamodel.Vector, ok bool)
	CheckForSnippet(hostname string, url string, referrer string) (string, bool)
	CheckForModel(hostname string, url string, referrer string) (string, bool)
}
type App struct {
	Service
	Storage
}

func (app *App) AppPage(vars map[string]string, w http.ResponseWriter, r *http.Request) {
	cid, found := app.CheckForModel(r.Host, r.URL.String(), r.Header.Get("Referer"))
	if found {
		http.Redirect(w, r, "/p/"+cid+"/", http.StatusFound)
		return
	}
	m := model.Model{}
	if vars["pflowCid"] != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		m = app.Storage.Model.GetByCid(vars["pflowCid"]).ToModel()
		if m.ID != 0 && m.IpfsCid == vars["pflowCid"] {
			m.MetaModel()
		}
	}
	_ = app.IndexPage().ExecuteTemplate(w, "index.html", m)
}

func (app *App) SvgHandler(vars map[string]string, w http.ResponseWriter, r *http.Request) {
	cid, found := app.CheckForModel(r.Host, r.URL.String(), r.Header.Get("Referer"))
	if found {
		http.Redirect(w, r, "/img/"+cid+".svg", http.StatusFound)
		return
	}
	if vars["pflowCid"] == "" {
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml ; charset=utf-8")
	m := app.Storage.Model.GetByCid(vars["pflowCid"]).ToModel()
	_, mm := m.MetaModel()
	if m.IpfsCid != vars["pflowCid"] {
		return
	}
	app.Event("viewSvg", map[string]interface{}{
		"id":      m.ID,
		"ipfsCid": m.IpfsCid,
	})
	x1, y1, width, height := mm.GetViewPort()
	i := image.NewSvg(w, width, height, x1, y1, width, height)

	state, stateOk := app.GetState(r)
	if !stateOk || len(state) != len(mm.Net().Places) {
		state = mm.Net().InitialVector()
	}
	i.Render(mm, state)
}

func (app *App) JsonHandler(vars map[string]string, w http.ResponseWriter, r *http.Request) {
	mm := metamodel.New()
	cid, found := app.CheckForModel(r.Host, r.URL.String(), r.Header.Get("Referer"))
	if found {
		http.Redirect(w, r, "/src/"+cid+".json", http.StatusFound)
	} else if vars["pflowCid"] != "" {
		m := app.Storage.Model.GetByCid(vars["pflowCid"])
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		mm.UnpackFromUrl("?z="+m.Base64Zipped, "model.json")
		app.Event("viewJson", map[string]interface{}{
			"id":      m.ID,
			"ipfsCid": m.IpfsCid,
		})
		data, _ := json.MarshalIndent(mm.ToDeclarationObject(), "", "  ")
		_, err := w.Write(data)
		if err != nil {
			panic(err)
		}
	}
}

func (app *App) SandboxHandler(vars map[string]string, w http.ResponseWriter, r *http.Request) {
	cid, found := app.CheckForSnippet(r.Host, r.URL.String(), r.Header.Get("Referer"))
	if found {
		http.Redirect(w, r, "/sandbox/"+cid+"/", http.StatusFound)
	} else if vars["pflowCid"] != "" {
		rec := app.Storage.Snippet.GetByCid(vars["pflowCid"])
		sourceCode, ok := metamodel.UnzipUrl("?z="+rec.Base64Zipped, "declaration.js")
		if !ok {
			http.Error(w, "Failed to unzip snippet", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templateData := struct {
			IpfsCid    string
			SourceCode string
		}{
			IpfsCid:    vars["pflowCid"],
			SourceCode: sourceCode,
		}
		_ = app.SandboxPage().ExecuteTemplate(w, "index.html", templateData)
	}
}
