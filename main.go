package main

import (
  "fmt"
  "os"
  "bufio"
)

type cliCommand struct{
  name string
  description string
  callback func() error
}

func commandMap() map[string]cliCommand{
  commandMap := map[string]cliCommand{
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
  }
  return commandMap
}

func commandHelp() error{
  helpText := "Welcome to the Pokedex!\nUsage:\n"
  commandMap := commandMap()
  fmt.Println(helpText)
  for _, command := range commandMap{
    fmt.Println(fmt.Sprintf("%s: %s",command.name, command.description))
  }
  fmt.Println()
  return nil
}


func commandExit() error{
  os.Exit(0)
  return nil
}

func main(){
  prompt := "pokedex > "
  commandMap := commandMap() 
  fmt.Printf(prompt)
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    fmt.Println()
    val, ok := commandMap[scanner.Text()] 
    if ok{
      val.callback()
    }

    fmt.Printf(prompt)
  }
}
