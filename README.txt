The goal is to manage the accountings for our beverages on our floor of a dorm conveniently

This repo consists of the server, wich exposes a restAPI
The definiton of said api can be found in another repo
a reference implementation of a webapp that uses said API can be found in another repo

Webapp made with Vuejs, JQuery and Bootstrap
Accounts/Beverages/Users are stored in json files

Testlogin: 
```
    export SES="$(curl -X GET 127.0.0.1:8080/api/session/getid)"
    curl -X POST 127.0.0.1:8080/api/session/login -H "sessionID: "$SES --form "name=Moritz" --form "password=test"
```
