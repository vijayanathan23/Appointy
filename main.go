package main

import (
	"net/http"
	"encoding/json"
	"log"
	"context"
	"strings"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/bson"
	"time"
)

//Creating type Meetings for APIs
type Meetings struct{
	Id	string	`json:"Id" 	bson:"Id"`
	Title	string	`json:"Title"	bson:"Title"`
	StartTime	string	`json:"StartTime"	bson:"StartTime"`
	EndTime	string	`json:"EndTime"	bson:"EndTime"`
	creation_TS time.Time `json:"Creation Timestamp" bson:"Creation Timestamp"`
	Participants []Participants	`json: "Participants" bson: "Participants"`
}

//Creating type Participants for APIs
type Participants struct{
	Name string `json: "Name" bson:"Name"`
	Email string `json: "Email" bson: "Email"`
	RSVP string `json: "RSVP" bson: "RSVP"`
}

func writeToDB(m Meetings){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second);
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err!=nil{
		log.Fatal(err);
	}
	if err!=nil{
		log.Fatal(err);
	}
	database := client.Database("appointy")
	meetingCollection := database.Collection("Meetings")
	meetingCollection.InsertOne(ctx,m)
	if err != nil {
    panic(err)
	}
}

func ReadDB1(ID string) Meetings{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second);
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err!=nil{
		log.Fatal(err);
	}
	if err!=nil{
		log.Fatal(err);
	}
	database := client.Database("appointy")
	meetingCollection := database.Collection("Meetings")
	curser , err := meetingCollection.Find(ctx, ID)
	if err != nil {
		log.Fatal(err)
	}
	l := Meetings{}
	er := curser.Decode(&l)
		if er != nil {
		}
		
	return l;
	
}

func main(){



	http.HandleFunc("/",func (w http.ResponseWriter,r *http.Request)  {
		w.Write([]byte("I am an API and I work in '/meetings'"));
	});
	

	//Get Meeting by ID
	http.HandleFunc("/meetings/",func (w http.ResponseWriter,r *http.Request)  {
		urlPath := strings.Split(r.URL.Path, "/");
		var ID string;
		ID = urlPath[2];
		userJson, err := json.Marshal( ReadDB1(ID))
			if err != nil{
				panic(err)
			}
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusOK)
			//Write json response back to response 
			w.Write(userJson)
		
	})


	
	http.HandleFunc("/meetings",func (w http.ResponseWriter,r *http.Request)  {
		
		if(r.Method == "GET"){
			w.Write([]byte("Send me a POST requset(JSON format) to create meeting or Ask the right question you will 'GET' the response"));
		}	else if(r.Method == "POST"){
			
			meetings := Meetings{}
			err := json.NewDecoder(r.Body).Decode(&meetings)
			if err != nil{
				panic(err)
			}

			meetings.creation_TS = time.Now().Local();
			userJson, err := json.Marshal(meetings)
			if err != nil{
				panic(err)
			}
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusOK)
			//Write json response back to response 
			w.Write(userJson)
			//write the post data into DB
			writeToDB(meetings);
		}
	});
	
	http.ListenAndServe(":8081",nil);
}