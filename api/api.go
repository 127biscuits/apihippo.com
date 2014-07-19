package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

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
		hippo.Populate()
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
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(doc.JSON())
}

// PutHippoHandler is going to increment the number of votes for a cerating
// hippo
func PutHippoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	change := bson.M{"$inc": bson.M{"votes": 1}}
	err := mongo.Collection.UpdateId(bson.ObjectIdHex(id), change)
	switch {
	case err == mgo.ErrNotFound:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	GetHippoHandler(w, r)
}

// PostHandler is able to receive hippo image and store them in our backend.
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: check that the posted file is an image

	const MAXSIZE = 32 << 10 // 32M

	if err := r.ParseMultipartForm(MAXSIZE); err != nil {
		errMessage := fmt.Sprintf(
			"Have you added the Content-Type: multipart/form-data header?"+
				"This is the detailed error: %s", err.Error())
		http.Error(w, errMessage, http.StatusBadRequest)
		return
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

	md5 := md5.New()
	io.Copy(md5, file)
	chksum := md5.Sum(nil)

	doc, err := mongo.GetHippoByMD5(string(chksum[:]))

	if doc == nil {
		doc, err := mongo.InsertHippo(file)
		if err != nil {
			http.Error(w, "Holy s*£%t! I couldn't store your hippo!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(doc.JSON())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(doc.JSON())
	}

}

// FakeCDNHandler will return the image stream for the hippo.
// TODO: this is just temporal until we have a proper CDN.
func FakeCDNHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["id"]

	w.Header().Set("Content-Type", "image/jpeg") // TODO: check the type of the image before adding this header

	file, err := mongo.GridFS.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	image := make([]byte, file.Size())
	file.Read(image)
	w.Write(image)
}

func RandomHippoHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure index on Random if we want efficience
	// TODO: not pretty sure if I should do this always
	err := mongo.Collection.EnsureIndexKey("random")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// TODO: wow, so much seed, wow, so random
	rand.Seed(time.Now().UnixNano())

	random := rand.Float32()
	hippo := &mongo.Hippo{}

	// TODO: just get the verified ones
	hippoQuerySet := mongo.Collection.Find(
		bson.M{
			"random": bson.M{"$gte": random}})
	if hippo == nil {
		hippoQuerySet = mongo.Collection.Find(
			bson.M{
				"random": bson.M{"$lte": random}})
	}
	hippoQuerySet.One(hippo)

	w.Header().Set("Content-Type", "application/json")
	w.Write(hippo.JSON())
}
