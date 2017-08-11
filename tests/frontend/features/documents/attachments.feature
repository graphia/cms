Feature: Creating documents
	So I can create beautiful, interesting documents
	As an author
	I want to add graphics and images

	Background:
		Given my user account exists
		And I have logged in
		And I am on the new document page

	Scenario: Dragging an image into a new document
		Given I have a 'jpeg' image
		When I drag the image into the editor
		Then the image placeholder should be added to the content
		And the image should be displayed in the gallery