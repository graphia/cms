Feature: Site Title

	So I feel as though the CMS is part of my company infrastructure
	As a user
	I want to see my company name in the header

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: The title should match the configured site title
		Given my company name is "Krusty Burger"
		Then the header in the main navigation should contain "Krusty Burger"
		And the title should be a link to the homepage