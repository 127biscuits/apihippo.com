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

TODO
----

	grep TODO . -R
