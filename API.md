# Interface for HaDiBar

The answer of the server always is structured like this:

    {"status": "_somestatus_", "response":_array, object or primitive_}   

**Response** specifies the _response_ field.  
**Status**  specifies the _status_ field  
If not stated otherwise **Status** contains  
1. "OK" if the action was executed successfully  
2. "ERROR" if the action failed  
3. "NOSES" if you havent send an sessionid as a header  
4. "NOAUTH" if you havent logged in yet  

All requests (besides /session/getid) must include a sessionid as a header like this:  

    $.ajax({
        url: "/beverage/all",
        type: 'GET',
        data: {},
        beforeSend: function (xhr) {
          xhr.setRequestHeader("sessionID", app.sessionid)
        },
        success: function (data) {
          status = JSON.parse(data).status
        }
    })  

There are two relevant models:
1. Beverage
    {
        ID: long,
        name: String,
        price: int
    }
2. Account
    {
        ID: long,
        name: String,
        values: Map<String,String>
    }
The current amount of money an account has is stored in values["value"] as an integer representing cents.


There are three groups of endpoints in the API:

1. /api/session
2. /api/f/beverage
3. /api/f/account

All calls to the /api/f/:floor routes require you to be logged in to check wether you are allowed to see the data

## /api/session

**Note**: If you get youself a new id and send this new id in the header your old login will *not* be connected with this new id 
**URL**: /api/session/getid  
**Method**: GET  
**QueryParameters**: none  
**FormParameters**: none  
**Response**: sessionid if successful, "ERROR" if not. 

**Note**: Connects the sessionid to an account on whose behalf the actions are executed by the server  
**URL**: /api/session/login  
**Method**: POST  
**QueryParameters**: none  
**FormParameters**: _name, password_   
**Response**: nothing

**Note**: Disconnects the sessionid from the old account. You may login with another account, or get a new sessionid  
**URL**: /api/session/logout  
**Method**: POST  
**QueryParameters**: none  
**FormParameters**: none   
**Response**: nothing

## /beverage Get all the Drinks!
**URL**: /api/f/beverage/all  
**Method**: Get
**QueryParameters**: none  
**FormParameters**: none  
**Response**: Array of beverages  
  
**URL**: /api/f/beverage/get  
**Method**: Get  
**QueryParameters**: _id_ ID of the beverage 
**FormParameters**: none  
**Response**: the beverage specified by :id  
  
**URL**: /api/f/beverage/update   
**Method**: POST  
**QueryParameters**: _id_ ID of the beverage  
**FormParameters**: _value, name_: the new values for the beverage with ID=id  
**Response**: the updated beverage object if successful, "ERROR" if not.   

**URL**: /api/f/beverage/delete   
**Method**: DELETE  
**QueryParameters**: _id_ ID of the beverage  
**FormParameters**: none
**Response**: nothing  
   
**URL**: /api/f/beverage/new  
**Method**: POST  
**QueryParameters**: none  
**FormParameters**: _value, name_: the values for the new beverage  
**Response**: the beverage object if successful, "ERROR" if not.  

## /accounts 

**URL**: /api/f/account/all  
**Method**: GET  
**QueryParameters**: none  
**FormParameters**: none   
**Response**: Array of all accounts  

**URL**: /api/f/account/get  
**Method**: GET  
**QueryParameters**: _id_ ID of the account 
**FormParameters**: none   
**Response**: the beverage object if successful, "ERROR" if not.  

**URL**: /api/f/account/update
**Method**: POST  
**QueryParameters**: _id_ ID of the account  
**FormParameters**: difference   
**Response**: the updated account object if successful, "ERROR" if not.  


