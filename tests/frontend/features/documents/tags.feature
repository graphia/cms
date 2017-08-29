Feature: Document tags
	So I can make documents easier to search for
	As an author
	I want to be able to add appropriate and informative tags

	Background:
		Given a repository has been initialised
		And my user account exists
		And I have logged in
		And I am on the new document page

	Scenario: Tag entry
		Given I have selected the tags editor
		When I enter the following tags and press "enter":
			| Sales     |
			| Marketing |
		Then I should see the following tags listed:
			| Sales     |
			| Marketing |