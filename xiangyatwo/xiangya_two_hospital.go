package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const host = "gzhceh.xyeyy.com"

var (
	mode     string
	token    string
	deptId   string
	date     string
	doctorId string
	trace    bool
)

var (
	errNoDeptID   = errors.New("xiangyatwo: error: 无deptId，请使用departments获取。")
	errNoDoctorID = errors.New("xiangyatwo: error: 无doctorId，请使用doctors获取。")
	errNoDate     = errors.New("xiangyatwo: error: 请指定出诊日期。")
	errNoToken    = errors.New("xiangyatwo: error: 请输入会话密钥")
	errNoMode     = errors.New("xiangyatwo: error: mode无匹配")
)

var modeToFn = map[string]func() error{
	"departments": departments,
	"doctors":     doctors,
	"schedules":   schedules,
	"visits":      visits,
	"order":       order,
	"patients":    patients,
}

// 1681044834099-7140E2A6C57409D057E8F5
func init() {
	flag.StringVar(&mode, "mode", "", "必须，操作类型，\ndepartments：列出所有部门\ndoctors：列出医生信息\nschedules：列出医生出诊安排\nvisits：医生出诊时间\norder：下单\npatients：列出出诊卡信息")
	flag.StringVar(&token, "token", "", "必须，会话密钥")
	flag.StringVar(&deptId, "deptId", "", "部门ID")
	flag.StringVar(&doctorId, "doctorId", "", "医生ID")
	flag.StringVar(&date, "date", "", "挂号日期，格式2006-01-02。")
	flag.BoolVar(&trace, "v", false, "打印请求结果。")
	flag.Parse()
}

func main() {
	fn, ok := modeToFn[mode]
	if !ok {
		fmt.Println(errNoMode)
		return
	}
	err := fn()
	if err != nil {
		fmt.Println(err)
	}
}

type departmentsRsp struct {
	DeptList *struct {
		DeptList []*struct {
			DeptId   string `json:"deptId"`
			DeptName string `json:"deptName"`
			DeptList []*struct {
				DeptId   string `json:"deptId"`
				DeptName string `json:"deptName"`
				DeptList []*struct {
					DeptId   string `json:"deptId"`
					DeptName string `json:"deptName"`
				}
			} `json:"deptList"`
		} `json:"deptList"`
	} `json:"deptList"`
}

