apihippo.com
============

- Do you like hippos?
+ Who doesn't like hippos?

And this is how it all started.

How to run it?
--------------

If you are brave enough you can try it, but (advice here) it's not finished
yet:

We're using MongoDB as our database. If you don't have a MongoDB database at hand, go
go MongoHQ and create a database there. There's a free layer and it will save all the fuss
of installing and configuring it.

If you go for a Database in MongoHQ, your connection string will be something like this:

	MONGODB_URL=mongodb://user:pass@kahana.mongohq.com:10009/apihippo \
	    go run main.go


Otherwise, if you're using a local install, you can run it as:

	MONGODB_URL=mongodb://localhost:27017/apihippo \
	    go run main.go

Doc?
----

### POST a hippo

	curl -H "Accept: application/json; \
                 Content-Type: multipart/form-data" \
            http://localhost:8000/ --form "data=@/tmp/hippo.jpg"

### PUT a hippo (vote)

	curl -H "Accept: application/json" -X PUT \
	    http://localhost:8000/[id_returned_by_the_previous_post]

### GET a hippo

	curl -H "Accept: application/json" \
	    http://localhost:8000/[id_returned_by_the_previous_post]

### GET all the hippos paginated

	curl -H "Accept: application/json" \
	    http://localhost:8000/?page=N  # page is optional

### GET the web

	curl http://localhost:8000

Concerns
--------

If you have read until here you are interested, so, I am not pretty sure if I
should add the URLs of the resources on the response, something like:

	"_links": {
		"self": {
			"href": "/123asd123"
		}
	},
	"hippo: {
		...
	}

Or in the case of a paginated response:

	"_links: {
		"self": { "href": "/" }
		"next": { "next": "/?page=1" }
		"previous": nil
	}
	"hippos": {
		...
	}

TODO
----

### Important

- Add a CDN or at least a way to serve the pictures.
- Properly test the pagination/limits with different hippos.
- Add a GET parameter to filter by verified/unverified hippos.
- Update verified at runtime (do not save the value on the DB, just have a
  tresshold on the code).

### A lot of things to do anyway...

	grep TODO . -R
