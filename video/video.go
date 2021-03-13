package video

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
)
type Video struct{
    Bvid string `json:"bvid"`
    Name string `json:"title"`
    Instroction string `json:"description"`
}

func GetVideoList(mid int) []Video {
    fmt.Println("正在获取用户: ", mid, " 的视频")
	res, _ := http.Get("https://api.bilibili.com/x/space/arc/search?order=cllick&mid=" + fmt.Sprint(mid))

	buf, _ := ioutil.ReadAll(res.Body)
	data, _ := simplejson.NewJson(buf)

	list, _ := data.Get("data").Get("list").Get("vlist").Encode()

    var videoList []Video
    json.Unmarshal(list, &videoList)
    return videoList
}
