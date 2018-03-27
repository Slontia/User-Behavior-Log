package main

import ( 
    "gopkg.in/mgo.v2" 
    "gopkg.in/mgo.v2/bson"
    //"encoding/base64"
    "encoding/hex"
    "encoding/json"
    "crypto/md5"
    //"github.com/gorilla/schema"
    "io"
    "io/ioutil"
    "net/http"
    "html/template"
    //"fmt"
    "time"
    "log"
    //"bytes"
    //"strings"
    //"sort"
    //"reflect"
    "strconv"
)

// database
const DATABASE_ADDR = "127.0.0.1:27017"
const DATABASE_NAME = "oo_test"
const COLLECTION_NAME = "behavior"

// frontend
const FRONT_ADDR = "127.0.0.1:8080"
const TEMPLATE_PATH = "./user.tpl"


func md5Encode(str string) string {
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(str))
    cipherStr := md5Ctx.Sum(nil)
    return hex.EncodeToString(cipherStr)
}

/*============================
|          Database          |
============================*/

var (
    mgoSession *mgo.Session
    dataBase   = DATABASE_NAME
)

func getSession() *mgo.Session {
    if mgoSession == nil {
        var err error
        mgoSession, err = mgo.Dial(DATABASE_ADDR)
        if err != nil {
            panic(err)
        }
    }
    return mgoSession.Clone()
}

func witchCollection(collection string, s func(*mgo.Collection) error) error {
    session := getSession()
    defer session.Close()
    c := session.DB(dataBase).C(collection)
    return s(c)
}

func removeAll(colName string) string {
    query := func(c *mgo.Collection) error {
        c.RemoveAll(nil)
        return nil
    }
    err := witchCollection(colName, query)
    if (err != nil) {
        return "false"
    }
    return "true"
}

/*============================
|          Request           |
============================*/

type Request struct {
	Timestp		int64     `json:"time,string"`
	Sign        string    `json:"sign,omitempty"`
	Action		string    `json:"action,omitempty"`
	Id			string	  `json:"id,omitempty"`
}

type UserLog struct {
	Time		int64
	Action		string   
	Id			string
}

func GetSign(request Request) string {
	id := request.Id
	action := request.Action
	timestp := request.Timestp
	dataStr := "id" + id
	dataStr += "action" + action
	dataStr += "time" + strconv.FormatInt(timestp, 10)
	log.Printf("Discode Sign: " + dataStr)
	return md5Encode(dataStr)
}

// Interface Entry
func entry(res http.ResponseWriter, req *http.Request) {
	// unmarshal json
	request := &Request{}
    resBody, err := ioutil.ReadAll(req.Body)
    if err != nil {
        io.WriteString(res, "Unreadable Request")
        return
    }
	err = json.Unmarshal([]byte(resBody), &request)
	if (err != nil) {
        log.Printf("error decoding sakura response: %v", err)
        if e, ok := err.(*json.SyntaxError); ok {
            log.Printf("syntax error at byte offset %d", e.Offset)
        }
        log.Printf("sakura response: %q", resBody)
        io.WriteString(res, "JSON Decoding Error")
		return
	}
	// sign check
	sign := GetSign(*request)
	curTimestp := time.Now().Unix()
	if (curTimestp - request.Timestp >= 20) {
		log.Printf("Overdue Request!")
		log.Printf("Send Time: %d", request.Timestp)
		log.Printf("Current Time: %d", curTimestp)
		io.WriteString(res, "Overdue Request")
		return
	}
	if (request.Sign != sign) {
		log.Printf("Sign Dismatch!")
		log.Printf("Got Sign: " + request.Sign)
		log.Printf("Expected Sign: " + sign)
		io.WriteString(res, "Sign Dismatch")
		return
	}
	// write to database
	ulog := &UserLog{request.Timestp, request.Action, request.Id}
    query := func(c *mgo.Collection) error {
        return c.Insert(ulog)
    }
    err = witchCollection(COLLECTION_NAME, query) 
    if (err != nil) {
    	io.WriteString(res, "Write Database Error")
        return
    }
    log.Printf("Success!")
}

func findAll(res http.ResponseWriter, req *http.Request) {
    var logs []UserLog
    query := func(c *mgo.Collection) error {
        return c.Find(bson.M{}).All(&logs)
    }
    err := witchCollection(COLLECTION_NAME, query)
    if (err != nil) {
        io.WriteString(res, "Cannot Read Database")
        return
    }
    jsons, errs := json.Marshal(logs)
    if errs != nil {  
      log.Println(errs.Error())  
    }
    log.Printf(string(jsons))
    res.Write(jsons)
}

// load template
func index(response http.ResponseWriter, request *http.Request) {
    tmpl, err := template.ParseFiles(TEMPLATE_PATH)
    if err != nil {
    	log.Printf("Template Loading Error")
    }
    tmpl.Execute(response, nil)
}

// route
func main() {
	http.HandleFunc("/index/", index)
    http.HandleFunc("/entry/", entry)
    http.HandleFunc("/show/", findAll)
	http.ListenAndServe(FRONT_ADDR, nil)
 }