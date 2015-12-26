# hugodeploy
Simple SFTP deployment tool for static websites (e.g. created by Hugo) with minification.

This was built to allow deployment to cheap and nasty hosting providers such as bluehost and namecheap with no dependencies on third party deployment systems.

## Basic Usage (Assuming you are using hugo)
1. Run `hugodeploy init` in the source directory of your hugo website.
2. This will create a hugodeploy.yaml (if one doesn't exist) file in that directory with some sensible defaults and create also create a tracking directory (Deployment Record Directory) for storing a copy of what's on your server.
3. Edit


## Warnings
SFTP username and password are stored in plaintext in the config file. Probably not a good idea to check your config file into a public repository.

## A bit more detail
```bash
hugodeply init
```
Initialises the configuration file and required directories.
If it has been run previously it won't create a configuration file, but it will create the deploymentRecordDir directory or clear it if it already exists.

If, for whatever reason, your deploymentRecordDir directory gets out of sync with the deployment target you can re-run this to force hugodeploy to resubmit all the files to the deployment target.


### Credits
Many thanks to [Steve Francia](http://github.com/spf13) for [hugo](http://github.com/spf13/hugo) - A Fast and Flexible static site generator, the reason hugodeploy is needed in the first place. Also, thanks for the supporting libraries such as [cobra](http://github.com/spf13/cobra) and [viper](http://github.com/spf13/viper) that made building this a whole lot easier.

SFTP library provided by [pkg](https://github.com/pkg/sftp).

Minification library from [tdewolff](https://github.com/tdewolff/minify).

