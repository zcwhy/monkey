package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "Global"
)

type SymbolTableEntry struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store     map[string]SymbolTableEntry
	numSymbol int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:     make(map[string]SymbolTableEntry),
		numSymbol: 0,
	}
}

func (s *SymbolTable) Set(name string, scope SymbolScope) SymbolTableEntry {
	entry := SymbolTableEntry{
		Name:  name,
		Scope: scope,
		Index: s.numSymbol,
	}

	s.store[name] = entry
	s.numSymbol += 1
	return entry
}

func (s *SymbolTable) Get(name string) (SymbolTableEntry, bool) {
	entry, ok := s.store[name]
	return entry, ok
}
