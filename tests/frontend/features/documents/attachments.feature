Feature: Creating documents
	So I can create beautiful, interesting documents
	As an author
	I want to add graphics and images

	Background:
		Given a repository has been initialised
		And the CMS is running with the "default" config
		And my user account exists
		And I have logged in

	Scenario: Images visible in the gallery
		Given the document I'm working on already has an attachment
		When I am on the edit page for my document
		And I click the "Images" tab in the sidebar
		Then I should see the image in the gallery

	@wip
	Scenario: Dragging an existing image into a document
		Given the document I'm working on already has an attachment
		And I am on the edit page for my document
		When I drag an image from the gallery to the editor
		Then the markdown image placeholder should be added to the editor