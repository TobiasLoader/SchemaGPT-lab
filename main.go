package main

import (
    "SchemaGPT-lab/dbs"
    "SchemaGPT-lab/req"
    "SchemaGPT-lab/res"
    "fmt"
    "io/ioutil"
    "net/http"
)

// GET & POST

type HandlerGET func(http.ResponseWriter,*http.Request)
type HandlerPOST func([]byte,http.ResponseWriter,*http.Request)

func GET(path string, handler HandlerGET){
    http.HandleFunc(path, handler);
}

func POST(path string, handler HandlerPOST){
    http.HandleFunc(path, func (w http.ResponseWriter, r *http.Request){
        // Check if the request method is POST
        if r.Method != http.MethodPost {
            http.Error(w, "Must be a POST request", http.StatusNotFound)
            return
        }
        
        // Read the body of the request
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Error reading request body", http.StatusInternalServerError)
            return
        }
        defer r.Body.Close()
        
        handler(body, w, r);
    });
}

// read database produce response

func db(fileName string, success res.SuccessDBType, failure res.FailureDBType) res.Response {
    maybeDB := dbs.MaybeReadDB(fileName);
    if maybeDB.Success {
        db := maybeDB.DB;
        return success(db);
    } else {
        return failure(fileName);
    }
}

// GET api endpoints

func getDB(w http.ResponseWriter, r *http.Request) {
    obj := db("db.json",func (db dbs.AnimalData) res.Response {
        return res.ConstructResponse(db);
    }, res.DefaultFailure);
    
    res.SendResponse(obj, w, r);
}

// POST api endpoints

type PostAuthor struct{
    Author string
}

func (body *PostAuthor) Valid() bool {
    return body.Author != ""
}

func postAuthor(body []byte, w http.ResponseWriter, r *http.Request) {
    bodyReq := req.UnmarshalBody(body, &PostAuthor{});
    if bodyReq.Success {
        if typeCheckedPostBody, ok := bodyReq.Body.(*PostAuthor); ok {
            postAuthorReq(typeCheckedPostBody, w, r);
        } else {
            res.SendResponse(res.ErrorResponse("Incorrect type"), w, r);
        }
    } else {
        res.SendResponse(res.ErrorResponse(bodyReq.Error), w, r);
    }
}

func postAuthorReq(body *PostAuthor, w http.ResponseWriter, r *http.Request) {
    obj := db("db.json",func (db dbs.AnimalData) res.Response {
        db.Author = body.Author;
        written := dbs.WriteDB("db.json",db);
        if written {
            return res.SuccessResponse();
        } else {
            return res.ErrorResponse("Author not written to database.");
        }
    }, res.DefaultFailure);
    
    res.SendResponse(obj, w, r);
}

// main

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        errRes := res.ErrorResponse("Nothing at '"+r.URL.Path+"' instead go to '/getDB'");
        res.SendResponse(errRes, w, r);
    })
    
    GET("/getDB", getDB);
    POST("/setAuthor", postAuthor);

    fmt.Println("Listen on http://localhost:8080");
    http.ListenAndServe("localhost:8080", nil);
}
