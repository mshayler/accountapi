# Account API

RESTful Account API written in *Go* implementing the [GoKit framework](https://gokit.io/) utilizing [Redis](https://redis.io/) for persistence.

#### Usage
* run ``` docker-compose up``` from inside directory to start redis
* ``` docker run --rm --name managerapi --network accountapi_default -p 8081:8081 amanager```
* Generate an account using the ``` /create ``` endpoint, and your desired username and password
* Obtain a verification token using ``` /login ``` with the valid credentials
* Use the verification token with ``` /verify ```

Go with Docker can finicky, if necessary adjust the connection line for the persistance to just point to localhost and run the go file directly.

##### Utilize the postman tests, or generate your own requests using a json format.

#### Port Information

* Redis: 6379
* API: 8081

#### Routes

|   | Params |
| ------------- |:-------------:|
| /create     | { user : string, pass : string }     |
| /login     | { user : string, pass : string }      |
| /verify      | { user : string, token : string }      |

##### Example Requests
* POST @ localhost:8081/create

```
{
    "user":"test",
    "pass":"1234"
}
```

* POST @ localhost:8081/login

```
{
    "user":"test",
    "pass":"1234"
}
```

* GET @ localhost:8081/verify

```
{
    "user":"test",
    "token":{GENERATED_TOKEN}
}
```

## Further Improvements
* Primarily the biggest room for improvement in this project is logging and error collection. Utilizing HTTP codes and providing more detailed errors is the next best step for this project.
* Secondly, building a more robust Authentication system, or more standardized. It would be best to implement a better flow using authentication headers to better mediate who is accessing this api, and at what times.
* Finally, setting up timeouts, and a cache flushing mechanism for the tokens is necessary, as the current implementation means a user just has to log in the one time. Redis supports a key-value flush on a timed basis, this could be extremely useful.
* 