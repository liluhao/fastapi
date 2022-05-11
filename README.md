# fastapi
## 项目介绍

### 功能一

* 用于快速构建 API ，相比于传统的“type HandlerFunc func(*Context)”处理函数Handler，本项目在使用处理函数Handler时候，不仅仅可以传入Context，还可以传入结构体struct，即如下

```go
func(ctx *gin.Context, in struct{}) (out struct{}, err error) { return }
```

### 功能二

* 创建一个 API后可以 获取 Swagger 定义的JSON序列,即可以无需任何标记的自动生成API的 OpenAPI/Swagger 架构

## 项目使用例子

### 展示功能一

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

> 正确请求如下：

```go
$ curl -H "Content-Type: application/json" -X POST --data '{"phrase": "hello"}' localhost:8080/api/echo
//     {"response":{"original_input":{"phrase":"hello"}}}
```

> 任举一种错误请求

```GO
$  curl -H "Content-Type: application/json" -X POST --data '{"phrase": "hello"}' localhost:8080/api/echoaDASD
//  {"error":"handler not found"}

```

### 展示功能二

> 运行如下代码：

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
	//func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error): MarshalIndent类似Marshal但会使用缩进将输出格式化。
	jsonBytes, _ := json.MarshalIndent(swagger, "", "    ")
	fmt.Println(string(jsonBytes))
}

```

> 控制台输出如下：

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

## 项目开发依赖

* go 1.18
* github.com/gin-gonic/gin v1.7.7

* github.com/go-openapi/spec v0.20.4

