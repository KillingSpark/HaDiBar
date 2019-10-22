# HaDiBar #
The goal is to manage the accountings for our beverages on our floor of a dorm conveniently

This is a Server for the HaDiBar-API. It exposes functionality to manage accounts and beverages.

The definiton of api can be found at [Hadibar-API](https://github.com/killingspark/Hadibar-api)
a reference implementation of a webapp that uses said API can be found at [Hadibar-Webapp](https://github.com/killingspark/Hadibar-Webapp).
The Webapp also contains some conveniance functions capsulating the ajax-calls and session management to use the api from javascript 

Webapp made with Vuejs, JQuery and Bootstrap
Accounts/Beverages/Users are stored in bolt key-value stores


## General usecase
This is meant for one person in a group of people to manage the accounts of all of them (in most cases that will be the one that manages the physical beverages as well). It is more of a convenient way of book keeping, not a way so that everyone is managing their own accounts.

## Users ##
There is no explicit user management right now. Usernames are aquired by logging in with the name and the password for the first time.
This will hopefully be improved in the future (with password-resets/registering with an email etc,etc)

## Test the server without the webapp ##
Make calls with curl like those:

Testlogin: 
```
    export SES="$(curl -X GET 127.0.0.1:8080/api/session/getid)"
    curl -X POST 127.0.0.1:8080/api/session/login -H "sessionID: "$SES --form "name=Moritz" --form "password=test"
```

## Admin Server
Besides the rest-api for the normal users there is an additional admin-server. It can listen on either
1. A unix socket
2. A tcp socket
3. A tls socket (requires keys provided by you)

The admin-server uses a JSON based RPC that allows to manage all things. There is support for 
1. listing/removing users
1. perform cleanup after deleting users
1. listing accounts
1. listing beverages
1. perform backups

The most comfortable way is probably to use the admin-client in the cmd directory. It is has support for unix-socket/tcp/tls but is not
yet able to identify itself with it's own certificate. The whole 'adding a CA' to the admin server is not tested yet but should work.

Example for performing backups and put them in a directory with the current date

```
go run src/cmd/admin-client/main.go -s sockets/control.socket backup "backups/$(date -u +"%Y-%m-%d__%H:%M:%S")"
```


Example of how to use the client if all tls options are enabled:
1. Client certificate required
2. Check server certificate against own root-ca
3. custom servername (needed if your server only runs on an ip or on localhost not on an actual domain)

```
go run cmd/admin-client/main.go --tls lsusrs --cert tlsexample/zertifikat-pub.pem --key tlsexample/zertifikat-key.pem --cacert tlsexample/ca-root.pem --servername test-server.test
```

## Generate your own CA for the server/client
* generate CA private key: openssl genrsa -out ca-key.pem 2048
* generate CA public cert: openssl req -x509 -new -nodes -extensions v3_ca -key ca-key.pem -days 1024 -out ca-root.pem -sha512

## Generate new certs for the server/client
* create key:               openssl genrsa -out server-key.pem 4096
* create signing request:   openssl req -new -key server-key.pem -out server.csr -sha51
* create cert:              openssl x509 -req -in server.csr -CA ca-root.pem -CAkey ca-key.pem -CAcreateserial -out server-pub.pem -days 365 -sha512
