# Top level app component
Chrome = require './chrome.coffee'  # Render a chrome
BranchList = require './branchList.coffee'  # Render a list of branches
BranchPage = require './branchPage.coffee'
Help = require './help.coffee'
module.exports = App = React.createClass

	getInitialState: ->
		indexHTMLName: "index.html" # name of the repos index html file
		branchToShow: null # if user clicks to show a branch we store it here
		ajaxDone: false

	# wrap BranchList in a container so the css for BranchList is portable (the container will have the positioning css)
	render: ->
		{UIPort, ajaxDone, indexHTMLName, apiHost, mainPort, repoName, branchToShow, allBranches} = @state
		return null if not ajaxDone

		content =
			if branchToShow
				`<BranchPage
					branch={branchToShow}
					apiHost={apiHost}
					mainPort={mainPort}
					urlForBranch={this.urlForBranch}
				/>
				`
			else
				`<div className="branch-manager-container">
					<BranchList
						indexHTMLName={indexHTMLName}
						apiHost={apiHost}
						mainPort={mainPort}
						UIPort={UIPort}
						repoName={repoName}
						urlForBranch={this.urlForBranch}
						getBranchUrl={this.getBranchUrl}
						alteredBranch={this.alteredBranch}
						allBranches={allBranches}
					/>
				</div>
				`
		help = if not branchToShow then `<Help repoName={repoName}/>` else null

		`<div className="app-container">
			<Chrome/>
			{help}
			{content}
		</div>
		`

	# Load the config here just once and pass params down to branch list
	componentDidMount: ->
		# load the config then get branches
		def = @getConfig()
		def
			.then @getAllBranches
			.then =>
				@setBranchParam()
				@setState {ajaxDone: true}

	getConfig: (cb) ->
		d = $.ajax
			type: "GET"
			url: "./config.json"
			success: (results) =>
				return if not results?.indexHTMLName
				{UIPort, indexHTMLName, apiHost, mainPort, repoName} = results
				@setState {UIPort, indexHTMLName, apiHost, mainPort, repoName}, cb
			error: (results) ->
				console.warn "error loading config: ", results
		d

	#
	# Ajax
	#
	getAllBranches: ->
		$.ajax
			url: @getBranchUrl()
			type: "GET"
			data: {whichBranches: "all"}
			success: (msg) =>
				@setState {allBranches: msg?.data or []}
			error: (msg) ->
				console.log "fail: ", msg

	alteredBranch: ->
		@getAllBranches()

	#
	# Utils
	#

	urlForBranch: (branch) ->
		"http://#{@state.apiHost}:#{branch.Port}/#{@state.indexHTMLName}"

	# Util to get the "branch" param from url.  Store it, or null, in state.
	setBranchParam: ->
		for pair in location.search.substring(1).split("&") # e.g. ["branch=scott-feature", "other_param=something"]
			[key, value] = pair.split("=")
			if key is "branch" and value isnt ""
				branchName = decodeURIComponent value
			else
				continue

		@setState
			branchToShow: @state.allBranches[branchName]

	getBranchUrl: ->
		"http://#{@state.apiHost}:#{@state.mainPort}/branch"
