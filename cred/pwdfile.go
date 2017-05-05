package cred

import (	
	"errors"
	"github.com/spf13/viper"
)

/**
	CredentialProvider which reads the credentials from a separate file
**/
type PwdFileCredProvider struct {			
}

func (cstore PwdFileCredProvider) GetCredentials(connectionType string) (credentials *Credentials, err error) {	

	pwdfile := viper.GetString(connectionType + ".pwdfile")
	
	pwdConfig := viper.New()
    pwdConfig.SetConfigName(pwdfile)
    pwdConfig.AddConfigPath(".")
    err = pwdConfig.ReadInConfig()
	if err != nil {
        return nil, err
    }
	
	uid := pwdConfig.GetString(connectionType + ".user")
	pwd := pwdConfig.GetString(connectionType + ".pwd")

	serr := ""
	
	if uid == "" {
		serr = serr + "UID not found. Define " + connectionType + ".user in credential config file. "
	}
	if uid == "" {
		serr = serr + "PWD not found. Define " + connectionType + ".pwd in credential config file. "
	}
	
	if (serr != "") {
		return nil, errors.New(serr)
	}

	return &Credentials{ UID: uid, PWD: pwd}, nil
}

