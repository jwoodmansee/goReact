package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ampData struct {
	ID             string `json:"id"`
	Model          string `json:"Model"`
	Bom            string `json:"Bom"`
	Band           string `json:"Band"`
	Direction      string `json: "Direction"`
	TestName       string `json: "TestName"`
	LowerLimit     string `json: "LowerLimit"`
	UpperLimit     string `json: "UpperLimit"`
	LowerEeprom    string `json: "LowerEeprom"`
	UpperEeprom    string `json: "UpperEeprom"`
	IP             string `json: "Ip"`
	Target         string `json: "Target"`
	LowerFrequency string `json: "LowerFrequency"`
	UpperFrequency string `json: "UpperFrequency"`
	Voltage        string `json: "Voltage"`
	Connectors     string `json: "Connectors"`
}

type Specs []ampData

type queryParams struct {
	model string
	bom   string
}

func GetAmpSpecs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	//Change the varible to something better
	serial := strings.Split(r.URL.Path, "/")
	//ampInfo := r.URL.Query().Get("params")

	if len(serial[2]) > 0 {
		ampBom := serial[2]
		fmt.Println(ampBom)
		w.Write(sqlQuery(ampBom))
	}

}

func sqlQuery(ampBom string) []byte {
	queryParams := queryParams{
		model: ampBom[0 : len(ampBom)-1],
		bom:   ampBom[len(ampBom)-1:],
	}
	log.Println(queryParams)

	log.Println("Enter SqlQuery")
	dsn := "server=dbcluster1;user id=testlogin;password=t3stin9;database=WeBooster"
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return nil
	}
	defer db.Close()

	data, err := db.Query(`
		SELECT am.model AS Model, bm.name AS Bom, bi.bandnumber AS Band, d.name AS Direction, ts.name AS TestName, ts.lowerlimit AS LowerLimit, ts.upperlimit AS UpperLimit, ts.eeprom_lowerlimit AS LowerEeprom, ts.eeprom_upperlimit AS UpperEeprom, ts.inputpower AS IP, ts.target AS Target, ts.lowerfrequency AS Lfrequency, ts.upperfrequency AS Ufrequency, ts.operatingvoltage AS Voltage, ts.connectortype AS Connectors
		FROM [WeBooster].[dbo].[TestSpec] AS ts
		LEFT JOIN [WeBooster].[dbo].Band AS b ON b.band_id = ts.band_id
		LEFT JOIN [WeBooster].[dbo].BandInfo AS bi ON bi.bandinfo_id = b.bandinfo_id
		LEFT JOIN [WeBooster].[dbo].Direction AS d ON b.direction_id = d.direction_id
		LEFT JOIN [WeBooster].[dbo].Bom AS bm ON bm.bom_id = ts.bom_id
		LEFT JOIN [WeBooster].[dbo].AmpModel AS am ON am.ampmodel_id = bm.ampmodel_id
		WHERE am.model=?1 AND bm.name=?2;`, queryParams.model, queryParams.bom)
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
			LowerLimit:     " ",
			UpperLimit:     " ",
			LowerEeprom:    " ",
			UpperEeprom:    " ",
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
		row.LowerLimit = parseValue(vals[4].(*interface{}))
		row.UpperLimit = parseValue(vals[5].(*interface{}))
		row.LowerEeprom = parseValue(vals[6].(*interface{}))
		row.UpperEeprom = parseValue(vals[7].(*interface{}))
		row.TestName = parseValue(vals[8].(*interface{}))
		row.IP = parseValue(vals[9].(*interface{}))
		row.Target = parseValue(vals[10].(*interface{}))
		row.LowerFrequency = parseValue(vals[11].(*interface{}))
		row.UpperFrequency = parseValue(vals[12].(*interface{}))
		row.Voltage = parseValue(vals[13].(*interface{}))
		row.Connectors = parseValue(vals[14].(*interface{}))
		fmt.Println(i)
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

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}
