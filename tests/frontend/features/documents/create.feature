Feature: Creating documents
	So I can add content to the CMS
	As an author
	I want to be able to create new documents

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page title is dynamic
		Given I am on the new document page
		And the page's title should be "New Document"
		When I set the "title" to "Boo-urns"
		Then now the title is "Boo-urns"

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
		And I should see a tags editing field
		And I should see a text area for the commit message

	Scenario: Default date
		Given I am on the new document page
		Then the date should be set to today

	Scenario: Creating a file
		Given I am on the new document page
		When I enter some text into the editor
		And I fill in the document metadata
		And I submit the form
		Then I should see my correctly-formatted document

	Scenario: Creating a file with tags
		Given I am on the new document page
		When I add content, tags and a commit message
		And I submit the form
		Then I should see my document with the correct tags

	Scenario: Redirection to new document after creation
		Given I am on the new document page
		When I have created a new document titled "sample document 2"
		Then I should be redirected to "/cms/documents/sample-document-2"

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

	Scenario: Submit button disabled by default
		Given I am on the new document page
		When I haven't interacted with the form
		Then the submit button should be disabled

	Scenario: Correctly dealing with conflicts
		Given I am on the new document page
		And a repository update has taken place in the background
		When I add my document's details and submit the form
		Then I should see the conflict modal box

	Scenario: Breadcrumbs without metadata
		Given I am on the new appendix page
		Then I should see the following breadcrumbs:
			| Text                | Reference       |
			| Dashboard           | /cms            |
			| appendices          | /cms/appendices |
			| New Document        | None            |

	Scenario: Breadcrumbs with metadata
		Given I am on the new document page
		Then I should see the following breadcrumbs:
			| Text                | Reference       |
			| Dashboard           | /cms            |
			| Important Documents | /cms/documents  |
			| New Document        | None            |