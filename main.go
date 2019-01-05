package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ampData struct {
	ID             string `json:"id"`
	Model          string `json:"model"`
	Bom            string `json:"bom"`
	Band           string `json:"band"`
	Direction      string `json:"direction"`
	TestName       string `json:"testName"`
	LowLimit       string `json:"lowLimit"`
	UpLimit        string `json:"upLimit"`
	LEeprom        string `json:"lEeprom"`
	UEeprom        string `json:"uEeprom"`
	IP             string `json:"ip"`
	Target         string `json:"target"`
	LowerFrequency string `json:"lowerFrequency"`
	UpperFrequency string `json:"upperFrequency"`
	Voltage        string `json:"voltage"`
	Connectors     string `json:"connectors"`
}

type queryParams struct {
	model string
	bom   string
}

//GetAmpSpecs gets the params sent by the GET request
func GetAmpSpecs(w http.ResponseWriter, r *http.Request) {
	//Change the varible to something better
	serial := strings.Split(r.URL.Path, "/")
	//ampInfo := r.URL.Query().Get("params")

	if len(serial[2]) > 0 {
		ampBom := serial[2]
		w.Write(sqlQuery(ampBom))
	}

}

func sqlQuery(ampBom string) []byte {
	dsn := "Data source"
	queryParams := queryParams{
		model: ampBom[0 : len(ampBom)-1],
		bom:   ampBom[len(ampBom)-1:],
	}

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return nil
	}
	defer db.Close()

	data, err := db.Query(`
		SELECT am.[model] AS Model,
			   bm.[name] AS Bom,
			   bi.[bandnumber] AS Band,
			   d.[name] AS Direction,
			   ts.[name] AS TestName,
			   ts.[lowerlimit] AS LowLimit,
			   ts.[upperlimit] AS UpLimit,
			   ts.[eeprom_lowerlimit] AS LEeprom,
			   ts.[eeprom_upperlimit] AS UEeprom,
			   ts.[inputpower] AS IP, ts.target AS Target,
			   ts.[lowerfrequency] AS Lfrequency,
			   ts.[upperfrequency] AS Ufrequency,
			   ts.[operatingvoltage] AS Voltage,
			   ts.[connectortype] AS Connectors
		FROM [Manufacturing].[dbo].[TestSpec] AS ts
		LEFT JOIN [Manufacturing].[dbo].[Band] AS b ON b.[band_id] = ts.[band_id]
		LEFT JOIN [Manufacturing].[dbo].[BandInfo] AS bi ON bi.[bandinfo_id] = b.[bandinfo_id]
		LEFT JOIN [Manufacturing].[dbo].[Direction] AS d ON b.[direction_id] = d.[direction_id]
		LEFT JOIN [Manufacturing].[dbo].[Bom] AS bm ON bm.[bom_id] = ts.[bom_id]
		LEFT JOIN [Manufacturing].[dbo].[AmpModel] AS am ON am.[ampmodel_id] = bm.[ampmodel_id]
		WHERE am.[model]=?1 AND bm.[name]=?2;`, queryParams.model, queryParams.bom)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	cols, err := data.Columns()
	if err != nil {
		return nil
	}
	if cols == nil {
		return nil
	}

	vals := make([]interface{}, len(cols))
	var columns = make([]string, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
		columns[i] = cols[i]
	}

	var jsonData []ampData
	i := 0
	for data.Next() {
		row := ampData{
			ID:             " ",
			Model:          " ",
			Bom:            " ",
			Band:           " ",
			Direction:      " ",
			TestName:       " ",
			LowLimit:       " ",
			UpLimit:        " ",
			LEeprom:        " ",
			UEeprom:        " ",
			IP:             " ",
			Target:         " ",
			LowerFrequency: " ",
			UpperFrequency: " ",
			Voltage:        " ",
			Connectors:     " ",
		}

		err = data.Scan(vals...)
		if err != nil {
			fmt.Println(err)
			continue
		}
		//for i := 0; i < len(vals); i += 15 {
		row.ID = strconv.Itoa(i)
		row.Model = parseValue(vals[0].(*interface{}))
		row.Bom = parseValue(vals[1].(*interface{}))
		row.Band = parseValue(vals[2].(*interface{}))
		row.Direction = parseValue(vals[3].(*interface{}))
		row.TestName = parseValue(vals[4].(*interface{}))
		row.LowLimit = parseValue(vals[5].(*interface{}))
		row.UpLimit = parseValue(vals[6].(*interface{}))
		row.LEeprom = parseValue(vals[7].(*interface{}))
		row.UEeprom = parseValue(vals[8].(*interface{}))
		row.IP = parseValue(vals[9].(*interface{}))
		row.Target = parseValue(vals[10].(*interface{}))
		row.LowerFrequency = parseValue(vals[11].(*interface{}))
		row.UpperFrequency = parseValue(vals[12].(*interface{}))
		row.Voltage = parseValue(vals[13].(*interface{}))
		row.Connectors = parseValue(vals[14].(*interface{}))
		//fmt.Println(row.Model)
		//APPEND ROW TO MAIN SLICE
		jsonData = append(jsonData, row)

		//}
		i++
	}
	if data.Err() != nil {
		return nil
	}

	//CONVERT TO JSON
	bytes, err := json.Marshal(jsonData)
	if err != nil {
		panic(err)
	}

	return bytes
}

func parseValue(pval *interface{}) string {
	switch v := (*pval).(type) {
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05.999")
	case string:
		return v
	case int64:
		return strconv.Itoa(int(v))
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)

	}

	return ""
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		fmt.Println("Port:", p)
		return ":" + p
	}
	fmt.Println("Port:8080")
	return ":8080"
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/search/{params}", GetAmpSpecs).Methods("GET")

	staticFileDirectory := http.Dir("./root")
	staticFileHandler := http.StripPrefix("/root/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/root/").Handler(staticFileHandler).Methods("GET")
	return r
}

func main() {

	fmt.Println("Server is running")
	r := newRouter()

	log.Fatal(http.ListenAndServe(getPort(), handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
