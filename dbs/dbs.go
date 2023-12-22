package dbs

import (
    "SchemaGPT-lab/utils"
    "encoding/json"
    "os"
)

type AnimalData struct {
    Author string            `json:"author"`
    Schema Schema            `json:"schema"`
    Data   map[string]Animal `json:"data"`
}

type Schema struct {
    ID             int `json:"id"`
    Name           string `json:"name"`
    Species        string `json:"species"`
    Diet           string `json:"diet"`
    Habitat        string `json:"habitat"`
    Characteristics []string `json:"characteristics"`
}

type Animal struct {
    Schema
}

type MaybeAnimalData struct {
    Success bool
    DB AnimalData
}

func MaybeReadDB(fileName string) MaybeAnimalData {
    // Open the JSON file
    file, err := os.Open(fileName)
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        return MaybeAnimalData{Success:false}
    }
    defer file.Close()
    
    // Decode the JSON file
    var local_db AnimalData
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&local_db)
    if err != nil {
        errStr := err.Error();
        utils.Error(errStr);
        return MaybeAnimalData{Success:false}
    }
    
    return MaybeAnimalData{Success:true,DB:local_db}
}
