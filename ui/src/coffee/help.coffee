# Display a help message to user
# Toggle visible not visible

module.exports = Help = React.createClass
	QUICK_DESCRIPTION: """
		Iteraid helps you iterate on your web software.   
		It lets you build all of the versions your team is working on, giving you a url that points to that built version.
		Click on 'build' next to a branch to make it live.  When you no longer need the live version, just click 'delete'.
	"""	
	


	getInitialState: ->
		visible: false

	render: ->
		{visible} = @state
		label = if visible then "Hide Help" else "Show Help"
		messageJSX =
			if visible
				`<div className="help-message">
					<div className="quick-description">{this.QUICK_DESCRIPTION}</div>
				</div>
				`
			else
				null

		`<div className="help-user">
			{messageJSX}
			<div className="help-toggle" onClick={this.toggleVisible}>
				{label}
			</div>
		</div>
		`
	toggleVisible: -> @setState {visible: !@state.visible}

