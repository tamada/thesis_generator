package thesis_generator

type Thesis struct {
	Id          string      `json:"-"` // generate from author name and title
	Title       string      `json:"title"`
	Degree      string      `json:"degree"`
	Year        int         `json:"year"`
	Supervisors []string    `json:"supervisors"`
	Author      *Author     `json:"author"`
	Repository  *Repository `json:"repository"`
}

type Author struct {
	StudentId   string       `json:"student-id"`
	Name        string       `json:"name"`
	Affiliation *Affiliation `json:"affiliation"`
	Email       string       `json:"email"`
}

type Affiliation struct {
	University string `json:"university"`
	Department string `json:"department"`
}

type Format int

const (
	LaTeX Format = iota + 1
	HTML
	Markdown
	MicrosoftWord
)
