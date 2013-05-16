//核心方法,加解密字符串
package discuz

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var UCenterUrl string
var Appid string
var uckey string
var API_RETURN_SUCCEED int = 1

//注册AppID
//id:AppID
//k:Key
//url:UCenter所在的URL
//listen:例如uc,那么在ucenter中的应用接口名称就必须是uc
func Register(id, k, url string) {
	Appid = id
	uckey = k
	UCenterUrl = url
	//激活当前目标的所有代理
}

type Handler struct {
}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server(w, r)
}
func DiscuzHandler(w http.ResponseWriter, r *http.Request) {
	server(w, r)
}
func server(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	if v := r.FormValue("code"); v != "" {
		ac := DiscuzDecode(v)
		val, err := url.ParseQuery(ac)
		if err == nil {
			switch val.Get("action") {
			case "test":
				{
					fmt.Fprint(w, API_RETURN_SUCCEED)
				}
			case "synlogin":
				{
					if bindLogin != nil {
						v, err := strconv.Atoi(val.Get("uid"))
						if err == nil {
							bindLogin(v, w, r)
						}
					}
				}
			case "synlogout":
				{
					if bindLogout != nil {
						bindLogout(w)
					}
				}
			case "deleteuser":
				{
					if bindDeleteUser != nil {
						if bindDeleteUser(strings.Split(val.Get("ids"), ",")) {
							fmt.Fprint(w, API_RETURN_SUCCEED)
						}
					}
				}
			case "renameuser":
				{
					if bindRenameUser != nil {
						v, err := strconv.Atoi(val.Get("uid"))
						if err == nil {
							if bindRenameUser(v, val.Get("oldusername"), val.Get("newusername")) {
								fmt.Fprint(w, API_RETURN_SUCCEED)
							}
						}
					}
				}
			case "updatepw":
				{
					if bindUpdatepw != nil {
						if bindUpdatepw(val.Get("username"), val.Get("password")) {
							fmt.Fprint(w, API_RETURN_SUCCEED)
						}
					}
				}
			case "gettag":
				{
					if bindGettag != nil {
						v, err := strconv.Atoi(val.Get("id"))
						if err == nil {
							data, b := bindGettag(v)
							if b {

								be, err := xml.Marshal(data)
								if err == nil {
									fmt.Fprint(w, be)
								}
							}
						}
					}
				}
			case "updatecredit":
				{
					if bindUpdatecredit != nil {
						v, err := strconv.Atoi(val.Get("uid"))
						credit, err := strconv.Atoi(val.Get("credit"))
						amount, err := strconv.Atoi(val.Get("amount"))
						if err == nil {
							if bindUpdatecredit(v, credit, amount) {
								fmt.Fprint(w, API_RETURN_SUCCEED)
							}
						}
					}
				}
			case "getcreditsettings":
				{
					if bindGetcreditsettings != nil {
						arr, b := bindGetcreditsettings()
						if b {
							st := XmlSerialize(arr, false, 1)
							fmt.Fprint(w, st)
						}
					}
				}
			case "updatecreditsettings":
				{
					if bindUpdatecreditsettings != nil {
						st := val.Get("credit")
						var kk map[int]CreditSettings
						err := xml.Unmarshal([]byte(st), &kk)
						if err == nil {
							if bindUpdatecreditsettings(kk) {
								fmt.Fprint(w, API_RETURN_SUCCEED)
							}
						}
					}
				}
			case "getcredit":
				{
					if bindGetcredit != nil {
						uid, err := strconv.Atoi(val.Get("uid"))
						credit, err := strconv.Atoi(val.Get("credit"))
						if err == nil {
							i := bindGetcredit(uid, credit)
							fmt.Fprint(w, i)
						}
					}
				}
			}
		}
	}

}

var bindGetcredit func(uid, credit int) int

func BindGetcredit(sum func(uid, credit int) int) {
	bindGetcredit = sum
}

type CreditSettings struct {
	//积分兑换的目标应用程序 ID
	Appiddesc int `xml:"appiddesc"`
	//积分兑换的目标积分编号
	Creditdesc int `xml:"creditdesc"`
	//积分兑换的源积分编号
	Creditsrc int `xml:"creditsrc"`
	//积分名称
	Title string `xml:"title"`
	//积分单位
	Unit string `xml:"unit"`
	//积分兑换比率
	Ratio int `xml:"ratio"`
}

var bindUpdatecreditsettings func(mp map[int]CreditSettings) bool

func BindUpdatecreditsettings(sum func(mp map[int]CreditSettings) bool) {
	bindUpdatecreditsettings = sum
}

var bindGetcreditsettings func() (arr [][]string, b bool)

func BindGetcreditsettings(sum func() (arr [][]string, b bool)) {
	bindGetcreditsettings = sum
}

var bindUpdatecredit func(uid, credit, amount int) bool

//uid:用户ID
//credit:积分编号,
//amount:数量
func BindUpdatecredit(sum func(uid, credit, amount int) bool) {
	bindUpdatecredit = sum
}

//TODO;updatebadwords,updatehosts,updateapps,updateclient,暂未实现

type Tag struct {
	Name     string `xml:"name"`
	Uid      int    `xml:"uid"`
	Username string `xml:"username"`
	Dateline string `xml:"dateline"`
	Url      string `xml:"url"`
	Image    string `xml:"image"`
}

var bindGettag func(id int) (tag []Tag, b bool)

func BindGettag(sum func(id int) (tag []Tag, b bool)) {
	bindGettag = sum
}

var bindUpdatepw func(username, passwd string) bool

func BindUpdatepw(sum func(username, passwd string) bool) {
	bindUpdatepw = sum
}

