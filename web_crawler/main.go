package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	MainPageUrl    = "https://bet.hkjc.com/football/schedule/schedule.aspx?lang=en"
	DataUrl        = "https://bet.hkjc.com/football/getJSON.aspx?jsontype=schedule.aspx"
	CookieFilePath = "./ck.txt"
)

type GameList struct {
	GameList []Game `json:"GameList"`
}

type Game struct {
	MatchID           string `json:"matchID"`
	MatchIDinofficial string `json:"matchIDinofficial"`
	MatchNum          string `json:"matchNum"`
	MatchDate         string `json:"matchDate"`
	MatchDay          string `json:"matchDay"`
	Coupon            struct {
		CouponID        string `json:"couponID"`
		CouponShortName string `json:"couponShortName"`
		CouponNameCH    string `json:"couponNameCH"`
		CouponNameEN    string `json:"couponNameEN"`
	} `json:"coupon"`
	League struct {
		LeagueID        string `json:"leagueID"`
		LeagueShortName string `json:"leagueShortName"`
		LeagueNameCH    string `json:"leagueNameCH"`
		LeagueNameEN    string `json:"leagueNameEN"`
	} `json:"league"`
	HomeTeam struct {
		TeamID     string `json:"teamID"`
		TeamNameCH string `json:"teamNameCH"`
		TeamNameEN string `json:"teamNameEN"`
	} `json:"homeTeam"`
	AwayTeam struct {
		TeamID     string `json:"teamID"`
		TeamNameCH string `json:"teamNameCH"`
		TeamNameEN string `json:"teamNameEN"`
	} `json:"awayTeam"`
	MatchStatus       string    `json:"matchStatus"`
	MatchTime         time.Time `json:"matchTime"`
	Statuslastupdated time.Time `json:"statuslastupdated"`
	Inplaydelay       string    `json:"inplaydelay"`
	LiveEvent         struct {
		IlcLiveDisplay  bool          `json:"ilcLiveDisplay"`
		HasLiveInfo     bool          `json:"hasLiveInfo"`
		IsIncomplete    bool          `json:"isIncomplete"`
		MatchIDbetradar string        `json:"matchIDbetradar"`
		Matchstate      string        `json:"matchstate"`
		StateTS         string        `json:"stateTS"`
		Liveevent       []interface{} `json:"liveevent"`
	} `json:"liveEvent"`
	Cornerresult      string `json:"cornerresult"`
	Cur               string `json:"Cur"`
	HasWebTV          bool   `json:"hasWebTV"`
	HasOdds           bool   `json:"hasOdds"`
	HasExtraTimePools bool   `json:"hasExtraTimePools"`
	Results           struct {
	} `json:"results"`
	DefinedPools []string      `json:"definedPools"`
	InplayPools  []interface{} `json:"inplayPools"`
}

func main() {

	cookieList := readCookie()
	if len(cookieList) == 0 {
		if getCookie() {
			cookieList = readCookie()
		}
	}
	if len(cookieList) == 0 {
		log.Println("get cookie failed")
		os.Exit(0)
	}

	responseString := sendRequestWithCookie(cookieList)
	log.Println(responseString)
	var gameList []*Game
	if err := json.Unmarshal([]byte(responseString), &gameList); err != nil {
		log.Println("error:", err)
		return
	}
	if len(gameList) == 0 {
		log.Println("game list is empty")
		return
	}
	writeCVS(gameList)
}

/**
 * @description: write data to cvs
 * @param {[]*Game} gameList
 * @return {*}
 */
func writeCVS(gameList []*Game) {
	csvFile, err := os.Create("./" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	line1 := []string{
		"MatchID",
		"MatchIDinofficial",
		"MatchDay",
		"MatchNum",
		"HomeTeam.TeamNameEN",
		"AwayTeam.TeamNameEN",
		"MatchTime",
	}
	err = writer.Write(line1)
	if err != nil {
		panic(err)
	}
	for _, game := range gameList {
		line := []string{
			game.MatchID,
			game.MatchIDinofficial,
			game.MatchDay,
			game.MatchNum,
			game.HomeTeam.TeamNameEN,
			game.AwayTeam.TeamNameEN,
			game.MatchTime.Format("2006-01-02 15:04:05"),
		}
		// 将切片类型行数据写入 csv 文件
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}

/**
 * @description: 获取cookie
 * @param {*}
 * @return {*}
 */
func readCookie() []*http.Cookie {

	list := []*http.Cookie{}

	pathExists, err := pathExists(CookieFilePath)
	if err != nil {
		log.Println("check cookie file path failed ", err.Error())
		return list
	}
	if pathExists == false {
		file, err := os.Create(CookieFilePath)
		if err != nil {
			fmt.Println("create cookie file path = ", err)
		}
		defer file.Close()
		return list
	}

	contents, err := ioutil.ReadFile(CookieFilePath)
	if err != nil {
		log.Printf("Failed to read file %s,err:%s", CookieFilePath, err.Error())
		return list
	}
	cookieContent := string(contents)
	if len(cookieContent) == 0 {
		return list
	}
	splitList := strings.Split(cookieContent, ";")
	for _, v := range splitList {
		cookie := strings.Split(v, "=")
		if len(cookie) != 2 {
			continue
		}
		cookie1 := &http.Cookie{
			Name:  cookie[0],
			Value: cookie[1],
		}
		list = append(list, cookie1)
	}
	return list
}

/**
 * @description: 发送请求
 * @param {string}
 * @return {*}
 */
func sendRequestWithCookie(cookieList []*http.Cookie) string {
	req, err := http.NewRequest("GET", DataUrl, nil)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
		return ""
	}
	for _, v := range cookieList {
		req.AddCookie(v)
	}
	for _, c := range req.Cookies() {
		log.Println(c)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occured. Error is: %s", err.Error())
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return string(body)
}

/**
 * @description: get cookie
 * @param {*}
 * @return {*}
 */
func getCookie() bool {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("mute-audio", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx, openMainPage(MainPageUrl)); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func openMainPage(url string) chromedp.Tasks {
	return chromedp.Tasks{
		//跳转到页面
		chromedp.Navigate(url),
		chromedp.WaitReady(`#printDiv`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			var c string
			for _, v := range cookies {
				c = c + v.Name + "=" + v.Value + ";"
			}
			log.Println(c)
			if err != nil {
				return err
			}
			if err := saveCookie(c); err != nil {
				return err
			}
			return nil
		}),
	}
}

/**
 * @description: saveCookie
 * @param {*}
 * @return {*}
 */
func saveCookie(ck string) error {
	filePath := CookieFilePath
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)

	write.WriteString(ck)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
