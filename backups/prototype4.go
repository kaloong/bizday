package main

import (
    "net/http"
    "fmt"
    "log"
    "time"
)

func main() {
    http.HandleFunc("/PrevBizDay", getPrevBizDay)
    http.HandleFunc("/NextBizDay", getNextBizDay)
    log.Fatal(http.ListenAndServe("localhost:8000",nil))

}

//Wrapper function of getBizDay
func getLastBizDay(t time.Time) time.Time  {
    offset_list := []int { -2,-3,-1,-1,-1,-1,-1,0 }
    return getBizDay( offset_list, t )
}

//Wrapper function of getBizDay
func getNextBizDay(t time.Time) time.Time  {
    offset_list := []int { 1,1,1,1,1,3,2,0 }
    return getBizDay( offset_list, t )
}

func getPrevBizDay( w http.ResponseWriter, req *http.Request ) {
    formats := [] string{"20060102", "2006/01/02", "2006-01-02", "2006.01.02"}
    offset_list := []int { -2,-3,-1,-1,-1,-1,-1,0 }
    enquiry_date := req.URL.Query().Get("date")
    tz := req.URL.Query().Get("tz")
    loc, err := time.LoadLocation(tz)
    if err != nil {
         panic(err)
    }
    for i := range formats  {
        if timeStamp, err := time.Parse(formats[i], enquiry_date);  err  == nil {
            var bizDay = getBizDay( offset_list, timeStamp.In(loc) )
            fmt.Fprintf(w, "The previous business day of: %v was: %v %v\n", enquiry_date, bizDay.Weekday(), bizDay.Format("2006-01-02 MST"))
        }
    }
}

func getNextBizDay( w http.ResponseWriter, req *http.Request ) {
    formats := [] string{"20060102", "2006/01/02", "2006-01-02", "2006.01.02"}
    offset_list := []int { 1,1,1,1,1,3,2,0 }
    enquiry_date := req.URL.Query().Get("date")
    tz := req.URL.Query().Get("tz")
    loc, err := time.LoadLocation(tz)
    if err != nil {
         panic(err)
    }
    for i := range formats  {
        if timeStamp, err := time.Parse(formats[i], enquiry_date);  err  == nil {
            var bizDay = getBizDay( offset_list, timeStamp.In(loc) )
            fmt.Fprintf(w, "The next business day of: %v is: %v %v\n", enquiry_date, bizDay.Weekday(), bizDay.Format("2006-01-02 MST"))
        }
    }
}

//Function to Business day base on offset
func getBizDay( offset_list []int , t time.Time ) time.Time {
    var offset = 0
    switch enquiry_date := time.Time(t).Weekday() ; {
    //switch enquiry_date := time.Now().Weekday() ; {
    // if it is Sunday
    case enquiry_date == 0 :
        log.Println("- case 0",enquiry_date+5)
        offset = offset_list[0]
    // if it is Monday
    case enquiry_date == 1 :
        log.Println("- case 1",enquiry_date+4)
        offset = offset_list[1]
    // if it is Tuesday
    case enquiry_date == 2 :
        log.Println("- case 2",enquiry_date-1)
        offset = offset_list[2]
    // if it is Wednesday
    case enquiry_date == 3 :
        log.Println("- case 3",enquiry_date-1)
        offset = offset_list[3]
    // if it is Thursday
    case enquiry_date == 4 :
        log.Println("- case 4",enquiry_date-1)
        offset = offset_list[4]
    // if it is Friday
    case enquiry_date == 5 :
        log.Println("- case 5",enquiry_date-1)
        offset = offset_list[5]
    // if it is Saturday
    case enquiry_date == 6 :
        log.Println("- case 6",enquiry_date-1)
        offset = offset_list[6]
    default:
        log.Println("- case d",enquiry_date+1)
        offset = offset_list[7]
    }
    return t.AddDate(0,0,offset)
}

