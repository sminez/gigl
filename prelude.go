package gigl

// This is a LISP prelude of procedures to be defined at runtime
// Probably more
var prelude = []string{
	"(defn list l l)",
	"(defn NOT (x) (if x #f #t))",
	"(defn OR (lst) (if (null? lst) #f (if (car lst) #t (OR (cdr lst)))))",
	"(defn AND (lst) (if (null? (cdr lst)) (car lst) (if (car lst) (AND (cdr lst)) #f)))",
	"(defn compose (f g) (λ (x) (f (g x))))",
	"(defn repeat (f) (compose f f))",
	"(defn abs (n) ((if (> n 0) + -) 0 n))",
	"(defn combine (f) (λ (x y) (if (null? x) (quote ()) (f (list (car x) (car y)) ((combine f) (cdr x) (cdr y))))))",
	"(define zip (combine cons))",
	"(defn caar (lst) (car (car lst)))",
	"(defn cadr (lst) (car (cdr lst)))",
	"(defn cdar (lst) (cdr (car lst)))",
	"(defn cddr (lst) (cdr (cdr lst)))",
	"(defn foldl (f acc lst) (if (null? lst) acc (foldl f (f acc (car lst)) (cdr lst))))",
	"(defn foldr (f acc lst) (if (null? lst) acc (f (car lst) (foldr f acc (cdr lst)))))",
	"(defn map (f lst) (foldr (λ (x y) (cons (f x) y)) (list) lst))",
	"(defn filter (f lst) (foldr (λ (x y) (if (f x) (cons x y) y)) (list) lst))",
	"(defn even? (n) (= (% n 2) 0))",
	"(defn odd? (n) (= (% n 2) 1))",
	"(defn flip (f) (λ (a b) (f b a)))",
	"(defn curry (f a) (λ (b) (f a b)))",
}
