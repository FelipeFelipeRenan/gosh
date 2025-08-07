package history

var commands []string
var index int = -1

type History struct {
	entries []string
	pos     int
}

func New() *History {
	return &History{}
}

func (h *History) Add(command string) {
	if command == ""{
		return
	}
	h.entries = append(h.entries, command)
	h.pos = len(h.entries)
}

func (h *History) Prev() string {
	if len(h.entries) == 0 || h.pos <= 0 {
		return ""
	}
	h.pos--
	return h.entries[h.pos]
}

func (h *History) Next() string {
	if len(h.entries) == 0 || h.pos >= len(h.entries)-1 {
		h.pos = len(h.entries)
		return ""
	}

	h.pos++
	return h.entries[h.pos]
}

func (h *History) ResetPos(){
	h.pos = len(h.entries)
}