package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"test1/mongo"
	"time"
)

func main() {
	mongo.GetMongoManager().Init("mongodb://127.0.0.1:27017", "", "", "rts")
	http.HandleFunc("/genOrder", handle_GenOrder)
	http.ListenAndServe(":8001", nil)
}

type gen_order_param struct {
	OrderID  int64   `json:"orderID" bson:"orderID"`
	ItemID   int64   `json:"itemID" bson:"itemID"`
	Count    int32   `json:"count" bson:"count"`
	Pay      float32 `json:"pay" bson:"pay"`
	PayState int32   `json:"payState" bson:"payState"`
	CreateAt string  `json:"createAt" bson:"createAt"`
}

func handle_GenOrder(w http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	c := &gen_order_param{}
	err = json.Unmarshal(body, c)
	if err != nil {
		return
	}
	write_to_db(c)
	ResultSuccess(w, nil)
}
func write_to_db(c *gen_order_param) {
	coll := mongo.GetMongoManager().GetCollection("order")
	if coll != nil {
		c.CreateAt = time.Now().Format("2006-01-02 15:04:05")
		coll.InsertOne(context.Background(), c)
	}
}
func CommonResult(w http.ResponseWriter, code int32, msg string) {
	Result(w, code, msg, nil)
}

func ResultSuccess(w http.ResponseWriter, obj interface{}) {
	Result(w, 0, "success", obj)
}

func Result(w http.ResponseWriter, code int32, msg string, obj interface{}) {
	// 当返回json时, 添加json头
	w.Header().Add("Content-Type", "application/json")
	result := map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	if obj != nil {
		vi := reflect.ValueOf(obj)
		if vi.Kind() == reflect.Interface {
			if !vi.IsNil() {
				result["data"] = obj
			}
		} else {
			result["data"] = obj
		}
	}
	_ = json.NewEncoder(w).Encode(result)
}
