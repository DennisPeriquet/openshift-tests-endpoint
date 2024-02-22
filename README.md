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
curl -X GET http://localhost:80/health -H "Audit-ID: 12345"
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
$ echo $(curl -sk -w "%{http_code}" -o response.txt -H "Audit-ID: 12345" "$url")
200
```

## Run as a container

Examples for quay.io and docker.io:

I've been using these and pushing to docker.io instead of quay.io:

```bash
QUAY_USER=`whoami`
podman build -t quay.io/${QUAY_USER}/openshift_tests_endpoint:$(cat VER) .
podman tag quay.io/${QUAY_USER}/openshift_tests_endpoint:$(cat VER) ${QUAY_USER}/openshift_tests_endpoint:$(cat VER)
podman push  ${QUAY_USER}/openshift_tests_endpoint:$(cat VER)
```

But here are examples for quay.io:

```bash
podman push quay.io/${QUAY_USER}/openshift_tests_endpoint:$(cat VER)
podman run --name openshift_endpoint_server -d -p 80:80 quay.io/${QUAY_USER}/openshift_tests_endpoint:$(cat VER)
podman run --name openshift_endpoint_server -d -p 80:80 ${QUAY_USER}/openshift_tests_endpoint:$(cat VER)
podman logs -f openshift_endpoint_server
```

Using the [Makefile](./Makefile):

```bash
make clean
make build
make podman-build
make podman-tag
make podman-push
```