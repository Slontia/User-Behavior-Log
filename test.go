package main

import ( 
    "gopkg.in/mgo.v2" 
    "gopkg.in/mgo.v2/bson"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "crypto/md5"
    "github.com/gorilla/schema"
    "io/ioutil"
    "net/http"
    "html/template"
    "fmt"
    "time"
    "log"
    //"bytes"
    //"strings"
    "sort"
    "reflect"
    "strconv"
)

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var coder = base64.NewEncoding(base64Table)

func base64Encode(encode_byte []byte) []byte {
    return []byte(coder.EncodeToString(encode_byte))
}

func base64Decode(decode_byte []byte) ([]byte, error) {
    return coder.DecodeString(string(decode_byte))
}

func md5Encode(str string) string {
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(str))
    cipherStr := md5Ctx.Sum(nil)
    return hex.EncodeToString(cipherStr)
}

// 用户登录注销信息
type SignLog struct {
    Id          bson.ObjectId   `bson:"_id"`            // Id
    UserId      string          `bson:"user_id"`        // 用户id
    LoginTime   string          `bson:"Login_time"`     // 登录时间
    LogoutTime  string          `bson:"Logout_time"`    // 登出时间
} 

// 代码行数信息
type LinesLog struct {
    Id          bson.ObjectId   `bson:"_id"`        // Id
    UserId      string          `bson:"user_id"`    // 用户id
    LineNum     int          	`bson:"line_num"`   // 代码行数
    Time        string          `bson:"time"`       // 时间戳
    SignId      bson.ObjectId          `bson:"sign_id"`    // 对应的sign信息id
}

// Debug信息
type DebugLog struct {
    Id          bson.ObjectId   `bson:"_id"`        // Id
    UserId      string          `bson:"user_id"`    // 用户id
    BeginTime   string          `bson:"b_time"`     // 开始时间
    EndTime     string          `bson:"e_time"`     // 结束时间
    SignId      bson.ObjectId   `bson:"sign_id"`    // 对应的sign信息id   
}

// 运行情况信息
type RunLog struct {
    Id          bson.ObjectId   `bson:"_id"`        // Id
    UserId      string          `bson:"user_id"`    // 用户id
    BeginTime   string          `bson:"b_time"`     // 开始时间
    EndTime     string          `bson:"e_time"`     // 结束时间
    SignId      bson.ObjectId   `bson:"sign_id"`    // 对应的sign信息id   
}


const URL = "127.0.0.1:27017"
var (
    mgoSession *mgo.Session
    dataBase   = "oo_test"
)

func getSession() *mgo.Session {
    if mgoSession == nil {
        var err error
        mgoSession, err = mgo.Dial(URL)
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

func GetSignLog(userId string) []SignLog {
    var ols []SignLog
    // find logs
    query := func(c *mgo.Collection) error {
        return c.Find(bson.M{"user_id": userId, "Logout_time": ""}).All(&ols)
    }
    witchCollection("sign", query)
    return ols
}

func getDebugLog(userId string) []DebugLog {
	var dbgs []DebugLog
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"user_id": userId, "e_time": ""}).All(&dbgs)
	}
	witchCollection("debug", query)
	return dbgs
}

func getRunLog(userId string) []RunLog {
	var runs []RunLog
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"user_id": userId, "e_time": ""}).All(&runs)
	}
	witchCollection("run", query)
	return runs
}

/*===========================
|            API            |
============================*/

type dbAPI struct {} 


// 用户登录
func (dbAPI) Login(data map[string]string) string {
	userId := data["userId"]
	ols := GetSignLog(userId)
	if (len(ols) > 0) {
		return "Logined"
	}
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    fmt.Println(time)
    ol := SignLog{id, userId, time, ""}
    // define 
    query := func(c *mgo.Collection) error {
        return c.Insert(&ol)
    }
    err := witchCollection("sign", query)
    if err != nil {
        return "failed"
    }
    return ol.Id.Hex()
}

// 用户登出
func (dbAPI) Logout(data map[string]string) int {
    userId := data["userId"]
    ols := GetSignLog(userId)
    if (len(ols) == 0) {    // user offline
        return -1
    }
    // add Logout time
    ols[0].LogoutTime = time.Now().Format("2006-01-02 15:04:05")
    query := func(c *mgo.Collection) error {
        return c.Update(bson.M{"Logout_time": ""}, &(ols[0]))
    }
    err := witchCollection("sign", query)  
    if (err != nil) {
        return -2
    }
    return 0
}

