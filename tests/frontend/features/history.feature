Feature: Commits
	So I can review recent commits
	As an author
	I want to be able to see a list of commits with messages

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Commits are visible
		Given I have added a new file
		When I navigate to the history page
		Then I should see my commit

	Scenario: Commits messages are displayed
		Given I have added a new file with a multiline commit message
		When I navigate to the history page
		Then the first line of my commit message should be a header
		And all subsequent lines should be paragraphs

	Scenario: Commits author names are displayed
		Given I have added a new file
		When I navigate to the history page
		Then the commit should include the committer name

	Scenario: Breadcrumbs
		When I am on the history page
		Then I should see the following breadcrumbs:
			| Text      | Reference   |
			| Dashboard | /cms        |
			| History   | None        |