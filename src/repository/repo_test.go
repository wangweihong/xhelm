package repository

import (
	"fmt"
	"testing"
)

func Test_Repo(t *testing.T) {
	t.Log("创建仓库test")
	err := RM.CreateRepo("test")
	if err != nil {
		t.Error("创建仓库失败:", err)
		t.Fail()
	}

	repos := RM.ListRepos()
	for _, v := range repos {
		fmt.Printf("%v %v %v\n", v.Name, v.State, v.CreateTime)
	}

	err = RM.DeleteRepo("test")
	if err != nil {
		t.Error("删除仓库失败:", err)
		t.Fail()
	}
}

func Test_AddRemote(t *testing.T) {
	opt := CreateOption{
		URL: "http://127.0.0.1:8879",
	}

	err := RM.AddRemoteRepo("localdddd", opt)
	if err != nil {
		t.Error("添加远程仓库失败:", err)
		t.Fail()
	}

	cs, err := RM.ListCharts("localdddd")
	if err != nil {
		t.Error("查看chart失败")
		t.Fail()
	}

	for _, v := range cs {
		fmt.Println(v.Name)
		for _, j := range v.Versions {
			fmt.Println(j.Version)
		}
	}
}
