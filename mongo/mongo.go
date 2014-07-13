package mongo

import (
	"bufio"
	"io"
	"log"
	"mime/multipart"
	"os"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	// DB is the current mongo database connection
	DB *mgo.Database

	// Collection is the only collection that we use with mongo (hippos ATM)
	Collection *mgo.Collection
)

// Hippo is the struct used to store the information that we save in mongo and
// we return to the user in JSON format
type Hippo struct {
	ID bson.ObjectId `bson:"_id"json:"id"`

	File *mgo.GridFile `json:"-"`

	URL      string `json:"url"`
	Verified bool   `json:"verified"`
	Votes    int    `json:"votes"`
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
}

// InsertHippo will store the Hippo on GridFS and return the Hippo document
// created
func InsertHippo(file multipart.File) (*Hippo, error) {
	// A random bson id as filename
	gridFSImage, _ := DB.GridFS("fs").Create(bson.NewObjectId().Hex())
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
		ID:    docID,
		File:  gridFSImage,
		Votes: 0,
	}

	if err := Collection.Insert(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
