package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "Global"
	LocalScope   SymbolScope = "Local"
	BuiltinScope SymbolScope = "Builtin"
)

type SymbolTableEntry struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	outer     *SymbolTable // 链式嵌套结构
	store     map[string]SymbolTableEntry
	numSymbol int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:     make(map[string]SymbolTableEntry),
		numSymbol: 0,
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		store:     make(map[string]SymbolTableEntry),
		numSymbol: 0,
		outer:     outer,
	}
}

func (s *SymbolTable) Define(name string) SymbolTableEntry {
	scope := LocalScope
	if s.outer == nil {
		scope = GlobalScope
	}

	entry := SymbolTableEntry{
		Name:  name,
		Scope: scope,
		Index: s.numSymbol,
	}

	s.store[name] = entry
	s.numSymbol += 1
	return entry
}

func (s *SymbolTable) DefineBuiltin(index int, name string) SymbolTableEntry {
	entry := SymbolTableEntry{
		Name:  name,
		Scope: BuiltinScope,
		Index: index,
	}

	s.store[name] = entry

	return entry
}

func (s *SymbolTable) Resolve(name string) (SymbolTableEntry, bool) {
	entry, ok := s.store[name]

	if !ok && s.outer != nil {
		entry, ok = s.outer.Resolve(name)
	}
	return entry, ok
}
