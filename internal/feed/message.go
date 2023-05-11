package feed

import (
	"math/rand"
	"strings"
)

// TODO: Remplace this by CHATGPT one day
var templates []string = []string{
	`
Hey! 👋

Check out this article by _$AUTHOR$_ on _$PLATFORM$_ about *$TITLE$*.

If you're interested in $TAGS$ this is a must-read!

🔗 $LINK$
`,
	`
_$PLATFORM$_ just published an intriguing piece by _$AUTHOR$_.

The article highlights *$TITLE$*.

If you're into $TAGS$, give it a read!

🔗 $LINK$
`,
	`
I came across this fascinating 🤔 article by _$AUTHOR$_ on _$PLATFORM$_ about *$TITLE$*.

If you're curious about $TAGS$, take a look!

🔗 $LINK$
`,
	`
❗ Attention❗

Check out this article by _$AUTHOR$_ on _$PLATFORM$_ about *$TITLE$*.

Don't miss out! $TAGS$

🔗 $LINK$
`,
	`
Want to learn more 👨‍🏫 about *$TITLE$*?

_$AUTHOR$'s_ _$PLATFORM$_ article dives deep into $TAGS$.

Check it out!

🔗 $LINK$
`,
	`
Hey everyone 👋,

I found this interesting _$PLATFORM$_ article by _$AUTHOR$_ about *$TITLE$*.

$TAGS$

🔗 $LINK$
`, `
👀 If you're interested in $TAGS$,

you won't want to miss this _$PLATFORM$_ article by _$AUTHOR$_ about *$TITLE$*.

Check it out!

🔗 $LINK$
`, `
👀 💰 💡💡💡

👉🏼 📝👩🏻‍💼 _$AUTHOR$_ on _$PLATFORM$_ about *$TITLE$* for $TAGS$

🔗 $LINK$
`, `
*$TITLE$*,

_$AUTHOR's$_ _$PLATFORM$_ article explores this topic $TAGS$.

If you're curious, give it a read!

🔗 $LINK$
`, `
_$PLATFORM$_ just published an insightful piece by _$AUTHOR$_ about *$TITLE$*.

Don't miss out on this one!

$TAGS$

🔗 $LINK$
`,
}

// Will inject the parameters into the template choose randomly.
func BuildMessage(title string, platform string, author string, tags []string, link string) string {
	text := templates[rand.Intn(len(templates))]
	replacer := strings.NewReplacer("$TITLE$", title, "$AUTHOR$", author, "$PLATFORM$", platform, "$TAGS$", "#"+strings.Join(tags, ", #"), "$LINK$", link)

	return replacer.Replace(text)
}
