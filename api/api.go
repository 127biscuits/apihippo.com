package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"

	"github.com/127biscuits/apihippo.com/cdn"
	"github.com/127biscuits/apihippo.com/mongo"
	"github.com/gorilla/mux"
)

// GetHandler is a JSON endpoint that returns ALL the hippos paginated
func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO: return all the hippos paginated")
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

	// A random bson id as filename
	gridFSImage, _ := mongo.DB.GridFS("fs").Create(bson.NewObjectId().Hex())
	defer gridFSImage.Close()

	io.Copy(gridFSImage, r.Body)

	docID := bson.NewObjectId()
	doc := mongo.Hippo{
		ID:    docID,
		File:  gridFSImage,
		Votes: 0,
	}
	if err := mongo.Collection.Insert(doc); err != nil {
		http.Error(w, "Holy s*£%t! I couldn't store your hippo!", http.StatusInternalServerError)
		return
	}

	// We don't want to store this on the DB
	doc.Verified = false
	doc.URL = cdn.GetHippoURL(docID.Hex())

	js, _ := json.Marshal(doc)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
