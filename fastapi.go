package fastapi

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Router struct {
	routesMap map[string]interface{}
}

func NewRouter() *Router {
	return &Router{
		routesMap: make(map[string]interface{}),
	}
}

func (r *Router) AddCall(path string, handler interface{}) {
	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() != 2 { // NumIn() ：返回func类型的参数个数，如果不是函数，将会panic
		panic("Wrong number of arguments")
	}
	if handlerType.NumOut() != 2 { //NumOut() int:返回func类型的返回值个数，如果不是函数，将会panic
		panic("Wrong number of return values")
	}

	ginCtxType := reflect.TypeOf(&gin.Context{})      //一定要传入指针类型，方便传入ConvertibleTo方法内进行比较
	if !handlerType.In(0).ConvertibleTo(ginCtxType) { //In(i int) Type:返回func类型的第i个参数的类型即返回Type接口类型，如非函数或者i不在[0, NumIn())内将会panic;  ConvertibleTo(u Type) bool:如该类型的值可以转换为u代表的类型，返回真
		panic("First argument should be *gin.Context!")
	}

	if handlerType.In(1).Kind() != reflect.Struct {
		panic("Second argument must be a struct")
	}

	if handlerType.Out(0).Kind() != reflect.Struct {
		panic("First return value be a struct")
	}
	//注意：(*error)(nil)的写法，使用其他形式会出错
	errorInterface := reflect.TypeOf((*error)(nil)).Elem() //Elem() Type:返回该类型的元素类型，如果该类型的Kind不是Array、Chan、Map、Ptr或Slice，会panic
	if !handlerType.Out(1).Implements(errorInterface) {    //  Out(i int) Type: 返回func类型的第i个返回值的类型，如非函数或者i不在[0, NumOut())内将会panic; Implements(u Type) bool:如果该类型实现了u代表的接口，会返回真
		panic("Second return value should be an error")
	}

	r.routesMap[path] = handler //注意：存入的path的带/的，即path=/echo
}

func (r *Router) GinHandler(c *gin.Context) {
	path := c.Param("path") //获取路径参数
	log.Printf(path)        //用于在控制台打印,即打印出 :   2022/05/11 18:12:09 /echo
	handlerFuncPtr, present := r.routesMap[path]
	if !present {
		c.JSON(http.StatusNotFound, gin.H{"error": "handler not found"}) //404
		return
	}

	inputType := reflect.TypeOf(handlerFuncPtr).In(1)
	inputVal := reflect.New(inputType).Interface() //New(typ Type) Value: 返回一个Value类型值，该值持有一个指向类型为typ的新申请的零值的指针，返回值的Type为PtrTo(typ);Interface() (i interface{}):返回v当前持有的值
	err := c.BindJSON(inputVal)                    //inputVal： &{hello}

	fmt.Println(inputVal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"}) //400
		return
	}

	toCall := reflect.ValueOf(handlerFuncPtr)
	// Call(in []Value) []Value: Call方法使用输入的参数in调用v持有的函数
	outputVal := toCall.Call(
		[]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(inputVal).Elem(), //Elem() Value：Elem返回v持有的接口保管的值的Value封装，或者v持有的指针指向的值的Value封装。如果v的Kind不是Interface或Ptr会panic；如果v持有的值为nil，会返回Value零值。
		},
	)
	returnedErr := outputVal[1].Interface()

	if returnedErr != nil || !outputVal[1].IsNil() { //保证outputVal[1].Interface()==nil
		c.JSON(http.StatusInternalServerError, gin.H{"error": returnedErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": outputVal[0].Interface()})
}
