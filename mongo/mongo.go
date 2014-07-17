package mongo

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"time"
	"crypto/md5"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/127biscuits/apihippo.com/cdn"
	"github.com/127biscuits/apihippo.com/settings"
)

var (
	// DB is the current mongo database connection
	DB *mgo.Database

	// Collection is the only collection that we use with mongo (hippos ATM)
	Collection *mgo.Collection

	// GridFS is were we store our images in Mongo
	GridFS *mgo.GridFS
)

func init() {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		log.Panic("Please, set $MONGODB_URL as \"mongodb://user:pass@host/db_name\"")
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		log.Panic("Can't connect to mongo, go error %v\n", err)
	}
	// TODO: we can not defer here, perhaps on the main?
	// defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	DB = sess.DB("apihippo")
	Collection = DB.C("hippos")
	GridFS = DB.GridFS("fs")
}

// Hippo is the struct used to store the information that we save in mongo and
// we return to the user in JSON format
type Hippo struct {
	ID bson.ObjectId `bson:"_id"json:"id"`

	Filename string `json:"-"`

	URL      string `json:"url"`
	Verified bool   `json:"verified"`
	Votes    int    `json:"votes"`

	// Weird way of getting a random doc, but:
	// http://cookbook.mongodb.org/patterns/random-attribute/
	Random float32 `json:"-"`
}

// Populate is going to set the "calculated" fields to the struct
func (h *Hippo) Populate() {
	h.URL = cdn.GetHippoURL(h.Filename)
	h.Verified = h.Votes > settings.Config.NeededVotesToVerify
}

// JSON is going to return the marshalled version of the struct
func (h Hippo) JSON() []byte {
	h.Populate()

	js, _ := json.Marshal(h)
	return js
}

// InsertHippo will store the Hippo on GridFS and return the Hippo document
// created
func InsertHippo(file multipart.File) (*Hippo, error) {
	// A random bson id as filename
	filename := bson.NewObjectId().Hex()
	gridFSImage, _ := GridFS.Create(filename)
	defer gridFSImage.Close()

	reader := bufio.NewReader(file)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := gridFSImage.Write(buf[:n]); err != nil {
			return nil, err
		}
	}

	docID := bson.NewObjectId()
	// TODO: I don't know if there is a better way to seed this
	rand.Seed(time.Now().UnixNano())

	doc := &Hippo{
		ID:			docID,
		Filename:	filename,
		Votes:		0,
		Random:		rand.Float32(),
	}

	if err := Collection.Insert(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// checks if there's an entry with the received MD5
// if there's it returns the Hippo
func getHippoByMD5(string md5checksum) (*Hippo, error){
	doc := Hippo{}
	err = GridFS.Find(bson.M{"md5": "Ale"}).One(&doc)
	if err != nil {
		panic(err)
	}
	return doc, nil
}