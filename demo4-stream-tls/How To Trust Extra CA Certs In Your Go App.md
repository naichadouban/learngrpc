**Edit:** Thanks to [/u/epiris](https://www.reddit.com/r/golang/comments/6w3lcn/adding_private_ssl_ca_cert_to_trust_pool_for_only/dm77ccl/) for pointing out that I actually posted a server-side code example, not client side. I think I got distracted while finishing off the post… completely missing the point! Whoops. Code sample is fixed now.

* * *

At work this week, I was tasked with updating a couple of older internal go applications currently serving HTTP to serve HTTPS instead.

Generally, go makes this a pretty simple task and there are plenty of existing guides on the web that cover the process, so I won’t bore you with it here.

Our system is (now) modeled as follows:

```
             +---------+
             | Backend |
             +----^----+
                  |
              +---+---+
         +->  |  API  |  <-+
         |    +-------+    |
   https |                 |
         v                 |
     +---+---+             |
     | WebUI |             | https
     +---+---+             |
         ^                 |
   https |                 |
         v                 v
    +----+----+        +---+---+
    | Browser |        |  CLI  |
    +---------+        +-------+

```

*   The Backend and API always run on the same server
*   The WebUI typically runs on the same server, but can run anywhere
*   The CLI typically runs on a client workstation, but can run on the server
*   The Browser (e.g. Chrome, Firefox, etc) is not our product.

The WebUI and API share some packages, including a “settings” package which handles some basic shared config for when the applications are both running in the same environment. When running on the same host, the applications also share the same SSL certificate.
On application start, both the WebUI and API check for existence of a cert/key pair on the filesystem; If one does not exist, a self-signed CA cert is generated (using code extracted from [here](https://golang.org/src/crypto/tls/generate_cert.go)) with a Subject Common Name: `localhost`.

## The Problem With Self-Signed Certs

… is trust. By default, most clients will not trust a self-signed certificate, because they don’t recognise the signer as a trusted Root CA. The browser will warn you, but let you choose to ignore the warning and continue accessing the resource; so for our WebUI, simply serving HTTPS with the self-signed cert “out of the box” is enough - an administrator deploying and managing the product can then simply replace the certificate with one signed by a client-trusted CA, if they choose to do so.

The API is a different story, because its client is our WebUI service written in go. So when the self-signed cert is presented, we will see the well known error: `x509: certificate signed by unknown authority`.

## So What Now?

Now we needed to establish a trust between the WebUI and API. I said earlier that they can be running on the same host, so we have a number of options available, here’s some examples, in order from least to most favourable:

1.  Lazily ignore server certificates for API calls with the [InsecureSkipVerify](https://golang.org/pkg/crypto/tls/#Config.InsecureSkipVerify) option.
    *   but then why use SSL in the first place?
2.  Start an HTTP listener on a different port with local-only IP, so we don’t have to bother with certs.
    *   at least it’s obvious that the security is missing
3.  Replace the `RootCAs` in our client `tls.Config{}` with the self-signed CA cert
    *   but the WebUI might be connecting to multiple API servers, not just localhost
4.  Append the self-signed cert to the host system trust store
    *   this requires specifically ordered steps and manual intervention because we generate the cert on first start
    *   it also means all clients on the host will trust that cert
    *   not to mention that you might not have permissions on the host to pull this off
5.  *Append the self-signed cert to an in-app copy of the host system trust store*
    *   finally, what you came here for!!
6.  Deploy certs signed by a trusted CA (did someone say [Lets Encrypt](https://letsencrypt.org/)?) 

## And The Winner Is…

Since the release of [go1.7](https://golang.org/doc/go1.7#crypto_x509), the `crypto/x509` package provides a handy function called `SystemCertPool()`. This allows us to take a copy of the host system trusted CA certs, to which we can append our self-signed cert (in memory) without affecting any other clients on the host; and without removing the ability for our client to trust certs from other external resources.

Here’s some sample code:

```
package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	localCertFile = "/usr/local/internal-ca/ca.crt"
)

func main() {
	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(localCertFile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: *insecure,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	// Uses local self-signed cert
	req := http.NewRequest(http.MethodGet, "https://localhost/api/version", nil)
	resp, err := client.Do(req)
	// Handle resp and err

	// Still works with host-trusted CAs!
	req = http.NewRequest(http.MethodGet, "https://example.com/", nil)
	resp, err = client.Do(req)
	// Handle resp and err

	// ...
}

```

## Problem Solved

The method above solved our key problem: When the applications were installed together on a single host, they needed to provide a good “out of the box experience” (OOBE), but also remain secure. The WebUI can now communicate securely ‘out of the box’ with the locally-running API instance using the shared self-signed cert.

This change, along with an optional `-insecure-ssl` flag (which falls back to `InsecureSkipVerify: true`) means our applications can still be deployed simply for trials or into otherwise-secure lab environments, but can also now be used securely in public cloud or other multi-host environments, by deploying with CA-signed certificates.
