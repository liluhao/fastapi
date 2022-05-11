# go-fastapi
## Introduction

go-fastapi is a library to quickly build APIs. 

Create an API and get Swagger definition for free

## Features

* Auto-generated OpenAPI/Swagger schema without any markup
* Declare handlers using types, not just `Context`
* Based on [gin](https://github.com/gin-gonic/gin) framework

> [OpenAPI 规范 (中文版) (apifox.cn)](https://openapi.apifox.cn/)

## Example

### 1.Declare handlers using types, not just `Context`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/llh/fastapi"
)

type EchoInput struct {
	Phrase string `json:"phrase"`
}

type EchoOutput struct {
	OriginalInput EchoInput `json:"original_input"`
}

func EchoHandler(ctx *gin.Context, in EchoInput) (out EchoOutput, err error) {
	out.OriginalInput = in
	return
}

func main() {
	r := gin.Default()

	myRouter := fastapi.NewRouter()
	myRouter.AddCall("/echo", EchoHandler)

	r.POST("/api/*path", myRouter.GinHandler) // must have *path parameter
	r.Run()
}

```

```GO
$ curl -H "Content-Type: application/json" -X POST --data '{"phrase": "hello"}' localhost:8080/api/echo
//     {"response":{"original_input":{"phrase":"hello"}}}
```

```GO
$  curl -H "Content-Type: application/json" -X POST --data '{"phrase": "hello"}' localhost:8080/api/echoaDASD
//  {"error":"handler not found"}

```

### 2.To generate OpenAPI/Swagger schema:

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/llh/fastapi"
)

type EchoInput struct {
	Phrase string `json:"phrase"`
}

type EchoOutput struct {
	OriginalInput EchoInput `json:"original_input"`
}

func EchoHandler(ctx *gin.Context, in EchoInput) (out EchoOutput, err error) {
	out.OriginalInput = in
	return
}

func main() {
	myRouter := fastapi.NewRouter()
	myRouter.AddCall("/echo", EchoHandler)

	swagger := myRouter.EmitOpenAPIDefinition()
	swagger.Info.Title = "My awesome API"
	jsonBytes, _ := json.MarshalIndent(swagger, "", "    ")
	fmt.Println(string(jsonBytes))
}

```

```go
{
    "swagger": "2.0",
    "info": {
        "title": "My awesome API",
        "version": "1.0"
    },
    "paths": {
        "/echo": {
            "post": {
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/EchoInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "$ref": "#/definitions/EchoOutput"
                    }
                }
            }
        }
    },
    "definitions": {
        "EchoInput": {
            "type": "object",
            "properties": {
                "phrase": {
                    "type": "string"
                }
            }
        },
        "EchoOutput": {
            "type": "object",
            "properties": {
                "original_input": {
                    "$ref": "#/definitions/EchoInput"
                }
            }
        }
    }
}
进程 已完成，退出代码为 0
```



<img src="https://mdmdmdmd.oss-cn-beijing.aliyuncs.com/img/146807480-be53b3fb-6de8-451f-8373-e8d6da54a032.png" width="400px" height="auto">

## Dependencies

* go 1.18
* github.com/gin-gonic/gin v1.7.7

* github.com/go-openapi/spec v0.20.4

