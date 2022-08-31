package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jison/uni/core/valuer"
	"github.com/jison/uni/internal/location"
	"github.com/stretchr/testify/assert"
)

type testStructForFuncProviderTest struct {
	a int
	b string
	c []int
	d rune
}

func funcForFuncProviderTest(a int, b string, c []int, d rune) (*testStructForFuncProviderTest, int, error) {
	return &testStructForFuncProviderTest{a, b, c, d}, 0, nil
}

func TestFunc(t *testing.T) {
	t.Run("normal value", func(t *testing.T) {
		type testInterface interface{}

		tag1 := NewSymbol("tag1")
		tag2 := NewSymbol("tag2")
		scope1 := NewScope("scope1")

		baseLoc := location.GetCallLocation(0)
		fp := Func(funcForFuncProviderTest,
			InScope(scope1),
			Param(0, ByName("a1"), ByTags(tag1), Optional(true)),
			Param(1, ByName("b1")),
			Param(2, ByName("c1"), AsCollector(true)),
			Param(3),
			Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil)))),
			Return(1),
			nil,
		)

		pro := fp.Provider()
		err := pro.Validate()
		assert.Nil(t, err)

		deps := dependencyIteratorToArray(pro.Dependencies())
		assert.Equal(t, 4, len(deps))
		for _, dep := range deps {
			if dep.Name() == "a1" {
				assert.Equal(t, TypeOf(1), dep.Type())
				assert.Equal(t, "a1", dep.Name())
				assert.Equal(t, newSymbolSet(tag1), dep.Tags())
				assert.Equal(t, valuer.Param(0), dep.Valuer())
				assert.True(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "b1" {
				assert.Equal(t, TypeOf(""), dep.Type())
				assert.Equal(t, "b1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(1), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "c1" {
				assert.Equal(t, TypeOf(0), dep.Type())
				assert.Equal(t, "c1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(2), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.True(t, dep.IsCollector())
			} else {
				assert.Equal(t, TypeOf('d'), dep.Type())
				assert.Equal(t, "", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(3), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			}

			assert.Same(t, pro, dep.Consumer())
		}

		coms := pro.Components().ToArray()
		assert.Equal(t, 2, len(coms))
		for _, com := range coms {
			if com.Name() == "r1" {
				assert.Equal(t, TypeOf(&testStructForFuncProviderTest{}), com.Type())
				assert.Equal(t, "r1", com.Name())
				assert.Equal(t, newSymbolSet(tag2), com.Tags())
				assert.Equal(t, newTypeSet(TypeOf((*testInterface)(nil))), com.As())
				assert.Equal(t, true, com.Ignored())
				assert.Equal(t, true, com.Hidden())
				assert.Equal(t, valuer.Index(0), com.Valuer())
			} else {
				assert.Equal(t, TypeOf(1), com.Type())
				assert.Equal(t, "", com.Name())
				assert.Equal(t, (*symbolSet)(nil), com.Tags())
				assert.Equal(t, (*typeSet)(nil), com.As())
				assert.Equal(t, false, com.Ignored())
				assert.Equal(t, false, com.Hidden())
				assert.Equal(t, valuer.Index(1), com.Valuer())
			}

			assert.Same(t, pro, com.Provider())
		}

		assert.Equal(t, valuer.Func(reflect.ValueOf(funcForFuncProviderTest)), pro.Valuer())
		assert.Equal(t, baseLoc.FileName(), pro.Location().FileName())
		assert.Equal(t, baseLoc.FileLine()+1, pro.Location().FileLine())
		assert.Equal(t, scope1, pro.Scope())
	})

	t.Run("nil function", func(t *testing.T) {
		type testInterface interface{}

		tag1 := NewSymbol("tag1")
		tag2 := NewSymbol("tag2")
		scope1 := NewScope("scope1")

		fp := Func(nil,
			InScope(scope1),
			Param(0, ByName("a1"), ByTags(tag1), Optional(true)),
			Param(1, ByName("b1")),
			Param(2, ByName("c1"), AsCollector(true)),
			Param(3),
			Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil)))),
			Return(1),
			nil,
		)

		pro := fp.Provider()
		err := pro.Validate()
		assert.NotNil(t, err)
	})
}