// 记录当前代码行数
func (dbAPI) RecordLine(data map[string]string) string {
    userId := data["userId"]
    lineNum, err := strconv.Atoi(data["lineNum"])
    if (err != nil) {
    	return "type error"
    }
    ols := GetSignLog(userId)
    if (len(ols) == 0) {    // user offline
        return "offline"
    }
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    linesLog := LinesLog{id, userId, lineNum, time, ols[0].Id}
    query := func(c *mgo.Collection) error {
        return c.Insert(&linesLog)
    }
    err = witchCollection("lines", query) 
    if (err != nil) {
        return "failed"
    }
    return linesLog.Id.Hex()
}

// debug开始
func (dbAPI) DebugBegin(data map[string]string) string {
    userId := data["userId"]
    ols := GetSignLog(userId)
    if (len(ols) == 0) {
        return "offline"
    }
    dbgs := getDebugLog(userId)
    if (len(dbgs) > 0) {
    	return "debugging"
    }
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    debugLog := DebugLog{id, userId, time, "", ols[0].Id}
    query := func(c *mgo.Collection) error {
        return c.Insert(&debugLog)
    }
    err := witchCollection("debug", query)
    if (err != nil) {
        return "failed"
    }
    return debugLog.Id.Hex()
}

// debug结束
func (dbAPI) DebugOver(data map[string]string) int {
	userId := data["userId"]
    debugs := getDebugLog(userId)
    if (len(debugs) == 0) {
        return -1
    }
    debugs[0].EndTime = time.Now().Format("2006-01-02 15:04:05")
    query := func(c *mgo.Collection) error {
        return c.Update(bson.M{"user_id": userId, "e_time": ""}, &debugs[0])    
    }
    err := witchCollection("debug", query)
    if (err != nil) {
        return -2
    }
    return 0
}

// 运行开始
func (dbAPI) RunBegin(data map[string]string) string {
	userId := data["userId"]
    ols := GetSignLog(userId)
    if (len(ols) == 0) {
        return "offline"
    }
    runs := getRunLog(userId)
    if (len(runs) > 0) {
    	return "running"
    }
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    runLog := RunLog{id, userId, time, "", ols[0].Id}
    query := func(c *mgo.Collection) error {
        return c.Insert(&runLog)
    }
    err := witchCollection("run", query)
    if (err != nil) {
        return "null"
    }
    return runLog.Id.Hex()
}

// 运行结束
func (dbAPI) RunOver(data map[string]string) int {
	userId := data["userId"]
    runs := getRunLog(userId)
    if (len(runs) == 0) {
        return -1	// not running
    }
    runs[0].EndTime = time.Now().Format("2006-01-02 15:04:05")
    query := func(c *mgo.Collection) error {
        return c.Update(bson.M{"user_id": userId, "e_time": ""}, &runs[0])    
    }
    err := witchCollection("run", query)
    if (err != nil) {
        return -2	// failed
    }
    return 0
}


/* 以下为测试用的函数 */

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

func removeAllCol() {
    removeAll("sign")
    removeAll("lines")
    removeAll("debug")
    removeAll("run")
}

/*===========================
|          Entry            |
============================*/

func GetSign(request Request) string {
	dataMap := request.Data
	timestp := request.Timestp
	api := request.Api

	keySlice := make([]string, len(dataMap))
	
	// Sort key
	i := 0
	for k, _ := range dataMap {
		keySlice[i] = k
		i++
	}
	sort.Strings(keySlice)
	
	dataStr := ""
	// create data string
	for _, key := range keySlice {
		value := dataMap[key]
		dataStr += key + value
	}
	dataStr += "timestp" + strconv.FormatInt(timestp, 10) + "api" + string(api)
	fmt.Println(dataStr)

	return md5Encode(dataStr)
}

type Request struct {
	Api			string               `json:"api,omitempty"`
	Timestp		int64                 `json:"timestp,string,omitempty"`
	Sign        string                 `json:"sign,omitempty"`
	Data		map[string]string        `json:"data,string,omitempty"`
}

type Request_Send struct {
	Api			string
	Timestp		int64
	Sign        string
	Data		map[string]interface{}
}



