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

### subdomains

It's quite important in case that you want to access to the hippo images that
you configure the subdomains. For doing that, the easiest way is to edit your
`/etc/hosts` and add the following line:

		127.0.0.1 	cdn.apihippo.com

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

### GET a hippo image

Please, remember to configure the subdomain on your `/etc/hosts` before
accessing to this URL, if not, it's not going to work:

    curl http://cdn.apihippo.com:8000/[filename]

### GET a random hippo

Please, read the advice above about `/etc/hosts` and then:

    curl http://random.apihippo.com:8000/

TODO
----

### Important

- Check that the file is an image before storing it in our CDN.
- Return the proper time of the image on the Header (we are always returning
  image/jpeg)
- Configure a proper CDN.
- Properly test the pagination/limits with different hippos.
- Add a GET parameter to filter by verified/unverified hippos.
- Move the documentation to some place where it can be demoed: Swagger could be
  a good option.
- Change paginated output to something more HAL style.

### A lot of things to do anyway...

	grep TODO . -R
