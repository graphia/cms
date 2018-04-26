Feature: Learning how to clone the theme
	So I can customise the CMS's output
	As a designer
	I want to be able to clone the theme repository

	Background:
		Given a repository has been initialised
		And the CMS is running with the "ssh_enabled" config
		And my user account exists
		And I have logged in

	Scenario: Page title
		Given I visit the theme settings page
		Then the page's title should be "Settings: Theme"
		And the main heading should be "Customising the theme"

	Scenario: Page contents
		Given I visit the theme settings page
		Then the page should contain the following dynamic snippets:
			| HostName docs.somecompany.com        |
			| Port 2223                            |
			| git clone docs.somecompany.com:theme |
		And the following subheadings:
			| Prerequisites |
			| Cloning       |