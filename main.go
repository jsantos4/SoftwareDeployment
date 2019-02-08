package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "github.com/jamespearly/loggly"
)

type Player struct {
	Data struct {
		Type string `json:"type"`
		Attributes struct {
			GameModeStats struct {
				SoloFpp struct {
                    Wins                int     `json:"wins"`
                    RoundsPlayed        int     `json:"roundsPlayed"`
                    Top10S              int     `json:"top10s"`
                    LongestTimeSurvived float64 `json:"longestTimeSurvived"`
                    TimeSurvived        float64 `json:"timeSurvived"`
                    Kills               int     `json:"kills"`
                    HeadshotKills       int     `json:"headshotKills"`
					DamageDealt         float64 `json:"damageDealt"`
                    LongestKill         float64 `json:"longestKill"`
                    MaxKillStreaks      int     `json:"maxKillStreaks"`
					Heals               int     `json:"heals"`
				} `json:"solo-fpp"`
            }`json: "GameModeStats"`
        }`json: "attributes"`
    }`json: "data"`
}

func apiRequest() *http.Response {
    url := "https://api.pubg.com/shards/steam/players/account.ab3ebdd2cbe44e5996bad678dd15ac3b/seasons/lifetime"

   var bearer = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI0NTU1YTNhMC0wYWNlLTAxMzctZDgzZS01OTE5NGVkMTExNzMiLCJpc3MiOiJnYW1lbG9ja2VyIiwiaWF0IjoxNTQ5MzAwNDc2LCJwdWIiOiJibHVlaG9sZSIsInRpdGxlIjoicHViZyIsImFwcCI6ImJhZ3VubmVyMzAwIn0.7SZyFCrSf9_LaoSdmqSJ-IqGj9ONcTs1ehPFlTp05Rw"
   var the = "application/vnd.api+json"
   req, err := http.NewRequest("GET", url, nil)

   req.Header.Add("Authorization", bearer)
   req.Header.Add("Accept", the)

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
       fmt.Println("Error on response.\n[ERRO] -", err)
   }

   return resp
}

func printStats(player *Player) {
    fmt.Println("Stats")
    fmt.Println("Wins: ", player.Data.Attributes.GameModeStats.SoloFpp.Wins)
    fmt.Println("Rounds played: ", player.Data.Attributes.GameModeStats.SoloFpp.RoundsPlayed)
    fmt.Println("Top 10s: ", player.Data.Attributes.GameModeStats.SoloFpp.Top10S)
    fmt.Println("Longest time survived: ", player.Data.Attributes.GameModeStats.SoloFpp.LongestTimeSurvived)
    fmt.Println("Total time survived: ", player.Data.Attributes.GameModeStats.SoloFpp.TimeSurvived)
    fmt.Println("Kills: ", player.Data.Attributes.GameModeStats.SoloFpp.Kills)
    fmt.Println("Headshots: ", player.Data.Attributes.GameModeStats.SoloFpp.HeadshotKills)
    fmt.Println("Damage dealt: ", player.Data.Attributes.GameModeStats.SoloFpp.DamageDealt)
    fmt.Println("Longest Kill: ", player.Data.Attributes.GameModeStats.SoloFpp.LongestKill)
    fmt.Println("Highest killstreak", player.Data.Attributes.GameModeStats.SoloFpp.MaxKillStreaks)
    fmt.Println("Times healed", player.Data.Attributes.GameModeStats.SoloFpp.Heals)
}


func main() {

    response := apiRequest()

    defer response.Body.Close()

    byteValue, _ := ioutil.ReadAll(response.Body)

    var stats Player

    json.Unmarshal(byteValue, &stats)

    client := loggly.New("PUBG api")

    logContent := client.Send("info", string(stats.Data.Attributes.GameModeStats.SoloFpp.RoundsPlayed))

}
