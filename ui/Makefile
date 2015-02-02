GULP=./node_modules/gulp/bin/gulp.js

all:
	@$(GULP) default

js:
	@$(GULP) jsSrc

install:
	npm install
	bower install

uninstall:
	rm -rf ./node_modules
	rm -rf ./bower_components

clean:
	rm -rf ./public

watch:
	@make
	@$(GULP) watch
	
