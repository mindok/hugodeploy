package cred

type Credentials struct {
	UID string
	PWD string	
}

type CredentialProvider interface {	
	// connectionType e.g. ftp, sftp
	GetCredentials(connectionType string) (*Credentials, error)
}

func GetCredentialProvider(identifier string) CredentialProvider {
	switch identifier {
		case "classic":
			return ClassicCredProvider{}
		case "interactive":
			return InteractiveCredProvider{}		
		case "pwdfile":
			return PwdFileCredProvider{}
		case "wincred":
			return WinCredProvider{}
		default:
			return nil;
	}
}
