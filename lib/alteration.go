package lib

// Alteration is general type of all SQL alterations, providing interface to get statements and diff strings.
type Alteration interface {
	// Statements returns SQL statements executing the alteration
	Statements() []string

	// Diff returns diff strings between before and after alteration
	Diff() []string

	// Id returns alteration ID.
	// ID must be unique among the same type of alterations (for example, table name for table alterations)
	Id() string

	// SeqNum returns sequential number of alterations.
	// This is supposed to be used on sort.
	SeqNum() int

	// SetSeqNum sets sequential number of alterations.
	SetSeqNum(int)

	// DependsOn returns dependent alterations.
	// This is supposed to be used on topological sort.
	DependsOn() []Alteration

	// AddDependsOn adds a dependent alteration.
	AddDependsOn(Alteration)

	Prefix() string

	SetPrefix(string)
}

type Prefixable struct {
	prefix string
}

func (r Prefixable) Prefix() string {
	return r.prefix
}

func (r *Prefixable) SetPrefix(s string) {
	r.prefix = s
}

// Sequential provides default implementation of Alteration.SeqNum
type Sequential struct {
	seqNum int
}

func (r Sequential) SeqNum() int {
	return r.seqNum
}

func (r *Sequential) SetSeqNum(n int) {
	r.seqNum = n
}

// Dependent provides default implementation of Alteration.DependsOn
type Dependent struct {
	dependsOn []Alteration
}

func (r Dependent) DependsOn() []Alteration {
	return r.dependsOn
}

func (r *Dependent) AddDependsOn(a Alteration) {
	r.dependsOn = append(r.dependsOn, a)
}
