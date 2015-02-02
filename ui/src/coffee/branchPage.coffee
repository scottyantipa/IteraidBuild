# Manages a page shownig details of a branch, including iframe and comments

module.exports = BranchPage = React.createClass

	render: ->
		{branch, urlForBranch} = @props
		{Name, Status} = branch
		url = urlForBranch(branch)

		`<div className="branch-page">
			<div className="iframe-container">
				<a
					className="branch-url"
					href={url}
				>
					Show "{Name}" in a separate window
				</a>
				<iframe ref="iframe" src={url}></iframe>
			</div>
		</div>
		`

	componentDidMount: ->
		$(window).on "resize", _.debounce @onResize, 100
		@onResize()


	onResize: ->
		if $frame = $ @refs["iframe"].getDOMNode()
			$frame.css "height", $(window).height() * .9
			$frame.css "width", $(window).width() * .98