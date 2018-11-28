# HaDiBar #
The goal is to manage the accountings for our beverages on our floor of a dorm conveniently

This is a Server for the HaDiBar-API. It exposes functionality to manage accounts and beverages.

The definiton of api can be found at [Hadibar-API](https://github.com/killingspark/Hadibar-api)
a reference implementation of a webapp that uses said API can be found at [Hadibar-Webapp](https://github.com/killingspark/Hadibar-Webapp).
The Webapp also contains some conveniance functions capsulating the ajax-calls and session management to use the api from javascript 

Webapp made with Vuejs, JQuery and Bootstrap
Accounts/Beverages/Users are stored in json files

## Users ##
There is no explicit user management right now. Usernames are aquired by logging in with the name and the password for the firs time.
This will hopefully be improved in the future (with password-resets/registering with an email etc,etc)

## Test the server without the webapp ##
Make calls with curl like those:

Testlogin: 
```
    export SES="$(curl -X GET 127.0.0.1:8080/api/session/getid)"
    curl -X POST 127.0.0.1:8080/api/session/login -H "sessionID: "$SES --form "name=Moritz" --form "password=test"
```
