# Basic chrome wrapper

module.exports = Chrome = React.createClass
	render: ->
		`<div className="chrome">
			<span className="logo"><a href="/">Iteraid</a></span>
			<span className="catch-phrase">The aid to product iteration</span>
			<div className="git-link">
				<a href="https://github.com/scottyantipa/Iteraid">
					<i className="fa fa-github"/>
				</a>
			</div>
		</div>
		`
