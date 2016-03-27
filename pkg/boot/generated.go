// Released under an MIT-style license. See LICENSE.

// Generated by generate.oh

package boot

var Script string = `
define caar: method (l) =: car: car l
define cadr: method (l) =: car: cdr l
define cdar: method (l) =: cdr: car l
define cddr: method (l) =: cdr: cdr l
define caaar: method (l) =: car: caar l
define caadr: method (l) =: car: cadr l
define cadar: method (l) =: car: cdar l
define caddr: method (l) =: car: cddr l
define cdaar: method (l) =: cdr: caar l
define cdadr: method (l) =: cdr: cadr l
define cddar: method (l) =: cdr: cdar l
define cdddr: method (l) =: cdr: cddr l
define caaaar: method (l) =: caar: caar l
define caaadr: method (l) =: caar: cadr l
define caadar: method (l) =: caar: cdar l
define caaddr: method (l) =: caar: cddr l
define cadaar: method (l) =: cadr: caar l
define cadadr: method (l) =: cadr: cadr l
define caddar: method (l) =: cadr: cdar l
define cadddr: method (l) =: cadr: cddr l
define cdaaar: method (l) =: cdar: caar l
define cdaadr: method (l) =: cdar: cadr l
define cdadar: method (l) =: cdar: cdar l
define cdaddr: method (l) =: cdar: cddr l
define cddaar: method (l) =: cddr: caar l
define cddadr: method (l) =: cddr: cadr l
define cdddar: method (l) =: cddr: cdar l
define cddddr: method (l) =: cddr: cddr l
define _connect_: syntax (conduit name) = {
	set conduit: eval conduit
	syntax (left right) e = {
		define p: conduit
		spawn {
			e::eval: quasiquote: block {
				public (unquote name) = (unquote p)
				eval (unquote left)
			}
			p::writer-close
		}
		block {
			e::eval: quasiquote: block {
				public _stdin_ = (unquote p)
				eval (unquote right)
			}
			p::reader-close
		}
	}
}
define _redirect_: syntax (name mode closer) = {
	syntax (c cmd) e = {
		define c: e::eval c
		define f = ()
		if (not: or (is-channel c) (is-pipe c)) {
			set f: open mode c
			set c = f
		}
		e::eval: quasiquote: block {
			public (unquote name) (unquote c)
			eval (unquote cmd)
		}
		if (not: is-null f): eval: quasiquote: f::(unquote closer)
	}
}
define ...: method (: args) = {
	cd _origin_
	define path: car args
	if (eq 2: length args) {
		cd: car args
		set path: cadr args
	}
	while true {
		define abs: symbol: "/"::join $CWD path
		if (exists abs): return abs
		if (eq $CWD /): return path
		cd ..
	}
}
define and: syntax (: lst) e = {
	define r = false
	while (not: is-null: car lst) {
		set r: e::eval: car lst
		if (not r): return r
		set lst: cdr lst
	}
	return r
}
define _append-stderr_: _redirect_ _stderr_ "a" writer-close
define _append-stdout_: _redirect_ _stdout_ "a" writer-close
define apply: method (f: args) =: f @args
define _backtick_: syntax (cmd) e = {
	define p: pipe
	spawn {
		e::eval: quasiquote: block {
			public _stdout_ = (unquote p)
			eval (unquote cmd)
		}
		p::writer-close
	}
	define r: cons () ()
	define c = r
	while (define l: p::readline) {
		set-cdr c: cons l ()
		set c: cdr c
	}
	p::reader-close
	return: cdr r
}
define catch: syntax (name: clause) e = {
	define args: list name (quote throw)
	define body: list (quote throw) name
	if (is-null clause) {
		set body: list body
	} else {
		set body: append clause body
	}
	define handler: e::eval {
		list (quote method) args (quote =) @body
	}
	define _return: e::eval (quote return)
	define _throw = throw
	e::public throw: method (condition) = {
		_return: handler condition _throw
	}
}
define _channel-stderr_: _connect_ channel _stderr_
define _channel-stdout_: _connect_ channel _stdout_
define echo: builtin (: args) = {
	if (is-null args) {
		_stdout_::write: symbol ""
	} else {
		_stdout_::write @(for args symbol)
	}
}
define error: builtin (: args) =: _stderr_::write @args
define for: method (l m) = {
	define r: cons () ()
	define c = r
	while (not: is-null l) {
		set-cdr c: cons (m: car l) ()
		set c: cdr c
		set l: cdr l
	}
	return: cdr r
}
define glob: builtin (: args) =: return args
define import: syntax (name) e = {
	set name: e::eval name
	define m: module name
	if (or (is-null m) (is-object m)) {
		return m
	}

	e::eval: quasiquote: _root_::define (unquote m): object {
		source (unquote name)
	}
}
define is-list: method (l) = {
	if (is-null l): return false
	if (not: is-cons l): return false
	if (is-null: cdr l): return true
	is-list: cdr l
}
define is-text: method (t) =: or (is-string t) (is-symbol t)
define list-ref: method (k x) =: car: list-tail k x
define list-tail: method (k x) = {
	if k {
		list-tail (sub k 1): cdr x
	} else {
		return x
	}
}
define object: syntax (: body) e = {
	e::eval: cons (quote block): append body (quote: context)
}
define or: syntax (: lst) e = {
	define r = false
	while (not: is-null: car lst) {
		set r: e::eval: car lst
		if r: return r
		set lst: cdr lst 
	}
	return r
}
define _pipe-stderr_: _connect_ pipe _stderr_
define _pipe-stdout_: _connect_ pipe _stdout_
define printf: method (f: args) =: echo: f::sprintf @args
define _process-substitution_: syntax (:args) e = {
	define fifos = ()
	define procs = ()
	define cmd: for args: method (arg) = {
		if (not: is-cons arg): return arg
		if (eq (quote substitute-stdin) (car arg)) {
			define fifo: temp-fifo
			define proc: spawn {
				e::eval: quasiquote {
					_redirect-stdin_ {
						unquote fifo
						unquote: cdr arg
					}
				}
			}
			set fifos: cons fifo fifos
			set procs: cons proc procs
			return fifo
		}
		if (eq (quote substitute-stdout) (car arg)) {
			define fifo: temp-fifo
			define proc: spawn {
				e::eval: quasiquote {
					_redirect-stdout_ {
						unquote fifo
						unquote: cdr arg
					}
				}
			}
			set fifos: cons fifo fifos
			set procs: cons proc procs
			return fifo
		}
		return arg
	}
	e::eval cmd
	wait @procs
	rm @fifos
}
define quasiquote: syntax (cell) e = {
	if (not: is-cons cell): return cell
	if (is-null cell): return cell
	if (eq (quote unquote): car cell): return: e::eval: cadr cell
	cons {
		e::eval: list (quote quasiquote): car cell
		e::eval: list (quote quasiquote): cdr cell
	}
}
define quote: syntax (cell) =: return cell
define read: builtin () =: _stdin_::read
define readline: builtin () =: _stdin_::readline
define _redirect-stderr_: _redirect_ _stderr_ "w" writer-close
define _redirect-stdin_: _redirect_ _stdin_ "r" reader-close
define _redirect-stdout_: _redirect_ _stdout_ "w" writer-close
define source: syntax (name) e = {
	define basename: e::eval name
	define paths = ()
	define name = basename

	if (has $OHPATH): set paths: (string $OHPATH)::split ":"
	while (and (not: is-null paths) (not: exists name)) {
		set name: "/"::join (car paths) basename
		set paths: cdr paths
	}

	if (not: exists name): set name = basename

	define f: open r- name

	define r: cons () ()
	define c = r
	while (define l: f::read) {
		set-cdr c: cons (cons (get-line-number) l) ()
		set c: cdr c
	}
	set c: cdr r
	f::close

	define rval: status 0
	define eval-list: syntax (first rest) o = {
		set first: o::eval first
		set rest: o::eval rest
		if (is-null first): return rval
		set-line-number: car first
		set rval: e::eval: cdr first
		eval-list (car rest) (cdr rest)
	}
	eval-list (car c) (cdr c)
	return rval
}
define write: method (: args) =: _stdout_::write @args
_sys_::public exception: method (type message status file line) = {
	object {
		public type = type
		public status = status
		public message = message
		public line = line
		public file = file
	}
}
_sys_::public get-prompt: method self (suffix) = {
	catch unused {
		return suffix
	}
	self::prompt suffix
}
_sys_::public prompt: method (suffix) = {
	define dirs: (string $CWD)::split "/"
	define last: sub (length dirs) 1
	return (list-ref last dirs) ^ suffix
}
_sys_::public throw: method (c) = {
	error: ": "::join c::file c::line c::type c::message
	fatal c::status
}

exists ("/"::join $HOME .ohrc) && source ("/"::join $HOME .ohrc)

`

//go:generate ./generate.oh
//go:generate go fmt generated.go
