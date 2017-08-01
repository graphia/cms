Feature: Creating documents
	So I can add content to the CMS
	As an author
	I want to be able to create new documents

	Background:
		Given my user account exists
		And I have logged in

	Scenario: The editor
		Given I am on the edit document page for "document_1.md"
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
			| Synopsis |
			| Author   |
			| Tags     |
		And I should not see the "Filename" field
		And I should see a text area for the commit message

	Scenario: Updating a file
		Given I am on the edit document page for "document_1.md"
		When I amend the text in the editor
		And I submit the form
		Then I should see my updated document

	Scenario: Redirection to modified document after post update
		Given I am on the edit document page for "document_1.md"
		When I set the "title" to "updated document"
		And I have edited the document and commit message
		And I submit the form
		Then I should be redirected to "/cms/documents/document_1.md"

	Scenario: Default page heading
		Given I am on the edit document page for "document_1.md"
		When the "title" is "document 1"
		Then the page heading should be "document 1"

	Scenario: Updating the page heading
		Given I am on the edit document page for "document_1.md"
		When I clear the "title"
		Then the page heading should be "No title"