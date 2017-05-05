package cred

import (	
	"errors"
	"github.com/spf13/viper"
)

/**
	CredentialProvider which reads the credentials directly from the config file
	(old behavior)
**/
type ClassicCredProvider struct {			
}

func (cstore ClassicCredProvider) GetCredentials(connectionType string) (credentials *Credentials, err error) {	

	uid := viper.GetString(connectionType + ".user")
	pwd := viper.GetString(connectionType + ".pwd")

	serr := ""
	
	if uid == "" {
		serr = serr + "UID not found. Define " + connectionType + ".user in config file. "
	}
	if uid == "" {
		serr = serr + "PWD not found. Define " + connectionType + ".pwd in config file. "
	}
	
	if (serr != "") {
		return nil, errors.New(serr)
	}

	return &Credentials{ UID: uid, PWD: pwd}, nil
}

