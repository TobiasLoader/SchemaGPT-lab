package res

import (
    "SchemaGPT-lab/utils"
    "SchemaGPT-lab/dbs"
    "fmt"
    "net/http"
    "encoding/json"
)

type Response struct {
    Success bool
    Error string
    Data json.RawMessage
}

type SuccessDBType func(dbs.AnimalData) Response
type FailureDBType func(string) Response


func DefaultFailure(fileName string) Response {
    err := "could not read database '"+fileName+"'";
    utils.Error(err);
    return Response{Success:false,Error:err};
}

func ConstructResponse(db dbs.AnimalData) Response {
    data, err := json.Marshal(db)
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        return Response{Success:false,Error:errStr};
    }
    return Response{Success:true,Data:data};
}

func SendResponse(res Response, w http.ResponseWriter, r *http.Request){
    resData, err := json.Marshal(res);
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        w.WriteHeader(http.StatusInternalServerError);
        fmt.Fprintf(w, "Internal Server Error");
    } else {
        w.Header().Set("Content-Type", "application/json");
        w.Write(resData);
    }
}