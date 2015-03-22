package gotojava

type JBlockStmt struct {
	List []JStmt
}

func (j *JBlockStmt) Add(s ...JStmt) {
	j.List = append(j.List, s...)
}
