package utils

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func AppendPostfix(filename string, addPostfix string) string {
	postfix := filename[strings.Index(filename, "."):len(filename)]
	justname := filename[0:strings.Index(filename, ".")]
	newfile := justname + addPostfix + postfix
	log.Println("-----old:%s just:%s addPost:%s newfile %s-----", filename, justname, addPostfix, newfile)
	return newfile
}

func ReadLine(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				return result, nil
			}
			return nil, err
		}
		//	log.Println("-----ReadLine:%s", line)
		result = append(result, line)
	}
	return result, nil
}

func SaveFile(filename string, data [][]string, append bool) {
	if !append {
		os.Remove(filename)
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600) //创建文件

	if err != nil {
		panic(err)
	}
	defer f.Close()

	//	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)      //写入数据
	w.Flush()
}

// fileName:文件名字(带全路径)
// content: 写入的内容
func AppendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

func WriteToFile(fileName string, content []string) error {
	// 以只写的模式，打开文件
	os.Remove(fileName)
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("%s file create failed. err: %s", fileName, err.Error())
	} else {
		// 查找文件末尾的偏移量

		// 从末尾的偏移量开始写入内容
		for i := 0; i < len(content); i++ {
			n, _ := f.Seek(0, os.SEEK_END)
			_, err = f.WriteAt([]byte(content[i]+"\n"), n)
		}

	}
	defer f.Close()
	return err
}
