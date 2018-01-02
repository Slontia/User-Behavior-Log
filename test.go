package main
import ( 
    "encoding/base64"
    "fmt"
    "gopkg.in/mgo.v2" 
    "gopkg.in/mgo.v2/bson"
    "time"
    "crypto/md5"
    "encoding/hex"
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

func getSign(dataMap map[string]string, appKey string) string {
	keySlice := make([]string, len(dataMap))
	
	// Sort key
	i := 0
	for k, _ := range dataMap {
		keySlice[i] = k
		i++
	}	
	sort.Strings(keySlice)
	
	// create data string
	dataStr := ""
	for _, key := range keySlice {
		dataStr += key + dataMap[key]
	}
	dataStr += appKey

	return md5Encode(dataStr)
}

// 用户登录注销信息
type SignLog struct {
    Id          bson.ObjectId   `bson:"_id"`            // Id
    UserId      string          `bson:"user_id"`        // 用户id
    LoginTime   string          `bson:"login_time"`     // 登录时间
    LogoutTime  string          `bson:"logout_time"`    // 登出时间
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

func getSignLog(userId string) []SignLog {
    var ols []SignLog
    // find logs
    query := func(c *mgo.Collection) error {
        return c.Find(bson.M{"user_id": userId, "logout_time": ""}).All(&ols)
    }
    witchCollection("sign", query)
    return ols
}

/*===========================
|            API            |
============================*/

// Interface Entry


// 用户登录
func login(userId string) string {
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
        return "null"
    }
    return ol.Id.Hex()
}

// 用户登出
func logout(userId string) int {
    ols := getSignLog(userId)
    if (len(ols) == 0) {    // user offline
        return -1
    }
    // add logout time
    ols[0].LogoutTime = time.Now().Format("2006-01-02 15:04:05")
    query := func(c *mgo.Collection) error {
        return c.Update(bson.M{"logout_time": ""}, &(ols[0]))
    }
    err := witchCollection("sign", query)  
    if (err != nil) {
        return -2
    }
    return 0
}

// 记录当前代码行数
func recordLine(lineNum int, userId string) string {
    ols := getSignLog(userId)
    if (len(ols) == 0) {    // user offline
        return "null"
    }
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    linesLog := LinesLog{id, userId, lineNum, time, ols[0].Id}
    query := func(c *mgo.Collection) error {
        return c.Insert(&linesLog)
    }
    err := witchCollection("lines", query) 
    if (err != nil) {
        return "null"
    }
    return linesLog.Id.Hex()
}

// debug开始
func debug_begin(userId string) string {
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    ols := getSignLog(userId)
    if (len(ols) == 0) {
        return "null"
    }
    debugLog := DebugLog{id, userId, time, "", ols[0].Id}
    query := func(c *mgo.Collection) error {
        return c.Insert(&debugLog)
    }
    err := witchCollection("debug", query)
    if (err != nil) {
        return "null"
    }
    return debugLog.Id.Hex()
}

// debug结束
func debug_over(userId string) int {
    var debugs []DebugLog
    query := func(c *mgo.Collection) error {
        return c.Find(bson.M{"user_id": userId, "e_time": ""}).All(&debugs)       
    }
    err := witchCollection("debug", query)
    if (len(debugs) == 0) {
        return -1
    }
    debugs[0].EndTime = time.Now().Format("2006-01-02 15:04:05")
    query = func(c *mgo.Collection) error {
        return c.Update(bson.M{"user_id": userId, "e_time": ""}, &debugs[0])    
    }
    err = witchCollection("debug", query)
    if (err != nil) {
        return -2
    }
    return 0
}

// 运行开始
func run_begin(userId string) string {
    id := bson.NewObjectId()
    time := time.Now().Format("2006-01-02 15:04:05")
    ols := getSignLog(userId)
    if (len(ols) == 0) {
        return "null"
    }
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
func run_over(userId string) int {
    var runs []RunLog
    query := func(c *mgo.Collection) error {
        return c.Find(bson.M{"user_id": userId, "e_time": ""}).All(&runs)       
    }
    err := witchCollection("run", query)
    if (len(runs) == 0) {
        return -1
    }
    runs[0].EndTime = time.Now().Format("2006-01-02 15:04:05")
    query = func(c *mgo.Collection) error {
        return c.Update(bson.M{"user_id": userId, "e_time": ""}, &runs[0])    
    }
    err = witchCollection("run", query)
    if (err != nil) {
        return -2
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

func main() {
    //removeAllCol()
    //fmt.Println("HW")
    //fmt.Println(login("999"))
    //fmt.Println(recordLine(10, "999"))
    //fmt.Println(run_begin("999"))
    //fmt.Println(debug_begin("999"))
    //fmt.Println(debug_over("999"))
    //fmt.Println(run_over("999"))
    //fmt.Println(logout("999"))
 }