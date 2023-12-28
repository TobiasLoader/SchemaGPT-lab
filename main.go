package main

import (
    "SchemaGPT-lab/dbs"
    "SchemaGPT-lab/req"
    "SchemaGPT-lab/res"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "github.com/joho/godotenv"
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

func getAnimal(w http.ResponseWriter, r *http.Request) {
    queryParams := r.URL.Query()
    animals, exists := queryParams["animal"]
    // fmt.Println(queryParams)
    if !exists || len(animals) == 0 {
        res.SendResponse(res.ErrorResponse("No 'animal' query parameter supplied"), w, r);
        return
    }
    animal := animals[0]
    obj := db("db.json",func (db dbs.AnimalData) res.Response {
        animalData, exists := db.Data[animal]
        if exists {
            return res.ConstructResponse(animalData);
        } else {
            return res.ErrorResponse("Animal provided does not exist in database.");
        }
    }, res.DefaultFailure);
    
    res.SendResponse(obj, w, r);
}


// POST api endpoints

//// POST AUTHOR

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


/// POST ANIMAL CHARACTERISTIC

type PostAnimalCharacteristic struct{
    Animal string
    Characteristic string
}

func (body *PostAnimalCharacteristic) Valid() bool {
    return body.Animal != "" && body.Characteristic != ""
}

func postAnimalCharacteristic(body []byte, w http.ResponseWriter, r *http.Request) {
    bodyReq := req.UnmarshalBody(body, &PostAnimalCharacteristic{});
    if bodyReq.Success {
        if typeCheckedPostBody, ok := bodyReq.Body.(*PostAnimalCharacteristic); ok {
            postAnimalCharacteristicReq(typeCheckedPostBody, w, r);
        } else {
            res.SendResponse(res.ErrorResponse("Incorrect type"), w, r);
        }
    } else {
        res.SendResponse(res.ErrorResponse(bodyReq.Error), w, r);
    }
}

func postAnimalCharacteristicReq(body *PostAnimalCharacteristic, w http.ResponseWriter, r *http.Request) {
    obj := db("db.json",func (db dbs.AnimalData) res.Response {
        animalData, exists := db.Data[body.Animal]
        if exists {
            animalData.Characteristics = append(animalData.Characteristics,body.Characteristic);
            db.Data[body.Animal] = animalData;
            written := dbs.WriteDB("db.json",db);
            if written {
                return res.SuccessResponse();
            } else {
                return res.ErrorResponse("Characteristic not written to database.");
            }
        } else {
            return res.ErrorResponse("Animal provided does not exist in database.");
        }
    }, res.DefaultFailure);
    
    res.SendResponse(obj, w, r);
}

// deploy

func deploy(){
    godotenv.Load() // This will load the .env file
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    deployEnv, exists := os.LookupEnv("DEPLOY_ENV");
    if exists {
        if deployEnv == "local" {
            fmt.Println("Listen on http://localhost:"+port);
            http.ListenAndServe("localhost:"+port, nil);
        } else if deployEnv == "prod" {
            fmt.Println("Listen on port "+port);
            http.ListenAndServe(":"+port, nil);
        } else {
            fmt.Println("Unrecognised deploy environment '"+deployEnv+"'");
        }
    } else {
        fmt.Println("Could't find environment variable 'DEPLOY_ENV'");
    }
}

// main

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        errRes := res.ErrorResponse("Nothing at '"+r.URL.Path+"' instead go to '/getDB'");
        res.SendResponse(errRes, w, r);
    })
    
    GET("/getDB", getDB);
    GET("/getAnimal", getAnimal);
    POST("/postAuthor", postAuthor);
    POST("/postAnimalCharacteristic", postAnimalCharacteristic);

    deploy();
}
