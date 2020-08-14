![Run Gosec](https://github.com/juliankoehn/wetter-service/workflows/Run%20Gosec/badge.svg)
![Test](https://github.com/juliankoehn/wetter-service/workflows/Test/badge.svg)

# Wetter Dienst regiocast

Wetter-Dienst basierend auf Analyse von anforderungen der Privater Radio-Sender
https://meteor.loverad.io/v1?lat=53.86893&lon=10.68729

Beim Debuggen der API fiel mir auf das diese zwar alle lat/longs akzeptiert, wie hier einen Ort in Libyen, allerdings dafür einen Fallback verwendet. In dem Fall wird das Wetter für "Oberstdorf" angezeigt.

Eine Alternative verwendet diese API, neben dem Fallback Config Flag der erlaubt direkt die OWM API zu Pollen, kann ein weiterer Fallback verwendet werden der die Distanz zum Ziel (haversine) berechnet und dem User die nächstmögliche Geo-Location ausgibt. Das macht dann Sinn wenn der Service für alle Standorte läuft, so kann Deutschlandweit das Wetter über eine API ausgegeben werden.

Im vergleich zur derzeitigen API dürfte ebenfalls ein Verbesserung der ResponseTimes (Latency) zu merken sein. Beim Benchen der bestehenden API kam ich teilweise auf 1.3-2.5 Sekunden. Die neue API sollte real-world Latencies liefern von 2-25ms / request.

Ebenfalls fiel mir auf das ich die API ohne Rate-Limit / CORS beschränkung spammen kann.
https://meteor.loverad.io/v1?lat=28.879878&lon=11.561894

Eine CORS implementierung wäre hier innerhalb von wenigen Minuten erledigt. Sollte hier bedarf zum nachbessern sein.


Quelle für Wetterdaten:
Die Wetterdaten dieser API stammen von https://openweathermap.org/api/one-call-api und kommen mit etwas mehr "infos". Die derzeitge API scheint 2 Requests pro Location zu verwenden um "Current" sowie "Forecast" auszugeben. Über die One-Call-API kann Current sowie Forecast in einem Zug requested werden (5 Tage, sowie Stündliche Updates)


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