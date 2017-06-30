package gigl

// This is a LISP prelude of procedures to be defined at runtime
// Probably more
var prelude = []string{
	"(defn list l l)",
	"(defn compose (f g) (lambda (x) (f (g x))))",
	"(defn repeat (f) (compose f f))",
	"(defn abs (n) ((if (> n 0) + -) 0 n))",
	"(defn combine (f) (lambda (x y) (if (null? x) (quote ()) (f (list (car x) (car y)) ((combine f) (cdr x) (cdr y))))))",
	"(define zip (combine cons))",
}
