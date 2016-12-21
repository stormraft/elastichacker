package redisc

import (
    "bytes"
    "encoding/gob"
    "fmt"
//	"github.com/garyburd/redigo/redis"
)

//func Process_json_bytes(index, itype, id string, byteArray []byte) {
//}

type P struct {
    Itype string
    Id int
    Json []byte
}

func Process_json_test(index, itype string, id int) error {
    c := getRedisConn()
    defer c.Close()

    _, err := c.Do("HSET", index, id, itype)
    return err
}

func Process_json_bytes(index, itype string, id int, byteArray []byte) error {
    c := getRedisConn()
    defer c.Close()

    nbytearray := process_bytes(itype, id, byteArray)

    _, err := c.Do("HSET", index, id, nbytearray)
    return err
}

func process_bytes(itype string, id int, byteArray []byte) ([]byte) {
    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err := enc.Encode(P{itype, id, byteArray})
    if err != nil {
        fmt.Println("process_bytes error in Encoder")
    }
    return buf.Bytes()
}
