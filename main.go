package main

import (
    "fmt"
	"log"
    "net/http"
	"context"
	"time"
    "encoding/json"
    "strings"
    "io/ioutil"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "strconv"

)

var client *mongo.Client
var MyMap map[string]*joined_participant


type Meeting struct {
	ID string  `json:"ID" bson:"ID"`
    Title string  `json:"Title" bson:"Title"`
    arrayParticipants [] struct{
	Name string `json:"Name" bson:"Name"`
    Email string `json:"Email" bson:"Email"`
    RSVP string `json:"RSVP" bson:"RSVP"`
	}`json:"arrayParticipants" bson:"arrayParticipants"`
	start int `json:"start" bson:"start"`
	end int `json:"end" bson:"end"`
}
type joined_participant struct {
    
    Start []int
    End []int

}

//To test the HTTP Handler
/*func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")

    
    io.WriteString(w, `{"alive": true}`)
}*/

//To add values to map for keeping track of participants
func (data *joined_participant) AppendValues(s int, e int) {
    data.Start = append(data.Start, s)
    data.End = append(data.End, s)
}

//Function to redirect the URL pased on parameters
func appropriate_function(w http.ResponseWriter, r *http.Request){
    key1 := r.URL.Query()["ID"]
    key2 := r.URL.Query()["start"]
    key3 := r.URL.Query()["Email"]
    if key1!= nil{
        find_meeting_id(w,r)
    }else if key2!= nil{
        find_meeting_time(w,r)
    }else if key3!= nil{
        find_meeting_email(w,r)
    }else {
        addMeeting(w,r)
    }
}


//Adding New meeting via push request
func addMeeting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "invalid_http_method")
        return
    }
	
	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    log.Println(string(body))

    if err != nil {
    panic(err)
    }

    //Checking RACE Conditions
    res1 := strings.Split(string(body), "[{") 
    res2 := strings.Split(res1[1], "}]") 
    res3 := strings.ReplaceAll(res2[0], "} , {", ",") 
    res5 := strings.Split(res2[1],":")
    res6 := strings.Split(res5[1],",")
    end_time,err:= strconv.Atoi(strings.Trim(res5[2],"}"))
    start_time,err:= strconv.Atoi(strings.Trim(res6[0]," "))
    res4 := strings.Split(res3, ",") 
    for i, s := range res4 {
        a := strings.Split(s, ":")
        e1:=" "
        if (i+1)%2==0{
            e1 =strings.Trim(a[1],"\"")
        }
        if (i+1)%3==0 {
            
            if a[1]== "\"YES\""{
                obj := &joined_participant{[]int{}, []int{}}
                obj.AppendValues(start_time, end_time)
                MyMap[e1] = obj
                fmt.Fprintf(w,"ADDED WITHOUT COLLISION")
            }


        
        }
        fmt.Println(MyMap)
    }
    
    var m interface{}
    errr := bson.UnmarshalExtJSON([]byte(body),true,&m)
    if errr != nil {
    log.Println(errr)
    }
    log.Println(m)
 
    collection := client.Database("Meetings").Collection("Meets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    result, _ := collection.InsertOne(ctx, m)
    
    json.NewEncoder(w).Encode(result)
}


//To find meetiongs of Participants via Email ID
func find_meeting_email(w http.ResponseWriter, r *http.Request) {
    
    email, ok := r.URL.Query()["Email"]
    
    if !ok || len(email[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return
    }
    

    
    e := email[0]
    
    fmt.Println(e)
   
    collection := client.Database("Meetings").Collection("Meets")
    pipe := bson.M{
        "arrayParticipants": bson.M{
           "Email":  e},}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    showLoadedCursor, err := collection.Aggregate(ctx, pipe)
    if err != nil {
        panic(err)
    }
    var showsLoaded []bson.M
    if err = showLoadedCursor.All(ctx, &showsLoaded); err != nil {
        panic(err)
    }
    fmt.Println(showsLoaded)
    json.NewEncoder(w).Encode(showsLoaded)
}



//To find meetiongs of Participants By time of the meeting
func find_meeting_time(w http.ResponseWriter, r *http.Request) {
    
    start, ok := r.URL.Query()["start"]
    
    if !ok || len(start[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return
    }
    end, ok := r.URL.Query()["end"]
    
    if !ok || len(end[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return
    }

    s := start[0]
    e := end[0]
    fmt.Println(s)
    fmt.Println(e)
    st, err := strconv.Atoi(s)
   if err != nil {
      log.Println(err)
      }
    ed, err := strconv.Atoi(e)
   if err != nil {
    log.Println(err)
   }
	collection := client.Database("Meetings").Collection("Meets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    cursor, err := collection.Find(ctx, bson.M{
        "$and": bson.A{ // you can try this in []interface
            bson.M{"start": bson.M{
                "$gt": st,},
            },
            bson.M{"end": bson.M{
                "$lt": ed,
            },
        },
    },
},
    )
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    
    defer cursor.Close(ctx)
    var jsonDocuments []map[string]interface{}
    var bsonDocument bson.D
    var jsonDocument map[string]interface{}
    var temporaryBytes []byte
    defer cursor.Close(ctx)
	for cursor.Next(ctx) {
       
		err = cursor.Decode(&bsonDocument)

    temporaryBytes, err = bson.MarshalExtJSON(bsonDocument, true, true)

    err = json.Unmarshal(temporaryBytes, &jsonDocument)

    jsonDocuments = append(jsonDocuments, jsonDocument)
    fmt.Println(jsonDocuments)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(jsonDocuments)
}

////To find meetiongs of Participants by ID
func find_meeting_id(w http.ResponseWriter, r *http.Request) {
    
    keys, ok := r.URL.Query()["ID"]
    
    if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        return
    }

    key := keys[0]
   
	collection := client.Database("Meetings").Collection("Meets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    cursor, err := collection.Find(ctx, bson.M{"ID":key})
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    var jsonDocuments []map[string]interface{}
    var bsonDocument bson.D
    var jsonDocument map[string]interface{}
    var temporaryBytes []byte
    defer cursor.Close(ctx)
	for cursor.Next(ctx) {
       
		err = cursor.Decode(&bsonDocument)

    temporaryBytes, err = bson.MarshalExtJSON(bsonDocument, true, true)

    err = json.Unmarshal(temporaryBytes, &jsonDocument)

    jsonDocuments = append(jsonDocuments, jsonDocument)
    fmt.Println(jsonDocuments)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(jsonDocuments)
    
}
func main() {
    fmt.Println("Starting the application...")
    MyMap = make(map[string]*joined_participant)
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
    http.HandleFunc("/meeting/", appropriate_function)
	http.ListenAndServe(":8000", nil)
}