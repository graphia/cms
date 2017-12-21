Feature: Displaying documents
	So I can view CMS contents
	As an author
	I want to be able to read the fully formatted documents

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And a document called 'appendix_1.md' exists
		And I have logged in

	Scenario: Viewing a document
		Given I navigate to that document's 'show' page
		Then I should see the correctly-formatted document

	Scenario: Page title is properly set
		Given I am on the document's show page
		And the page's title should be "Appendix 1"

	Scenario: Frontmatter
		Given I am on the document's show page
		When the document has some frontmatter set up
		Then I should see the following frontmatter items with correct values:
			| Title       |
			| Author      |
			| Tags        |
			| Date        |
			| Synopsis    |
			| Version     |

	Scenario: Toolbar
		Given I am on the document's show page
		Then I should see a toolbar with the following buttons:
			| Edit    |
			| History |
			| Delete  |

	Scenario: Clicking the Edit button
		Given I am on the document's show page
		When I click the toolbar's 'Edit' button
		Then I should be on the document's edit page

	Scenario: Clicking the Delete button
		Given I am on the document's show page
		When I click the toolbar's 'Delete' button
		And I click the "Confirm deletion" button
		Then I should be on the directory's index page
		And the document should have been deleted

	Scenario: Clicking the History button
		Given I am on the document's show page
		When I click the toolbar's 'History' button
		Then I should be on the document's history page

	Scenario: Breadcrumbs without metadata
		Given I am on the document's show page
		Then I should see the following breadcrumbs:
			| Text                | Reference           |
			| Dashboard           | /cms                |
			| appendices          | /cms/appendices     |
			| appendix_1.md       | None                |

	Scenario: Breadcrumbs with metadata
		Given I am on the show page for a document with metadata
		Then I should see the following breadcrumbs:
			| Text                | Reference           |
			| Dashboard           | /cms                |
			| Important Documents | /cms/documents      |
			| document 1          | None                |