Feature: Listing documents
	So I can be informed about a directory's contents
	As an author
	I want to see the directory's metadata

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Displaying document metadata
		Given I am on the "documents" index page
		Then I should be able to see the directory's description
		And I should be able to see the directory's introduction

	Scenario: When the introductory text is missing
		Given I am on the "appendices" index page
		Then there should be no ".directory-info-text" section