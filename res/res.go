package res

import (
    "SchemaGPT-lab/utils"
    "SchemaGPT-lab/dbs"
    "fmt"
    "net/http"
    "encoding/json"
)

type ResData interface {};

type Response struct {
    Success bool
    Error string
    Data json.RawMessage
}

func ErrorResponse(err string) Response {
    return Response{Success:false,Error:err};
}

func SuccessResponse() Response {
    return Response{Success:true};
}

type SuccessDBType func(dbs.AnimalData) Response
type FailureDBType func(string) Response

func DefaultFailure(fileName string) Response {
    err := "could not read database '"+fileName+"'";
    utils.Error(err);
    return ErrorResponse(err);
}

func ConstructResponse(resData ResData) Response {
    data, err := json.Marshal(resData)
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        return ErrorResponse(errStr);
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

func WriteDBRes(fileName string, data dbs.AnimalData) Response {
    written := dbs.WriteDB("db.json",data);
    if written {
        return SuccessResponse();
    } else {
        return ErrorResponse("Failed to write to database.");
    }
}