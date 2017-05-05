// +build !windows

package cred

import "errors"

/**
	Credential Provider which retrieves the credentials from the Windows Credential Manager
	(Windows only, naturally)
**/
type WinCredProvider struct {			
}

func (cstore WinCredProvider) GetCredentials(connectionType string) (credentials *Credentials, err error) {		
	return nil, errors.New("The wincred credential provider is only supported on Windows");	
}