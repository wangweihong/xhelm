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

	/*
		err = RM.DeleteRepo("test")
		if err != nil {
			t.Error("删除仓库失败:", err)
			t.Fail()
		}
	*/

}