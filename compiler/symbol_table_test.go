package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]SymbolTableEntry{
		"a": SymbolTableEntry{Name: "a", Scope: GlobalScope, Index: 0},
		"b": SymbolTableEntry{Name: "b", Scope: GlobalScope, Index: 1},
	}
	global := NewSymbolTable()
	a := global.Set("a", GlobalScope)
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}
	b := global.Set("b", GlobalScope)
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}
}
func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Set("a", GlobalScope)
	global.Set("b", GlobalScope)
	expected := []SymbolTableEntry{
		SymbolTableEntry{Name: "a", Scope: GlobalScope, Index: 0},
		SymbolTableEntry{Name: "b", Scope: GlobalScope, Index: 1},
	}
	for _, sym := range expected {
		result, ok := global.Get(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				sym.Name, sym, result)
		}
	}
}
