package mongo

import (
	"log"
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
	ID bson.ObjectId `bson:"_id"`

	File *mgo.GridFile `json:"-"`

	URL      string `json:"url"`
	Verified bool   `json:"verified"`
	Votes    int    `json:"votes"`
}

func init() {
	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		log.Panic("Please, set $MONGOHQ_URL as \"mongodb://user:pass@server.mongohq.com/db_name\"")
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
