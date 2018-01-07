Feature: Publishing drafts
	So I can work on documents without publishing them
	As an author
	I want to be able to set a document's draft status

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Draft checkbox on new document page
		Given I am on the new document page
		Then there should be a checkbox called "Draft"

	Scenario: Draft checkbox on the edit document page
		Given I am on the edit document page for "document_1.md"
		Then there should be a checkbox called "Draft"

	Scenario: New documents should default to draft
		Given I am on the new document page
		Then the "Draft" checkbox should be checked

	Scenario: Creating a non-draft document
		Given I am on the new document page
		When I uncheck the "Draft" checkbox
		And I fill in the rest of the document form and submit it
		Then my document should not be a draft

	Scenario: Marking an existing document as a draft
		Given I am on the edit document page for "document_1.md"
		When I uncheck the "Draft" checkbox
		And I enter 'the first document' in the 'Title' field
		And I enter 'commit-message' in the 'Commit Message' field
		And I submit the form
		Then my document should be a draft