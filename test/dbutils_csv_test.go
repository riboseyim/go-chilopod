package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test_CSV_Read(t *testing.T) {
	data := read_csv("./data/Data_GMOF_CaiPiao_List2.csv", ",")
	if data != nil {
		for i := range data {
			var url = "http://zhs.mof.gov.cn/zhuantilanmu/caipiaoguanli/"
			row := strings.Replace(data[i], "./", "", -1)
			url = url + row
			fmt.Printf("========%s \n", url)
		}
	}

}
