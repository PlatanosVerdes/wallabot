# Using the bot

## The bot

1. `/start` to join.
2. Write what you are looking for and it is saved (a greeting gets a joke instead). A price and a radius are picked out of the text:

   | Written | Searches |
   | :--- | :--- |
   | `iphone 13` | "iphone 13", all of Spain |
   | `kallax hasta 40` · `máx 40` · `menos de 40` | up to 40 € |
   | `moto desde 1.500` · `más de 1500` | from 1500 € |
   | `bici 100-300` · `entre 100 y 300` | 100 to 300 € |
   | `sofá en Sant Cugat` | 30 km around Sant Cugat |
   | `sofá en Sant Cugat a 10 km` | 10 km around Sant Cugat |

   The town after the last `en` is looked up on OpenStreetMap; when it is no town (`funda en piel`), it stays part of the text.
3. Optional: send a location (📎 → Location, any point on the map) to move the latest search there, 30 km unless the text said otherwise.

A listing shared from the app (`wallapop.com/item/...`) becomes a search for things like it: the first three words of its title, in its category, up to 20% dearer.

For filters the text cannot say (category, condition, brand and model), make the search on es.wallapop.com in a phone browser (typed in the address bar, so it does not jump to the app) and paste its address. Whatever is written next to it is the name, e.g. `coches top https://es.wallapop.com/search?...`

| Command | What it does |
| :--- | :--- |
| `/nueva <search>` | Same as writing it, or pasting the address |
| `/busquedas` | Your searches: 🔔 silence, ✏️ rename, 🗑 delete, and when the last round ended |
| `/ahora` | Check one search, or all, right now |
| `/baja` | Leave and delete your searches |
| `/ayuda` | Help |

Each listing arrives with its photo, price, town, the details Wallapop shows for it (year and km for a motorbike, condition, colour and material for furniture, whether it ships), and two buttons: open it, or silence that search. Details already in the title are left out, and a listing whose details cannot be read is sent without them.

Each chat only sees its own searches. Up to `WALLA_MAX_USERS` chats, `WALLA_MAX_SEARCHES` searches each.