var bindRenameUser func(userId int, oldname, newname string) bool

func BindRenameUser(sum func(userId int, oldname, newname string) bool) {
	bindRenameUser = sum
}

var bindDeleteUser func(userId []string) bool

//注册删除用户方法
func BindDeleteUser(sum func(userId []string) bool) {
	bindDeleteUser = sum
}

var bindLogin func(userId int, w http.ResponseWriter, r *http.Request) bool

//注册同步登陆处理方法
func BindLogin(sum func(userId int, w http.ResponseWriter, r *http.Request) bool) {
	bindLogin = sum
}

var bindLogout func(w http.ResponseWriter)

//注册同步退出处理方法
func BindLogout(sum func(w http.ResponseWriter)) {
	bindLogout = sum
}

//加密字符串
func DiscuzEncode(source string) (str string) {
	return DiscuzAuthcode(source, 0)
}

//解密字符串
func DiscuzDecode(source string) (str string) {
	return DiscuzAuthcode(source, 1)
}
func DiscuzAuthcode(source string, t int) (str string) {
	ckeyLength := 4
	key := MD5(uckey)
	keya := MD5(key[:16])
	keyb := MD5(key[16:])
	var keyc string
	if t == 1 {
		keyc = source[:ckeyLength]
	} else {
		keyc = RandomString(ckeyLength)
	}
	cryptkey := keya + MD5(keya+keyc)
	key_length := len(cryptkey)
	var st []byte
	if t == 1 {
		var err error
		st, err = base64.StdEncoding.DecodeString(source[ckeyLength:])
		if err != nil {
			st, err = base64.StdEncoding.DecodeString(source[ckeyLength:] + "==")
			if err != nil {
				st, err = base64.StdEncoding.DecodeString(source[ckeyLength:] + "=")
				if err != nil {
					return ""
				}
			}
		}
	} else {
		st = []byte("0000000000" + MD5(source + keyb)[:16] + source)
	}
	st_length := len(st)
	result := bytes.Buffer{}
	box := make([]int, 256)
	for i := 0; i < 256; i++ {
		box[i] = i
	}
	rndkey := make([]int, 256)
	for i := 0; i < 256; i++ {
		l := i % key_length
		rndkey[i] = int([]byte(cryptkey[l : l+1])[0])
	}
	j := 0
	for i := 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp := box[i]
		box[i] = box[j]
		box[j] = tmp
	}

	a := 0
	j = 0
	for i := 0; i < st_length; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp

		fir := st[i]
		nex := byte(box[(box[a]+box[j])%256])
		result.WriteByte(fir ^ nex)
	}
	if t == 1 {
		res := result.String()
		if res[10:26] == MD5(res[26:] + keyb)[:16] {
			return res[26:]
		} else {
			return ""
		}
	} else {
		return keyc + strings.Replace(base64.StdEncoding.EncodeToString(result.Bytes()), "=", "", -1)
	}
	return
}
func MD5(str string) (st string) {
	m := md5.New()
	io.WriteString(m, str)
	return fmt.Sprintf("%x", m.Sum(nil))
}
func RandomString(len int) (str string) {
	ch := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	sb := bytes.Buffer{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < len; i++ {
		sb.WriteRune(ch[rand.Intn(len)])
	}
	return sb.String()
}

//核心方法调用
//参数args第一位为UserAgent
func ApiPost(moudle, action, arg string, args ...string) (str string) {
	var ua string
	if len(args) > 0 {
		ua = args[0]
	} else {
		ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.22 (KHTML, like Gecko) Chrome/25.0.1364.160 Safari/537.22"
	}

	arg = fmt.Sprintf("%s&agent=%s&time=%d", arg, MD5(ua), time.Now().Unix())
	arg = DiscuzEncode(arg)
	arg = url.QueryEscape(arg)
	postData := fmt.Sprintf("m=%s&a=%s&inajax=2&release=20110501&input=%s&appid=%s", moudle, action, arg, Appid)
	client := &http.Client{}
	req, err := http.NewRequest("POST", UCenterUrl+"/index.php?__times__=1", strings.NewReader(postData))
	req.Header = GetHeader()
	req.Header.Set("User-Agent", ua)
	resp, err := client.Do(req)

	if err == nil {
		data, err := ioutil.ReadAll(resp.Body) //取出主体的内容
		if err != nil {
			return ""
		}
		return string(data)
	} else {
		return ""
	}
}

var header *http.Header

func GetHeader() http.Header {
	if header == nil {
		header = &http.Header{}
		header.Set("Accept", "*/*")
		header.Set("Connection", "Close")
		header.Set("Accept-Language", "zh-cn")
		header.Set("Cache-Control", "no-cache")
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return *header
}
func XmlSerialize(arr interface{}, html bool, level int) (str string) {
	buff := bytes.Buffer{}
	if level == 1 {
		buff.WriteString("<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?>\r\n<root>\r\n")
	}
	switch arr.(type) {
	case []string:
		{
			for i, v := range arr.([]string) {
				buff.WriteString(fmt.Sprintf("<item id=\"%v\">", i))
				if html {
					buff.WriteString(v)
					buff.WriteString("]]>")
				} else {
					buff.WriteString(v)
				}
				buff.WriteString("</item>\r\n")
			}
		}
	case [][]string:
		{
			for i, v := range arr.([][]string) {
				buff.WriteString(fmt.Sprintf("<item id=\"%v\">\r\n", i))
				buff.WriteString(XmlSerialize(v, html, level+1))
				buff.WriteString("</item>\r\n")
			}
		}
	}
	if level == 1 {
		buff.WriteString("</root>")
	}
	return buff.String()
}
