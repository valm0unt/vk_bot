package notes

import (
	"strconv"
)

type notesList []string

var (
	notes notesList
)

func Init() {
	notes = make(notesList, 0)
}

func Add(text string) {
	if len(text) > 200 {
		text = text[:200]
	}
	notes = append(notes, text)
}

func List() string {
	result := ""
	for id, text := range notes {
		result += strconv.Itoa(id) + ": " + text + "\n"
	}

	return result
}

func Del(id int) bool {
	if len(notes) > 0 && id >= 0 && id <= len(notes) {
		notes = append(notes[:id], notes[id+1:]...)
		return true
	}

	return false
}
