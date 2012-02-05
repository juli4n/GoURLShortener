// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package linky

import (
	"fmt"
	"http"
  "template"
  "appengine/datastore"
  "time"
  "strings"
)

const (
  MainPage = "main.html"
  StatsPage = "stats.html"
)

type PageModel struct {
  LinkyName string
  Visits map[int64] int64
  ClickMap map[string] uint
  FooterMessage string
  LinkyURL string
  StatsURL string
}

func init() {
  http.HandleFunc("/add/", addLinkHandler)
  http.HandleFunc("/stats/", statsHandler)
	http.HandleFunc("/", mainHandler)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

  if r.URL.Path == "/" {
    writeTemplate(w, MainPage, &PageModel{})
  } else {
    linky := getStorage(r).Get(r.URL.Path[1:len(r.URL.Path)])
    if linky.URL == "" {
      writeTemplate(w, MainPage, &PageModel{})
    } else {
      location := r.Header.Get("X-AppEngine-country")
      referer := r.Header.Get("Referer")
      getStorage(r).AddVisit(Visit{LinkyName: linky.Name, Location: location, VisitedDate: datastore.SecondsToTime(time.Seconds()), Referer: referer })
      w.Header().Add("Location", linky.URL)
      w.WriteHeader(301)
    }
  }
}

func addLinkHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  name := r.FormValue("name")
  url := r.FormValue("url")
  if ! strings.HasPrefix(url, "http://") {
    url = "http://" + url
  }
  location := r.Header.Get("X-AppEngine-country")
  getStorage(r).Save( Linky { Name: name, URL: url, Location: location, Created: datastore.SecondsToTime(time.Seconds())} )
  writeTemplate(w, MainPage, &PageModel { FooterMessage: "Link added", LinkyURL: r.Host + "/" + name, StatsURL: r.Host + "/stats/" + name, LinkyName: name})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
  linkyName := r.URL.Path[7: len(r.URL.Path)]
  visits := getStorage(r).GetVisit(linkyName)
  vd := make(map[int64] int64)
  clicksMap := make(map[string] uint)
  for v := range visits {
    et := visits[v].VisitedDate.Time().Seconds()
    ret := et - (et % 86400) + 86400
    vd[ret] += 1
    if visits[v].Location != "" {
      clicksMap[visits[v].Location] += 1
    }
  }
  fmt.Print(vd)
  fmt.Print(clicksMap)
  writeTemplate(w, StatsPage, &PageModel{LinkyName: linkyName, Visits: vd, ClickMap : clicksMap})
}

func writeTemplate(w http.ResponseWriter, tempFile string, model interface{}) {
    if tpl, err := template.ParseFile(tempFile); err != nil {
      fmt.Fprintf(w, "Internal Server Error")
      w.WriteHeader(500)
    } else {
      tpl.Execute(w, model)
    }
}
