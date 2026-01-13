// Package test provides benchmarks for smog v0.4.0 features.
package test

import (
	"testing"

	"github.com/kristofer/smog/pkg/compiler"
	"github.com/kristofer/smog/pkg/parser"
)

// BenchmarkSuperMessageSend benchmarks super message send parsing and compilation
func BenchmarkSuperMessageSend(b *testing.B) {
	b.Run("ParseSuperUnary", func(b *testing.B) {
		input := "super initialize"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseSuperKeyword", func(b *testing.B) {
		input := "super at: 5 put: 10"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("CompileSuperMessage", func(b *testing.B) {
		input := "super initialize"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})
}

// BenchmarkDictionaryLiterals benchmarks dictionary literal operations
func BenchmarkDictionaryLiterals(b *testing.B) {
	b.Run("ParseSmallDictionary", func(b *testing.B) {
		input := "#{'a' -> 1. 'b' -> 2}"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseMediumDictionary", func(b *testing.B) {
		input := "#{'a' -> 1. 'b' -> 2. 'c' -> 3. 'd' -> 4. 'e' -> 5}"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseLargeDictionary", func(b *testing.B) {
		input := "#{'a' -> 1. 'b' -> 2. 'c' -> 3. 'd' -> 4. 'e' -> 5. 'f' -> 6. 'g' -> 7. 'h' -> 8. 'i' -> 9. 'j' -> 10}"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("CompileSmallDictionary", func(b *testing.B) {
		input := "#{'x' -> 10. 'y' -> 20}"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})

	b.Run("CompileLargeDictionary", func(b *testing.B) {
		input := "#{'a' -> 1. 'b' -> 2. 'c' -> 3. 'd' -> 4. 'e' -> 5. 'f' -> 6. 'g' -> 7. 'h' -> 8}"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})
}

// BenchmarkCascadingMessages benchmarks cascading message operations
func BenchmarkCascadingMessages(b *testing.B) {
	b.Run("ParseTwoCascades", func(b *testing.B) {
		input := "point x: 10; y: 20"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseFiveCascades", func(b *testing.B) {
		input := "obj m1; m2; m3; m4; m5"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseTenCascades", func(b *testing.B) {
		input := "obj m1; m2; m3; m4; m5; m6; m7; m8; m9; m10"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("CompileTwoCascades", func(b *testing.B) {
		input := "point x: 10; y: 20"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})

	b.Run("CompileFiveCascades", func(b *testing.B) {
		input := "obj m1; m2; m3; m4; m5"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})
}

// BenchmarkCascadeVsNormalMessages compares cascade performance to normal message sends
func BenchmarkCascadeVsNormalMessages(b *testing.B) {
	b.Run("NormalMessages", func(b *testing.B) {
		// Separate message sends to simulate setting x and y
		input := "point x: 10. point y: 20"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})

	b.Run("CascadedMessages", func(b *testing.B) {
		// Same operations but cascaded
		input := "point x: 10; y: 20"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})
}

// BenchmarkSelfKeyword benchmarks self keyword usage
func BenchmarkSelfKeyword(b *testing.B) {
	b.Run("ParseSelf", func(b *testing.B) {
		input := "self"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseSelfMessageSend", func(b *testing.B) {
		input := "self initialize"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})

	b.Run("ParseSelfWithCascade", func(b *testing.B) {
		input := "self x: 10; y: 20; z: 30"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			_, _ = p.Parse()
		}
	})
}

// BenchmarkComplexExpressions benchmarks complex expressions using v0.4.0 features
func BenchmarkComplexExpressions(b *testing.B) {
	b.Run("DictionaryInCascade", func(b *testing.B) {
		input := "obj data: #{'x' -> 10. 'y' -> 20}; process"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})

	b.Run("NestedCascades", func(b *testing.B) {
		input := "outer inner; process: (point x: 1; y: 2)"
		for i := 0; i < b.N; i++ {
			p := parser.New(input)
			program, _ := p.Parse()
			c := compiler.New()
			_, _ = c.Compile(program)
		}
	})
}
