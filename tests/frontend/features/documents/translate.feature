Feature: Translating documents
	So I can make CMS contents more accessible
	As an author
	I want to be able to facilitate their translation

	Background:
		Given a repository has been initialised
		And the CMS is running with the "multilingual" config
		And my user account exists
		And I have logged in

	Scenario: Translation dropdown
		Given my document is untranslated
		And I navigate to my document's 'show' page
		When I click the 'Translate' dropdown button
		Then I should see a list of available languages:
			| Swedish |
			| Finnish |
		And the existing language 'English' should not be included

	Scenario: Translation dropdown when translations exist
		Given there is already a 'Swedish' translation of my document
		And I navigate to my document's 'show' page
		When I click the 'Translate' dropdown button
		Then I should see a list of available languages:
			| Finnish |
		And the existing languages 'English' and 'Swedish' should not be included


	Scenario: Initiating a translation
		Given I navigate to my document's 'show' page
		When I click the 'Translate' dropdown button
		And I click the 'Swedish' translation option
		Then I should be on the new 'Swedish' document
		And there should be an alert with the message "This is the placeholder for your translation"