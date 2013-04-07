package main

var chAdd chan string
var chGet chan []*FileInfo
var chGetOver chan bool
var mapStat map[string]int64

func StatStart() {
	chAdd = make(chan string, 1000)
	chGet = make(chan []*FileInfo)
	chGetOver = make(chan bool)
	mapStat = make(map[string]int64)
	go statLoop()
}
func statLoop() {
	for {
		select {
		case key := <-chAdd:
			count, ok := mapStat[key]
			if ok {
				mapStat[key] = 1 + count
			} else {
				mapStat[key] = 1
			}
		case keys := <-chGet:
			for _, v := range keys {
				v.Count = mapStat[v.Name]
			}
			chGetOver <- true
		}
	}
}

func StatAdd(key string) {
	chAdd <- key
}
func StatGet(keys []*FileInfo) {
	chGet <- keys
	<-chGetOver
}
