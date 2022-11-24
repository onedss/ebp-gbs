package client

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/onedss/ebp-gbs/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type httpServer struct {
	httpPort   int
	httpServer *http.Server
}

func NewOneHttpServer(httpPort int) (server *httpServer) {
	return &httpServer{
		httpPort: httpPort,
	}
}

func (p *httpServer) Start() (err error) {
	redis_addr := "192.168.0.101:26379"
	redis_pass := "livegbs@2019"
	url := "http://127.0.0.1:8081/async/alarm?method=fireCameraAlarm"
	redisAddr := utils.Conf().Section("redis").Key("address").MustString(redis_addr)
	redisPass := utils.Conf().Section("redis").Key("password").MustString(redis_pass)
	redisDB := utils.Conf().Section("redis").Key("database").MustInt(0)
	httpUrl := utils.Conf().Section("ebp").Key("httpUrl").MustString(url)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass, // no password set
		DB:       redisDB,   // use default DB
	})
	ctx := context.Background()
	cn := rdb.Conn(ctx)
	defer cn.Close()
	name := fmt.Sprintf("OneGBS_%s", time.Now().Format("2006.0102.150405"))
	if err := cn.ClientSetName(ctx, name).Err(); err != nil {
		log.Printf("Connect to redis error!!! %v", err)
		return err
	}

	name, err = cn.ClientGetName(ctx).Result()
	if err != nil {
		log.Printf("Visit to redis error!!! %v", err)
		return err
	}
	fmt.Println("Client Name:", name)

	// There is no error because go-redis automatically reconnects on error.
	pubSub := rdb.Subscribe(ctx, "alarm")

	// Close the subscription when we are done.
	defer pubSub.Close()

	ch := pubSub.Channel()

	for msg := range ch {
		alarm := msg.Payload
		log.Println("收到：", msg.Channel, msg.Payload)
		sendAlarm(alarm, httpUrl)
	}

	log.Println("Done.")
	return nil
}

func (p *httpServer) Stop() (err error) {
	return nil
}

func (p *httpServer) GetPort() int {
	return p.httpPort
}

func sendAlarm(alarm string, url string) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(alarm))
	if err != nil {
		log.Println("请求失败，错误原因：", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	//resp, err := client.Post("http://192.168.101.219:58080/async/alarm?method=fireCameraAlarm", "application/json", strings.NewReader(alarm))
	if err != nil {
		log.Println("请求失败，错误原因：", err)
		return
	}
	defer resp.Body.Close()
	// 200 OK
	log.Println("返回码：", resp.Status, "请求内容", alarm)
	//fmt.Println("返回头：", resp.Header)
	if resp.StatusCode != 200 {
		log.Println("请求失败，返回码：", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("返回数据失败，错误原因：", err)
		return
	}
	content := string(body)
	log.Println("处理完毕。返回数据", content)
	//buf := make([]byte, 1024)
	//for {
	//	// 接收服务端信息
	//	n, err := resp.Body.Read(buf)
	//	if err != nil && err != io.EOF {
	//		fmt.Println(err)
	//		return
	//	} else {
	//		fmt.Println("处理完毕")
	//		res := string(buf[:n])
	//		fmt.Println(res)
	//		break
	//	}
	//}
}
