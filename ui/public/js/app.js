(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
/** @jsx React.DOM */var App, BranchList, BranchPage, Chrome, Help;

Chrome = require('./chrome.coffee');

BranchList = require('./branchList.coffee');

BranchPage = require('./branchPage.coffee');

Help = require('./help.coffee');

module.exports = App = React.createClass({displayName: "App",
  getInitialState: function() {
    return {
      indexHTMLName: "index.html",
      branchToShow: null,
      ajaxDone: false
    };
  },
  render: function() {
    var UIPort, ajaxDone, allBranches, apiHost, branchToShow, content, help, indexHTMLName, mainPort, repoName, _ref;
    _ref = this.state, UIPort = _ref.UIPort, ajaxDone = _ref.ajaxDone, indexHTMLName = _ref.indexHTMLName, apiHost = _ref.apiHost, mainPort = _ref.mainPort, repoName = _ref.repoName, branchToShow = _ref.branchToShow, allBranches = _ref.allBranches;
    if (!ajaxDone) {
      return null;
    }
    content = branchToShow ? React.createElement(BranchPage, {
					branch: branchToShow, 
					apiHost: apiHost, 
					mainPort: mainPort, 
					urlForBranch: this.urlForBranch}
				)
				 : React.createElement("div", {className: "branch-manager-container"}, 
					React.createElement(BranchList, {
						indexHTMLName: indexHTMLName, 
						apiHost: apiHost, 
						mainPort: mainPort, 
						UIPort: UIPort, 
						repoName: repoName, 
						urlForBranch: this.urlForBranch, 
						getBranchUrl: this.getBranchUrl, 
						alteredBranch: this.alteredBranch, 
						allBranches: allBranches}
					)
				)
				;
    help = !branchToShow ? React.createElement(Help, {repoName: repoName}) : null;
    return React.createElement("div", {className: "app-container"}, 
			React.createElement(Chrome, null), 
			help, 
			content
		)
		;
  },
  componentDidMount: function() {
    var def;
    def = this.getConfig();
    return def.then(this.getAllBranches).then((function(_this) {
      return function() {
        _this.setBranchParam();
        return _this.setState({
          ajaxDone: true
        });
      };
    })(this));
  },
  getConfig: function(cb) {
    var d;
    d = $.ajax({
      type: "GET",
      url: "./config.json",
      success: (function(_this) {
        return function(results) {
          var UIPort, apiHost, indexHTMLName, mainPort, repoName;
          if (!(results != null ? results.indexHTMLName : void 0)) {
            return;
          }
          UIPort = results.UIPort, indexHTMLName = results.indexHTMLName, apiHost = results.apiHost, mainPort = results.mainPort, repoName = results.repoName;
          return _this.setState({
            UIPort: UIPort,
            indexHTMLName: indexHTMLName,
            apiHost: apiHost,
            mainPort: mainPort,
            repoName: repoName
          }, cb);
        };
      })(this),
      error: function(results) {
        return console.warn("error loading config: ", results);
      }
    });
    return d;
  },
  getAllBranches: function() {
    return $.ajax({
      url: this.getBranchUrl(),
      type: "GET",
      data: {
        whichBranches: "all"
      },
      success: (function(_this) {
        return function(msg) {
          return _this.setState({
            allBranches: (msg != null ? msg.data : void 0) || []
          });
        };
      })(this),
      error: function(msg) {
        return console.log("fail: ", msg);
      }
    });
  },
  alteredBranch: function() {
    return this.getAllBranches();
  },
  urlForBranch: function(branch) {
    return "http://" + this.state.apiHost + ":" + branch.Port + "/" + this.state.indexHTMLName;
  },
  setBranchParam: function() {
    var branchName, key, pair, value, _i, _len, _ref, _ref1;
    _ref = location.search.substring(1).split("&");
    for (_i = 0, _len = _ref.length; _i < _len; _i++) {
      pair = _ref[_i];
      _ref1 = pair.split("="), key = _ref1[0], value = _ref1[1];
      if (key === "branch" && value !== "") {
        branchName = decodeURIComponent(value);
      } else {
        continue;
      }
    }
    return this.setState({
      branchToShow: this.state.allBranches[branchName]
    });
  },
  getBranchUrl: function() {
    return "http://" + this.state.apiHost + ":" + this.state.mainPort + "/branch";
  }
});


},{"./branchList.coffee":2,"./branchPage.coffee":3,"./chrome.coffee":4,"./help.coffee":5}],2:[function(require,module,exports){
/** @jsx React.DOM */
/*
Component for managing which branches are active
	-- add branch
	-- delete branch
	-- search brancehs
NOTE: This probably does too many things at the moment
 */
var BranchList;

module.exports = BranchList = React.createClass({displayName: "BranchList",
  refForInput: "branchInputField",
  BUILDING_CLASS: "building-branch",
  WAITING_CLASS: "waiting-branch",
  getInitialState: function() {
    return {
      allBranches: this.props.allBranches || {},
      searchText: "",
      postedToBuild: {},
      postedToDelete: {}
    };
  },
  componentWillReceiveProps: function(newProps) {
    return this.setState({
      allBranches: newProps.allBranches || []
    }, this.updatePostedBranches);
  },
  render: function() {
    var branchInfo, branchName, builtJSX, inputPlaceHolder, jsx, status, unbuiltJSX, _ref;
    if (_.isEmpty(this.state.allBranches)) {
      return React.createElement("div", {className: "main-page-spinner"}, React.createElement("i", {className: "fa fa-spinner fa-spin fa-5x"}));
    }
    builtJSX = [];
    unbuiltJSX = [];
    _ref = this.state.allBranches;
    for (branchName in _ref) {
      branchInfo = _ref[branchName];
      if (!branchInfo.isInSearch) {
        continue;
      }
      status = this.state.postedToBuild[branchName] ? "Building" : this.state.postedToDelete[branchName] ? "Unbuilt" : branchInfo.Status;
      jsx = this.formatBranch(branchName, branchInfo, status);
      if (status === "Built") {
        builtJSX.push(jsx);
      } else {
        unbuiltJSX.push(jsx);
      }
    }
    inputPlaceHolder = "Search remote branches of " + this.props.repoName + "...";
    return React.createElement("div", {className: "branch-manager"}, 
			React.createElement("input", {
				type: "text", 
				className: "enter-branch", 
				onChange: this.onChangeSearch, 
				placeholder: inputPlaceHolder, 
				ref: this.refForInput}
			), 
			React.createElement("ul", {className: "branch-list"}, 
				React.createElement("div", {className: "list-section-header"}, "BUILT BRANCHES (", builtJSX.length, ")"), 
				builtJSX, 
				React.createElement("div", {className: "list-section-header"}, "UNBUILT BRANCHES (", unbuiltJSX.length, ")"), 
				unbuiltJSX
			)
		)
		;
  },
  formatBranch: function(branchName, branchInfo, status) {
    var className, directUrl;
    switch (status) {
      case "Built":
        directUrl = this.props.urlForBranch(branchInfo);
        return React.createElement("li", {key: branchName, className: "built-branch"}, 
					React.createElement("span", {className: "build-button-container"}, 
						React.createElement("a", {
							className: "delete-branch", 
							onClick: this.deleteBranch.bind(this, branchInfo.Name)
						}, 
							"remove"
						)
					), 
					React.createElement("a", {className: "branch-name", href: this.urlForBranchPage(branchInfo)}, branchInfo.Name), 
					React.createElement("a", {
						className: "direct-url", 
						href: directUrl
						}, 
							"direct", 
							React.createElement("i", {className: "fa fa-external-link"})
					)

				);
      case "Building":
      case "Waiting":
        className = status === "Building" ? this.BUILDING_CLASS : this.WAITING_CLASS;
        return React.createElement("li", {key: branchName, className: className}, 
					this.getBuildingLabelJSX(status), 
					React.createElement("span", {className: "branch-name"}, branchInfo.Name)
				);
      case "Unbuilt":
        className = "not-built-branch";
        return React.createElement("li", {key: branchName, className: className}, 
					this.getUnbuiltLabelJSX(branchName), 
					React.createElement("span", {className: "branch-name"}, branchName)
				);
      default:
        return null;
    }
  },
  componentDidMount: function() {
    return this.setPrunedBranches();
  },
  getBuildingLabelJSX: function(status) {
    return React.createElement("span", {className: "build-status"}, 
			React.createElement("i", {className: "fa fa-spinner fa-spin status-spinner"}), 
			React.createElement("span", {className: "build-button-container"}, 
				React.createElement("span", null, status)
			)
		)
		;
  },
  getUnbuiltLabelJSX: function(branchName) {
    return React.createElement("span", null, 
			React.createElement("span", {className: "build-button-container"}, 
				React.createElement("a", {
					className: "build-branch", 
					onClick: this.postNewBranch.bind(this, branchName)
				}, 
					"create"
				)
			)
		)
		;
  },
  focusSearch: function() {
    return this.refs[this.refForInput].getDOMNode().focus();
  },
  onChangeSearch: function(e) {
    return this.setState({
      searchText: e.target.value.trim().toLowerCase()
    }, this.setPrunedBranches);
  },
  setPrunedBranches: function() {
    var allBranches, branchName, info;
    allBranches = this.state.allBranches;
    for (branchName in allBranches) {
      info = allBranches[branchName];
      info.isInSearch = branchName.toLowerCase().indexOf(this.state.searchText) !== -1;
    }
    return this.setState({
      allBranches: allBranches
    });
  },
  postNewBranch: function(branchName) {
    var postedToBuild;
    postedToBuild = this.state.postedToBuild;
    postedToBuild[branchName] = branchName;
    this.setState({
      postedToBuild: postedToBuild
    });
    return $.ajax({
      type: "POST",
      url: this.props.getBranchUrl(),
      data: {
        name: branchName
      },
      success: (function(_this) {
        return function(results) {
          return _this.props.alteredBranch();
        };
      })(this),
      error: function(results) {
        return console.log('post error with results: ', results);
      }
    });
  },
  deleteBranch: function(branchName) {
    var postedToDelete;
    postedToDelete = this.state.postedToDelete;
    postedToDelete[branchName] = branchName;
    this.setState({
      postedToDelete: postedToDelete
    });
    return $.ajax({
      type: "DELETE",
      url: "" + (this.props.getBranchUrl()) + "?name=" + branchName,
      success: (function(_this) {
        return function(results) {
          return _this.props.alteredBranch();
        };
      })(this),
      error: function(results) {
        return console.log('delete error with results: ', results);
      }
    });
  },
  updatePostedBranches: function() {
    var deleted, match, posted, postedToBuild, postedToDelete, _ref;
    _ref = this.state, postedToBuild = _ref.postedToBuild, postedToDelete = _ref.postedToDelete;
    for (posted in postedToBuild) {
      match = _.filter(this.state.allBranches, function(branch) {
        return branch.Name === posted && branch.Status === "Built";
      });
      if (match.length !== 0) {
        delete postedToBuild[posted];
      }
    }
    for (deleted in postedToDelete) {
      match = _.filter(this.state.allBranches, function(branch) {
        return branch.Name === deleted && branch.Status === "Unbuilt";
      });
      if (match.length === 0) {
        delete postedToDelete[deleted];
      }
    }
    return this.setState({
      postedToDelete: postedToDelete,
      postedToBuild: postedToBuild
    }, this.setPrunedBranches);
  },
  urlForBranchPage: function(branch) {
    var encodedBranch;
    encodedBranch = encodeURIComponent(branch.Name);
    return "http://" + this.props.apiHost + ":" + this.props.UIPort + "/?branch=" + encodedBranch;
  }
});


},{}],3:[function(require,module,exports){
/** @jsx React.DOM */var BranchPage;

module.exports = BranchPage = React.createClass({displayName: "BranchPage",
  render: function() {
    var Name, Status, branch, url, urlForBranch, _ref;
    _ref = this.props, branch = _ref.branch, urlForBranch = _ref.urlForBranch;
    Name = branch.Name, Status = branch.Status;
    url = urlForBranch(branch);
    return React.createElement("div", {className: "branch-page"}, 
			React.createElement("div", {className: "iframe-container"}, 
				React.createElement("a", {
					className: "branch-url", 
					href: url
				}, 
					"Show \"", Name, "\" in a separate window"
				), 
				React.createElement("iframe", {ref: "iframe", src: url})
			)
		)
		;
  },
  componentDidMount: function() {
    $(window).on("resize", _.debounce(this.onResize, 100));
    return this.onResize();
  },
  onResize: function() {
    var $frame;
    if ($frame = $(this.refs["iframe"].getDOMNode())) {
      $frame.css("height", $(window).height() * .9);
      return $frame.css("width", $(window).width() * .98);
    }
  }
});


},{}],4:[function(require,module,exports){
/** @jsx React.DOM */var Chrome;

module.exports = Chrome = React.createClass({displayName: "Chrome",
  render: function() {
    return React.createElement("div", {className: "chrome"}, 
			React.createElement("span", {className: "logo"}, React.createElement("a", {href: "/"}, "Iteraid")), 
			React.createElement("span", {className: "catch-phrase"}, "The aid to product iteration"), 
			React.createElement("div", {className: "git-link"}, 
				React.createElement("a", {href: "https://github.com/scottyantipa/Iteraid"}, 
					React.createElement("i", {className: "fa fa-github"})
				)
			)
		)
		;
  }
});


},{}],5:[function(require,module,exports){
/** @jsx React.DOM */var Help;

module.exports = Help = React.createClass({displayName: "Help",
  QUICK_DESCRIPTION: "Iteraid helps you iterate on your web software.   \nIt lets you build all of the versions your team is working on, giving you a url that points to that built version.\nClick on 'build' next to a branch to make it live.  When you no longer need the live version, just click 'delete'.",
  getInitialState: function() {
    return {
      visible: false
    };
  },
  render: function() {
    var label, messageJSX, visible;
    visible = this.state.visible;
    label = visible ? "Hide Help" : "Show Help";
    messageJSX = visible ? React.createElement("div", {className: "help-message"}, 
					React.createElement("div", {className: "quick-description"}, this.QUICK_DESCRIPTION)
				)
				 : null;
    return React.createElement("div", {className: "help-user"}, 
			messageJSX, 
			React.createElement("div", {className: "help-toggle", onClick: this.toggleVisible}, 
				label
			)
		)
		;
  },
  toggleVisible: function() {
    return this.setState({
      visible: !this.state.visible
    });
  }
});


},{}],6:[function(require,module,exports){
/** @jsx React.DOM */var App;

App = require("./app.coffee");

$(function() {
  return React.renderComponent(App({}), $("body")[0]);
});


},{"./app.coffee":1}]},{},[6]);
