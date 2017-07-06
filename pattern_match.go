package gigl

import "fmt"

/*
This is the pattern matching language for syntax-rules based macros.

(defmacro <name>
  (syntax-rules (<syms>)
    ((pattern1) (template1)
	 (pattern2) (template2)
	 .
	 .
	 .)))


If everything works, this should parse! (Taken from the R7RS spec)

(defmacro let
  (syntax-rules ()
    ((let ((name val) ...) body1 body2 ...)
     ((lambda (name ...) body1 body2 ...)
       val ...))
    ((let tag ((name val) ...) body1 body2 ...)
     ((letrec ((tag (lambda (name ...)
       body1 body2 ...)))
       tag)
     val ...))))
*/

type matcher interface {
	Matches(lispVal) bool
	isSubpattern() bool
	makeRepeating()
	isRepeating() bool
}

type pVar struct {
	symbol    SYMBOL
	value     lispVal
	rawMatch  bool
	repeating bool
}

type matchPattern struct {
	value       lispVal
	repeating   bool
	hasEllipsis bool
	pVars       *LispList
	pattern     *LispList
	varMap      *map[SYMBOL][]lispVal
}

// Check to see if a pVar matches a target
func (p *pVar) Matches(other lispVal) bool {
	// Underscores are wildcards that match anything
	if p.value == "_" {
		return true
	}
	// Raw matches only ever match their value
	if p.rawMatch {
		return p.value == other
	}
	// Otherwise we match if we are not yet bound...
	if p.value == nil {
		p.value = other
		return true
	}
	// ...or if our value is equal to the target
	return p.value == other
}

// Make sure that repeated vars are the same
func (p *pVar) propagateMatch(attempt *map[SYMBOL][]lispVal) bool {
	if p.rawMatch || p.value == "_" {
		// Raw matches and underscores don't propagate
		return true
	}

	existingVal, present := (*attempt)[p.symbol]

	if present && !p.isRepeating() {
		return p.value == existingVal[0]
	}

	(*attempt)[p.symbol] = append((*attempt)[p.symbol], p.value)
	return true
}

func (p *pVar) makeRepeating() {
	p.repeating = true
}

func (m *matchPattern) makeRepeating() {
	m.repeating = true
}

func (p *pVar) isRepeating() bool {
	return p.repeating
}

func (m *matchPattern) isRepeating() bool {
	return m.repeating
}

func (p *pVar) isSubpattern() bool {
	return false
}

func (m *matchPattern) isSubpattern() bool {
	return true
}

// Initialise a new template
func NewMatchPattern(ptrn *LispList) (*matchPattern, error) {
	varMap := make(map[SYMBOL][]lispVal)
	p := matchPattern{
		repeating:   false,
		hasEllipsis: false,
		varMap:      &varMap,
		pattern:     ptrn,
	}

	// Collecting things up in a slice is conceptually easier to think about
	// when compared the repeated appends of lists or cons -> reverse.
	pvarSlice := make([]lispVal, 0)
	pvar, remaining := ptrn.popHead()

	for {
		// If we've reached the end of the pattern s-expression then bind it
		// onto the matchPattern and return
		if remaining.Len() == 0 && (pvar == nil || isEmptyList(pvar)) {
			p.pVars = List(pvarSlice...)
			return &p, nil
		}

		switch pvar.(type) {
		case *LispList:
			// Add a nested pattern
			subPattern, err := NewMatchPattern(pvar.(*LispList))
			if err != nil {
				return &p, err
			}
			pvarSlice = append(pvarSlice, subPattern)

		case SYMBOL:
			switch pvar {
			case SYMBOL("_"):
				underscore := &pVar{
					symbol:   "_",
					value:    nil,
					rawMatch: false,
				}
				pvarSlice = append(pvarSlice, underscore)

			case SYMBOL("..."):
				if p.hasEllipsis {
					return &p, fmt.Errorf("Can only have one ... per template")
				}
				p.hasEllipsis = true
				pvarSlice[len(pvarSlice)-1].(matcher).makeRepeating()

			default:
				newPvar := &pVar{
					symbol:   pvar.(SYMBOL),
					value:    nil,
					rawMatch: false,
				}
				pvarSlice = append(pvarSlice, newPvar)
			}
		}
		pvar, remaining = remaining.popHead()
	}
}

