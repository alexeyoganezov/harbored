// Manages presentations stored in the database
// and broadcasts their changes to audience via WebSocket rooms
package presentations

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"harbored/config"
	"harbored/database"
	"harbored/models/presentation"
	"harbored/utils"
	"harbored/webserver"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// Get an array of all presentations from database
func Get() []*presentation.Presentation {
	txn := database.DB.Txn(false)
	defer txn.Abort()
	iterator, err := txn.Get("presentations", "id")
	if err != nil {
		panic(err)
	}
	var presentations []*presentation.Presentation
	for obj := iterator.Next(); obj != nil; obj = iterator.Next() {
		p := obj.(presentation.Presentation)
		presentations = append(presentations, &p)
	}
	return presentations
}

// Find pdf-files in path and create "presentation" records in database
func Load(path string) []presentation.Presentation {
	if path == "" {
		return make([]presentation.Presentation, 0)
	}
	presentations := make([]presentation.Presentation, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	var index = 0
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			pdfPath := config.Config.StaticDir + "/" + file.Name()
			utils.CopyFile(path+"/"+file.Name(), pdfPath)
			index += 1
			presentations = append(presentations, presentation.Presentation{
				ID:                index,
				Filename:          file.Name(),
				CurrentPageNumber: 0,
			})
		}
	}
	txn := database.DB.Txn(true)
	for _, presentation := range presentations {
		if err := txn.Insert("presentations", presentation); err != nil {
			panic(err)
		}
	}
	txn.Commit()
	return presentations
}

// Mark a presentation as started and broadcast a message
// Called on presentation opening.
func Start(presentationId int) {
	room, _ := webserver.RoomList.Get("global")
	var presentation presentation.Presentation
	txn := database.DB.Txn(true)
	raw, err := txn.First("presentations", "id", presentationId)
	if err != nil {
		panic(err)
	}
	mapstructure.Decode(raw, &presentation)
	presentation.IsOnline = true
	presentation.CurrentPageNumber = 0
	if err := txn.Insert("presentations", presentation); err != nil {
		fmt.Println(err)
	}
	txn.Commit()
	message := webserver.WSMessage{
		Command: "presentation:start",
		Payload: make(map[string]interface{}),
	}
	message.Payload["presentationId"] = presentationId
	wsjson, _ := json.Marshal(message)
	room.Broadcast <- wsjson
}

// Broadcast a message about switched slide
func ChangeSlide(presentationId int, pageNumber int) {
	room, _ := webserver.RoomList.Get("global")
	var presentation presentation.Presentation
	txn := database.DB.Txn(true)
	raw, err := txn.First("presentations", "id", presentationId)
	if err != nil {
		panic(err)
	}
	mapstructure.Decode(raw, &presentation)
	presentation.CurrentPageNumber = pageNumber
	if err := txn.Insert("presentations", presentation); err != nil {
		fmt.Println(err)
	}
	txn.Commit()
	message := webserver.WSMessage{
		Command: "presentation:slide-change",
		Payload: make(map[string]interface{}),
	}
	message.Payload["presentationId"] = presentationId
	message.Payload["pageNumber"] = pageNumber
	wsjson, _ := json.Marshal(message)
	room.Broadcast <- wsjson
}

// Mark a presentation as ended and broadcast a message
// Called on presentation switching.
func Stop(presentationId int) {
	room, _ := webserver.RoomList.Get("global")
	var presentation presentation.Presentation
	txn := database.DB.Txn(true)
	raw, err := txn.First("presentations", "id", presentationId)
	if err != nil {
		panic(err)
	}
	mapstructure.Decode(raw, &presentation)
	presentation.IsOnline = false
	if err := txn.Insert("presentations", presentation); err != nil {
		fmt.Println(err)
	}
	txn.Commit()
	message := webserver.WSMessage{
		Command: "presentation:stop",
		Payload: make(map[string]interface{}),
	}
	message.Payload["presentationId"] = presentationId
	wsjson, _ := json.Marshal(message)
	room.Broadcast <- wsjson
}

// Broadcast a message about presentation end.
// Called on application exit.
func End() {
	room, _ := webserver.RoomList.Get("global")
	message := webserver.WSMessage{
		Command: "presentation:end",
		Payload: make(map[string]interface{}),
	}
	wsjson, _ := json.Marshal(message)
	room.Broadcast <- wsjson
}
