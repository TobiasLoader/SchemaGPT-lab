package main

import (
    "SchemaGPT-lab/dbs"
    "SchemaGPT-lab/res"
    "fmt"
    "io/ioutil"
    "net/http"
)

func db(fileName string, success res.SuccessDBType, failure res.FailureDBType) res.Response {
    maybeDB := dbs.MaybeReadDB(fileName);
    if maybeDB.Success {
        db := maybeDB.DB;
        return success(db);
    } else {
        return failure(fileName);
    }
}

// api endpoints

func getDB(w http.ResponseWriter, r *http.Request) {
    obj := db("db.json",func (db dbs.AnimalData) res.Response {
        fmt.Println(db);
        return res.ConstructResponse(db);
    }, res.DefaultFailure);
    
    res.SendResponse(obj, w, r);
}

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
            http.Error(w, "Method is not supported", http.StatusNotFound)
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


func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
    })
    
    GET("/getDB", getDB);
    

    http.ListenAndServe(":8080", nil)
}
