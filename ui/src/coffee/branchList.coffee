###
Component for managing which branches are active
	-- add branch
	-- delete branch
	-- search brancehs
NOTE: This probably does too many things at the moment
###

module.exports = BranchList = React.createClass
	refForInput: "branchInputField"
	BUILDING_CLASS: "building-branch"
	WAITING_CLASS: "waiting-branch"

	getInitialState: ->
		allBranches: @props.allBranches or {}
		searchText: ""
		# when user clicks 'build' or 'delete, store in these two objs so we can show 'building...' (this status is also stored on back end)
		postedToBuild: {}
		postedToDelete: {}

	# take builtBranches from parent and make sure to update the ones you've posted to add/delete
	componentWillReceiveProps: (newProps) ->
		@setState
			allBranches: newProps.allBranches or []
		, @updatePostedBranches

	render: ->
		if _.isEmpty @state.allBranches
			return `<div className="main-page-spinner"><i className="fa fa-spinner fa-spin fa-5x"/></div>`

		# loop through all available branches and put list items into two lists, built and unbuilt
		builtJSX = []
		unbuiltJSX = []
		for branchName, branchInfo of @state.allBranches
			continue if not branchInfo.isInSearch
			status =
				if @state.postedToBuild[branchName] # we've marked it on client side as 'building'
					"Building"
				else if @state.postedToDelete[branchName]
					"Unbuilt"
				else
					branchInfo.Status

			jsx = @formatBranch branchName, branchInfo, status
			if status is "Built"
				builtJSX.push jsx
			else
				unbuiltJSX.push jsx

		inputPlaceHolder = "Search remote branches of #{@props.repoName}..."

		`<div className="branch-manager">
			<input
				type="text"
				className="enter-branch"
				onChange={this.onChangeSearch}
				placeholder={inputPlaceHolder}
				ref={this.refForInput}
			/>
			<ul className="branch-list">
				<div className="list-section-header">BUILT BRANCHES ({builtJSX.length})</div>
				{builtJSX}
				<div className="list-section-header">UNBUILT BRANCHES ({unbuiltJSX.length})</div>
				{unbuiltJSX}
			</ul>
		</div>
		`

	formatBranch: (branchName, branchInfo, status) ->

		switch status
			when "Built"
				directUrl = @props.urlForBranch branchInfo
				`<li key={branchName} className="built-branch">
					<span className="build-button-container">
						<a
							className="delete-branch"
							onClick={this.deleteBranch.bind(this, branchInfo.Name)}
						>
							remove
						</a>
					</span>
					<a className="branch-name" href={this.urlForBranchPage(branchInfo)}>{branchInfo.Name}</a>
					<a
						className="direct-url"
						href={directUrl}
						>
							direct
							<i className="fa fa-external-link"/>
					</a>

				</li>`

			when "Building", "Waiting" # the server has marked as Building
				className = if status is "Building" then @BUILDING_CLASS else @WAITING_CLASS
				`<li key={branchName} className={className}>
					{this.getBuildingLabelJSX(status)}
					<span className="branch-name">{branchInfo.Name}</span>
				</li>`

			when "Unbuilt"
				className = "not-built-branch"
				`<li key={branchName} className={className}>
					{this.getUnbuiltLabelJSX(branchName)}
					<span className="branch-name">{branchName}</span>
				</li>`


			else
				null

	componentDidMount: ->
		@setPrunedBranches()

	getBuildingLabelJSX: (status) ->
		`<span className="build-status">
			<i className="fa fa-spinner fa-spin status-spinner"/>
			<span className="build-button-container">
				<span>{status}</span>
			</span>
		</span>
		`

	# returns jsx for the list item of an unbuilt branch
	getUnbuiltLabelJSX: (branchName) ->
		`<span>
			<span className="build-button-container">
				<a
					className="build-branch"
					onClick={this.postNewBranch.bind(this, branchName)}
				>
					create
				</a>
			</span>
		</span>
		`

	focusSearch: ->
		@refs[@refForInput].getDOMNode().focus() # focus the search box

	#
	# Input methods
	#
	onChangeSearch: (e) ->
		@setState {searchText: e.target.value.trim().toLowerCase()}, @setPrunedBranches

	# calculate which branches to return for text search
	setPrunedBranches: ->
		allBranches = @state.allBranches
		for branchName, info of allBranches
			info.isInSearch = branchName.toLowerCase().indexOf(@state.searchText) isnt -1
		@setState {allBranches}

	postNewBranch: (branchName) ->
		{postedToBuild} = @state
		postedToBuild[branchName] = branchName
		@setState {postedToBuild}
		$.ajax
			type: "POST"
			url: @props.getBranchUrl()
			data: name: branchName
			success: (results) =>
				@props.alteredBranch()
			error: (results) ->
				console.log 'post error with results: ', results


	deleteBranch: (branchName) ->
		{postedToDelete} = @state
		postedToDelete[branchName] = branchName
		@setState {postedToDelete}
		$.ajax
			type: "DELETE"
			url: "#{@props.getBranchUrl()}?name=#{branchName}" # cant pass as data option, must be full URI for DELETEs
			success: (results) =>
				@props.alteredBranch()
			error: (results) ->
				console.log 'delete error with results: ', results

	# loop through the postedToBuild and postedToDelete and if any have been reported as built or deleted, remove them
	updatePostedBranches: ->
		{postedToBuild, postedToDelete} = @state

		for posted of postedToBuild
			match = _.filter @state.allBranches, (branch) ->
				branch.Name is posted and branch.Status is "Built"

			if match.length isnt 0
				delete postedToBuild[posted]

		for deleted of postedToDelete
			match = _.filter @state.allBranches, (branch) ->
				branch.Name is deleted and branch.Status is "Unbuilt"

			if match.length is 0
				delete postedToDelete[deleted]

		@setState {postedToDelete, postedToBuild}, @setPrunedBranches

	urlForBranchPage: (branch) ->
		encodedBranch = encodeURIComponent branch.Name
		"http://#{@props.apiHost}:#{@props.UIPort}/?branch=#{encodedBranch}"
