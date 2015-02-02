# Iteraid
#### *The aid to software iteration*

Iteraid lets you and anyone working on your web application view a live version of that app for every outstanding git branch.  This lets you iterate crazy fast.  No more pulling/building code all the time just to see someone's changes. Just open your browser, go to your iteraid instance, and click on the branch name and you've got a live app.  

Athletes use Gatorade, developers use Iteraid.


Iteraid requires a web application that can be served as a static, single page app.  It works by building that version of the application, then setting up a dedicated static server which serves that version of the app on a new port.

### Dependencies
1. Golang, with a GOPATH including this Iteraid directory
2. Mongodb running locally
3. npm installed globally
4. bower installed globally

### To configure, you must provide 2 files in ./user dir:

1. ./user/config.json file like this:

```
	repoName: name of the repo to serve
	repoSSH: the SSH (important) url to the repo for 'git clone the_url'
	indexHTMLName: the path and name to your main html file after build
	startPort: first port ot start the branch servers on (e.g. 2010).
	mongoPort: port on which your local mongo is running,
	hostName: host name of your machine,
	mainPort: port for setup the webserver
    UIPort: port for the ui to run on
```

2.  A build script in ./user/build_branch.sh for packaging your app into a stand alone directory which can be served by a static file server.  The easiest script to put here is simply to recursively copy you're whole directory.  However, for speed and better storage, it is best to create a build that has only the necessary assets.


### You must build both the server and the client

To build both client and server
  1. git clone git@git.soma.salesforce.com:santipa/Iteraid.git
  2. cd Iteraid
  3. make install
  4. git submodule init (better way to do this?)
  5. git submodule update
  6. make
  7. make repo
  8. make serve

To watch/recompile UI files
  1. $ cd ui
  2. $ make watch

