package controllers

import (
	//	"ecode"
	"encoding/json"
	//	"errno"
	"fmt"
	"net"
	"runtime"
	//	"server/models/etcd"
	"strings"

	"github.com/astaxie/beego"
	"github.com/golang/glog"
)

type baseController struct {
	beego.Controller
}

type ErrStruct struct {
	Err  string      `json:"message"`
	Code int         `json:"ecode"`
	Data interface{} `json:"data"`
}

func debugPrintFunc(err string) string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)

	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	file, line := fun.FileLine(fpcs[0])
	return fmt.Sprintf("File(%v) Line(%v) Func(%v): %v", file, line, fun.Name(), err)
}

func (this *baseController) errReturn(data interface{}, statusCode int) {
	this.AllowCross()

	var errStruct ErrStruct
	errStruct.Code = statusCode
	switch v := data.(type) {
	case string:
		errStruct.Err = v
	case error:
		errStruct.Err = v.Error()
	case ErrStruct:
		errStruct = v

	}

	debugErr := fmt.Errorf("RequestIP:%v,Error:%v", this.Ctx.Request.RemoteAddr, errStruct.Err)
	glog.V(4).Info(debugPrintFunc(debugErr.Error()))
	//	uerr.PrintAndReturnError(err)

	this.Ctx.Output.SetStatus(statusCode)
	this.Data["json"] = errStruct

	this.ServeJSON()
}

func (this *baseController) normalReturn(data interface{}, statusCode ...int) {
	this.AllowCross()
	var normal ErrStruct
	normal.Data = data
	this.Data["json"] = normal

	if len(statusCode) != 0 {
		this.Ctx.Output.SetStatus(statusCode[0])
	}

	this.ServeJSON()
}

/*
func (c *baseController) Prepare() {
	if !c.checkLeader() {
		return
	}
}
*/

// AllowCross 跨域
func (c *baseController) AllowCross() {
	beego.Debug("AllowCross|Origin:", c.Ctx.Input.Header("Origin"))
	c.Ctx.Output.Header("Access-Control-Allow-Origin", c.Ctx.Input.Header("Origin"))
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")                  //允许method
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,Cookie,UserName,Token") //header的类型
	c.Ctx.Output.Header("Access-Control-Max-Age", "1728000")
	c.Ctx.Output.Header("content-type", "application/json") //返回数据格式是json
}

/*
func (c *baseController) Finish() {
	c.Data["json"] = c.MyOutput
	c.ServeJSON()
	//	c.Finish()
}
*/

func getRouteControllerName() string {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(3, fpcs)

	if n == 0 {
		return "n/a"
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	sl := strings.Split(fun.Name(), ".")
	return sl[len(sl)-1]

}

// GetClientIP 获取客户端ip
func (c *baseController) GetClientIP() string {
	clientIP := c.Ctx.Request.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		return ""
	}

	return strings.Split(clientIP, ",")[0]
}

// 检查请求地址是否为外网地址
func (c *baseController) IsPublicAddress() bool {
	clientIP := c.GetClientIP()
	if clientIP == "" {
		return false
	}

	IP := net.ParseIP(clientIP)
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

/*
func (c *baseController) checkLeader() bool {
	// check leader
	if !etcd.IsLeader(c.IsPublicAddress()) {
		err := errno.NewError(ecode.EcodeCurrentIsNotMaster, "非master")
		data := etcd.GetLeaderCandidate(c.IsPublicAddress())

		var es ErrStruct
		es.Code = err.Code
		es.Data = data
		es.Err = err.Error()
		c.Data["json"] = es
		c.ServeJSON()
		c.StopRun()
		return false
	}
	return true
}
*/

func (c *baseController) unmarshalRequestBodyJSON(req interface{}) error {
	if req == nil {
		return fmt.Errorf("Request Body Is Nil")
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		return fmt.Errorf("Unmarshal Request Body JSON Error: " + err.Error())
	}

	return nil
}

func (c *baseController) setOutput(data interface{}, err error) {
	if err != nil {
		c.errReturn(err, 500)
	} else {
		c.normalReturn(data)
	}
}
