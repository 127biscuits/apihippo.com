package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"labix.org/v2/mgo/bson"

	"github.com/127biscuits/apihippo.com/cdn"
	"github.com/127biscuits/apihippo.com/mongo"
	"github.com/gorilla/mux"
)

// PaginatedResponse is the struct used for paginated JSON responses
type PaginatedResponse struct {
	Meta struct {
		HasPrevious bool `json:"hasPrevious"`
		HasNext     bool `json:"hasNext"`
		Pages       int  `json:"pages"`
	} `json:"meta"`
	Hippos []*mongo.Hippo `json:"hippos"`
}

// GetHandler is a JSON endpoint that returns ALL the hippos paginated
func GetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: move it to a setting
	const PAGESIZE = 10

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	all := mongo.Collection.Find(nil)
	sliceAll := all.Limit(PAGESIZE)
	if page > 0 {
		sliceAll = sliceAll.Skip(PAGESIZE * (page - 1))
	}

	count, _ := all.Count()
	response := &PaginatedResponse{}

	response.Meta.Pages = count / PAGESIZE
	response.Meta.HasPrevious = page > 0
	response.Meta.HasNext = page < response.Meta.Pages

	sliceAll.All(&response.Hippos)

	// Add URLs
	for _, hippo := range response.Hippos {
		hippo.URL = cdn.GetHippoURL(hippo.ID.Hex())
	}

	js, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// GetHippoHandler is going to find a hippo by Mongo ID and return it in JSON
// format.
// In case that the hippo is not found, we are going to return a 404.
func GetHippoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	doc := &mongo.Hippo{}
	if err := mongo.Collection.FindId(bson.ObjectIdHex(id)).One(doc); err != nil {
		log.Panic(err)
		http.NotFound(w, r)
		return
	}

	doc.URL = cdn.GetHippoURL(id)

	js, _ := json.Marshal(doc)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// PostHandler is able to receive hippo image and store them in our backend.
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: check that the posted file is an image
	// TODO: check that the md5 of the upload file doesn't match with anything that we have (mongo created md5s for us)

	const MAXSIZE = 10 * 1024 // 10M

	if err := r.ParseMultipartForm(MAXSIZE); err != nil {
		http.Error(w, "Not Multipart?", http.StatusBadRequest)
	}

	// TODO: support multiple file upload, for now, we return after the first insertion
	var key string
	for key, _ = range r.MultipartForm.File {
		break
	}
	files := r.MultipartForm.File[key]

	if !strings.HasPrefix(files[0].Header["Content-Type"][0], "image/jpeg") {
		http.Error(w, "I will just accept an \"image/*\" here!", http.StatusBadRequest)
		return
	}

	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doc, err := mongo.InsertHippo(file)
	if err != nil {
		http.Error(w, "Holy s*£%t! I couldn't store your hippo!", http.StatusInternalServerError)
		return
	}

	// We don't want to store this on the DB
	doc.Verified = false
	doc.URL = cdn.GetHippoURL(doc.ID.Hex())

	js, _ := json.Marshal(doc)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return
}