// Check that a subcomponent matches correctly
func (m *matchPattern) componentMatches(mvar matcher, t lispVal) bool {
	// NOTE :: this will bind the match if true
	if mvar.Matches(t) {
		if mvar.isSubpattern() {
			// update the master map with the results of the subPattern
			subMap := mvar.(*matchPattern).varMap
			for k, v := range *subMap {
				(*m.varMap)[k] = v
			}
			return true
		}
		return mvar.(*pVar).propagateMatch(m.varMap)
	}
	return false
}

// Check to see if a pVar matches a target
func (m *matchPattern) Matches(target lispVal) bool {
	_, ok := target.(*LispList)
	// A matchPattern can only match against an s-expression that is
	// less than or equal to it in length
	if !ok {
		return false
	}

	other := target.(*LispList)

	if m.pVars.Len() > other.Len() {
		return false
	}

	pRemaining, tRemaining := m.pVars, other
	var p matcher
	var t lispVal

MATCHLOOP:
	for i := 0; i < other.Len(); i++ {
		p, pRemaining, t, tRemaining, ok = getNext(pRemaining, tRemaining)

		// If we run out of pvars then the match failed
		if !ok {
			return false
		}

		// This does the binding match here
		if m.componentMatches(p, t) {
			if p.isRepeating() {
				var matchedVals []lispVal
				var subMap map[SYMBOL][]lispVal

				switch p.(type) {
				case *pVar:
					matchedVals = []lispVal{p.(*pVar).value}

				case *matchPattern:
					subMap = make(map[SYMBOL][]lispVal)
					for k, v := range *p.(*matchPattern).varMap {
						if val, alreadyMatched := (*m.varMap)[k]; alreadyMatched {
							if val[0] != v[0] {
								return false
							}
						}
						subMap[k] = v
					}
				}

				// Now consume the rest of the input
			REPLOOP:
				for {
					t, tRemaining = tRemaining.popHead()

					if tRemaining.Len() == 0 && (t == nil || isEmptyList(t)) {
						// we've consumed all of the input
						break REPLOOP
					}

					// reset the pVars and try the match again
					switch p.(type) {
					case *pVar:
						p.(*pVar).value = nil

					case *matchPattern:
						vars := p.(*matchPattern).pVars
						var v lispVal
						for {
							v, vars = vars.popHead()
							if vars.Len() == 0 && (v == nil || isEmptyList(v)) {
								break
							}
							v.(*pVar).value = nil
						}
					}

					// Try the match again
					if !m.componentMatches(p, t) {
						return false
					}

					// Update the map
					switch p.(type) {
					case *pVar:
						matchedVals = append(matchedVals, p.(*pVar).value)
					case *matchPattern:
						for _, v := range p.(*matchPattern).pVars.toSlice() {
							subMap[v.(*pVar).symbol] = append(subMap[v.(*pVar).symbol], v.(*pVar).value)
						}
					}
				}
				// Update the master map
				switch p.(type) {
				case *pVar:
					(*m.varMap)[p.(*pVar).symbol] = []lispVal{List(matchedVals...)}
				case *matchPattern:
					for k, v := range subMap {
						(*m.varMap)[k] = []lispVal{List(v...)}
					}
				}
				break MATCHLOOP
			}
		}
	}
	// Make sure everything actually matched
	for _, v := range *(m.varMap) {
		if v == nil {
			return false
		}
	}
	m.value = m.pVars
	return true
}

func (m *matchPattern) getMatches() map[SYMBOL]lispVal {
	matches := make(map[SYMBOL]lispVal)
	for k, v := range *(m.varMap) {
		matches[k] = v[0].(lispVal)
	}
	return matches
}

func (m *matchPattern) PrintMatch() {
	fmt.Println("match map:")
	for k, v := range m.getMatches() {
		fmt.Printf("\t%v --> %v\n", k, v)
	}
}

func getNext(pvars, targets *LispList) (matcher, *LispList, lispVal, *LispList, bool) {
	t, remainingTargets := targets.popHead()
	p, remainingPvars := pvars.popHead()
	okP := !(remainingPvars.Len() == 0) || !(p == nil || isEmptyList(p))
	okT := !(remainingTargets.Len() == 0) || !(t == nil || isEmptyList(t))
	if okP && okT {
		switch p.(type) {
		case *pVar:
			return p.(*pVar), remainingPvars, t, remainingTargets, true
		case *matchPattern:
			return p.(*matchPattern), remainingPvars, t, remainingTargets, true
		}
	}
	return &pVar{}, remainingPvars, t, remainingTargets, true
}
