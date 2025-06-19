package main

type Item struct {
	Id           StringInt `json:"id"`
	WeekNumber   StringInt `json:"weekNumber"`
	WeekDay      string    `json:"weekDay"`
	Time         string    `json:"time"`
	Group        string    `json:"group"`
	Teacher      string    `json:"teacher"`
	Subject      string    `json:"subject"`
	Subject_lvl2 string    `json:"subject_lvl2"`
	ClassRoom    string    `json:"class_room"`
}
