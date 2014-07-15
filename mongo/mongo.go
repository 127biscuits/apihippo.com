package mongo

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/127biscuits/apihippo.com/cdn"
)

var (
	// DB is the current mongo database connection
	DB *mgo.Database

	// Collection is the only collection that we use with mongo (hippos ATM)
	Collection *mgo.Collection

	// GridFS is were we store our images in Mongo
	GridFS *mgo.GridFS
)

// Hippo is the struct used to store the information that we save in mongo and
// we return to the user in JSON format
type Hippo struct {
	ID bson.ObjectId `bson:"_id"json:"id"`

	Filename string `json:"-"`

	URL      string `json:"url"`
	Verified bool   `json:"verified"`
	Votes    int    `json:"votes"`
}

// Populate is going to set the "calculated" fields to the struct
func (h *Hippo) Populate() {
	const NEEDED_VOTES_TO_VERIFY = 1 // TODO: move to a setting

	h.URL = cdn.GetHippoURL(h.Filename)
	h.Verified = h.Votes > NEEDED_VOTES_TO_VERIFY
}

// JSON is going to return the marshalled version of the struct
func (h Hippo) JSON() []byte {
	h.Populate()

	js, _ := json.Marshal(h)
	return js
}

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
	doc := &Hippo{
		ID:       docID,
		Filename: filename,
		Votes:    0,
	}

	if err := Collection.Insert(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
