package dbs

import (
    "SchemaGPT-lab/utils"
    "encoding/json"
    "os"
)

type AnimalData struct {
    Author string            `json:"author"`
    Data   map[string]Animal `json:"data"`
}

type Animal struct {
    Population             int `json:"population"`
    Name           string `json:"name"`
    Species        string `json:"species"`
    Diet           string `json:"diet"`
    Habitat        string `json:"habitat"`
    Characteristics []string `json:"characteristics"`
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

func WriteDB(fileName string, data AnimalData) bool {
    file, errOpen := os.Create(fileName)
    if errOpen != nil {
        utils.Error("Error could not open file. Error: " + errOpen.Error());
        return false;
    } else {
        jsonData, errMarshal := json.MarshalIndent(data, "", "    ");
        if errMarshal != nil {
            utils.Error("Error could not marshal data. Error: " + errMarshal.Error());
            return false;
        } else {
            _, errWrite := file.Write(jsonData)
            if errWrite != nil {
                utils.Error("Error occurred during writing to the file. Error: " + errWrite.Error());
                return false;
            }
        }
        defer file.Close()
    }
    return true;
}


func (body *Animal) Valid() bool {
    return body.Population >= 0 && body.Name != "" && body.Species != "" && body.Diet != "" && body.Habitat != "";
}