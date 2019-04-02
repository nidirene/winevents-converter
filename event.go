package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type EventData struct {
	Data []Data `xml:"Data"`
}

type Data struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

type Provider struct {
	ID   string `xml:"Guid,attr"`
	Name string `xml:"Name,attr"`
}

type TimeCreated struct {
	SystemTime string `xml:"SystemTime,attr"`
}

type Correlation struct {
	ActivityID string `xml:"ActivityID,attr"`
}

type Execution struct {
	ProcessID string `xml:"ProcessID,attr"`
	ThreadID  string `xml:"ThreadID,attr"`
}
type Security struct {
	UserID string `xml:"UserID,attr,omitempty"`
}

type System struct {
	Provider      Provider    `xml:"Provider"`
	EventID       int         `xml:"EventID"`
	Version       int         `xml:"Version"`
	Level         int         `xml:"Level"`
	Task          int         `xml:"Task"`
	Opcode        int         `xml:"Opcode"`
	Keywords      string      `xml:"Keywords"`
	TimeCreated   TimeCreated `xml:"TimeCreated"`
	EventRecordID int         `xml:"EventRecordID"`
	Correlation   Correlation `xml:"Correlation"`
	Execution     Execution   `xml:"Execution"`
	Channel       string      `xml:"Channel"`
	Computer      string      `xml:"Computer"`
	Security      Security    `xml:"Security"`
}

type Event struct {
	System    System    `xml:"System"`
	EventData EventData `xml:"EventData"`
}

var levelTable = []string{"Verbose", "Informational", "Warning", "Error", "Critical"}
var loginTypeTable = []string{"", "", "Interactive", "Network", "Batch", "Service", "", "Unlock", "NetworkCleartext", "NewCredentials", "RemoteInteractive", "CachedInteractive"}

func convertToWinC(ev *Event) ([]byte, error) {

	kw, _ := strconv.ParseInt(ev.System.Keywords, 0, 64)
	jsys := map[string]string{
		"EventID":        strconv.Itoa(ev.System.EventID),
		"Channel":        string(ev.System.Channel),
		"Version":        strconv.Itoa(ev.System.Version),
		"ProviderName":   string(ev.System.Provider.Name),
		"ProviderID":     strings.Trim(ev.System.Provider.ID, "{}"),
		"Computer":       string(ev.System.Computer),
		"EventRecordID":  strconv.Itoa(ev.System.EventRecordID),
		"Keywords":       strconv.FormatInt(kw, 10),
		"Level":          levelTable[ev.System.Level],
		"Opcode":         strconv.Itoa(ev.System.Opcode),
		"ProcessID":      string(ev.System.Execution.ProcessID),
		"ThreadID":       string(ev.System.Execution.ThreadID),
		"Task":           strconv.Itoa(ev.System.Task),
		"RelatedAcivity": "",
		"Qualifiers":     "",
		"TimeCreated":    string(ev.System.TimeCreated.SystemTime),
		"UserId":         string(ev.System.Security.UserID),
	}

	sdata := make(map[string]string)
	for _, e := range ev.EventData.Data {
		sdata[e.Name] = e.Value
	}

	rc := map[string]map[string]string{
		"System":    jsys,
		"EventData": sdata,
	}
	return json.MarshalIndent(rc, "", "\t")
}
