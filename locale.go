package clickcounteriframe

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type LocalePriority struct {
	Locale   string
	Priority float64
}

func (p LocalePriority) String() string {
	return fmt.Sprintf("%s: %f", p.Locale, p.Priority)
}

// ByPriority implements sort.Interface for []LocalePriority based on
// the Age field.
type ByPriority []LocalePriority

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority < a[j].Priority }

// Parse the value in a Accept-Language Header
func GetClientLocales(acceptLanguage string) []string {
	if len(acceptLanguage) == 0 {
		locales := []string{"en"}
		return locales
	}

	languages := strings.Split(acceptLanguage, ",")
	prioLocales := make([]LocalePriority, len(languages))
	for i, localeQPair := range languages {
		localeQ := strings.Split(localeQPair, ";")
		prioLocale := new(LocalePriority)
		if localeQ[0] == localeQPair {
			prioLocale.Locale = localeQ[0]
			prioLocale.Priority = 1.0
		} else {
			qSplit := strings.Split(localeQ[1], "=")
			q, err := strconv.ParseFloat(qSplit[1], 32)
			if err != nil {
				q = 0
			}
			prioLocale.Locale = localeQ[0]
			prioLocale.Priority = q
		}
		prioLocales[i] = *prioLocale
	}

	sort.Sort(sort.Reverse(ByPriority(prioLocales)))

	locales := make([]string, len(languages))
	for i, l := range prioLocales {
		locales[i] = l.Locale
	}

	return locales
}
