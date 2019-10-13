package model

// 站点元信息

type Pair struct {
	Key   string `meddler:"kv_key"`
	Value string `meddler:"kv_value"`
}
