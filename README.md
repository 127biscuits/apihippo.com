apihippo.com
============

- Do you like hippos?
+ Who doesn't like hippos?

And this is how it all started.

How to run it?
--------------

If you are brave enough you can try it, but (advice here) it's not finished
yet:

	MONGOHQ_URL=mongodb://user:pass@kahana.mongohq.com:10009/apihippo \
	    go run main.go

Doc?
----

POST a hippo
~~~~~~~~~~~~

	curl -H "Accept: application/json; \
                 Content-Type: multipart/form-data" \
            http://localhost:8000/ --form "data=@/tmp/hippo.jpg"

TODO
----

grep TODO . -R
