package spam

type SpamChecker interface {
	Validate(text string) (ok bool, word string)
}

// uniform checker
type spamChecker struct {
	checkers []SpamChecker
}

func (ck *spamChecker) Validate(text string) (ok bool, word string) {
	for _, checker := range ck.checkers {
		if ok, _ := checker.Validate(text); !ok {
			return false, ""
		}
	}
	return true, ""
}

func New(checkers ...SpamChecker) SpamChecker  {
	return &spamChecker {
		checkers: checkers,
	}
}