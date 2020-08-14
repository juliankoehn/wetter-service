# Wetter Dienst regiocast

Wetter-Dienst basierend auf Analyse von anforderungen der Privater Radio-Sender
https://meteor.loverad.io/v1?lat=53.86893&lon=10.68729

Beim Debuggen der API fiel mir auf das diese zwar alle lat/longs akzeptiert, wie hier einen Ort in Libyen, allerdings dafür einen Fallback verwendet. In dem Fall wird weiterhin das Wetter für "Oberstdorf" angezeigt.

Die ResponseTimes sind btw verdammt krass (teilweise 1500-2400ms).

Ebenfalls fiel mir auf das ich die API ohne Rate-Limit / CORS beschränkung spammen kann.
https://meteor.loverad.io/v1?lat=28.879878&lon=11.561894


Quelle für Wetterdaten:
Entsprechend: https://openweathermap.org/api/hourly-forecast - die API ermöglicht eine abfrage "By city ID" und stellt entsprechend eine file [`city.list.json`](http://bulk.openweathermap.org/sample/) bereit.


Die City.list aktualisiert sich alle paar Tage bis hin zu alle paar Monate, entsprechend macht es sinn einen Bootup in der App zu verwenden und alle 24 Stunden nach einem Update zu prüfen, die JSON Daten auszulesen und in eine Datenbank zu Speichern. Da hier kein frequentierter Zugriff drauf sein muss genügt dafür bereits eine SQLite.

Ein API Call sieht dann entsprechend so aus: `pro.openweathermap.org/data/2.5/forecast/hourly?id={city ID}&appid={your api key}`

oder alternativ: `pro.openweathermap.org/data/2.5/forecast/hourly?lat={lat}&lon={lon}&appid={your api key}` für lat/long.

Bei Regiocast werden diese Daten in einer Redis-DB gespeichert, um offenbar schnelle Zugriffe zu gewähren (0-latency). Wir können diese Anforderung in dieser Go-Application reduzieren und weichen auf einen Mem-Cache aus. Dies reduziert a) die Laufzeit Kosten der Applikation und b) wir verkleinern den Service-Layer. Ebenfalls wird ein Deployment deutlich leichter.


### Warum nicht `github.com/briandowns/openweathermap`

Das Package hat leider kein `Exception` handling, so führen 401 Errors z.B. zu Unmarshal Errors (go structs).
Mögliche alternativen:
* Contributing: Repo wirkt weitgehend inaktiv die ältesten PR sind von 2017
* Integration: eigenes Package `/omw` basierend auf gegebenen Package.

## RSH.de Geo-Service

RSH-Wetter bietet Informationen zu 12 Standorten an, zusätzlich besteht die möglichkeit "Mein Standort" zu verwenden. (siehe /rsh.de/Screen Shot 2020-08-14 at 12.45.53.png - uBlock disabled)

# Third Parties
* github.com/spf13/cobra
* github.com/kelseyhightower/envconfig
* github.com/joho/godotenv
* github.com/sirupsen/logrus
* github.com/jinzhu/gorm (why not V2?: V2 is still under developement, current state: public testing)
* github.com/labstack/echo/v4

## Testing Third Parties
* github.com/tj/assert
* github.com/stretchr/testify

Missing? see go.mod


## CLI Commands

| CLI  | Description  |
|---|---|
| `weather refresh-cities` |  Updates Cities in our Storage Engine |
| `weather search -n cityName` |  Returns a list of Cities matching name |
| `weather enable -i cityID` |  Enabled a City by ID to be parsed |
| `weather disable -i cityID` |  Disables a City by ID to be parsed |
| `weather show` |  shows a list of all enabled cities |

## Endpoints

| Endpoint | Parameters | Description  |
|---|---|---|
| `/v1?` | `lat={LATITUDE}&lon={LONGITUDE}` |  Gets weather data by lat long |
| `/v1/metrics` | | Returns metrics from our cache |