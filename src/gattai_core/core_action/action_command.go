package core_action

import (
	"fmt"
	"log"
)

func RunCmdAction(action *Action) (string, error){

	return action.Exec.RunAction(action.Name)
}

func ValidateCmdAction(action *Action) (string, error){
	result, err := action.Test.RunAction(action.Name)
	if err != nil {
		return "", fmt.Errorf("ValidateCmdAction error: %v",err)
	}
	log.Println(result)
	return action.Exec.RunAction(action.Name)
}
