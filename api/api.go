package api

import (
	"net/http"

	"github.com/ArcticStrap/glacierglide/api/edit"
	"github.com/ArcticStrap/glacierglide/api/history"
	"github.com/ArcticStrap/glacierglide/api/search"
	"github.com/ArcticStrap/glacierglide/api/source"
	"github.com/ArcticStrap/glacierglide/api/user"
	"github.com/ArcticStrap/glacierglide/api/wiki"
	"github.com/ArcticStrap/glacierglide/api/wikido"
	"github.com/ArcticStrap/glacierglide/appsignals"
	"github.com/ArcticStrap/glacierglide/data"
)

func Setup(rt *http.ServeMux, db data.Datastore, sc *appsignals.SignalConnector) {
	// Core api routes
	wikido.SetupDoRoute(rt, db)
	wiki.SetupWikiRoute(rt, db)
	edit.SetupEditRoute(rt, db, sc)
	history.SetupHistoryRoute(rt, db)
	search.SetupSearchEngine(rt, db)
	source.SetupSourceRoute(rt, db)
	user.SetupUserRoute(rt, db, sc)
}
