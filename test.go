package main

import ( 
    "encoding/base64"
    "fmt"
    "gopkg.in/mgo.v2" 
    "gopkg.in/mgo.v2/bson"
    "time"
    "sort"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
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
	Api			string
	Timestp		int64
	Sign        string
	Data		map[string]string
}

type Request_Send struct {
	Api			string
	Timestp		int64
	Sign        string
	Data		map[string]interface{}
}

// Interface Entry
func entry(reqByte []byte) {
	// unmarshal json
	request := &Request{}
	err := json.Unmarshal(reqByte, &request)
	if (err != nil) {
		fmt.Println("Unmarshal failed!")
		return
	}

	// sign check
	sign := GetSign(*request)
	fmt.Println("Got: " + request.Sign)
	fmt.Println("Act: " + sign)
	curTimestp := time.Now().Unix()
	fmt.Println(curTimestp)
	if (curTimestp - request.Timestp >= 20) {
		fmt.Println("Over Time!")
		return
	}
	if (request.Sign != sign) {
		fmt.Println("Invalid!")
		return
	}

	// call interface with reflection
	var dbapi dbAPI
	v := reflect.ValueOf(&dbapi)
	args := []reflect.Value{ reflect.ValueOf(request.Data) }
	out := v.MethodByName(request.Api).Call(args)
    for _, v := range out {
        fmt.Println(v)
    }
}


func main() {
	data := make(map[string]interface{})
	data["userId"] = "001"
	data["lineNum"] = "50"
	r := &Request_Send{
		"Login", 
		int64(1432710115), 
		"000", 
		data}

	jsonByte, err := json.Marshal(r)
	if (err != nil) {
		fmt.Println("ERROR!")
		return
	}
	entry(jsonByte)
    //removeAllCol()
 }