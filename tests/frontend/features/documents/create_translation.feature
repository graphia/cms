Feature: Creating documents
	So I can add content to the CMS
	As an author
	I want to be able to create new documents

	Background:
		Given a repository has been initialised
		And the CMS is running with the "multilingual" config
		And my user account exists
		And I have logged in

	Scenario: Creating a new file with a non-default language code
		Given I am on the new document page
		When I enter some text into the editor
		And I fill in the document metadata
		And I select "Finnish" from the languages dropdown
		And I submit the form
		Then the "Finnish" file should be created and contain the correct information

	Scenario: Creating a new file with a non-default language code and custom document name
		Given I am on the new document page
		And I have filled in details for a new "Finnish" document
		And I have entered the custom document "Top 10 Moomins"
		When I submit the form
		Then the "Finnish" custom file should be created and contain the correct information