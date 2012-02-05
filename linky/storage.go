package linky

import (
  "appengine"
  "appengine/datastore"
  "http"
)

type Linky struct {
  Name string
  URL string
  Location string
  Created datastore.Time
}

type Visit struct {
  LinkyName string
  VisitedDate datastore.Time
  Location string
  Referer string
}

//An interface for storage and retrieval of Linkys
type URLStorage interface {
  // Persists a given Linky
  Save(link Linky)
  // Returns a Linky for a given name
  Get(name string) Linky
  // Adds a new visit entry (Just for statistics)
  AddVisit(visit Visit)
  // Returns all the visit instances for a given linky ordered by date
  GetVisit(linkyName string) []*Visit
}

// A persistent URLStorage implementation backed by GAE DataStore

type EntityURLStorage struct {
  context appengine.Context
}

func getStorage(r *http.Request) URLStorage {
  return &EntityURLStorage{context : appengine.NewContext(r)}
}
func (self *EntityURLStorage) GetVisit(linkyName string) (visits []*Visit) {
  q := datastore.NewQuery("visit").Filter("LinkyName =", linkyName).Order("-VisitedDate")
  q.GetAll(self.context, &visits)
  return visits
}

func (self *EntityURLStorage) Save(link Linky) {
  key := datastore.NewKey(self.context, "linky", link.Name, 0, nil)
  datastore.Put(self.context, key, &link)
}

func (self *EntityURLStorage) AddVisit(visit Visit) {
  key := datastore.NewIncompleteKey(self.context, "visit", nil)
  datastore.Put(self.context, key, &visit)
}

func (self *EntityURLStorage) Get(name string) (link Linky) {
  key := datastore.NewKey(self.context, "linky", name, 0, nil)
  datastore.Get(self.context, key, &link)
  return link
}

// A Non-persistent in-memory URLStorage implementation
/*
type MapURLStorage map[string] Linky

func (self *MapURLStorage) Save(link Linky) {
  (*self)[link.Name] = link
}

func (self *MapURLStorage) Get(name string) Linky {
  return (*self)[name]
}

var storage MapURLStorage = make(map[string] Linky)

*/
