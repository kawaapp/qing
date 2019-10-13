package wxmp

import (
	"testing"
	"encoding/json"
	"log"
)

func TestClient_Decrypt(t *testing.T) {
	data := data1
	text, err := decrypt(data.Key, data.Data, data.IV)
	if err != nil {
		t.Error(err)
	}

	log.Println("decrypt text:", text)

	m := make(map[string]interface{}, 0)
	err = json.Unmarshal(([]byte)(text), &m)
	if err != nil {
		t.Error(err)
	}
	if v := m["purePhoneNumber"]; v != "18610342257" {
		t.Error("expect phone:", "18610342257", "get:", v)
	}
}

type data struct {
	Key string
	Data string
	IV string
}

var data1 = data {
	Key: "nZ64ibIf59p1Z21VqT1fwg==",
	Data: "x3TsJzievTcCo/IUQbPKEKxO8HJF7AFVxQzoS01ASdjzxgj89DKSONA3opOvc8gGGoQTxbeG6tH1VNj1398ElIgthhEoAEHYFtDbOwrP0PwNvbVyto8Oa3HNug0t2HOabryCblemTJZZ7HigsCiZmLU2FYSG3YWwUzWopZVLs6SMGgLUZIqARYg4BAqqcZRMht44mtzJLQtNnX7PTZvgBA==",
	IV: "z0XpJ5Se5NmIu4J9ThovwA==",
}

var data2 = data{
	Key: "3bK7o6+JSjndruyEjweM9g==",
	Data: "ruIjN94RNb3fSklCtz3Yn2bGy70dOm2Unu9rja32x4f++g1D7s73Ikb2FpJPzeLd0txCXR2Sino5ma8K7zCMtw9JvTRdajfEsXTrBxOFiq0sP7vrPkAhrvbWNXZOZdatFotANPdTj8Cj9rRtkJ26bvLgBaIN7jryz7/Cwf3UTqh2hJWaAVJqNMwOA9EnR3vB0nrvC32rQh0ChCeY7M8sDA==",
	IV: "ZPbsSihoj9Ol96cdW4E5WQ==",
}

var data3 = data{
	Key: "tiihtNczf5v6AKRyjwEUhQ==",
	Data: "" +
	"CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZM" +
	"QmRzooG2xrDcvSnxIMXFufNstNGTyaGS" +
	"9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+" +
	"3hVbJSRgv+4lGOETKUQz6OYStslQ142d" +
	"NCuabNPGBzlooOmB231qMM85d2/fV6Ch" +
	"evvXvQP8Hkue1poOFtnEtpyxVLW1zAo6" +
	"/1Xx1COxFvrc2d7UL/lmHInNlxuacJXw" +
	"u0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn" +
	"/Hz7saL8xz+W//FRAUid1OksQaQx4CMs" +
	"8LOddcQhULW4ucetDf96JcR3g0gfRK4P" +
	"C7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB" +
	"6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns" +
	"/8wR2SiRS7MNACwTyrGvt9ts8p12PKFd" +
	"lqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYV" +
	"oKlaRv85IfVunYzO0IKXsyl7JCUjCpoG" +
	"20f0a04COwfneQAGGwd5oa+T8yO5hzuy" +
	"Db/XcxxmK01EpqOyuxINew==",
	IV: "r7BXXKkLb8qrSNn05n0qiA==",
}