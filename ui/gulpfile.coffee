gulp = require 'gulp'
browserify = require 'browserify'
concat = require 'gulp-concat'
coffeeify = require 'coffeeify'
plumber = require 'gulp-plumber'
stylus = require 'gulp-stylus'
clean = require 'gulp-clean'
reactify = require 'reactify'
transform = require 'vinyl-transform'
rename = require 'gulp-rename'

gulp.task 'clean', ->
	gulp.src "./public/**/*.*", read: false
		.pipe plumber()
		.pipe clean()

gulp.task 'jsSrc', ->
	browserifyCoffeeReact = ->
		transform (file) ->
			browserify debug: false
				.add file
				.transform coffeeify
				.transform everything: true, reactify
				.bundle()

	gulp.src './src/coffee/start.coffee'
		.pipe plumber()
		.pipe browserifyCoffeeReact()
		.pipe concat 'app.js'
		.pipe gulp.dest 'public/js/'

# dependency javascripts (unminified for now)
jsDepNames = [
	"jquery/dist/jquery.js"
	"react/react-with-addons.js"
	"underscore/underscore.js"
]
jsDepPaths = ("./bower_components/#{depName}" for depName in jsDepNames)
gulp.task 'jsDeps', ->
	gulp.src jsDepPaths
		.pipe concat 'dependencies.js'
		.pipe gulp.dest './public/js'

# copy static style_assets over from src to public
gulp.task 'style_assets', ->
	gulp.src ["./src/styles/**/*.*", "!./src/styles/css/app.styl"]
		.pipe plumber()
		.pipe gulp.dest './public/styles/'

# compile the stylus
gulp.task 'stylus', ->
	gulp.src "./src/styles/css/app.styl"
		.pipe plumber()
		.pipe stylus()
		.pipe gulp.dest './public/styles/css'

gulp.task 'config', ->
	gulp.src "../user/config.json"
		.pipe plumber()
		.pipe gulp.dest './public/'

gulp.task 'html', ->
	gulp.src "./src/index.html"
		.pipe plumber()
		.pipe gulp.dest './public/'

# watch
gulp.task 'watch', ->
	gulp.watch './src/coffee/*.coffee', ['jsSrc']
	gulp.watch '../user/config.json', ['config']
	gulp.watch './src/styles/**/*.*', ['stylus'] # if any style src file changes, build them and copy into ./public

gulp.task 'default', ['jsSrc', 'jsDeps', 'stylus', 'style_assets', 'config', 'html']
