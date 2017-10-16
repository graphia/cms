Given %r{^a repository update has taken place in the background$} do
  step "I have added a new file"
end

When %r{^I make my changes and submit the form$} do
  steps %{
    When I set the "title" to "updated document"
		And I have edited the document and commit message
    And I submit the form
  }
end

When %r{^I add my document's details and submit the form$} do
  steps %{
    When I enter some text into the editor
		And I fill in the document metadata
		And I submit the form
  }
end

Then %r{^I should see the conflict modal box$} do
  expect(page).to have_css("#conflict-warning.modal")
end