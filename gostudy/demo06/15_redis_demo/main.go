package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type News struct {
	ID    int64
	Title string
	Link  url.URL
	Ctime time.Time
}

var wg sync.WaitGroup

const (
	Dir         = "D:\\news"
	CachePrefix = "gonews"
	Host        = "192.168.1.61:6379"
	CacheDB     = 2
)

var client *redis.Client

func init() {
	if client == nil {
		client = newClient()
	}
}

func newClient() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     Host,
		Password: "",
		DB:       CacheDB,
	})
	return client
}

func main() {
	//loadData()
	getData()
}

func getData() {
	pageNum, pageSize := int64(1), int64(5)
	news, count := getAllNews(pageNum, pageSize)
	data := map[string]interface{}{
		"total":    count,
		"pagesize": pageSize,
		"items":    news,
	}
	jData, _ := json.Marshal(data)
	fmt.Println(string(jData))
}

// 获取全部新闻
func getAllNews(pageNum int64, pageSize int64) ([]map[string]string, int64) {
	key0 := CachePrefix
	keys1, _ := client.SMembers(key0).Result()
	newsList := map[string]map[string]string{}
	var count int
	for _, key1 := range keys1 {
		keys2, _ := client.SMembers(key0 + ":" + key1).Result()
		for _, key2 := range keys2 {
			news, err := getNewsCache(key0 + ":" + key1 + ":" + key2)
			if err == nil {
				newsList[news["ctime"]+news["id"]] = news
			}
			count++
		}
	}
	fmt.Println(count)
	allNews := sortNews(newsList)
	pageNews := []map[string]string{}
	var i int64 = 0
	for _, item := range allNews {
		if i >= (pageNum-int64(1))*pageSize && i < pageNum*pageSize {
			pageNews = append(pageNews, item)
		}
		i++
	}
	return pageNews, i
}

// 对新闻进行排序
func sortNews(raw map[string]map[string]string) []map[string]string {
	keys := []string{}
	data := []map[string]string{}
	for key := range raw {
		keys = append(keys, key)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	for _, key := range keys {
		data = append(data, raw[key])
	}
	return data
}

// 获取新闻缓存
func getNewsCache(key string) (map[string]string, error) {
	return client.HGetAll(key).Result()
}

func loadData() {
	files := getFileList(Dir)
	for _, file := range files {
		wg.Add(1)
		go cacheNews(file)
	}
	wg.Wait()
}

// 并发缓存数据
func cacheNews(path string) {
	defer wg.Done()
	newsList, _ := getNews(path)
	for _, item := range newsList {
		cache := map[string]interface{}{
			"id":    strconv.Itoa(int(item.ID)),
			"title": item.Title,
			"link":  item.Link.String(),
			"ctime": item.Ctime.Format("20060102"),
		}
		setNewsCache(cache) // 缓存数据
	}
}

func setNewsCache(cache map[string]interface{}) error {
	key0 := CachePrefix
	var key1, key2 string
	if value, ok := cache["ctime"].(string); ok {
		key1 = key0 + ":" + value
		err := client.SAdd(key0, value).Err()
		if err != nil {
			return err
		}
	}
	if value, ok := cache["title"].(string); ok {
		key2 = key1 + ":" + value
		err := client.SAdd(key1, value).Err()
		if err != nil {
			return err
		}
		return client.HMSet(key2, cache).Err()
	}
	return nil
}

func getNews(path string) (newsList []News, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	reg := regexp.MustCompile(`(?i)#{0,3}\s*GoCN每日新闻\(([\d-]*)\)\n+\D*((.*\n)+?\n)`)
	allMatch := reg.FindAllSubmatch(data, -1)
	for _, item := range allMatch {
		loc, _ := time.LoadLocation("Local")
		ctime, err := time.ParseInLocation("2006-01-02", string(item[1]), loc)
		if err != nil {
			continue
		}
		subReg := regexp.MustCompile(`(\d)\.\s*(.*)\s+(http.*)\n?`)
		subAll := subReg.FindAllSubmatch(item[2], -1)
		for _, subItem := range subAll {
			id, err := strconv.ParseInt(string(subItem[1]), 10, 64)
			if err != nil {
				continue
			}
			title := strings.TrimSpace(string(subItem[2]))
			sUrl, err := url.Parse(string(subItem[3]))
			if err != nil {
				continue
			}
			singleNews := News{id, title, *sUrl, ctime}
			newsList = append(newsList, singleNews)
		}
	}
	return
}

func getFileList(dir string) []string {
	files := []string{}
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if strings.Contains(path, ".git") {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		baseName := filepath.Base(path)
		if !strings.Contains(baseName, "-") {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files
}