func Test_funcProvider_Provider(t *testing.T) {
	type testInterface interface{}

	tag1 := NewSymbol("tag1")
	tag2 := NewSymbol("tag2")
	scope1 := NewScope("scope1")

	fp0 := funcProviderOf(funcForFuncProviderTest,
		InScope(scope1),
		Param(0, ByName("a1"), ByTags(tag1), Optional(true)),
		Param(1, ByName("b1")),
		Param(2, ByName("c1"), AsCollector(true)),
		Param(3),
		Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil)))),
		Return(1),
	)

	t.Run("Dependencies", func(t *testing.T) {
		fp := fp0.clone()
		pro := fp.Provider()

		deps := dependencyIteratorToArray(pro.Dependencies())
		assert.Equal(t, 4, len(deps))
		for _, dep := range deps {
			if dep.Name() == "a1" {
				assert.Equal(t, TypeOf(1), dep.Type())
				assert.Equal(t, "a1", dep.Name())
				assert.Equal(t, newSymbolSet(tag1), dep.Tags())
				assert.Equal(t, valuer.Param(0), dep.Valuer())
				assert.True(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "b1" {
				assert.Equal(t, TypeOf(""), dep.Type())
				assert.Equal(t, "b1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(1), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "c1" {
				assert.Equal(t, TypeOf(0), dep.Type())
				assert.Equal(t, "c1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(2), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.True(t, dep.IsCollector())
			} else {
				assert.Equal(t, TypeOf('d'), dep.Type())
				assert.Equal(t, "", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(3), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			}
			assert.True(t, dep.Valuer() == dep.Valuer())
			assert.Same(t, pro, dep.Consumer())
		}
	})

	t.Run("Valuer", func(t *testing.T) {
		fp := fp0.clone()
		pro := fp.Provider()
		assert.Equal(t, valuer.Func(reflect.ValueOf(funcForFuncProviderTest)), pro.Valuer())
		assert.True(t, pro.Valuer() == pro.Valuer())
	})

	t.Run("Location", func(t *testing.T) {
		fp := fp0.clone()
		loc1 := location.GetCallLocation(0)
		fp.SetLocation(loc1)
		assert.Equal(t, loc1, fp.Provider().Location())
	})

	t.Run("Components", func(t *testing.T) {
		fp := fp0.clone()
		pro := fp.Provider()
		coms := pro.Components().ToArray()
		assert.Equal(t, 2, len(coms))
		for _, com := range coms {
			if com.Name() == "r1" {
				assert.Equal(t, TypeOf(&testStructForFuncProviderTest{}), com.Type())
				assert.Equal(t, "r1", com.Name())
				assert.Equal(t, newSymbolSet(tag2), com.Tags())
				assert.Equal(t, newTypeSet(TypeOf((*testInterface)(nil))), com.As())
				assert.Equal(t, true, com.Ignored())
				assert.Equal(t, true, com.Hidden())
				assert.Equal(t, valuer.Index(0), com.Valuer())
			} else {
				assert.Equal(t, TypeOf(1), com.Type())
				assert.Equal(t, "", com.Name())
				assert.Equal(t, (*symbolSet)(nil), com.Tags())
				assert.Equal(t, (*typeSet)(nil), com.As())
				assert.Equal(t, false, com.Ignored())
				assert.Equal(t, false, com.Hidden())
				assert.Equal(t, valuer.Index(1), com.Valuer())
			}
			assert.True(t, com.Valuer() == com.Valuer())
			assert.Same(t, pro, com.Provider())
		}
	})

	t.Run("Scope", func(t *testing.T) {
		t.Run("scope", func(t *testing.T) {
			fp := fp0.clone()
			pro := fp.Provider()
			assert.Equal(t, scope1, pro.Scope())
		})

		t.Run("scope", func(t *testing.T) {
			fp := fp0.clone()
			fp.SetScope(nil)
			pro := fp.Provider()
			assert.Equal(t, GlobalScope, pro.Scope())
		})
	})

	t.Run("Validate", func(t *testing.T) {
		t.Run("val is not a function", func(t *testing.T) {
			intVal := 1
			strVal := ""
			runeVal := 'a'

			tests := []struct {
				name string
				t    reflect.Type
			}{
				{"int", TypeOf(1)},
				{"string", TypeOf("")},
				{"rune", TypeOf('a')},
				{"[]int", TypeOf([]int{})},
				{"[]string", TypeOf([]string{})},
				{"[]rune", TypeOf([]rune{})},
				{"*int", TypeOf(&intVal)},
				{"*string", TypeOf(&strVal)},
				{"*rune", TypeOf(&runeVal)},
				{"map[int]string", TypeOf(map[int]string{})},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					err := Func(tt.t).Provider().Validate()
					assert.NotNil(t, err)
				})
			}
		})

		t.Run("parameter not exist", func(t *testing.T) {
			con := Func(funcForFuncProviderTest, Param(4, ByName("p4"))).Provider()
			err := con.Validate()
			assert.NotNil(t, err)
		})

		t.Run("error in parameters", func(t *testing.T) {
			con := Func(funcForFuncProviderTest, Param(0, AsCollector(true))).Provider()
			err := con.Validate()
			assert.NotNil(t, err)
		})

		t.Run("return index not exist", func(t *testing.T) {
			con := Func(funcForFuncProviderTest, Return(3, Name("p4"))).Provider()
			err := con.Validate()
			assert.NotNil(t, err)
		})

		t.Run("return index invalid", func(t *testing.T) {
			con := Func(funcForFuncProviderTest, Return(2, Name("p2"))).Provider()
			err := con.Validate()
			assert.NotNil(t, err)
		})

		t.Run("no valid value return", func(t *testing.T) {
			func1 := func() error { return nil }
			func2 := func() {}
			con1 := Func(func1).Provider()
			err1 := con1.Validate()
			assert.NotNil(t, err1)

			con2 := Func(func2).Provider()
			err2 := con2.Validate()
			assert.NotNil(t, err2)
		})

		t.Run("error in components", func(t *testing.T) {
			con := Func(funcForFuncProviderTest, Return(0, As((*error)(nil)))).Provider()
			err := con.Validate()
			assert.NotNil(t, err)
		})
	})

	t.Run("Format", func(t *testing.T) {
		fp := fp0.clone()
		pro := fp.Provider()

		t.Run("not verbose", func(t *testing.T) {
			expected := fmt.Sprintf("Function[%v] in %v", TypeOf(funcForFuncProviderTest), fp.Scope())
			assert.Equal(t, expected, fmt.Sprintf("%v", pro))
		})

		t.Run("verbose", func(t *testing.T) {
			expected := fmt.Sprintf("Function[%v] in %v at %v", TypeOf(funcForFuncProviderTest),
				fp.Scope(), pro.Location())
			assert.Equal(t, expected, fmt.Sprintf("%+v", pro))
		})
	})
}

func Test_funcProvider_StructProviderBuilder(t *testing.T) {
	type testInterface interface{}

	tag1 := NewSymbol("tag1")
	tag2 := NewSymbol("tag2")
	scope1 := NewScope("scope1")
	loc1 := location.GetCallLocation(0)

	fp0 := funcProviderOf(funcForFuncProviderTest)

	t.Run("ApplyModule", func(t *testing.T) {
		fp := fp0.clone()
		fp.Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil))))
		fp.Return(1)

		mb := NewModuleBuilder()
		fp.ApplyModule(mb)

		coms := mb.Module().AllComponents().ToArray()
		assert.Equal(t, 2, len(coms))
		for _, com := range coms {
			if com.Name() == "r1" {
				assert.Equal(t, TypeOf(&testStructForFuncProviderTest{}), com.Type())
				assert.Equal(t, "r1", com.Name())
				assert.Equal(t, newSymbolSet(tag2), com.Tags())
				assert.Equal(t, newTypeSet(TypeOf((*testInterface)(nil))), com.As())
				assert.Equal(t, true, com.Ignored())
				assert.Equal(t, true, com.Hidden())
				assert.Equal(t, valuer.Index(0), com.Valuer())
			} else {
				assert.Equal(t, TypeOf(1), com.Type())
				assert.Equal(t, "", com.Name())
				assert.Equal(t, (*symbolSet)(nil), com.Tags())
				assert.Equal(t, (*typeSet)(nil), com.As())
				assert.Equal(t, false, com.Ignored())
				assert.Equal(t, false, com.Hidden())
				assert.Equal(t, valuer.Index(1), com.Valuer())
			}
		}
	})

	t.Run("Provide", func(t *testing.T) {
		fp := fp0.clone()
		fp.Param(0, ByName("a1"), ByTags(tag1), Optional(true))
		fp.Param(1, ByName("b1"))
		fp.Param(2, ByName("c1"), AsCollector(true))
		fp.Param(3)
		fp.Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil))))
		fp.Return(1)
		fp.SetScope(scope1)
		fp.SetLocation(loc1)

		pro := fp.Provider()
		err := pro.Validate()
		assert.Nil(t, err)

		deps := dependencyIteratorToArray(pro.Dependencies())
		assert.Equal(t, 4, len(deps))
		for _, dep := range deps {
			if dep.Name() == "a1" {
				assert.Equal(t, TypeOf(1), dep.Type())
				assert.Equal(t, "a1", dep.Name())
				assert.Equal(t, newSymbolSet(tag1), dep.Tags())
				assert.Equal(t, valuer.Param(0), dep.Valuer())
				assert.True(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "b1" {
				assert.Equal(t, TypeOf(""), dep.Type())
				assert.Equal(t, "b1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(1), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "c1" {
				assert.Equal(t, TypeOf(0), dep.Type())
				assert.Equal(t, "c1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(2), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.True(t, dep.IsCollector())
			} else {
				assert.Equal(t, TypeOf('d'), dep.Type())
				assert.Equal(t, "", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(3), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			}

			assert.Same(t, pro, dep.Consumer())
		}

		coms := pro.Components().ToArray()
		assert.Equal(t, 2, len(coms))
		for _, com := range coms {
			if com.Name() == "r1" {
				assert.Equal(t, TypeOf(&testStructForFuncProviderTest{}), com.Type())
				assert.Equal(t, "r1", com.Name())
				assert.Equal(t, newSymbolSet(tag2), com.Tags())
				assert.Equal(t, newTypeSet(TypeOf((*testInterface)(nil))), com.As())
				assert.Equal(t, true, com.Ignored())
				assert.Equal(t, true, com.Hidden())
				assert.Equal(t, valuer.Index(0), com.Valuer())
			} else {
				assert.Equal(t, TypeOf(1), com.Type())
				assert.Equal(t, "", com.Name())
				assert.Equal(t, (*symbolSet)(nil), com.Tags())
				assert.Equal(t, (*typeSet)(nil), com.As())
				assert.Equal(t, false, com.Ignored())
				assert.Equal(t, false, com.Hidden())
				assert.Equal(t, valuer.Index(1), com.Valuer())
			}

			assert.Same(t, pro, com.Provider())
		}

		assert.Equal(t, valuer.Func(reflect.ValueOf(funcForFuncProviderTest)), pro.Valuer())
		assert.Equal(t, loc1, pro.Location())
		assert.Equal(t, scope1, pro.Scope())
	})

	t.Run("Param", func(t *testing.T) {
		t.Run("Optional", func(t *testing.T) {
			fp := fp0.clone()
			fp.Param(0, ByName("p0"), Optional(true))
			meet := false
			fp.Provider().Dependencies().Iterate(func(dep Dependency) bool {
				if dep.Name() == "p0" {
					meet = true
					assert.True(t, dep.Optional())
				}
				return true
			})
			assert.True(t, meet)
		})

		t.Run("AsCollector", func(t *testing.T) {
			fp := fp0.clone()
			fp.Param(2, ByName("p2"), AsCollector(true))
			meet := false
			fp.Provider().Dependencies().Iterate(func(dep Dependency) bool {
				if dep.Name() == "p2" {
					meet = true
					assert.True(t, dep.IsCollector())
				}
				return true
			})
			assert.True(t, meet)
		})

		t.Run("Name", func(t *testing.T) {
			fp := fp0.clone()
			fp.Param(0, ByName("p0"))
			meet := false
			fp.Provider().Dependencies().Iterate(func(dep Dependency) bool {
				if dep.Name() == "p0" {
					meet = true
				}
				return true
			})
			assert.True(t, meet)
		})

		t.Run("Tags", func(t *testing.T) {
			fp := fp0.clone()
			fp.Param(0, ByName("p0"), ByTags(tag1))
			meet := false
			fp.Provider().Dependencies().Iterate(func(dep Dependency) bool {
				if dep.Name() == "p0" {
					meet = true
					assert.Equal(t, newSymbolSet(tag1), dep.Tags())
				}
				return true
			})
			assert.True(t, meet)
		})

		t.Run("param not exist", func(t *testing.T) {
			fp := fp0.clone()
			fp.Param(4, ByName("p4"))
			meet := false
			fp.Provider().Dependencies().Iterate(func(dep Dependency) bool {
				if dep.Name() == "p0" {
					meet = true
				}
				return true
			})
			assert.False(t, meet)
		})
	})

	t.Run("Return", func(t *testing.T) {
		t.Run("SetIgnore", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(0, Name("r1"), Ignore())
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r1" {
					meet = true
					assert.True(t, com.Ignored())
				}
				return true
			})

			assert.True(t, meet)
		})

		t.Run("SetHidden", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(0, Name("r1"), Hide())
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r1" {
					meet = true
					assert.True(t, com.Hidden())
				}
				return true
			})

			assert.True(t, meet)
		})

		t.Run("AddAs", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(0, Name("r1"), As(TypeOf((*testInterface)(nil))))
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r1" {
					meet = true
					assert.Equal(t, newTypeSet(TypeOf((*testInterface)(nil))), com.As())
				}
				return true
			})

			assert.True(t, meet)
		})

		t.Run("SetName", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(0, Name("r1"))
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r1" {
					meet = true
				}
				return true
			})

			assert.True(t, meet)
		})

		t.Run("AddTags", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(0, Name("r1"), Tags(tag1))
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r1" {
					meet = true
					assert.Equal(t, newSymbolSet(tag1), com.Tags())
				}
				return true
			})

			assert.True(t, meet)
		})

		t.Run("return not exist", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(3, Name("r3"))
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r3" {
					meet = true
				}
				return true
			})

			assert.False(t, meet)
		})

		t.Run("return not valid", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(2, Name("r2"))
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r2" {
					meet = true
				}
				return true
			})

			assert.False(t, meet)
		})

		t.Run("option is nil", func(t *testing.T) {
			fp := fp0.clone()
			fp.Return(0, Name("r1"), nil)
			meet := false
			fp.Provider().Components().Iterate(func(com Component) bool {
				if com.Name() == "r1" {
					meet = true
				}
				return true
			})

			assert.True(t, meet)
		})
	})

	t.Run("SetScope", func(t *testing.T) {
		t.Run("set scope", func(t *testing.T) {
			fp := fp0.clone()
			fp.SetScope(scope1)
			assert.Equal(t, scope1, fp.Provider().Scope())
		})

		t.Run("set nil scope", func(t *testing.T) {
			fp := fp0.clone()
			fp.SetScope(nil)
			assert.Equal(t, GlobalScope, fp.Provider().Scope())
		})

	})

	t.Run("SetLocation", func(t *testing.T) {
		fp := fp0.clone()
		fp.SetLocation(loc1)
		assert.Equal(t, loc1, fp.Provider().Location())
	})

	t.Run("UpdateCallLocation", func(t *testing.T) {
		t.Run("location have been set", func(t *testing.T) {
			fp := fp0.clone()
			fp.SetLocation(loc1)
			fp.UpdateCallLocation(nil)
			assert.Equal(t, loc1, fp.Location())
		})

		t.Run("location have not been set", func(t *testing.T) {
			fp := fp0.clone()
			baseLoc := location.GetCallLocation(0)
			func() {
				fp.UpdateCallLocation(nil)
			}()
			assert.Equal(t, baseLoc.FileName(), fp.Location().FileName())
			assert.Equal(t, baseLoc.FileLine()+3, fp.Location().FileLine())
		})

		t.Run("location is not nil", func(t *testing.T) {
			loc := location.GetCallLocation(0)
			fp := fp0.clone()
			fp.UpdateCallLocation(loc)
			assert.Equal(t, loc, fp.Location())
		})
	})
}

