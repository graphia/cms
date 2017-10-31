Feature: Creating documents
	So I can add content to the CMS
	As an author
	I want to be able to create new documents

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
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
		And I should see a tags editing field
		And I should not see the "Filename" field
		And I should see a text area for the commit message

	Scenario: Updating a file
		Given I am on the edit document page for "document_1.md"
		When I amend the text in the editor, modify the metadata and add a commit message
		And I submit the form
		Then I should see my updated document

	Scenario: Updating a file with tags
		Given I am on the edit document page for "document_1.md"
		When I add tags for Sales and Marketing
		And I enter "added some tags" in the "Commit Message" field
		And I submit the form
		Then I should see my document with the correct tags

	Scenario: Redirection to modified document after post update
		Given I am on the edit document page for "document_1.md"
		When I set the "title" to "updated document"
		And I have edited the document and commit message
		And I submit the form
		Then I should see the document containing my recent changes
		And I should have been redirected to "/cms/documents/document_1.md"

	Scenario: Cancelling an edit
		Given I am on the edit document page for "document_1.md"
		When I click the "Cancel" button
		Then I should be redirected to "/cms/documents/document_1.md"

	Scenario: Default page heading
		Given I am on the edit document page for "document_1.md"
		When the "title" is "document 1"
		Then the page heading should be "document 1"

	Scenario: Updating the page heading
		Given I am on the edit document page for "document_1.md"
		When I clear the "title"
		Then the page heading should be "No title"

	Scenario: Submit button disabled by default
		Given I am on the edit document page for "document_1.md"
		When I haven't interacted with the form
		Then the submit button should be disabled

	Scenario: Correctly dealing with conflicts
		Given I am on the edit document page for "document_1.md"
		And a repository update has taken place in the background
		When I make my changes and submit the form
		Then I should see the conflict modal box

	Scenario: Breadcrumbs without metadata
		Given I am on the edit appendix page for "appendix_1.md"
		Then I should see the following breadcrumbs:
			| Text                | Reference                     |
			| Dashboard           | /cms                          |
			| appendices          | /cms/appendices               |
			| appendix_1.md       | /cms/appendices/appendix_1.md |
			| Edit                | None                          |

	Scenario: Breadcrumbs with metadata
		Given I am on the edit document page for "document_1.md"
		Then I should see the following breadcrumbs:
			| Text                | Reference                    |
			| Dashboard           | /cms                         |
			| Important Documents | /cms/documents               |
			| document 1          | /cms/documents/document_1.md |
			| Edit                | None                         |