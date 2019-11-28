package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"reflect"
	"strconv"
	"strings"
)

const (
	LevelNatUp      = "国家正"
	LevelNatDown    = "国家副"
	LevelProUp      = "省级正"
	LevelProDown    = "省级副"
	LevelHallUp     = "厅正"
	LevelHallDown   = "厅副"
	LevelCountyUp   = "县处级正"
	LevelCountyDown = "县处级副"
	LevelSecUp      = "乡科级正"
	LevelSecDown    = "乡科级副"

	DepCommittee  = "党委"
	DepGovernment = "政府"
	DepPeople     = "人大"
	DepPolitical  = "政协"
)

type Info struct {
	ProVince string
	City     string
	Time     string
	Official string
	Level    string
	Dep      string
	Justice  string
	Man      string
}

type result struct {
	Province string
	City     string
	Date     string

	CityCount         int // 城市人数
	NOTProCount       int
	CityOfficial      int // 城市官员
	NOTProOfficial    int
	CityNotOfficial   int // 非官员
	NOTProNotOfficial int

	NationalUp      int // 国家正
	NOTNationalUp   int
	NationalDown    int
	NOTNationalDown int

	ProvinceUp      int // 省级正
	NOTProvinceUp   int
	ProvinceDown    int
	NOTProvinceDown int

	HallUp      int // 厅正
	NOTHallUp   int
	HallDown    int
	NOTHallDown int

	CountyUp      int // 县级正
	NOTCountyUp   int
	CountyDown    int
	NOTCountyDown int

	SecUp      int // 科正
	NOTSecUp   int
	SecDown    int
	NOTSecDown int

	Committee     int // 党委
	NOTCommittee  int
	Government    int // 政府
	NOTGovernment int
	People        int // 人大
	NOTPeople     int
	Political     int // 政协
	NOTPolitical  int

	Justice    int //司法
	NOTJustice int

	Man    int // 男性
	NOTMan int
}

func readXls(file string) []Info {
	var data []Info
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return data
	}
	rows, err := f.GetRows("Sheet2")
	if err != nil {
		fmt.Println("read rows failed.", err)
		return data
	}
	for _, row := range rows {
		if len(row) < 12 {
			fmt.Println("this row is invalid.", row)
			continue
		}
		tmp := Info{
			ProVince: row[0],
			City:     row[1],
			Time:     row[4],
			Official: row[7],
			Level:    row[8],
			Dep:      row[9],
			Justice:  row[10],
			Man:      row[11],
		}
		data = append(data, tmp)
	}
	fmt.Println(data)
	return data
}

func GetResult(data []Info) {
	proCity := GetProCities(data)
	fmt.Println(proCity)

	res := make(map[string]map[string]result)
	for _, d := range data {
		date := parseDate(d.Time)
		if date == "" {
			continue
		}
		if d.City == "" { // 未匹配市
			for city := range proCity[d.ProVince] {
				for i := range res[city] {
					r := res[city][i]
					judge(d, &r)
					res[city][i] = r
				}
			}
		} else { // 属于城市
			if _, ok := res[d.City]; !ok {
				res[d.City] = map[string]result{}
			}
			if _, ok := res[d.City][date]; !ok {
				res[d.City][date] = result{Province: d.ProVince, City: d.City, Date: date}
			}
			r := res[d.City][date]
			judge(d, &r)
			res[d.City][date] = r
		}
	}
	a := false
	for k := range res {
		for k2 := range res[k] {
			t := reflect.TypeOf(res[k][k2])
			v := reflect.ValueOf(res[k][k2])
			if !a {
				for k := 0; k < t.NumField(); k++ {
					fmt.Printf("%s,", t.Field(k).Name)
				}
				fmt.Println()
				a = true
			}
			for k := 0; k < t.NumField(); k++ {
				fmt.Printf("%v,", v.Field(k).Interface())
			}
			fmt.Println()
		}
	}
}

func judge(d Info, r *result) {
	if d.City == "" {
		r.NOTProCount++
		if d.Official == "官员" {
			r.NOTProOfficial++
		} else {
			r.NOTProNotOfficial++
		}
		switch d.Level {
		case LevelNatUp:
			r.NOTNationalUp++
		case LevelNatDown:
			r.NOTNationalDown++
		case LevelProUp:
			r.NOTProvinceUp++
		case LevelProDown:
			r.NOTProvinceDown++
		case LevelHallUp:
			r.NOTHallUp++
		case LevelHallDown:
			r.NOTHallDown++
		case LevelCountyUp:
			r.NOTCountyUp++
		case LevelCountyDown:
			r.NOTCountyDown++
		case LevelSecUp:
			r.NOTSecUp++
		case LevelSecDown:
			r.NOTSecDown++
		}
		switch d.Dep {
		case DepCommittee:
			r.NOTCommittee++
		case DepGovernment:
			r.NOTGovernment++
		case DepPeople:
			r.NOTPeople++
		case DepPolitical:
			r.NOTPolitical++
		}
		if d.Justice == "是" {
			r.NOTJustice++
		}
		if d.Man == "是" {
			r.NOTMan++
		}
	} else {
		r.CityCount++
		if d.Official == "官员" {
			r.CityOfficial++
		} else {
			r.CityNotOfficial++
		}
		switch d.Level {
		case LevelNatUp:
			r.NationalUp++
		case LevelNatDown:
			r.NationalDown++
		case LevelProUp:
			r.ProvinceUp++
		case LevelProDown:
			r.ProvinceDown++
		case LevelHallUp:
			r.HallUp++
		case LevelHallDown:
			r.HallDown++
		case LevelCountyUp:
			r.CountyUp++
		case LevelCountyDown:
			r.CountyDown++
		case LevelSecUp:
			r.SecUp++
		case LevelSecDown:
			r.SecDown++
		}
		switch d.Dep {
		case DepCommittee:
			r.Committee++
		case DepGovernment:
			r.Government++
		case DepPeople:
			r.People++
		case DepPolitical:
			r.Political++
		}
		if d.Justice == "是" {
			r.Justice++
		}
		if d.Man == "是" {
			r.Man++
		}
	}
}

func GetProCities(data []Info) map[string]map[string]bool {
	x := make(map[string]map[string]bool)
	for _, d := range data {
		if d.City == "" {
			continue
		}
		if _, ok := x[d.ProVince]; !ok {
			x[d.ProVince] = make(map[string]bool)
		}
		x[d.ProVince][d.City] = true
	}
	return x
}

func parseDate(d string) string {
	ss := strings.Split(d, "-")
	if len(ss) != 3 {
		fmt.Println("this data valid.", d)
		return ""
	}
	month, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		fmt.Println("this data valid.", d)
		return ""
	}
	if month < 7 {
		return ss[2]
	}
	year, err := strconv.ParseInt(ss[2], 10, 64)
	if err != nil {
		fmt.Println("this data valid.", d)
		return ""
	}
	return strconv.FormatInt(year+1, 10)
}
