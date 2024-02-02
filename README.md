# openshift-tests-endpoint

This endpoint can serve as a simple endpoint on a cloud so we can test connectivity to
endpoints in specific clouds during openshift-test runs.  There are two versions:

* Run on a VM
* Run in a [cloud function](https://cloud.google.com/functions/docs/writing/write-http-functions)

## Building

```bash
make clean
make build
```

## Run as a VM

### Run the server

```go
$ ./openshift-tests-endpoint -mode server

# Get help
$ ./openshift-tests-endpoint -h
Usage of main:
  -cert string
    	TLS certificate file (default "./cert.pem")
  -count int
    	number of clients to run in client mode (default 1)
  -https
    	use HTTPS (default HTTP)
  -key string
    	TLS private key file (default "./key.pem")
  -mode string
    	run in 'client' or 'server' mode
```

### Run the (test) client

Start 20 clients:

```go
$ ./openshift-tests-endpoint -mode client -count 20
```

or single client with `curl`:

```bash
curl -X GET http://localhost:49888/health -H "Audit-ID: 12345" -H "Cluster-ID: build02"
```


## Run as a cloud function

See [cloud_function](cloud_function/cloud_function.go).

See [deployment doc](https://cloud.google.com/functions/docs/deploy).

```bash
$ gcloud functions deploy OpenshiftTestsEndpoint --runtime go121 --trigger-http --allow-unauthenticated --entry-point OpenshiftTestsEndpoint
```

### Test the cloud function like this

Once you determine your endpoint DNS entry, you can test the server like this:

```
$ url=<get this from your cloud>
$ echo $(curl -sk -w "%{http_code}" -o response.txt -H "Audit-ID: 12345" -H "Cluster-ID: build03" "$url")
200
```