package consul_kv_mapper

import (
	"github.com/hashicorp/consul/api"
	"strings"
)

func BuildMap(client *api.Client, prefix string) (*MapType, error) {
	kvMap := &MapType{Value: "root"}

	kvPairs, _, err := client.KV().List(prefix, nil)
	if err != nil {
		return nil, err
	}

	if len(kvPairs) != 0 {

		for _, pair := range kvPairs {

			key := strings.Replace(pair.Key, prefix+"/", "", -1)

			keyParts := strings.Split(key, "/")

			var parts []KeyType

			for _, keyPart := range keyParts {
				parts = append(parts, KeyType(keyPart))
			}

			for i := 0; i < len(parts); i++ {
				if len(parts) > i {
					if (kvMap.Get(parts[:(i+1)]...) == nil) && (parts[i] != "") {
						var val ValueType
						val = ValueType("")

						if string(pair.Value) != "" {
							val = ValueType(pair.Value)
						}

						addNextLevel(kvMap, parts[i], val, parts[:i]...)
					}
				}
			}
		}

	}

	return kvMap, nil
}

func addNextLevel(n *MapType, k KeyType, v ValueType, p ...KeyType) {

	if len(p) > 0 {
		n.Get(p...).Add(k, v)
	} else {
		n.Add(k, v)
	}
}

type KeyType string
type ValueType string

type MapType struct {
	Children map[KeyType]*MapType
	Value    ValueType
}

func (n *MapType) Add(key KeyType, v ValueType) {
	if n.Children == nil {
		n.Children = map[KeyType]*MapType{}
	}
	n.Children[key] = &MapType{Value: v}
}

func (n *MapType) Get(keys ...KeyType) *MapType {
	for _, key := range keys {
		n = n.Children[key]
	}
	return n
}

func (n *MapType) Set(v ValueType, keys ...KeyType) {
	n = n.Get(keys...)
	n.Value = v
}
