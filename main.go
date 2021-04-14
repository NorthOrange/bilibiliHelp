package main

import (
	"bilibiliHelp/user"
	"bilibiliHelp/video"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

func main() {

	user := user.GetUser("331991614") // 获取 b 站用户信息
	id := Register(user)              // 用 拿到的 用户信息注册
	UploadAvatar(user.Avatar, id)     // 更新头像

	list := video.GetVideoList(user.Mid)
	for _, v := range list {
		if strings.Contains(v.Name, "互动视频") {
			continue
		} else {
			if VideoSize(v.Bvid) { // 大小 超过 200 兆的视频跳过
				GetVideo(v.Bvid)                       // 根据 bvid 下载视频
				VideoUpload(v.Name, v.Instroction, id) // 上传视频
				fmt.Println(v.Name)
			} else {
				continue
			}

		}
	}
}

func Register(user user.User) string { // 用一个 user 进行注册, 返回一个 id
	client := &http.Client{}
	data := make(map[string]interface{})
	postUrl := socket + "/api/user/register"
	data["name"] = strings.Replace(user.Name, " ", "", -1)
	data["account"] = fmt.Sprint(user.Mid)
	data["password"] = fmt.Sprint(user.Mid)
	byteData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", postUrl, bytes.NewReader(byteData))
	res, _ := client.Do(req)
	byte, _ := ioutil.ReadAll(res.Body)
	jsonData, _ := simplejson.NewJson(byte)
	id, _ := jsonData.Get("id").Int()
	return strconv.Itoa(id)
}

func UploadAvatar(avatarUrl string, id string) { // 根据从 B 站获取的用户资料更新头像
	log.Println("获取头像")
	res, _ := http.Get(avatarUrl)
	avatar, _ := ioutil.ReadAll(res.Body)
	file, _ := os.Create("./temp.png")
	file.Write(avatar)

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("avatar", "./temp.png") // 添加一个表单文件
	f, _ := os.Open("./temp.png")
	io.Copy(fw, f)

	log.Println("用户 id: ", id, "头像上传开始")
	w.WriteField("id", id) // 添加一个表单字段
	w.Close()              // !!! 不关闭会失去终止边界

	// 发送请求
	postUrl := socket + "/api/avatar/upload"
	req, _ := http.NewRequest("POST", postUrl, buf)
	req.Header.Set("Content-Type", w.FormDataContentType()) // 设置 内容 类型
	client := &http.Client{}
	res2, _ := client.Do(req)

	res2Data, _ := ioutil.ReadAll(res2.Body)
	log.Println("头像上传结果: ", string(res2Data))
}

func GetVideo(Bvid string) { // 根据 bvid 下载视频
	cmdStr := "you-get -dash-flv https://www.bilibili.com/video/" + Bvid
	cmd := exec.Command("bash", "-c", cmdStr)

	log.Println("开始进行视频: " + Bvid + " 的下载")
	cmd.CombinedOutput()
	log.Println("视频: " + Bvid + " 下载成功")
}

func VideoUpload(VideoName string, introduction string, id string) { // 上传视频
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	fw, _ := w.CreateFormFile("video", "./"+VideoName+".mp4") // 添加一个表单文件
	f, _ := os.Open("./" + VideoName + ".mp4")
	io.Copy(fw, f)

	log.Println("视频: ", VideoName, " 上传开始")
	w.WriteField("id", id)                     // 添加一个表单字段
	w.WriteField("name", VideoName)            // 添加一个表单字段
	w.WriteField("introduction", introduction) // 添加一个表单字段
	w.WriteField("nocover", "true")            // 添加一个表单字段
	w.Close()                                  // !!! 不关闭会失去终止边界

	// 发送请求
	postUrl := socket + "/api/video/upload"
	req, _ := http.NewRequest("POST", postUrl, buf)
	req.Header.Set("Content-Type", w.FormDataContentType()) // 设置 内容 类型
	client := &http.Client{}
	res2, _ := client.Do(req)

	res2Data, _ := ioutil.ReadAll(res2.Body)
	log.Println("视频: ", VideoName, "上传结果: ", string(res2Data))

}

func VideoSize(Bvid string) bool {
	cmdStr := "you-get --json https://www.bilibili.com/video/" + Bvid
	cmd := exec.Command("bash", "-c", cmdStr)
	log.Println("开始进行视频: " + Bvid + " 的大小判断")
	json, _ := cmd.CombinedOutput()
	jsondata, _ := simplejson.NewJson(json)

	size, _ := jsondata.Get("streams").Get("dash-flv").Get("size").Int()
	if size <= 0 || size > 200*1024*1024 {
		fmt.Println("大小不正常")
		return false
	} else {
		fmt.Println("大小正常")
		return true
	}
}
