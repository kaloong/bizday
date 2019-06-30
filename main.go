package main

import (
    "time"
    "fmt"
    "os"
    "flag"
    "log"
    "net/http"
)

func getPrevBizDay( w http.ResponseWriter, req *http.Request ) {
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

//Wrapper function of getBizDay
func getLastBizDay_cli(t time.Time) time.Time  {
    offset_list := []int { -2,-3,-1,-1,-1,-1,-1,0 }
    return getBizDay( offset_list, t )
}

//Wrapper function of getBizDay
func getNextBizDay_cli(t time.Time) time.Time  {
    offset_list := []int { 1,1,1,1,1,3,2,0 }
    return getBizDay( offset_list, t )
}

//Function to Business day base on offset
func getBizDay( offset_list []int , t time.Time ) time.Time {
    var offset = 0
    switch enquiry_date := time.Time(t).Weekday() ; {
    //switch enquiry_date := time.Now().Weekday() ; {
    // if it is Sunday
    case enquiry_date == 0 :
        fmt.Println("- case 0",enquiry_date+5)
        offset = offset_list[0] 
    // if it is Monday
    case enquiry_date == 1 :
        fmt.Println("- case 1",enquiry_date+4)
        offset = offset_list[1] 
    // if it is Tuesday
    case enquiry_date == 2 :
        fmt.Println("- case 2",enquiry_date-1)
        offset = offset_list[2] 
    // if it is Wednesday
    case enquiry_date == 3 :
        fmt.Println("- case 3",enquiry_date-1)
        offset = offset_list[3] 
    // if it is Thursday
    case enquiry_date == 4 :
        fmt.Println("- case 4",enquiry_date-1)
        offset = offset_list[4] 
    // if it is Friday
    case enquiry_date == 5 :
        fmt.Println("- case 5",enquiry_date-1)
        offset = offset_list[5] 
    // if it is Saturday
    case enquiry_date == 6 :
        fmt.Println("- case 6",enquiry_date-1)
        offset = offset_list[6] 
    default:
        fmt.Println("- case d",enquiry_date+1)
        offset = offset_list[7] 
    }
    return t.AddDate(0,0,offset)
}

var formats = [] string{"20060102", "2006/01/02", "2006-01-02", "2006.01.02"}
func main() {
    if len(os.Args) <= 1 {
        fmt.Printf("No command specified.\n")
        return
    }

    getCmd := flag.NewFlagSet("get", flag.ExitOnError)
    webCmd := flag.NewFlagSet("web", flag.ExitOnError)
    tzFlag := getCmd.String("t","Greenwich","Timezone name")
    bdFlag := getCmd.String("b","","Business Date name")
    hpFlag := getCmd.String("h","","Get help cli")
    wFlag  := webCmd.String("h","","Get help webserver")

    //formats := [] string{"20060102", "2006/01/02", "2006-01-02", "2006.01.02"}
    switch os.Args[1] {
    case "get":
        getCmd.Parse(os.Args[2:])
        if getCmd.Parsed() { 
            if *hpFlag != "" {
                log.Println("Get help page")
                os.Exit(2)
            }
            if *tzFlag == "" {
                log.Println("Please provide timezone.")
                return
            }
            if *bdFlag == "" {
                log.Println("Please provide business date.")
                return
            }
            log.Printf("Timezone selected: %q\n",*tzFlag)
            log.Printf("Bus Date selected: %q\n",*bdFlag)

        }
        loc, err := time.LoadLocation (*tzFlag)
        if err != nil {
             panic(err)
        }
        for i := range formats {
            if timeStamp, err := time.Parse(formats[i], *bdFlag); err == nil {
                //t := time.Now().Add(*bdFlag)
                prev_bizday := getLastBizDay_cli(timeStamp.In(loc))
                next_bizday := getNextBizDay_cli(timeStamp.In(loc))
                log.Printf("- Previous business day : %v\t%v\n", prev_bizday.Weekday(), prev_bizday.Format("2006-01-02 MST"))
                log.Printf("-     Next business day : %v\t%v\n", next_bizday.Weekday(), next_bizday.Format("2006-01-02 MST"))
            }
        }
    case "web":
        webCmd.Parse(os.Args[2:])
        if webCmd.Parsed() { 
            if *wFlag != "" {
                log.Println("Get help page")
                os.Exit(2)
            }
        }
        http.HandleFunc("/PrevBizDay", getPrevBizDay)
        http.HandleFunc("/NextBizDay", getNextBizDay)
        log.Fatal(http.ListenAndServe("localhost:8000",nil))
    default:
        fmt.Printf("%q is not valid command.\n", os.Args[1])
        os.Exit(2)
    }
}
