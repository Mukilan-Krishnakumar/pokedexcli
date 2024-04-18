package main

import (
  "fmt"
  "os"
  "bufio"
  "net/http"
  "errors"
  "encoding/json"
)

type cliCommand struct{
  name string
  description string
  callback func(*config) error
}

type config struct{
  Next *string
  Previous *string
}

type pokeJSON struct{
  Count int `json:"count"`
  Next *string `json:"next"`
  Previous *string `json:"previous"`
  Results []struct{
    Name string `json:"name"`
    URL string `json:"url"`
  } `json:"results"`

}

func getMap() map[string]cliCommand{
  getMap := map[string]cliCommand{
    "help" : {
        name:        "help",
        description: "Displays a help message",
        callback:    commandHelp,
    },
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
    "map" : {
      name: "map",
      description: "Map forward",
      callback: commandMap,
    },
    "mapb" : {
      name: "mapb",
      description: "Map backward",
      callback: commandMapB,
    },


  }
  return getMap
}

func commandHelp(c *config) error{
  helpText := "Welcome to the Pokedex!\nUsage:\n"
  getMap := getMap()
  fmt.Println(helpText)
  for _, command := range getMap{
    fmt.Println(fmt.Sprintf("%s: %s",command.name, command.description))
  }
  fmt.Println()
  return nil
}


func commandExit(c *config) error{
  os.Exit(0)
  return nil
}

func commandMap(c *config) error{
  if c.Previous == nil{
    init_url := "https://pokeapi.co/api/v2/location-area/"
    res, err := http.Get(init_url)
    if err != nil{
        return errors.New("Problem hitting the API")
      }
    defer res.Body.Close()
    decoder := json.NewDecoder(res.Body)
    var pokejsondata pokeJSON
    decoder.Decode(&pokejsondata)
    c.Previous = &init_url
    c.Next = pokejsondata.Next
    for _, location := range pokejsondata.Results{
      fmt.Println(location.Name)
    } 
    fmt.Println()
    return nil

  }
  if c.Next == nil{
    return nil
  }
  res, err := http.Get(*c.Next)
  if err != nil{
    return errors.New("Problem hitting the API")
  }
  defer res.Body.Close()
  decoder := json.NewDecoder(res.Body)
  var pokejsondata pokeJSON
  decoder.Decode(&pokejsondata)
  c.Next = pokejsondata.Next
  c.Previous = pokejsondata.Previous
  for _, location := range pokejsondata.Results{
    fmt.Println(location.Name)
  }
  fmt.Println()
  return nil
}

func commandMapB(c *config) error{
  if c.Previous == nil{
    return nil
  }
  res, err := http.Get(*c.Previous)
  if err != nil{
    return errors.New("Problem hitting the API")
  }
  defer res.Body.Close()
  decoder := json.NewDecoder(res.Body)
  var pokejsondata pokeJSON
  decoder.Decode(&pokejsondata)
  c.Next = pokejsondata.Next
  c.Previous = pokejsondata.Previous
  for _, location := range pokejsondata.Results{
    fmt.Println(location.Name)
  }
  fmt.Println()
  return nil
}


func main(){
  prompt := "pokedex > "
  getMap := getMap() 
  var actualconf config
  config := &actualconf
  fmt.Printf(prompt)
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    fmt.Println()
    val, ok := getMap[scanner.Text()] 
    if ok{
      val.callback(config)
    }

    fmt.Printf(prompt)
  }
}
