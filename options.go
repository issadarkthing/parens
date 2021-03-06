package parens

// Option can be used with New() to customize initialization of Evaluator
// Instance.
type Option func(env *Env)

// WithGlobals sets the global variables during initialisation. If factory
// is nil, a mutex based concurrent map will be used.
func WithGlobals(globals map[string]Any, factory func() ConcurrentMap) Option {
	return func(env *Env) {
		if factory == nil {
			factory = newMutexMap
		}
		if env.globals == nil {
			env.globals = factory()
		}
		for k, v := range globals {
			env.globals.Store(k, v)
		}
	}
}

// WithMaxDepth sets the max depth allowed for stack.  Panics if depth == 0.
func WithMaxDepth(depth uint) Option {
	if depth == 0 {
		panic("maxdepth must be nonzero.")
	}
	return func(env *Env) {
		env.maxDepth = int(depth)
	}
}

// WithExpander sets the macro Expander to be used by the p. If nil, a builtin
// Expander will be used.
func WithExpander(expander Expander) Option {
	return func(env *Env) {
		if expander == nil {
			expander = &builtinExpander{}
		}
		env.expander = expander
	}
}

// WithAnalyzer sets the Analyzer to be used by the p. If replaceBuiltin is set,
// given analyzer will be used for all forms. Otherwise, it will be used only for
// forms not handled by the builtinAnalyzer.
func WithAnalyzer(analyzer Analyzer) Option {
	return func(env *Env) {
		if analyzer == nil {
			analyzer = &BuiltinAnalyzer{
				SpecialForms: map[string]ParseSpecial{
					"go":    parseGoExpr,
					"def":   parseDefExpr,
					"quote": parseQuoteExpr,
				},
			}
		}
		env.analyzer = analyzer
	}
}

func withDefaults(opts []Option) []Option {
	return append([]Option{
		WithAnalyzer(nil),
		WithExpander(nil),
		WithMaxDepth(10000),
	}, opts...)
}
