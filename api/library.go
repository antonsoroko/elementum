package api

import (
	"fmt"
	"strconv"

	"github.com/anacrolix/missinggo/perf"
	"github.com/gin-gonic/gin"

	"github.com/elgatito/elementum/bittorrent"
	"github.com/elgatito/elementum/config"
	"github.com/elgatito/elementum/library"
	"github.com/elgatito/elementum/library/uid"
	"github.com/elgatito/elementum/trakt"
	"github.com/elgatito/elementum/xbmc"
)

const (
	playLabel  = "LOCALIZE[30023]"
	linksLabel = "LOCALIZE[30202]"

	trueType  = "true"
	falseType = "false"

	movieType   = "movie"
	showType    = "show"
	seasonType  = "season"
	episodeType = "episode"
	searchType  = "search"

	multiType = "\nmulti"
)

// AddMovie ...
func AddMovie(ctx *gin.Context) {
	defer perf.ScopeTimer()()

	xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

	tmdbID := ctx.Params.ByName("tmdbId")
	force := ctx.DefaultQuery("force", falseType) == trueType

	movie, err := library.AddMovie(tmdbID, force)
	if err != nil {
		isErrored := true
		if err == library.ErrVideoRemoved {
			if xbmcHost.DialogConfirmFocused("Elementum", fmt.Sprintf("LOCALIZE[30279];;%s", movie.Title)) {
				movie, err = library.AddMovie(tmdbID, true)
				if err == nil {
					isErrored = false
				}
			}
		}
		if isErrored {
			ctx.String(200, err.Error())
			return
		}
	}
	if config.Get().TraktToken != "" && config.Get().TraktSyncAddedMovies {
		go trakt.SyncAddedItem("movies", tmdbID, config.Get().TraktSyncAddedMoviesLocation)
	}

	label := "LOCALIZE[30277]"
	logMsg := "%s (%s) added to library"
	if force {
		label = "LOCALIZE[30286]"
		logMsg = "%s (%s) merged to library"
	}

	log.Noticef(logMsg, movie.Title, tmdbID)
	if config.Get().LibraryUpdate == 0 || (config.Get().LibraryUpdate == 1 && xbmcHost.DialogConfirmFocused("Elementum", fmt.Sprintf("%s;;%s", label, movie.Title))) {
		xbmcHost.VideoLibraryScanDirectory(library.MoviesLibraryPath(), true)
	} else {
		if ctx != nil {
			ctx.Abort()
		}
		library.ClearPageCache(xbmcHost)
	}
}

// AddMoviesList ...
func AddMoviesList(ctx *gin.Context) {
	defer perf.ScopeTimer()()

	listID := ctx.Params.ByName("listId")
	updatingStr := ctx.DefaultQuery("updating", falseType)

	updating := false
	if updatingStr != falseType {
		updating = true
	}

	library.SyncMoviesList(listID, updating, updating)
}

// RemoveMovie ...
func RemoveMovie(ctx *gin.Context) {
	defer perf.ScopeTimer()()

	xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

	tmdbID, _ := strconv.Atoi(ctx.Params.ByName("tmdbId"))
	tmdbStr := ctx.Params.ByName("tmdbId")
	movie, paths, err := library.RemoveMovie(tmdbID, false)
	if err != nil {
		ctx.String(200, err.Error())
	}
	if config.Get().TraktToken != "" && config.Get().TraktSyncRemovedMovies {
		go trakt.SyncRemovedItem("movies", tmdbStr, config.Get().TraktSyncRemovedMoviesLocation)
	}

	if ctx != nil {
		if movie != nil && xbmcHost.DialogConfirmFocused("Elementum", fmt.Sprintf("LOCALIZE[30278];;%s", movie.Title)) {
			for _, path := range paths {
				xbmcHost.VideoLibraryCleanDirectory(path, "movies", false)
			}
			if m, err := uid.GetMovieByTMDB(movie.ID); err == nil && m != nil {
				xbmcHost.VideoLibraryRemoveMovie(m.XbmcUIDs.Kodi)
			}
		} else {
			ctx.Abort()
			library.ClearPageCache(xbmcHost)
		}
	}

}

//
// Shows externals
//

