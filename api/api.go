package api

import (
	"github.com/r3boot/go-rtbh/config"
	"github.com/r3boot/go-rtbh/lists"
	"github.com/r3boot/rlib/logger"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var Config *config.Config
var Log logger.Log
var Blacklist *lists.Blacklist
var Whitelist *lists.Whitelist
var History *lists.History

// Use a custom template delimiter because it crashes with Angular
var templateDelims = []string{"{{%", "%}}"}
var templates *template.Template

func Setup(l logger.Log, cfg *config.Config) (err error) {
	Config = cfg
	Log = l

	Blacklist = lists.NewBlacklist()
	Whitelist = lists.NewWhitelist()

	template_dir := Config.Api.Resources + "/templates"
	err = filepath.Walk(template_dir, func(p string, fs os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fs.IsDir() {
			return nil
		}

		name := p[len(template_dir)-1:]
		if templates == nil {
			templates = template.New(name)
			templates.Delims(templateDelims[0], templateDelims[1])
			_, err = templates.ParseFiles(p)
		} else {
			_, err = templates.New(name).ParseFiles(p)
		}
		Log.Debug("[api]: Processed template: " + name)

		return err
	})

	return
}

func RunTillDeath() {
	img_dir := Config.Api.Resources + "uimages"
	js_dir := Config.Api.Resources + "/app/node_modules"
	app_dir := Config.Api.Resources + "/app"
	bootstrap_dist_dir := Config.Api.Resources + "/app/node_modules/bootstrap/dist"

	// First, add handlers for static assets
	http.Handle("/imgs/", http.StripPrefix("/imgs/", http.FileServer(http.Dir(img_dir))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(js_dir))))
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(app_dir))))
	http.Handle("/css/", http.FileServer(http.Dir(bootstrap_dist_dir)))
	http.Handle("/fonts/", http.FileServer(http.Dir(bootstrap_dist_dir)))

	// API functions
	http.HandleFunc("/v1/blacklist", getBlacklist)

	// Application views
	http.HandleFunc("/blacklist", blacklistHandler)
	http.HandleFunc("/", homeHandler)

	// Finally, fire up the API
	url := Config.Api.BindIp + ":" + Config.Api.BindPort
	Log.Verbose("[api] Listening on " + url)
	http.ListenAndServe(url, nil)

	return
}
