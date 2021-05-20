# Estuary Agent

The agent is written in Go, and it executes low-level commands. Part of the Estuary stack.

It enables any use case which implies system commands:

- Controlling and configuring the machines (via REST API)
- Exposing CLI applications via REST API
- Testing support by enabling SUT control and automation framework control
- IoT
- Home control integrations

It supports command execution having several modes:

- Commands executed sequentially
- Commands executed in parallel
- Commands executed in background
- Commands executed synchronously

With the help of the agent the user can also do IO operations:

- File upload and download (binary / text)
- Folder download (as zip archive)

This code acts both as a microservice as well as a library.

## Fluentd logging

- FLUENTD_IP_PORT -> This env var sets the fluentd ip:port connection. Example: localhost:24224

## Authentication

- HTTP_AUTH_TOKEN -> This env var sets the auth token for the service. Will be matched with the header 'Token'  

[!!!]() Use this env variable, otherwise you will open a major security hole. The attacker will have access to your system. [!!!]()


## Eureka registration

To have all your Agent instances in a central location, Netflix Eureka is used. This means your client will discover all services used for your test and then spread the tests across all.

Start Eureka server with docker (example):

docker run -p 8080:8080 estuaryoss/netflix-eureka:1.10.5

-  EUREKA_SERVER -> Example: EUREKA_SERVER=http://10.10.15.30:8080/eureka/v2  
-  APP_IP_PORT-> Example: APP_IP_PORT=10.10.15.28:8081. This is the information on how the instance in seen externally. If your instance is inside docker, put the host&port of the host VM where the docker container resides. 

## Enable HTTPS

-  HTTPS_ENABLE -> The main env var.  Example: HTTPS_ENABLE=true
-  HTTPS_CER and HTTPS_KEY - > Set the certificate and the private key path with these env variables. If you do not set cert and private key file env vars, it defaults to a folder in the same path called https, and the default files https/cert.pem and https/key.pem.

! Please also change the app port by setting the env var called **PORT** to *8443*. Default is 8080.

## Environment variables injection
User defined environment variables will be stored in a 'virtual' environment. The extra env vars will be used by the process that executes system commands.  
There are two ways to inject user defined environment variables.    
-   call POST on **/env** endpoint. The body will contain the env vars in JSON format. E.g. {"FOO1":"BAR1"}  
-   create an **environment.properties** file with the extra env vars needed and place it in the same path as the binary. Example in this repo.  

## Example output
curl -X POST -d 'ls -lrt' http://localhost:8080/command

```json
{
    "code": 1000,
    "message": "Success",
    "description": {
        "finished": true,
        "started": false,
        "startedat": "2020-08-15 19:38:16.138962",
        "finishedat": "2020-08-15 19:38:16.151067",
        "duration": 0.012,
        "pid": 2315,
        "id": "none",
        "commands": {
            "ls -lrt": {
                "status": "finished",
                "details": {
                    "out": "total 371436\n-rwxr-xr-x 1 dinuta qa  13258464 Jun 24 09:25 main-linux\ndrwxr-xr-x 4 dinuta qa        40 Jul  1 11:42 tmp\n-rw-r--r-- 1 dinuta qa  77707265 Jul 25 19:38 testrunner-linux.zip\n-rw------- 1 dinuta qa   4911271 Aug 14 10:00 nohup.out\n",
                    "err": "",
                    "code": 0,
                    "pid": 6803,
                    "args": [
                        "/bin/sh",
                        "-c",
                        "ls -lrt"
                    ]
                },
                "startedat": "2020-08-15 19:38:16.138970",
                "finishedat": "2020-08-15 19:38:16.150976",
                "duration": 0.012
            }
        }
    },
    "timestamp": "2020-08-15 19:38:16.151113",
    "path": "/command?",
    "name": "estuary-agent",
    "version": "4.0.8"
}
```

References

1. https://www.tutorialspoint.com/go/
2. https://golang.org
3. https://github.com/sohamkamani/jwt-go-example
4. https://github.com/adigunhammedolalekan/go-contacts