// AddShow ...
func AddShow(ctx *gin.Context) {
	defer perf.ScopeTimer()()

	xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

	tmdbID := ctx.Params.ByName("tmdbId")
	force := ctx.DefaultQuery("force", falseType) == trueType

	show, err := library.AddShow(tmdbID, force)
	if err != nil {
		isErrored := true
		if err == library.ErrVideoRemoved {
			if xbmcHost.DialogConfirmFocused("Elementum", fmt.Sprintf("LOCALIZE[30279];;%s", show.Name)) {
				show, err = library.AddShow(tmdbID, true)
				if err == nil {
					isErrored = false
				}
			}
		}
		if isErrored {
			ctx.String(200, err.Error())
			return
		}
	}
	if config.Get().TraktToken != "" && config.Get().TraktSyncAddedShows {
		go trakt.SyncAddedItem("shows", tmdbID, config.Get().TraktSyncAddedShowsLocation)
	}

	label := "LOCALIZE[30277]"
	logMsg := "%s (%s) added to library"
	if force {
		label = "LOCALIZE[30286]"
		logMsg = "%s (%s) merged to library"
	}

	log.Noticef(logMsg, show.Name, tmdbID)
	if config.Get().LibraryUpdate == 0 || (config.Get().LibraryUpdate == 1 && xbmcHost.DialogConfirmFocused("Elementum", fmt.Sprintf("%s;;%s", label, show.Name))) {
		xbmcHost.VideoLibraryScanDirectory(library.ShowsLibraryPath(), true)
	} else {
		library.ClearPageCache(xbmcHost)
	}
}

// AddShowsList ...
func AddShowsList(ctx *gin.Context) {
	defer perf.ScopeTimer()()

	listID := ctx.Params.ByName("listId")
	updatingStr := ctx.DefaultQuery("updating", falseType)

	updating := false
	if updatingStr != falseType {
		updating = true
	}

	library.SyncShowsList(listID, updating, updating)
}

// RemoveShow ...
func RemoveShow(ctx *gin.Context) {
	defer perf.ScopeTimer()()

	xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

	tmdbID := ctx.Params.ByName("tmdbId")
	show, paths, err := library.RemoveShow(tmdbID, false)
	if err != nil {
		ctx.String(200, err.Error())
	}
	if config.Get().TraktToken != "" && config.Get().TraktSyncRemovedShows {
		go trakt.SyncRemovedItem("shows", tmdbID, config.Get().TraktSyncRemovedShowsLocation)
	}

	if ctx != nil {
		if show != nil && paths != nil && xbmcHost.DialogConfirmFocused("Elementum", fmt.Sprintf("LOCALIZE[30278];;%s", show.Name)) {
			for _, path := range paths {
				xbmcHost.VideoLibraryCleanDirectory(path, "tvshows", false)
			}
			if s, err := uid.GetShowByTMDB(show.ID); err == nil && s != nil {
				xbmcHost.VideoLibraryRemoveTVShow(s.XbmcUIDs.Kodi)
			}
		} else {
			ctx.Abort()
			library.ClearPageCache(xbmcHost)
		}
	}

}

// UpdateLibrary ...
func UpdateLibrary(ctx *gin.Context) {
	xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

	if err := library.Refresh(); err != nil {
		ctx.String(200, err.Error())
	}
	if config.Get().LibraryUpdate == 0 || (config.Get().LibraryUpdate == 1 && xbmcHost.DialogConfirmFocused("Elementum", "LOCALIZE[30288]")) {
		xbmcHost.VideoLibraryScan()
	}
}

// UpdateTrakt ...
func UpdateTrakt(ctx *gin.Context) {
	xbmcHost, _ := xbmc.GetXBMCHostWithContext(ctx)

	xbmcHost.Notify("Elementum", "LOCALIZE[30358]", config.AddonIcon())
	ctx.String(200, "")
	go func() {
		library.IsTraktInitialized = false
		library.RefreshTrakt()
		if config.Get().LibraryUpdate == 0 || (config.Get().LibraryUpdate == 1 && xbmcHost.DialogConfirmFocused("Elementum", "LOCALIZE[30288]")) {
			xbmcHost.VideoLibraryScan()
		}
	}()
}

// PlayMovie ...
func PlayMovie(s *bittorrent.Service) gin.HandlerFunc {
	if config.Get().ChooseStreamAutoMovie {
		return MovieRun("play", s)
	}
	return MovieRun("links", s)
}

// PlayShow ...
func PlayShow(s *bittorrent.Service) gin.HandlerFunc {
	if config.Get().ChooseStreamAutoShow {
		return ShowEpisodeRun("play", s)
	}
	return ShowEpisodeRun("links", s)
}
