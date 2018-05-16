Feature: Translations

	So I have an awareness of a document's translations
	As a reader
	I want to be able to navigate to translated versions of the document

	Background:
		Given a repository has been initialised
		And the CMS is running with the "multilingual" config
		And my user account exists
		And I have logged in

	Scenario: Translation links visible on the documents show page
		Given my document has been translated into 'Finnish' and 'Swedish'
		When I visit my document's 'English' page
		Then my document should have links to 'English', 'Finnish' and 'Swedish' in the Translations section

