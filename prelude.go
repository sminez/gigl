package gigl

// This is a LISP prelude of procedures to be defined at runtime
// Probably more efficient to define these in go but meh...this is more fun!
var prelude = []string{
	// Simple procedures that are just easier to define in LISP...!
	"(defn list l l)",
	"(defn abs (n) ((if (> n 0) + -) 0 n))",
	// Drop/take the first elements of a list
	"(defn drop (n lst) (if (= n 0) lst (drop-n (- n 1) (cdr lst))))",
	"(defn take (n lst) (if (= n 0) '() (cons (car lst) (take (- n 1) (cdr lst)))))",
	"(defn dropwhile (pred lst) (cond ((null? lst) '()) ((pred (car lst)) (dropwhile pred (cdr lst))) (:else lst)))",
	"(defn takewhile (pred lst) (cond ((null? lst) '()) ((pred (car lst)) (cons (car lst) (takewhile pred (cdr lst)))) (:else '())))",
	// Selectors for specific elements of a list
	"(defn caar (lst) (car (car lst)))",
	"(defn cadr (lst) (car (cdr lst)))",
	"(defn cdar (lst) (cdr (car lst)))",
	"(defn cddr (lst) (cdr (cdr lst)))",
	"(defn caddr (lst) (car (cdr (cdr lst))))",
	// TBH, these are a lot less archaic and easier to remember than c...r
	"(defn last (lst) (cond ((= 0 (len lst)) '()) ((= (len lst) 1) (car lst)) (:else (last (cdr lst)))))",
	"(defn nth (n lst) (if (= 0 (len lst)) '() (if (= n 0) (car lst) (nth (- n 1) (cdr lst)))))",
	// Higher order functions
	"(defn compose (f g) (λ (x) (f (g x))))",
	"(defn repeat (f) (compose f f))",
	"(defn map (f lst) (foldr (λ (x y) (cons (f x) y)) (list) lst))",
	// The f in map-append must return a list. The final result is a list of
	// all of the results of (f elem) appended together
	// (map-append (λ (n) (list n (* 10 n))) (range 5)) --> (0 0 1 10 2 20 3 30 4 40)
	"(defn map-append (f lst) (if (= 0 (len lst)) '() (append (f (car lst)) (map-append f (cdr lst)))))",
	"(defn amap (f lst) (if (= 0 (len lst)) '() (append (f (car lst)) (amap f (cdr lst)))))",
	// map-tail will build a list of lists: the result of calling f on first the entire
	// list, then the tail, tail of the tail...etc until we reach '()
	// (map-tail (λ (lst) (apply * lst)) (range 5)) --> (120 120 60 20 5)
	"(defn map-tail (f lst) (if (= 0 (len lst)) '() (cons (f lst) (map-tail f (cdr lst)))))",
	"(defn tmap (f lst) (if (= 0 (len lst)) '() (cons (f lst) (tmap f (cdr lst)))))",
	"(defn filter (f lst) (foldr (λ (x y) (if (f x) (cons x y) y)) (list) lst))",
	"(defn flip (f) (λ (a b) (f b a)))",
	"(defn curry (f a) (λ (b) (f a b)))",
	"(defn combine (f) (λ (x y) (if (null? x) '() (f (list (car x) (car y)) ((combine f) (cdr x) (cdr y))))))",
	"(define zip (combine cons))",
	// Simple short circuiting boolean logic
	"(defn not (x) (if x #f #t))",
	"(defn or (lst) (if (= 0 (len lst)) #f (if (car lst) #t (or (cdr lst)))))",
	"(defn and (lst) (if (null? (cdr lst)) (car lst) (if (car lst) (and (cdr lst)) #f)))",
	// Boolean checks
	"(defn zero? (n) (curry = 0))",
	"(defn positive? (n) (if (float? n) (> n 0) #f))",
	"(defn pos? (n) (if (float? n) (> n 0) #f))",
	"(defn negative? (n) (if (float? n) (< n 0) #f))",
	"(defn neg? (n) (if (float? n) (< n 0) #f))",
	"(defn even? (n) (if (int? n) (= (% n 2) 0) #f))",
	"(defn odd? (n) (if (int? n) (= (% n 2) 1) #f))",
	// these are useful for filters as otherwise the inequality is reversed and it
	// get confusing --> (filter (>than 4) lst) == (filter (curry < 4) lst)
	"(defn <than (n) (curry > n))",
	"(defn <=to (n) (curry >= n))",
	"(defn >than (n) (curry < n))",
	"(defn >=to (n) (curry <= n))",
	// Scans and folds: fold and scan are left based and use the first element
	// of their list argument as the accumulator.
	// NOTE :: scans require a list based accumulator!
	"(defn foldl (f acc lst) (if (= 0 (len lst)) acc (foldl f (f acc (car lst)) (cdr lst))))",
	"(defn foldr (f acc lst) (if (= 0 (len lst)) acc (f (car lst) (foldr f acc (cdr lst)))))",
	"(defn fold (f lst) (if (= 0 (len lst)) lst (foldl f (car lst) (cdr lst))))",
	"(defn reduce (f lst) (if (= 0 (len lst)) lst (foldl f (car lst) (cdr lst))))",
	"(defn scanl (f acc lst) (if (= 0 (len lst)) acc (scanl f (append acc (list (f (car lst) (last acc)))) (cdr lst))))",
	"(define scanr (λ (f acc lst) (scanl f acc (reverse lst))))",
	"(define scan (λ (f lst) (if (= 0 (len lst)) lst (scanl f (list (car lst)) (cdr lst)))))",
	"(defn reverse (lst) (foldl (flip cons) '() lst))",
	// More fun with maps and higher order functions
	"(defn concat-map (f lst) (fold append (map f lst)))",
	"(defn cmap (f lst) (fold append (map f lst)))",
	"(defn flatten (lst) (if (list? lst) (cmap flatten lst) (list lst)))",
	// Built-in macros
	// NOTE :: as I'm still working on the macro syntax, these may change...
	"(defmacro unless (arg body) `(if (not ~arg) ~body))",
}
