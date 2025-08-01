Much thanks to: https://deadlock-api.com/. You guys have created quite possibly the easiest API I've ever worked with in my life.

some basic CRUDslop

#TL;DR
-type in your steam ID
-API looks up to see if you're in a match
-API looks at what team you're on
-API parses heroes on the enemy team
-For every hero on the enemy team, you get a formatted reccomendation for what items to buy
-If the same item has gotten recommended multiple times, you'll see it highlighted.

Do not use this as the primary source of truth of what to buy. Item builds are HIGHLY complex, contextual and requires a lot of intuition.

This tool is for players who are struggling to find/buy counters for particular heroes. Use this tool for suggestions - not just for you, but for your team.

#Why make this?
I have played too many damn games where no one on my team goes counter items. I've also noticed people on the discord/reddit complain about characters without doing the bare minimum to counter them.


#FUTURE WORK:
-Parse through the character your playing and match items with opponents.
EX: Against high fire rate heroes, a character like Mcginnis loves Supressor, but a frontliner character like Shiv wants Juggernaut.
-Classify items into early/mid/late
EX: Against Grey Talon, in lane spirit shielding is fantastic against his all-ins. Late game, witchmail/spellbreaker would offer a lot more.
-Classify laning items vs midgame items
EX: Itemizing against sinclair is a lot different in lane vs in game. You shouldn't get anti-poke reccomendations if your lane doesn't have poke characters.
-In the middle of the game, identifying "top performers" and highlighting items to build against that hero.
EX: The current system will give you advice on how to build against a 0/10 grey talon that's AFK. This is because it can't parse current game state and give updated recommendations.
-Adding in reasoning for every item that's recommended.
EX: Counter items aren't always intuitive. A lot of them are conditional/require some setup before usage.  For counter spell, you can say something like "Can be used against Shiv's ultimate" or for slowing hex I want to add "Use before Pocket throws out coat/suitcase."
-Identifying build archetypes/item+hero combinations
EX: Sinclair has two build archetypes. He can go for spirit poke items (expansion on 1, extra charge/compress cooldown on 1) or for more utility/control (expansion on 3, duration extender)
-Identifying counter-counter items
EX: Knockdown is great against Seven. However, once he gets unstoppable - it's not as good. Knockdown should be recommended at first, then forgotten.

A lot of this requires an algorithmic reasoning/suggestion system to look at the game state, and make recommendations. I *could* go that far, bu

Any more ideas? Feel free to open an issue. I thrive off feedback.

#What I need from YOU (yes, you):
I need members of the community to parse through my "reccomended counter items" list, and give me hints on what to add/remove.
