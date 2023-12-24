package req

import (
    "SchemaGPT-lab/utils"
    "encoding/json"
)

type PostBody interface {
    Valid() bool
};

type Request struct {
    Success bool
    Error string
    Body PostBody
}

func UnmarshalBody(body []byte, data PostBody) Request {
    err := json.Unmarshal(body,data);
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        return Request{Success:false, Error:errStr};
    } else {
        if data.Valid() {
            return Request{Success:true, Body:data};
        } else {
            return Request{Success:false, Error:"POST body doesn't have all required fields."};
        }
    }
}

// func ExtractRequest()