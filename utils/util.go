/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2021-11-10 09:48:35
 * @LastEditTime: 2022-02-11 17:22:37
 */

package utils

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/syyongx/php2go"
)

// 获取两个 string 的相似度
func GetRatio(first string, second string) (ratio float64) {
	_ = php2go.SimilarText(first, second, &ratio)
	return ratio / 100
}

// HandleError 用于处理 error
func HandleError(action string, err error) {
	if err != nil {
		_ = fmt.Errorf(fmt.Sprintf("%s => %s\n", action, err))
	}
}

// Sha1 加密
func Sha1(content string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(content)))
}

// md5 加密
func Md5(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}

// 获取终端宽度
func GetTermWidth() int {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	width, _ := termbox.Size()
	termbox.Close()

	return width
}

// 反转 [][]string
func ReverseSlice(s [][]string) [][]string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// slice 转为 string
func SliceToString(slc []string) string {
	return "[" + strings.Join(slc, ", ") + "]"
}

// slice 排序
func SortSlice(t []string) {
	sort.Slice(t, func(i, j int) bool {
		return t[i] < t[j]
	})
}

// 判断是否在 slice, array, map 中
func InSlice(needle interface{}, haystack interface{}) bool {
	return php2go.InArray(needle, haystack)
}

// slice(string类型)元素去重
func UniqStringSlice(slc []string) []string {
	result := make([]string, 0)
	tempMap := make(map[string]bool, len(slc))
	for _, e := range slc {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

// slice(int类型)元素去重
func UniqIntSlice(slc []int) []int {
	result := make([]int, 0)
	tempMap := make(map[int]bool, len(slc))
	for _, e := range slc {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

// 获取区间中的随机整数
func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

// 获取指定长度的随机字符串
func RandomStr(length int) string {
	b_str := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		s := b_str[rand.Intn(len(b_str))]
		result = append(result, s)
	}
	return string(result)
}

// 获取文件内容
func FileGetContents(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(content)
}

// 按行读取文件内容
func FileEachLineRead(filename string) []string {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0664)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var datas []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		datas = append(datas, sc.Text())
	}
	return datas
}

// 判定文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// 时间戳转时间字符串 => 2006-01-02 15:04:05
func UnixToTime(timestamp interface{}) string {
	// 通过 i.(type) 来判断是什么类型,下面的 case 分支匹配到了则执行相关的分支
	switch timestamp.(type) {
	case int:
		t := int64(timestamp.(int)) // interface 转为 int 再转为 int64
		return time.Unix(t, 0).Format("2006-01-02 15:04:05")
	case int64:
		return time.Unix(timestamp.(int64), 0).Format("2006-01-02 15:04:05")
	case string:
		t, _ := strconv.ParseInt(timestamp.(string), 10, 64) // interface 转为 string 再转为 int64
		return time.Unix(t, 0).Format("2006-01-02 15:04:05")
	}
	return ""
}
