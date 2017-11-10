Feature: Translations

	So that I can see an overview of multilingual content
	As a reader
	I want to be able to view translated files

	Background:
		Given a repository has been initialised
		And the CMS is running with the "multilingual" config
		And my user account exists
		And I have logged in

	Scenario: Translation links visible on the documents index page
		Given my document has been translated into 'Finnish' and 'Swedish'
		When I visit the homepage
		Then I should see my document listed under 'Documents'
		And it should have 'Finnish' and 'Swedish' flags

	Scenario: Navigating to a translated page
		Given my document has been translated into 'Swedish'
		And I visit the homepage
		When I click the 'Swedish' link
		Then I should be on the document's 'Swedish' translation