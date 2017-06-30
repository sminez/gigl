package gigl

// This is a LISP prelude of procedures to be defined at runtime
// Probably more
var prelude = []string{
	"(defn list l l)",
	"(defn compose (f g) (λ (x) (f (g x))))",
	"(defn repeat (f) (compose f f))",
	"(defn abs (n) ((if (> n 0) + -) 0 n))",
	"(defn combine (f) (λ (x y) (if (null? x) (quote ()) (f (list (car x) (car y)) ((combine f) (cdr x) (cdr y))))))",
	"(define zip (combine cons))",
	// "(defn zip-with (f) (λ (iters) (map (λ (x) (f x)) (zip iters))))",
}
