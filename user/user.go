package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)


type User struct {
    Mid int `json:"mid"`
	Name string `json:"name"`
    Avatar string `json:"face"`
}




func GetUser(id string) User { // 根据 id 获取 bilibili 用户数据
    log.Println("开始获取: ", id, " 的信息.")
	res,_ := http.Get("https://api.bilibili.com/x/space/acc/info?mid=" + id  + "&jsonp=jsonp")
	buf, _ := ioutil.ReadAll(res.Body)
	data, _ := simplejson.NewJson(buf)

    byte, _ := data.Get("data").Encode()

    var user User
    json.Unmarshal(byte, &user)
    return user
}
