package main

import (
	"fmt"
	"os"
    "errors"
	"strconv"
    "strings" 
	"time"

	"github.com/BurntSushi/toml"
	//"github.com/sqweek/dialog"
    //"gopkg.in/toast.v1"
)

type Meeting struct {
    Name string
    Date string
    DayOfWeek string
    Time string
}

type Calendar struct {
    Meetings []Meeting
}

func ReadCalendar(calPath string) (Calendar, error) {
	var cal Calendar
	f, err := os.Open(calPath)

	if err != nil {
	    return Calendar{}, errors.New(fmt.Sprintf("Could not open calendar %s", calPath))
	}

	if _, err := toml.DecodeReader(f, &cal); err != nil {
	    return Calendar{}, errors.New(fmt.Sprintf("Could not decod calendar file %s", calPath))
	}

	f.Close()
	return cal, nil
}


func main() {
	//duration, err := time.ParseDuration(os.Args[1])
	//if err != nil {
	//	panic("error: invalid duration: " +  os.Args[1])
	//}

    cal, err := ReadCalendar(os.Args[1])

    if err != nil {
        panic(err)
    }

    t := time.Now()
    st := 3600*t.Hour() + 60*t.Minute() + t.Second()

    fmt.Printf("Seconds elapsed in day: %d\n", st)

    for _, meeting := range cal.Meetings {
        fmt.Printf("Name: %s\nTime: %s\n\n", meeting.Name, meeting.Time)
        parts := strings.Split(meeting.Time, ":")
        h, _ := strconv.Atoi(parts[0])
        m, _ := strconv.Atoi(parts[1])
        sm := 3600*h + 60*m
        fmt.Printf("%d, %d\n", h, m)

        if st >= sm {
            fmt.Printf("Meeting %s is starting\n", meeting.Name)
        }
    }


	//secs := int(duration.Seconds())
	//// fmt.Println(secs)

	//if err != nil {
	//	fmt.Println(os.Args)
	//	panic("Invalid time specification")
	//}

	//for secs > 0 {
	//	mm := secs / 60
	//	ss := secs - 60*mm
	//	fmt.Printf("\r%02d:%02d", mm, ss)
	//	secs -= 1
	//	time.Sleep(time.Second)
	//}

	//fmt.Printf("\r%02d:%02d\n", 0, 0)
	//dialog.Message("%s", "Time's Up!").Title("GoTime").Info()
}
