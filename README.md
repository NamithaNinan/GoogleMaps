# cmpe273-Assignment2
google map in go language

#What is the code about ?
This program aims at making use of the Google maps API in order to fetch the coordinate locations of a user given address. I have tried implementin basic CRUD operations like illustrated below :

1. Create

This uses the POST verb of a REST API . Here the user gives input in the following format through the link(http://localhost:8082/locations)

{

    "Name" : "namitha ninan",
    
    "Address" : "101 san fernando",
    
    "City" : "san jose",
    
    "State" : "CA",
    
    "Zip" : "95112"
    
}

The program reads this from the body of the request and provides the address to Google Map API which then returns the coordinates. This is then fed into MongoLab along with a unique ID to identify it by later

{

  "id": "5666312266bae326c4708949",
  
  "name": "namitha ninan",
  
  "address": "101 san fernando",
  
  "city": "san jose",
  
  "state": "CA",
  
  "zip": "95112",
  
  "coordinates": {
    "latitude": -121.887067,
    "longitude": 37.336259
  }
  
}

2. Read

Get command is used to retrieve the data of a specific id from the database and display it to the user. On giving the address http://localhost:8082/locations/5666312266bae326c4708949, we'll get the same output as shown above as the id's provided are the same.

3. Update

Put command is used to update/change the input given. Similar to get we specify the id of the record to be altered http://localhost:8082/locations/5666312266bae326c4708949

Now suppose we change the name and address and give the following input,

{

    "Name" : "Abhishek Gupta",
    
    "Address" : "villa torrino apartments",
    
    "City" : "san jose",
    
    "State" : "CA",
    
    "Zip" : "95112"
    
}

We now get the modified output for the same id

{

  "id": "5666312266bae326c4708949",
  
  "name": "Abhishek Gupta",
  
  "address": "villa torrino apartments",
  
  "city": "san jose",
  
  "state": "CA",
  
  "zip": "95112",
  
  "coordinates": {
    "latitude": -121.8847222,
    "longitude": 37.3456227
  }
  
}

4. Delete

This is as simple as using the Delete Verb and providing the same link as above where we specify the id to be deleted
