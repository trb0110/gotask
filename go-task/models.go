package gotask


type User struct {
	UserID       		string `json:"UserID"									xml:"UserID"`
	Username     		string `json:"Username"									xml:"Username"`
	Password     		string `json:"Password"									xml:"Password"`
	Role         		string `json:"Role"										xml:"Role"`
	PreferredHours 		string `json:"PreferredHours"							xml:"PreferredHours"`
	Token        		string `json:"Token"									xml:"Token"`
}

type Task struct {
	TaskId			string 								`json:"TaskId"								xml:"TaskId"`
	Timestamp 		string 								`json:"Timestamp"							xml:"Timestamp"`
	Description		string								`json:"Description"							xml:"Description"`
	Duration		string								`json:"Duration"							xml:"Duration"`
	UserID			string								`json:"UserID"								xml:"UserID"`
	Username		string								`json:"Username"							xml:"Username"`
}