// Interface Entry
func entry(response http.ResponseWriter, req *http.Request) {
	// unmarshal json
	request := &Request{}
    resBody, err := ioutil.ReadAll(req.Body)
    //reqByte := []byte("")
    if err != nil {
        fmt.Println("ReadAll Failed")
        return
    }
    //fmt.Println(string(resBody))
    //fmt.Println(reflect.TypeOf(resBody).String())

    //reqByte := bytes.NewBuffer(result)
	err = json.Unmarshal([]byte(resBody), &request)
	if (err != nil) {
        log.Printf("error decoding sakura response: %v", err)
        if e, ok := err.(*json.SyntaxError); ok {
            log.Printf("syntax error at byte offset %d", e.Offset)
        }
        log.Printf("sakura response: %q", resBody)
		return
	}
    //log.Printf("sakura response: %q", resBody)
	// sign check
	sign := GetSign(*request)
	fmt.Println("Got: " + request.Sign)
	fmt.Println("Act: " + sign)
	curTimestp := time.Now().Unix()
	//fmt.Println(request.Data)
	if (curTimestp - request.Timestp >= 20) {
		fmt.Println("Over Time!")
		//return
	}
	if (request.Sign != sign) {
		fmt.Println("Invalid!")
		//return
	}

	// call interface with reflection
	var dbapi dbAPI
    //fmt.Println(111)
	v := reflect.ValueOf(&dbapi)
    //fmt.Println(222)
	args := []reflect.Value{ reflect.ValueOf(request.Data) }
	//fmt.Println(333)
    out := v.MethodByName(request.Api).Call(args)
    //fmt.Println(444)
    for _, v := range out {
        fmt.Println(v)
    }
    //fmt.Println(555)
}


func Hello(response http.ResponseWriter, request *http.Request) {
    type person struct {
        Id      int
        Name    string
        Country string
    }
    fmt.Println("Yes")

    liumiaocn := person{Id: 1001, Name: "liumiaocn", Country: "China"}

    tmpl, err := template.ParseFiles("./user.tpl")
    if err != nil {
            fmt.Println("Error happened..")
    }
    tmpl.Execute(response, liumiaocn)
}

type BaseJsonBean struct {  
    Code    int         `json:"code"`  
    Data    interface{} `json:"data"`  
    Message string      `json:"message"`  
}  
  
func NewBaseJsonBean() *BaseJsonBean {  
    return &BaseJsonBean{}  
}

func AjaxTest(response http.ResponseWriter, req *http.Request) {
    fmt.Println(999)
    err := req.ParseForm()
    if err != nil {
        fmt.Println("解析失败")
    }
    param_api, found1 := req.Form["api"]
    param_timestp, found2 := req.Form["timestp"]
    param_sign, found3 := req.Form["sign"]
    param_data, found4 := req.Form["data"]
    var decoder = schema.NewDecoder()
    var request Request
    err = decoder.Decode(&request, req.PostForm)
    if err != nil {
        fmt.Println("decode失败")
        fmt.Println(err)
    }
    //param_data := req.FormValue("data")
    if !(found1 && found2 && found3 && found4) {
        fmt.Println(found1)
        fmt.Println(found2)
        fmt.Println(found3)
        fmt.Println(found4)
        fmt.Println("Error")
        return
    }
    result := NewBaseJsonBean()
    api := param_api[0]
    timestp := param_timestp[0]
    sign := param_sign[0]
    data := param_data[0]
    //data := strings.Split(param_data[0], ",")

    s := "api: " + api + ", timestp: " + timestp + ", sign: " + sign
    fmt.Println("info: " + s)
    fmt.Println("data: ", data)

    result.Code = 12450
    result.Message = "屌你老母哦❤" 

    bytes, _ := json.Marshal(result)
    fmt.Println(string(bytes))
}

func main() {
	/*data := make(map[string]interface{})
	data["userId"] = "001"
	data["lineNum"] = "50"
	r := &Request_Send{
		"RecordLine", 
		int64(1432710115), 
		"000",    
		data}

	jsonByte, err := json.Marshal(r)
	if (err != nil) {
		fmt.Println("ERROR!")
		return
	}
	entry(jsonByte)
    */
	http.HandleFunc("/what/", Hello)
    http.HandleFunc("/ajax/", AjaxTest)
    http.HandleFunc("/entry/", entry)
	http.ListenAndServe("127.0.0.1:8080", nil)
    //removeAllCol()
 }