package graphalg

import (
	"fmt"
	"log"
	"reflect"

	"github.com/golang/protobuf/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

var (
	typeRegistry = make(map[string]reflect.Type)
)

func RegisterMessage(i interface{}) string {
	t := reflect.TypeOf(i)
	ti := reflect.New(t).Interface()
	url := fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	if _, ok := ti.(proto.Message); !ok {
		panic("type does " + url + "not support proto.Message interface")
	}
	log.Println(url)
	typeRegistry[url] = t
	return url
}

func unmarshalAny(any *google_protobuf.Any) (proto.Message, error) {
	class := any.TypeUrl
	bytes := any.Value

	instance := reflect.New(typeRegistry[class]).Interface()
	err := proto.Unmarshal(bytes, instance.(proto.Message))
	if err != nil {
		return nil, err
	}
	log.Printf("instance: %v", instance)

	return instance.(proto.Message), nil
}
