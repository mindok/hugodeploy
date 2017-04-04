# hugodeploy
Simple FTP deployment tool for static websites (e.g. created by [Hugo](https://gohugo.io/)) with built-in minification.

## Why use hugodeploy?
This was built to allow easy deployment to el-cheapo hosting providers, such as bluehost and namecheap, with no dependencies on third party deployment systems. It keeps a local copy of what has already been deployed and figures out what's different each time it is run so it minimises transfers to your web server.

It was originally written to support the deployment of two or three websites I've been working on, starting with [A Philosopher and A Businessman](http://philosopherbusinessman.com/).

It is totally independent of any code repository - it just needs the output from whatever static site generator you are using and ftp details for where to put it.

## What alternatives are there?
### SaaS Continuous Deployment
[This article](http://jice.lavocat.name/blog/2015/hugo-deployment-via-codeship/) talks about using [Codeship](https://codeship.com/), and [this one](https://gohugo.io/tutorials/automated-deployments/) uses [Wercker](http://wercker.com) - both SaaS continous integration and deployment options. I decided against these as there are just more dynamic dependencies to worry about keeping up to date. They do have the advantage that they respond to git commits - i.e. your public site gets updated automatically when you commit changes.

### Bash scripts
Again, more dynamic dependencies, plus if I have to write code I'd rather write go.

### Rsync
Rsync should work ok, but some of the lower cost hosting providers don't support it as well as they should. I wanted to use what works out of the box.

## How does it work?
hugodeploy keeps a local copy of the latest version of all files successfully sent to to the deployment target. hugodeploy does a binary compare of the file contents ready to deploy with this local copy to determine whether a file needs to be deployed. This is handy where you have images, videos, bloated javascript libraries etc that are slow to send - they only get sent once.

hugodeploy minifies html, css, js, json and XML by default prior to deploying. You can disable this using the DontMinify option in the config file or the -m flag.

## TODOs
1. Fix up path handling for directories so they can be relative to working directory rather than absolute
2. Modify ftp invocation infrastructure so it is substitutable with another deployment method (e.g. sftp, scp). sftp is kinda done, but not tested AT ALL, and not plumbed in
3. Allow file ignores (like .gitignore) so we don't get random stuff like .DS_Store sent over the wire. Done at a naive level - good enough for me.
4. <del>Allow specification of website root in ftp client</del> DONE
5. Clean up some of the interaction between package level variables, command line flags & viper in cmd/root.go
6. Possible refactor to push down connection of DeployScanner to appropriate Deployer into deploy package rather than handling in push & preview commands.
7. Implement directory delete in ftp. This will need to be done in the source library first.

## Installation
Currently there are no pre-built binaries so you will need go installed. See [https://golang.org](https://golang.org) for instructions.
Run `go get github.com/mindok/hugodeploy`
Change directory to $GOPATH/src/github.com/mindok/hugodeploy
Run `go build`
You should now have a hugodeploy binary that you can put somewhere on your path.

## Basic Usage (Assuming you are using hugo)
1. Navigate to the source directory of your website
2. Run `hugodeploy init` in the source directory of your hugo website.
3. This will create a hugodeploy.yaml (if one doesn't exist) file in that directory with some sensible defaults.
4. Edit hugodeploy.yaml, in particular to set up your ftp host, username and password. Most other options should be ok. If you want to specify a directory for tracking the deployed files you can the deployRecordDir and create that directory manually. Otherwise...
5. [Optional] Run `hugodeploy init` again. This will now create a tracking directory (Deployment Record Directory) for storing a copy of what's on your server.
6. Run `hugodeploy preview` to see a list of what will be sent to your deployment target
7. Run `hugodeploy push` to upload files to your deployment target


## Warnings
FTP username and password are stored in plaintext in the config file. Probably not a good idea to check your config file into a public repository.

Paths currently should be absolute rather than relative to the working directory.

It is designed for web-scale files (i.e. images < 10Mb, HTML, JS, CSS files). Each file is loaded into memory in its entirety before it is processed. Transferring Gb size files probably isn't a good idea with this tool.

Not tested on any platform other than Mac.

## Commands

### init
```bash
hugodeploy init [flags]
```
Initialises the configuration file if one doesn't exist.

If there is a configuration file, it will create the deployRecordDir directory or clear it if it already exists.

If, for whatever reason, your deploymentRecordDir directory gets out of sync with the deployment target you can re-run this to force hugodeploy to resubmit all the files to the deployment target.

Run `hugo init -h` or `hugo init --help` for information on available flags

### preview
```bash
hugodeploy preview [flags]
```
Generates a list of deployment actions based on differences between sourceDir and deploymentRecordDir (and whether minification is disabled or not).

Run `hugo preview -h` or `hugo preview --help` for information on available flags

### push
```bash
hugodeploy push [flags]
```
Performs deployment actions based on differences between sourceDir and deploymentRecordDir (and whether minification is disabled or not).

Run `hugo push -h` or `hugo push --help` for information on available flags

## Options
Life is easier if you set all the options in the config file, call the config file hugodeploy.yaml and place it in the source directory for your hugo website. Then set the current working directory to the source directory for your hugo website before running the commands. However, if you want a little more control here are the available options

### ConfigFile
Specifies the location of the configuration file. By default hugodeploy looks for a file called hugodeploy.xxx in the current working directory, where xxx indicates a format supported by [viper](http://github.com/spf13/viper) - currently JSON, TOML, YAML and HCL.

**Note for TOML config files:** Global settings (e.g. sourcedir, deployRecordDir) need to be at the top, above grouped settings such as ftp. This is a [limitation of TOML](https://npf.io/2014/08/intro-to-toml/). Thanks to [Hernán](https://github.com/hfoffani) for finding that one.


You can override the name and location of the config file using the --config flag. e.g.
```bash
hugo push --config 'whydoihavetodothingsthehardway.toml'
```

### sourceDir
TODO: Currently these should be set as absolute paths.

Specifies the source location of the directory to be deployed to the deployment target. For a typical hugo installation this will be 'public'.

Should generally be set in the config file (sourceDir option), but you can also set on the command-line using --sourceDir or -s. I'm not sure why you want to do that, but I was having fun exploring [cobra](http://github.com/spf13/cobra) & [viper](http://github.com/spf13/viper) so thought I'd put it in.

### deployRecordDir
TODO: Currently these should be set as absolute paths.

Specifies the location of the directory used to track what has been deployed. It defaults to 'deployed'.

Should generally be set in the config file (deployRecordDir option), but you can also set on the command-line using --deployRecordDir or -d.

### DontMinify Option
Disables minification. Can be set in the config file (DontMinify), or on the command-line. Command flags are -m or --dontminify.

Note that changing this is likely to cause all minifiable files (HTML, CSS etc) to be resent as the file compare operates with what was previously sent, and is done post-minification.

Minification is performed by the [tdewolff/minify](https://github.com/tdewolff/minify) library.

### FTP Options
Sets the host, username, password and root directory for the FTP deployment target. Can only be set in the config file as follows:
```
ftp:
  host: <host ip or name>
  user: <username>
  pwd: <password>
  rootdir: <root directory of website relative to root of ftp server. e.g. / or /public_html/ >
  disabletls: <optional. Should be false unless troubleshooting, or if your ftp server is known not to support TLS.>
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

### Troubleshooting FTP connections
Most problems with hugodeploy are related to FTP connections and the widely differing implementation of the FTP specification in different servers. 

The first thing to try is disabling TLS. This is generally a bad idea as your password and data will be transmitted in clear text. However, some servers just don't have a TLS connection option. To disable TLS, use the ftp.disabletls option in the configuration file.

The next thing to do is turn on Verbose or Debug mode either in the config file or using the -v or -d command line switches. This activates goftp's debugging output and you will be able to see the commands going back and forth between hugodeploy and your ftp server.

If neither of those work, log an issue on github with some details from the logged output. Be sure to remove passwords and other sensitive details from the log - they appear in a couple of places.

## A few notes on code organisation
deploy.DeployScanner traverse all files in sourceDir and compares them with what's in deployRecordDir.
A new DeployCommand is created for each difference between the two containing the details of what needs to be done to update the deployment target.

The DeployCommands thus generated are passed to the selected Deployer (currently only an FTPDeployer) for execution at the deployment target (e.g. creation of a file). Once the DeployCommand has successfully executed it is passed to a FileDeployer to update the deployRecordDir.

Feel free to suggest changes or enhancements, or send PRs for proposed code mods.

## Credits
Many thanks to [Steve Francia](http://github.com/spf13) for [hugo](http://github.com/spf13/hugo) - A Fast and Flexible static site generator. It's awesomeness inspired me to cook up this simple deployment tool. Also, thanks for the supporting libraries such as [cobra](http://github.com/spf13/cobra) and [viper](http://github.com/spf13/viper) that made building this a whole lot easier.

FTP library provided by [DutchCoders-goftp](https://github.com/dutchcoders/goftp)
- Local copy held here to allow pushing of byte array rather than file

SFTP library provided by [pkg](https://github.com/pkg/sftp). (Not implemented as yet)

Minification library from [tdewolff](https://github.com/tdewolff/minify).
