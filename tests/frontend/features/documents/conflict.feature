Feature: Creating documents
	So I can be confident in not losing my work
	As an author
	I want to be able to recover from conflicts

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	@chrome
	Scenario: Saving my work
		Given my downloads directory is empty
		And I have tried to save a file but the repository has been updated
		When I click the "Download your copy" button in the modal
		Then I should download a copy of my file "sample-document.md"

	Scenario: Ignoring the warning
		Given I have tried to save a file but the repository has been updated
		When I click the "Close" button in the modal
		Then I should be on the documents directory's index page