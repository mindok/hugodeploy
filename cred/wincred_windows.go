package cred

import (	
	"errors"	
	"strings"
	"github.com/danieljoos/wincred"
	"github.com/spf13/viper"
)

/**
	Credential Provider which retrieves the credentials from the Windows Credential Manager
	(Windows only, naturally)
**/
type WinCredProvider struct {			
}

func (cstore WinCredProvider) GetCredentials(connectionType string) (credentials *Credentials, err error) {	
	identifier := viper.GetString(connectionType + ".wincred-identifier")
	
	if (identifier == "") {
		return nil, errors.New("Wincred identifier not found. Define " + connectionType +".wincred-identifier in config file.")
	}

	cred, err := wincred.GetGenericCredential(identifier)		
    if err != nil {        
		return nil, err			
    }
	
	// the returned password contains null bytes between
	// the characters for some reason -> remove these
	pwd := string(cred.CredentialBlob)
	pwd = strings.Replace(pwd, string(0), "", -1)
	
	return &Credentials{ UID: cred.UserName, PWD: pwd }, nil			
}