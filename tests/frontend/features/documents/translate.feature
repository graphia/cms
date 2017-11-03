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
		Given I navigate to my document's 'show' page
		When I click the 'Translate' dropdown button
		Then I should see a list of available languages:
			| Swedish |
			| Finnish |
		And the default language 'English' should not be included

	Scenario: Translation dropdown when translations exist
		Given there is already a 'Swedish' translation of my document
		And I navigate to that document's 'show' page
		When I click the 'Translate' dropdown button
		Then I should see a list of available languages:
			| Swedish |
			| Finnish |
		And 'Swedish' should be disabled

	Scenario: Initiating a translation
		Given I navigate to my document's 'show' page
		When I click the 'Translate' dropdown button
		And I click 'Swedish'
		Then I should be on the new 'Swedish' document
		And there should be a banner "This is the placeholder for your translation"