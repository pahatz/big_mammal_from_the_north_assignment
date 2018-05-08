package main

import (
    "net/http"
    "log"
    "fmt"
    "encoding/xml"
    "io/ioutil"
    "strconv"
    "time"
)

/*
 * struct representing xml structure of ELK response used to get the data 
 */
type Numbers struct {
	Number int `xml:"NUMBER"`
}

/* struct holding the data for one player
 * TotalCalls holds the total calls for each player
 * Calls variable holds a list of calls per player
 */
type Player struct {
	TotalCalls int
	Calls []int
}

var Player_map = make(map[int]Player)
var App_total_calls int //total calls for candidate app for all players
var Temp_number Numbers
var Last_num = 0 // holding query paramater of last call

func getXML(url string) ([]byte, error) {
    var client = &http.Client{
        Timeout: time.Second * 10,
    }
	resp, err := client.Get(url)
	if err != nil {
        return []byte{}, fmt.Errorf("GET error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
    	log.Fatal(err)
        return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
    }
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []byte{}, fmt.Errorf("Read body: %v", err)
    }
    return data, nil
}

/*
 * Function called on listening player's calls. First we get the
 * player's number from the url, if it's valid we make the call to
 * ELK and get the XML and we unmarshal to get the number in it.
 * if the map for the specific player doesn't exist, we create it 
 * and update Players calls and counter, if not, we create first.
 * we use a temporary struct to append data to global map Player_map
 */
func handler(w http.ResponseWriter, r *http.Request) {
	player_num := r.URL.Path[len("/player/"):]
	//fmt.Println(reflect.TypeOf(player_num))
    if len(player_num) > 0 { //if 
        pl_num, _ := strconv.Atoi(player_num)
        url := "https://testapi.elk-studios.com/test.php?n=" + strconv.Itoa(Last_num)
        //fmt.Println("url is: "+url)
        resp, err:= getXML(url)
        //fmt.Println("response from url: "+ string(resp))
        if (err!=nil){
        	log.Printf("Failed to get XML: %v", err)
        }
        xml.Unmarshal(resp, &Temp_number)
        //fmt.Println("this is the taken num: ", Temp_number.Number)
        Last_num = Temp_number.Number

        if tmp, ok := Player_map[pl_num]; !ok {
        	tmp.TotalCalls=1
        	tmp.Calls= []int{Temp_number.Number}
        	Player_map[pl_num] = tmp
        }else {
        	tmp.TotalCalls +=1
        	tmp.Calls = append(tmp.Calls, Temp_number.Number)
        	Player_map[pl_num] = tmp
        }  
        //fmt.Println("player calls: ",Player_map[pl_num].TotalCalls)
        //fmt.Println("player calls: ",Player_map[pl_num].Calls)
        App_total_calls +=1
        fmt.Fprintf(w, "number from response: %v ", Temp_number.Number)

    } else {
        //w.WriteHeader(http.StatusNotFound)
        log.Println(w, "Error 401: Invalid User")
    }
}

/*
 * function called on statistic endpoint. Prints all requested data
 */
func statistic (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Total calls:", App_total_calls)
	fmt.Fprintf(w, "Total calls: %v \n", App_total_calls)
	for k := range Player_map {
		fmt.Println("calls for user ", k, ":", Player_map[k].Calls)
		fmt.Fprintf(w, "calls for user %v:", k)
		fmt.Fprintf(w, "%v \n", Player_map[k].TotalCalls)
		fmt.Fprintf(w, "call list: %v \n", Player_map[k].Calls)
	}
	
}

func main() {
    http.HandleFunc("/player/", handler)
    http.HandleFunc("/statistic/", statistic)
    log.Println("Serving on 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}