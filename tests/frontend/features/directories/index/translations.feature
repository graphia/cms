Feature: Translations

	So that I can interact with multilingual content
	As a reader
	I want to be able to view translated files

	Background:
		Given a repository has been initialised
		And the CMS is running with the "multilingual" config
		And my user account exists
		And I have logged in

	Scenario: Translation links visible on the documents index page
		Given my document has been translated into 'Finnish' and 'Swedish'
		When I visit the documents index page
		Then I should see my document listed with 'Finnish' and 'Swedish' flags

	Scenario: When there is only one version of the file available
		Given my document has not been translated into any other languages
		When I visit the documents index page
		Then my document shouldn't have any alternative languages section

	Scenario: Navigating to the translation's show page
		Given my document has been translated into 'Finnish'
		And I visit the documents index page
		When I click the 'Finnish' link
		Then I should be on the document's 'Finnish' translation
