package main


import (
    "encoding/json"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/binary"
    "fmt"
    "bytes"
    "strings"
    "gopkg.in/mgo.v2"
    "os"
    "gopkg.in/mgo.v2/bson"
)
var i int=0
var j int=0
var b []byte
var t input
var c coordinate
var generic Generic
var msg,msguser outputMongo
var collection *mgo.Collection

type input struct {
    Name string
    Address string
    City string
    State string
    Zip string
}

type coordinate struct {
    Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
}

type outputMongo struct {
    ID  bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Name  string `json:"name"`
    Address     string `json:"address"`
    City        string `json:"city"`
    State string `json:"state"`
    Zip   string `json:"zip"`
    Coordinates struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    } `json:"coordinates"`
       
}

type Generic struct{
       Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PartialMatch bool     `json:"partial_match"`
        PlaceID      string   `json:"place_id"`
        Types        []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}

//Indenting into proper JSON format
func prettyprint(b []byte) ([]byte, error) {
    var out bytes.Buffer
    err := json.Indent(&out, b, "", "  ")
    return out.Bytes(), err
}

func connectdb(){
    uri := "mongodb://dbuser:dbuser@ds053090.mongolab.com:53090/location"
  if uri == "" {
    fmt.Println("no connection string provided")
    os.Exit(1)
  }
 
  sess, err := mgo.Dial(uri)
  if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
  }
  sess.SetSafe(&mgo.Safe{})
  
  collection = sess.DB("location").C("address")
}

func posting(rw http.ResponseWriter, req *http.Request) {
    
    if(req.Method=="POST"){//------------------------> POST---------------------------------->
        
    body, err := ioutil.ReadAll(req.Body)
    //fmt.Println(i,t[i].Name, t[i].Address)
    err = json.Unmarshal(body, &t)
    if err != nil {
        panic("run after unmarshal")
    }
    Addr:=strings.Split(t.Address," ")
    City1:=strings.Split(t.City," ")
    
    api:="http://maps.google.com/maps/api/geocode/json?address="+Addr[0]+"+"+Addr[1]+"+"+Addr[2]+","+City1[0]+"+"+City1[1]+","+t.State+"+"+t.Zip+"&sensor=false"
    link,err:=http.Get(api)
    if err!=nil{
        log.Fatal("there is an error")
        }
    
        head, _ := ioutil.ReadAll(link.Body) 
        err = json.Unmarshal(head, &generic) // unmarshalling parameters
        if err != nil {
            panic("error in unmarshalling api")
        }
    
    c.Longitude=generic.Results[0].Geometry.Location.Lat
    c.Latitude=generic.Results[0].Geometry.Location.Lng
    msg=outputMongo{bson.NewObjectId(),t.Name,t.Address,t.City,t.State,t.Zip,c}
    fmt.Println(msg)
    b,_=json.Marshal(msg)   
    bx, _ := prettyprint(b)
    n:=binary.Size(bx)
    s := string(bx[:n])
    fmt.Fprintf(rw,s)
    
  err = collection.Insert(msg)
  if err != nil {
    fmt.Printf("Can't insert document: %v\n", err)
    }
    i=i+1
    j=j+1
    
}else if(req.Method=="GET"){//---------------------------> GET------------------------------------>
    
    s1:=req.URL.Path[1:]
    st1:=string(s1[10:])
    oid := bson.ObjectIdHex(st1)
  
    err := collection.FindId(oid).One(&msg)
if err != nil {
    fmt.Printf("Error in searching! %v\n", err)
    os.Exit(1)
  }
   
    b,_=json.Marshal(msg)   
    bx, _ := prettyprint(b)
    n:=binary.Size(bx)
    s := string(bx[:n])
    fmt.Fprintf(rw,s)
     

  }else if(req.Method=="PUT"){//------------------------> PUT---------------------------------->
   
     s1:=req.URL.Path[1:]
     st1:=string(s1[10:])
      oid := bson.ObjectIdHex(st1)
  fmt.Println(oid)
    err := collection.FindId(oid).One(&msg)
if err != nil {
    fmt.Printf("Error in searching! %v\n", err)
    os.Exit(1)
  }
  body, _ := ioutil.ReadAll(req.Body)
    
    err = json.Unmarshal(body, &t)
    if err != nil {
        panic("run after unmarshal")
    }
    Addr:=strings.Split(t.Address," ")
    City1:=strings.Split(t.City," ")
    
    api:="http://maps.google.com/maps/api/geocode/json?address="+Addr[0]+"+"+Addr[1]+"+"+Addr[2]+","+City1[0]+"+"+City1[1]+","+t.State+"+"+t.Zip+"&sensor=false"
    link,_:=http.Get(api)
    if err!=nil{
        log.Fatal("there is an error")
        }
    
        head, _ := ioutil.ReadAll(link.Body) 
        err = json.Unmarshal(head, &generic) // unmarshalling parameters
        if err != nil {
            panic("error in unmarshalling api")
        }
    
    c.Longitude=generic.Results[0].Geometry.Location.Lat
    c.Latitude=generic.Results[0].Geometry.Location.Lng
    msguser=outputMongo{oid,t.Name,t.Address,t.City,t.State,t.Zip,c}
    msg.Name=msguser.Name
    msg.Address=msguser.Address
    msg.City=msguser.City
    msg.State=msguser.State
    msg.Zip=msguser.Zip
    msg.Coordinates=msguser.Coordinates
    err = collection.UpdateId(oid,msg)
    if err != nil {
    fmt.Printf("Error in searching! %v\n", err)
    os.Exit(1)
  }
  err = collection.FindId(oid).One(&msg)
  if err != nil {
    fmt.Printf("Error in searching! %v\n", err)
    os.Exit(1)
  }
   
    b,_=json.Marshal(msg)   
    bx, _ := prettyprint(b)
    n:=binary.Size(bx)
    s := string(bx[:n])
    fmt.Fprintf(rw,s)
     
   
}else if(req.Method=="DELETE"){//---------------------------> DELETE----------------------------------->
    s1:=req.URL.Path[1:]
    st1:=string(s1[10:])
    
  oid := bson.ObjectIdHex(st1)
  //fmt.Println(oid)
    err := collection.RemoveId(oid)
if err != nil {
    fmt.Printf("Error in searching! %v\n", err)
    os.Exit(1)
        }
    }
}
func main() {
   
    http.HandleFunc("/locations", posting)
    http.HandleFunc("/locations/",posting)
    connectdb()
    log.Fatal(http.ListenAndServe(":8082", nil))
}

/*
{
    "Name" : "namitha ninan",
    "Address" : "101 san fernando",
    "City" : "san jose",
    "State" : "CA",
    "Zip" : "95112"
}
*/