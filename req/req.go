package req

import (
    "SchemaGPT-lab/utils"
    "fmt"
    "encoding/json"
)

type PostBody interface {};

type Request struct {
    Success bool
    Error string
    Body PostBody
}

func DeconstructBody(body []byte, data PostBody) Request {
    err := json.Unmarshal(body,data);
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        return Request{Success:false, Error:errStr};
    } else {
        fmt.Println("data",data);
        return Request{Success:true, Body:data};
    }
}

// func ExtractRequest()