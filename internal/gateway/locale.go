package gateway

import (
	"golang.org/x/text/language"
)

// supported lists the locales the owner apps ship with. The first entry is
// the fallback.
var supported = []language.Tag{
	language.AmericanEnglish,
	language.CanadianFrench,
	language.German,
	language.Arabic,
}

var matcher = language.NewMatcher(supported)

var messages = map[language.Tag]map[LinkState]string{
	language.AmericanEnglish: {
		LinkOnline:   "Your vehicle is connected.",
		LinkAsleep:   "Your vehicle is asleep and will wake when needed.",
		LinkOffline:  "Your vehicle has been out of coverage.",
		LinkUpdating: "A software update is being installed.",
	},
	language.CanadianFrench: {
		LinkOnline:   "Votre véhicule est connecté.",
		LinkAsleep:   "Votre véhicule est en veille.",
		LinkOffline:  "Votre véhicule est hors de la zone de couverture.",
		LinkUpdating: "Une mise à jour logicielle est en cours.",
	},
	language.German: {
		LinkOnline:   "Ihr Fahrzeug ist verbunden.",
		LinkAsleep:   "Ihr Fahrzeug befindet sich im Ruhezustand.",
		LinkOffline:  "Ihr Fahrzeug hat keinen Empfang.",
		LinkUpdating: "Ein Software-Update wird installiert.",
	},
	language.Arabic: {
		LinkOnline:   "سيارتك متصلة.",
		LinkAsleep:   "سيارتك في وضع السكون.",
		LinkOffline:  "سيارتك خارج نطاق التغطية.",
		LinkUpdating: "جارٍ تثبيت تحديث البرنامج.",
	},
}

// negotiateLocale picks the best supported locale for the raw
// Accept-Language header sent by the client. The header comes straight from
// the request and is not length limited before it is parsed.
func negotiateLocale(acceptLanguage string) language.Tag {
	if acceptLanguage == "" {
		return supported[0]
	}
	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil || len(tags) == 0 {
		return supported[0]
	}
	_, idx, _ := matcher.Match(tags...)
	return supported[idx]
}

func localizedMessage(tag language.Tag, state LinkState) string {
	if byState, ok := messages[tag]; ok {
		if msg, ok := byState[state]; ok {
			return msg
		}
	}
	return messages[supported[0]][state]
}
