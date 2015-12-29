# hugodeploy
Simple FTP deployment tool for static websites (e.g. created by [Hugo](https://gohugo.io/)) with built-in minification.

## Why use hugodeploy?
This was built to allow easy deployment to el-cheapo hosting providers, such as bluehost and namecheap, with no dependencies on third party deployment systems. 

It was originally written to support the deployment of two or three websites I've been working on, starting with [A Philosopher and A Businessman](http://philosopherbusinessman.com/).

## How does it work?
hugodeply keeps a local copy of the latest version of all files successfully sent to to the deployment target. hugodeploy does a binary compare of the file contents ready to deploy with this local copy to determine whether a file needs to be deployed. This is handy where you have images, videos, bloated javascript libraries etc that are slow to send - they only get sent once.

hugodeploy minifies html, css, js, json and XML by default prior to deploying. You can disable this using the DontMinify option in the config file or the -m flag. 

## TODOs
1. Fix up path handling for directories so they can be relative to working directory rather than absolute
2. Modify ftp invocation infrastructure so it is substituable with another deployment method (e.g. sftp, scp). sftp is kinda done, but not test AT ALL, and not plumbed in
3. Allow file ignores (like .gitignore) so we don't get random stuff like .DS_Store sent over the wire. Done at a naive level - good enough for me.
4. Allow specification of website root in ftp client
5. Clean up some of the interaction between package level variables, command line flags & viper in cmd/root.go
6. Possible refactor to push down connection of DeployScanner to appropriate Deployer into deploy package rather than handling in push & preview commands.

## Basic Usage (Assuming you are using hugo)
1. Navigate to the source directory of your website
2. Run `hugodeploy init` in the source directory of your hugo website.
3. This will create a hugodeploy.yaml (if one doesn't exist) file in that directory with some sensible defaults.
4. Edit hugodeploy.yaml, in particular to set up your sftp host, username and password. Most other options should be ok. If you want to specify a directory for tracking the deployed files you can the deployRecordDir and create that directory manually. Otherwise...
5. [Optional] Run `hugodeploy init` again. This will now create a tracking directory (Deployment Record Directory) for storing a copy of what's on your server.
6. Run `hugodeploy preview` to see a list of what will be sent to your deployment target
7. Run `hugodeploy push` to upload files to your deployment target


## Warnings
FTP username and password are stored in plaintext in the config file. Probably not a good idea to check your config file into a public repository.

It is designed for web-scale files (i.e. images < 10Mb, HTML, JS, CSS files). Each file is loaded into memory in its entirety before it is processed. Transferring Gb size files probably isn't a good idea with this tool.

## Commands

### init
```bash
hugodeply init [flags]
```
Initialises the configuration file if one doesn't exist.

If there is a configuration file, it will create the deployRecordDir directory or clear it if it already exists.

If, for whatever reason, your deploymentRecordDir directory gets out of sync with the deployment target you can re-run this to force hugodeploy to resubmit all the files to the deployment target.

Run `hugo init -h` or `hugo init --help` for information on available flags

### preview
```bash
hugodeply preview [flags]
```
Generates a list of deployment actions based on differences between sourceDir and deploymentRecordDir (and whether minification is disabled or not).

Run `hugo preview -h` or `hugo preview --help` for information on available flags

### push
```bash
hugodeply push [flags]
```
Performs deployment actions based on differences between sourceDir and deploymentRecordDir (and whether minification is disabled or not).

Run `hugo push -h` or `hugo push --help` for information on available flags

## Options
Life is easier if you set all the options in the config file, call the config file hugodeply.yaml and place it in the source directory for your hugo website. Then set the current working directory to the source directory for your hugo website before running the commands. However, if you want a little more control here are the available options

### ConfigFile
Specifies the location of the configuration file. By default hugodeploy looks for a file called hugodeply.xxx in the current working directory, where xxx indicates a format supported by [viper](http://github.com/spf13/viper) - currently JSON, TOML, YAML and HCL.

You can override the name and location of the config file using the --config flag. e.g.
```bash
hugo push --config 'whydoihavetodothingsthehardway.toml'
```

### sourceDir
Specifies the source location of the directory to be deployed to the deployment target. For a typical hugo installation this will be 'public'.

Should generally be set in the config file (sourceDir option), but you can also set on the command-line using --sourceDir or -s. I'm not sure why you want to do that, but I was having fun exploring [cobra](http://github.com/spf13/cobra) & [viper](http://github.com/spf13/viper) so thought I'd put it in.

### deployRecordDir
Specifies the location of the directory used to track what has been deployed. It defaults to 'deployed'.

Should generally be set in the config file (deployRecordDir option), but you can also set on the command-line using --deployRecordDir or -d. 

### DontMinify Option
Disables minification. Can be set in the config file (DontMinify), or on the command-line. Command flags are -m or --dontminify.

Note that changing this is likely to cause all minifiable files (HTML, CSS etc) to be resent as the file compare with what was previously sent is done post-minification.

Minification is performed by the [tdewolff/minify](https://github.com/tdewolff/minify) library.

### FTP Options
Sets the host, username and password for the SFTP deployment target. Can only be set in the config file as follows:
```
ftp: 
  host: <host ip or name>
  user: <username>
  pwd: <password>
```
Note that if you are using YAML, the indent between ftp & host is 2 spaces, not a tab. 

### Skipping files
There is a naive file and directory skipping capability that currently just does a simple string.Contains test. Substrings matched are set in the SkipFiles section of the config file as follows:
```
skipfiles:
  - .DS_Store
  - .git
  - /tmp
```


### Credits
Many thanks to [Steve Francia](http://github.com/spf13) for [hugo](http://github.com/spf13/hugo) - A Fast and Flexible static site generator. It's awesomeness inspired me to cook up this simple deployment tool. Also, thanks for the supporting libraries such as [cobra](http://github.com/spf13/cobra) and [viper](http://github.com/spf13/viper) that made building this a whole lot easier.

FTP library provided by [DutchCoders-goftp](https://github.com/dutchcoders/goftp)
- Local copy held here to allow pushing of byte array rather than file

SFTP library provided by [pkg](https://github.com/pkg/sftp). (Not implemented as yet)

Minification library from [tdewolff](https://github.com/tdewolff/minify).

