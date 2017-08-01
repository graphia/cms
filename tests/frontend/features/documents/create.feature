Feature: Creating documents
	So I can add content to the CMS
	As an author
	I want to be able to create new documents

	Background:
		Given my user account exists
		And I have logged in

	Scenario: The editor
		Given I am on the new document page
		Then I should see an editor with the following buttons:
			| Bold           |
			| Italic         |
			| Heading        |
			| Quote          |
			| Generic List   |
			| Numbered List  |
			| Create Link    |
			| Insert Image   |
			| Toggle Preview |
			| Markdown Guide |
		And I should see the following fields for document metadata:
			| Title    |
			| Filename |
			| Synopsis |
			| Author   |
			| Tags     |
		And I should see a text area for the commit message

	Scenario: Creating a file
		Given I am on the new document page
		When I enter some text into the editor
		And I fill in the document metadata
		And I submit the form
		Then I should see my correctly-formatted document

	Scenario: Automatically setting the new file name
		Given I am on the new document page
		When I set the "title" to "the world's most amazing, fantastic file"
		Then the "filename" should equal "the-worlds-most-amazing-fantastic-file"
		And the "filename" field should be read only

	Scenario: Customising the filename
		Given I am on the new document page
		And I have entered my new document's details
		When I check the "custom-filename" checkbox
		And the "filename" field should not be read only

	Scenario: Redirection to new document after creation
		Given I am on the new document page
		When I set the "title" to "sample document 2"
		And I have edited the document and commit message
		And I submit the form
		Then I should be redirected to "/cms/documents/sample-document-2.md"

	Scenario: Cancelling document creation
		Given I am on the new document page
		When I click the "Cancel" button
		Then I should be redirected to the documents index

	Scenario: Default page heading
		Given I am on the new document page
		When the "title" is blank
		Then the page heading should be "New Document"

	Scenario: Updating the page heading
		Given I am on the new document page
		When I set the "title" to "sample document"
		Then the page heading should be "sample document"