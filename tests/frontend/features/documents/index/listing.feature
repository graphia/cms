Feature: Listing documents
	So I can see files of a certain category
	As an author
	I want to view listings for a directory

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Page title is properly set
		Given I am on the "documents" index page
		Then the page's title should be "Important Documents"

	Scenario: Documents are visible on the documents page
		Given there are directories called "documents" and "appendices"
		And they both contain Markdown files
		When I navigate to the "documents" index page
		Then I should see a list containing the contents of the "documents" directory

	Scenario: Navigating to a document's show page
		Given I am on the "documents" index page
		When I click the link to "document 1"
		Then I should be on the "document 1" show page

	Scenario: Creating a new document
		Given I am on the "documents" index page
		When I click the "New document" button
		Then I should be on the new document page for the 'documents' directory

	Scenario: Documents are visible on the documents page
		Given I am on the "documents" index page
		When I click the "Appendices" navigation link
		Then I should be on the "appendices" index page
		And I should see a list containing the contents of the "appendices" directory

	Scenario: The page title should match the directory name
		Given there are directories called "documents" and "appendices"
		Then each directory index page should have the correct title:
			| Directory  | Title      |
			| documents  | Documents  |
			| appendices | Appendices |

	Scenario: When the directory is empty
		Given the 'empty' directory contains no files
		When I am on the "empty" index page
		Then there should be an alert with the message "There's nothing here yet"
		And there should be a 'Create a new document' link

	Scenario: The create a new document button
		Given the 'empty' directory contains no files
		And I am on the "empty" index page
		When I click the "Create a new document" button
		Then I should be on the new document page for the 'empty' directory

	Scenario: When the directory doesn't exist
		Given there is no directory called "operating-procedures"
		When I am on the "operating-procedures" index page
		Then the main heading should be "404"
		And the page's title should be "Not found"
		And there should be an alert with the message "The item you were looking for cannot be found"

	Scenario: Breadcrumbs
		Given I am on the "documents" index page
		Then I should see the following breadcrumbs:
			| Text                | Reference |
			| Dashboard           | /cms      |
			| Important Documents | None      |

	Scenario: Primary navigation highlighting
		When I am on the "appendices" index page
		Then the primary navigation link "Appendices" should be active

	Scenario: Identifying draft documents
		Given I have some documents that are drafts
		When I am on the "appendices" index page
		Then the draft document should be highlighted