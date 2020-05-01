package models

type BusinessLogErrorCount struct {
	Id         int
	App        string
	ModuleName string
	Level      string
	Count      int
}

type LogAppAccess struct {
	Id            int
	App       string
	Url           string
	AccessCount   int
	AccessAvgTime float64
	AccessMinTime float64
	AccessMaxTime float64
}