func Test_funcProvider_clone(t *testing.T) {
	type testInterface interface{}

	tag1 := NewSymbol("tag1")
	tag2 := NewSymbol("tag2")
	scope1 := NewScope("scope1")
	loc1 := location.GetCallLocation(0)

	fp := funcProviderOf(funcForFuncProviderTest,
		InScope(scope1),
		Param(0, ByName("a1"), ByTags(tag1), Optional(true)),
		Param(1, ByName("b1")),
		Param(2, ByName("c1"), AsCollector(true)),
		Param(3),
		Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil)))),
		Return(1),
	)
	fp.SetLocation(loc1)

	verifyProvider := func(t *testing.T, pro Provider) {
		err := pro.Validate()
		assert.Nil(t, err)

		deps := dependencyIteratorToArray(pro.Dependencies())
		assert.Equal(t, 4, len(deps))
		for _, dep := range deps {
			if dep.Name() == "a1" {
				assert.Equal(t, TypeOf(1), dep.Type())
				assert.Equal(t, "a1", dep.Name())
				assert.Equal(t, newSymbolSet(tag1), dep.Tags())
				assert.Equal(t, valuer.Param(0), dep.Valuer())
				assert.True(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "b1" {
				assert.Equal(t, TypeOf(""), dep.Type())
				assert.Equal(t, "b1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(1), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			} else if dep.Name() == "c1" {
				assert.Equal(t, TypeOf(0), dep.Type())
				assert.Equal(t, "c1", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(2), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.True(t, dep.IsCollector())
			} else {
				assert.Equal(t, TypeOf('d'), dep.Type())
				assert.Equal(t, "", dep.Name())
				assert.Equal(t, (*symbolSet)(nil), dep.Tags())
				assert.Equal(t, valuer.Param(3), dep.Valuer())
				assert.False(t, dep.Optional())
				assert.False(t, dep.IsCollector())
			}

			assert.Same(t, pro, dep.Consumer())
		}

		coms := pro.Components().ToArray()
		assert.Equal(t, 2, len(coms))
		for _, com := range coms {
			if com.Name() == "r1" {
				assert.Equal(t, TypeOf(&testStructForFuncProviderTest{}), com.Type())
				assert.Equal(t, "r1", com.Name())
				assert.Equal(t, newSymbolSet(tag2), com.Tags())
				assert.Equal(t, newTypeSet(TypeOf((*testInterface)(nil))), com.As())
				assert.Equal(t, true, com.Ignored())
				assert.Equal(t, true, com.Hidden())
				assert.Equal(t, valuer.Index(0), com.Valuer())
			} else {
				assert.Equal(t, TypeOf(1), com.Type())
				assert.Equal(t, "", com.Name())
				assert.Equal(t, (*symbolSet)(nil), com.Tags())
				assert.Equal(t, (*typeSet)(nil), com.As())
				assert.Equal(t, false, com.Ignored())
				assert.Equal(t, false, com.Hidden())
				assert.Equal(t, valuer.Index(1), com.Valuer())
			}

			assert.Same(t, pro, com.Provider())
		}

		assert.Equal(t, valuer.Func(reflect.ValueOf(funcForFuncProviderTest)), pro.Valuer())
		assert.Equal(t, loc1, pro.Location())
		assert.Equal(t, scope1, pro.Scope())
	}

	t.Run("equality", func(t *testing.T) {
		fp2 := fp.clone()
		verifyProvider(t, fp2.Provider())
		assert.False(t, fp2.Valuer() == fp.Valuer())
	})

	t.Run("update isolation", func(t *testing.T) {
		scope2 := NewScope("scope2")
		loc2 := location.GetCallLocation(0)

		fp2 := fp.clone()
		fp2.Param(0, ByName("a2"), ByTags(tag2), Optional(false))
		fp2.Param(1, ByName("b2"))
		fp2.Param(2, ByName("c2"), AsCollector(false))
		fp2.Param(3, ByName("d2"))
		fp2.Return(0, Name("rr1"), Tags(tag1), As(TypeOf(1)))
		fp2.Return(1, Name("rr2"))
		fp2.SetScope(scope2)
		fp2.SetLocation(loc2)

		verifyProvider(t, fp.Provider())
	})

	t.Run("update isolation 2", func(t *testing.T) {
		scope2 := NewScope("scope2")
		loc2 := location.GetCallLocation(0)

		fp2 := fp.clone()
		fp3 := fp2.clone()
		fp2.Param(0, ByName("a2"), ByTags(tag2), Optional(false))
		fp2.Param(1, ByName("b2"))
		fp2.Param(2, ByName("c2"), AsCollector(false))
		fp2.Param(3, ByName("d2"))
		fp2.Return(0, Name("rr1"), Tags(tag1), As(TypeOf(1)))
		fp2.Return(1, Name("rr2"))
		fp2.SetScope(scope2)
		fp2.SetLocation(loc2)

		verifyProvider(t, fp3.Provider())
	})
}

func Test_funcProvider_Equal(t *testing.T) {
	type testInterface interface{}

	tag1 := NewSymbol("tag1")
	tag2 := NewSymbol("tag2")
	scope1 := NewScope("scope1")
	loc1 := location.GetCallLocation(0)

	fp := funcProviderOf(funcForFuncProviderTest,
		InScope(scope1),
		Param(0, ByName("a1"), ByTags(tag1), Optional(true)),
		Param(1, ByName("b1")),
		Param(2, ByName("c1"), AsCollector(true)),
		Param(3),
		Return(0, Name("r1"), Tags(tag2), Hide(), Ignore(), As(TypeOf((*testInterface)(nil)))),
		Return(1),
	)
	fp.SetLocation(loc1)

	t.Run("equal", func(t *testing.T) {
		fp2 := fp.clone()
		assert.True(t, fp2.Equal(fp))
		assert.True(t, fp.Equal(fp2))
	})

	t.Run("funcConsumer", func(t *testing.T) {
		fp2 := fp.clone()
		fp2.Param(0, ByName("a2"))
		assert.False(t, fp2.Equal(fp))
		assert.False(t, fp.Equal(fp2))
	})

	t.Run("component", func(t *testing.T) {
		fp2 := fp.clone()
		fp2.Return(0, Name("r2"))
		assert.False(t, fp2.Equal(fp))
		assert.False(t, fp.Equal(fp2))
	})

	t.Run("fakeComponent", func(t *testing.T) {
		fp2 := fp.clone()
		fp2.Param(0, ByName("a2"))
		fp2.Return(2, Name("r"))
		assert.False(t, fp2.Equal(fp))
		assert.False(t, fp.Equal(fp2))
	})
}
