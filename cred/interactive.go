package cred

import (		
	"fmt"
	"os"
	"bufio"
	"github.com/eiannone/keyboard"
	"github.com/spf13/viper"
)

/**
	CredentialProvider which asks the user to input their credentials on stdin
	once they are required
**/
type InteractiveCredProvider struct {			
}

func readPwd() (pwd string, err error) {
	s := ""

	for {
		char, key, err := keyboard.GetSingleKey()
	
		if (err != nil) {
			return "", err
		}
	
		if (key == keyboard.KeyEnter) {
			break
		}
		
		s += string(char)
	}
	
	return s, err
}

func (cstore InteractiveCredProvider) GetCredentials(connectionType string) (credentials *Credentials, err error) {	
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter credentials for connection type " + connectionType);
	
	// is the username preset?
	uname := viper.GetString(connectionType + ".user")	
	
	if (uname != "") {
		fmt.Println("Username: " + uname);
	} else {
		// user has to enter username
		fmt.Print("Enter username: ");				
		unameBy, _, err := reader.ReadLine()
		uname = string(unameBy)

		if (err != nil) {
			return nil, err
		}
	}
	
	fmt.Print("Enter password (input will be hidden): ");
	pwd, err := readPwd();	
	
	if (err != nil) {
		return nil, err
	}
		      
	return &Credentials{ UID: uname, PWD: pwd }, nil		    	
}

