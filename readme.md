***

*Much thanks to [deadlock-api.com](https://deadlock-api.com/). You guys have created quite possibly the easiest API I've ever worked with in my life.*

## TL;DR

1.  Type in your Steam ID.
2.  The API looks up to see if you're in a match.
3.  The API checks what team you're on.
4.  The API parses the heroes on the enemy team.
5.  For every enemy hero, you get a formatted recommendation for what items to buy.
6.  If the same item is recommended multiple times, it will be highlighted.

> **Disclaimer:** Do not use this as the primary source of truth for what to buy. Item builds are highly complex, contextual, and require a lot of intuition. This tool is for players who are struggling to find or buy counters for particular heroes. Use this tool for suggestions—not just for you, but for your entire team.

## Why I Made This

I have played too many games where no one on my team buys counter items. I've also noticed people on Discord and Reddit complain about characters without doing the bare minimum to counter them.

## Future Work

-   **Parse the character you're playing and match items with opponents.**
    -   **Example:** Against high-fire-rate heroes, a character like Mcginnis loves Suppressor, but a frontliner like Shiv wants Juggernaut.

-   **Classify items into early, mid, and late game.**
    -   **Example:** Against Grey Talon, Spirit Shielding is fantastic against his all-ins during the laning phase. Late game, Witchmail or Spellbreaker would offer a lot more value.

-   **Distinguish between laning items and mid-game items.**
    -   **Example:** Itemizing against Sinclair is very different in lane versus mid-game. You shouldn't get anti-poke recommendations if your lane opponents don't have poke characters.

-   **Identify "top performers" mid-game and highlight items to build against them.**
    -   **Example:** The current system will give you advice on how to build against a 0/10 Grey Talon that's AFK. It can't yet parse the current game state to give updated recommendations.

-   **Add reasoning for every recommended item.**
    -   **Example:** Counter items aren't always intuitive. For Counterspell, you could add a note like, "Can be used against Shiv's ultimate." For Slowing Hex, "Use before Pocket throws out his coat/suitcase."

-   **Identify build archetypes and item/hero combinations.**
    -   **Example:** Sinclair has two main build archetypes: spirit poke items (expansion on Q, extra charge/compress cooldown) or utility/control items (expansion on E, duration extender).

-   **Identify "counter-counter" items.**
    -   **Example:** Knockdown is great against Seven. However, once he gets Unstoppable, it's not as good. Knockdown should be recommended at first, then de-prioritized.

A lot of this requires an algorithmic reasoning/suggestion system to look at the game state and make recommendations. I *could* go that far, but it's probably not worth the effort at this time.

*Any more ideas? Feel free to open an issue. I thrive on feedback!*

## How You Can Help

I need members of the community to parse through my "recommended counter items" list and give me feedback on what to add or remove.
