# Get Database
curl "http://localhost:8080/getDB"

# Get Animal
curl "http://localhost:8080/getAnimal?animal=Chimpanzee"

# Post Author
curl -X POST -H "Content-Type: application/json" -d '{"author":"Toby"}' http://localhost:8080/postAuthor

# Post Animal Characteristic
curl -X POST -H "Content-Type: application/json" -d '{"animal":"Tiger","characteristic":"aawesome"}' http://localhost:8080/postAnimalCharacteristic

# Post New Animal
curl -X POST -H "Content-Type: application/json" -d '{"population":0,"name":"Velociraptor","Species":"V. mongoliensis","diet":"Carnivorous","habitat":"desert"}' http://localhost:8080/postNewAnimal