func departments() error {
	if len(token) <= 0 {
		return errNoToken
	}
	body := url.Values{
		"login_access_token": []string{token},
		"hisId":              []string{"173"},
		"platformId":         []string{"173"},
		"platformSource":     []string{"3"},
		"subSource":          []string{"1"},
		"requestSource":      []string{"hc"},
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/register/deptlistfull", host), strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	setReqHeaders(request)
	var rsp departmentsRsp
	err = doPost(request, &rsp)
	if err != nil {
		return err
	}
	fmt.Println("序号", "部门ID", "部门名称")
	for i, v := range rsp.DeptList.DeptList {
		fmt.Println(strconv.Itoa(i+1)+":", v.DeptId, v.DeptName)
		for j, v1 := range v.DeptList {
			fmt.Println("  ", strconv.Itoa(i+1)+"."+strconv.Itoa(j+1)+":", v1.DeptId, v1.DeptName)
			if v1.DeptList != nil && len(v1.DeptList) != 0 {
				for l, v2 := range v1.DeptList {
					fmt.Println("    ", strconv.Itoa(i+1)+"."+strconv.Itoa(j+1)+"."+strconv.Itoa(l+1)+":", v2.DeptId, v2.DeptName)
				}
			}
		}
	}
	return nil
}

type doctorsRsp struct {
	DoctorList []*struct {
		DoctorID    string `json:"doctorId"`
		DoctorName  string `json:"doctorName"`
		DoctorTitle string `json:"doctorTitle"`
	} `json:"doctorList"`
}

func doctors() error {
	if len(token) <= 0 {
		return errNoToken
	}
	if len(deptId) <= 0 {
		return errNoDeptID
	}
	body := url.Values{
		"login_access_token": []string{token},
		"deptId":             []string{deptId},
		"hisId":              []string{"173"},
		"platformId":         []string{"173"},
		"platformSource":     []string{"3"},
		"subSource":          []string{"1"},
		"requestSource":      []string{"hc"},
		"subHisId":           []string{},
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/register/doctorlist", host), strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	setReqHeaders(request)
	var rsp doctorsRsp
	err = doPost(request, &rsp)
	if err != nil {
		return err
	}
	fmt.Println("序号", "医生ID", "医生姓名", "职位")
	for i, v := range rsp.DoctorList {
		fmt.Println(strconv.Itoa(i+1)+":", v.DoctorID, v.DoctorName, v.DoctorTitle)
	}
	return nil
}

type schedulesRsp struct {
	DoctorList []*struct {
		DoctorID   string `json:"doctorId"`
		DoctorName string `json:"doctorName"`
		Fee        int    `json:"registerFee"`
		Left       int    `json:"leftSource"`
		Total      int    `json:"totalSource"`
	} `json:"doctorList"`
	Date string `json:"scheduleDate"`
}

func schedules() error {
	if len(token) <= 0 {
		return errNoToken
	}
	if len(deptId) == 0 {
		return errNoDeptID
	}
	if len(date) == 0 {
		return errNoDate
	}
	body := url.Values{
		"deptId":             []string{deptId},
		"scheduleDate":       []string{date},
		"hisId":              []string{"173"},
		"platformId":         []string{"173"},
		"platformSource":     []string{"3"},
		"subSource":          []string{"1"},
		"requestSource":      []string{"hc"},
		"login_access_token": []string{token},
		"doctorSchedule":     []string{},
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/register/scheduledoctorlist?_route=h173&k", host), strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	setReqHeaders(request)
	var rsp schedulesRsp
	err = doPost(request, &rsp)
	if err != nil {
		return err
	}
	fmt.Println("序号", "医生ID", "医生姓名", "挂号费", "剩余", "总数")
	for i, v := range rsp.DoctorList {
		fmt.Println(strconv.Itoa(i+1)+":", v.DoctorID, v.DoctorName, float64(v.Fee)/100, v.Left, v.Total)
	}
	return nil
}

type visitsRsp struct {
	DeptID   string       `json:"deptID"`
	DoctorID string       `json:"doctorId"`
	Date     string       `json:"scheduleDate"`
	ItemList []*visitInfo `json:"itemList"`
}

func visits() error {
	list, err := getVisits()
	if err != nil {
		return err
	}
	fmt.Println("序号", "ID", "开始时间", "结束时间", "挂号费", "剩余", "总数")
	for i, v := range list {
		fmt.Println(strconv.Itoa(i+1)+":", v.ScheduleID, v.VisitBeginTime, v.VisitEndTime, float64(v.Fee)/100, v.Left, v.Total)
	}
	return nil
}

func order() error {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	count := 1
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		log.Printf("第%d轮开抢\n", count)
		list, err := getVisits()
		if err != nil {
			log.Printf("%v\n", err)
			continue
		}
		wg.Add(len(list))
		for _, req := range list {
			go func(req *visitInfo) {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
				}
				orderRsp, err := generatorOrder(ctx, req)
				if err != nil {
					log.Printf("%v\n", err)
					return
				}
				if len(orderRsp.OrderId) > 0 {
					cancel()
					log.Printf("锁号成功，请及时前往小程序上缴费，%s-%s %s\n", req.VisitBeginTime, req.VisitEndTime, orderRsp.OrderId)
				} else {
					log.Printf("%s %s-%s\n", date, req.VisitBeginTime, req.VisitEndTime)
				}
			}(req)
		}
		wg.Wait()
		count++
	}
}

type visitInfo struct {
	Fee            int    `json:"registerFee"`
	Left           int    `json:"leftSource"`
	ScheduleID     string `json:"scheduleId"`
	Status         int    `json:"status"`
	VisitBeginTime string `json:"visitBeginTime"`
	VisitEndTime   string `json:"visitEndTime"`
	VisitPeriod    int    `json:"visitPeriod"`
	Total          int    `json:"totalSource"`
}

func getVisits() ([]*visitInfo, error) {
	if len(token) <= 0 {
		return nil, errNoToken
	}
	if len(deptId) == 0 {
		return nil, errNoDeptID
	}
	if len(doctorId) == 0 {
		return nil, errNoDoctorID
	}
	if len(date) == 0 {
		return nil, errNoDate
	}
	body := url.Values{
		"login_access_token": []string{token},
		"deptId":             []string{deptId},
		"scheduleDate":       []string{date},
		"type":               []string{},
		"doctorId":           []string{doctorId},
		"hisId":              []string{"173"},
		"platformId":         []string{"173"},
		"platformSource":     []string{"3"},
		"subSource":          []string{"1"},
		"requestSource":      []string{"hc"},
		"subHisId":           []string{""},
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/register/schedulelist?_route=h173&k", host), strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	setReqHeaders(request)
	var rsp visitsRsp
	err = doPost(request, &rsp)
	if err != nil {
		return nil, err
	}
	return rsp.ItemList, nil
}

type generatorOrderRsp struct {
	OrderId string `json:"orderId"`
}

func generatorOrder(ctx context.Context, req *visitInfo) (*generatorOrderRsp, error) {
	body := url.Values{
		"deptId":             []string{deptId},
		"doctorId":           []string{doctorId},
		"extFields":          []string{`{"_bdaiGuide":"","_doctorQrGuide":"","_deptQrGuide":"","_hcSource":""}`},
		"scheduleDate":       []string{date},
		"scheduleId":         []string{req.ScheduleID},
		"visitPeriod":        []string{strconv.Itoa(req.VisitPeriod)},
		"visitBeginTime":     []string{req.VisitBeginTime},
		"visitEndTime":       []string{req.VisitEndTime},
		"patientId":          []string{"1500748"},
		"payFlag":            []string{"1"},
		"transParam":         []string{`{"type":"hcTransParam","plat":"gzhc365zhyy"}`},
		"patCardNo":          []string{"000009309082"},
		"outExtFieldsFlag":   []string{"1"},
		"hisId":              []string{"173"},
		"platformId":         []string{"173"},
		"platformSource":     []string{"3"},
		"subSource":          []string{"1"},
		"requestSource":      []string{"hc"},
		"login_access_token": []string{token},
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/register/generatororder", host), strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	setReqHeaders(request)
	var rsp generatorOrderRsp
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	err = doPost(request, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

type patientsRsp struct {
	CardList []*struct {
		PatCardNo   string `json:"patCardNo"`
		PatientName string `json:"patientName"`
	} `json:"cardList"`
}

func patients() error {
	body := url.Values{
		"login_access_token": []string{token},
		"noAuthOn999":        []string{"false"},
		"hisId":              []string{"173"},
		"platformId":         []string{"173"},
		"platformSource":     []string{"3"},
		"subSource":          []string{"1"},
		"requestSource":      []string{"hc"},
	}
	request, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("https://%s/api/homepage/getpatientslist", host), strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	setReqHeaders(request)
	var rsp patientsRsp
	err = doPost(request, &rsp)
	if err != nil {
		return err
	}
	fmt.Println("序号", "姓名", "卡号")
	for i, v := range rsp.CardList {
		fmt.Println(strconv.Itoa(i+1)+":", v.PatientName, v.PatCardNo)
	}
	return nil
}

type commonRsp struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

func doPost(req *http.Request, v interface{}) error {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	rspBody, err := io.ReadAll(response.Body)
	if trace {
		fmt.Println(string(rspBody))
	}
	if err != nil {
		return err
	}
	var r commonRsp
	err = json.Unmarshal(rspBody, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return errors.New(r.Msg)
	}
	err = json.Unmarshal(r.Data, v)
	if err != nil {
		return err
	}
	return nil
}

func setReqHeaders(req *http.Request) {
	req.Header.Set("Hc-Proj-Info", "project/his-wxapp;type/miniapp;ch/wechat;ver/4.0;")
	req.Header.Set("Hc-Src-Hisid", "173")
	req.Header.Set("requestSource", "hc")
	req.Header.Set("Host", "gzhceh.xyeyy.com")
	req.Header.Set("Referer", "https://servicewechat.com/wx4317b293c2c78b38/85/page-frame.html")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("bqqgu", "aOMGBN43p42s3h.tBlQRtpTDK59LwogcJoiXELo9_AHmOTzSWTwb5h1Jyjrtixyk8Zn5I6pVe7JkZH0sJmPx6cSo8WjorlZtPzHhcgWQRxE0FYllGECrfkNpF9orJ2UO8rI7qZWV3UN0o1XVLv_YTYKs2fHVx3aYmVK1QtiDtpHoTfgwv6m7a61kpcnDfQK.7yfdg5U6N1lWPprg4PVDvb8CahlVR9onGLTGIWgyeaXxfwvC9GHBr0HWSlS9xlBVJf_O7DbRyXfc7PpaqFOuxrgHNVHlQU9zGLifs.o9tJc1vGfJHBHkkzWkSu3hmjBzkLaZGXIezf_Ov6BJYBFON4Y8._p6PMflGAYQpluaTvKx6Z54MmvepwZjXhZmykVZHIKPUCv6smCAZMr3qDi2lTyksf8ho.bJqnOZDyEjK9JJ3fAGEvwEZVoegGjIj.sHG7EI3xTa5rNUQSFWtJL0HrzGgRXRTTACs3M7TOH4D4_CpoyAp6Bsn4VljCK4UOvNl2CMnJ40ftyq4jaU_nKtthHQVAHiPxqsgIgr4Vmg0o.ooeK1FuhERDnV4BfIFgKIsgIZlIySI7RQAMsRUwLN6hxyYV1YYwLe6i4sNelQ2cjG")
}
