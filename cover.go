package hangulize

// cov is a tuple to identify a rule.
type cov struct {
	step Step
	id   int
}

// cover stores covered rule IDs.
type cover map[cov]bool

func (c cover) Mark(step Step, ruleID int) {
	c[cov{step, ruleID}] = true
}

func (c cover) Has(step Step, ruleID int) bool {
	return c[cov{step, ruleID}]
}

// Coverage measures the test coverage of a spec.
func Coverage(h *Hangulizer) float64 {
	cover := make(cover)

	for _, exm := range h.spec.Test {
		word := exm[0]
		_, traces := h.HangulizeTrace(word)

		for _, tr := range traces {
			if tr.Rule != nil {
				cover.Mark(tr.Step, tr.Rule.ID)
			}
		}
	}

	total := len(h.spec.Rewrite) + len(h.spec.Transcribe)

	if total == 0 {
		return 0.0
	}

	covered := 0

	for _, r := range h.spec.Rewrite {
		if cover.Has(Rewrite, r.ID) {
			covered++
		}
	}
	for _, r := range h.spec.Transcribe {
		if cover.Has(Transcribe, r.ID) {
			covered++
		}
	}

	return float64(covered) / float64(total)
}
