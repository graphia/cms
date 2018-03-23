Then %r{^I should be redirected to the parent directory's index$} do
  expect(page).to have_css("h1", text: "Appendices")
  expect(page.current_path).to eql("/cms/appendices")
end

Then %r{^the file should have been deleted$} do
  expect(File.exist?(File.join(REPO_PATH, "appendices", "appendix_1.md"))).to be false
end

Then %r{^the file and attachments directory should have been deleted$} do
  expect(File.exist?(File.join(REPO_PATH, "appendices", "appendix_1", "index.md"))).to be false
  expect(Dir.exist?(File.join(REPO_PATH, "appendices", "appendix_1"))).to be false
end

Then %r{^the file should have been deleted but not the attachments directory$} do
  expect(File.exist?(File.join(REPO_PATH, "appendices", "appendix_1", "index.md"))).to be false
  expect(Dir.exist?(File.join(REPO_PATH, "appendices", "appendix_1"))).to be true
end

Given %r{^I have deleted a single file$} do
  steps %{
    Given I am on the document's show page
    When I click the "Delete" button
    And I click the "Confirm deletion" button
    Then I should be redirected to the parent directory's index
  }
end

Then %r{^the last commit message should contain the file's name$} do
  g = Git.open(REPO_PATH)
  expect(g.log.first.message.to_s).to eql("File deleted appendices/appendix_1/index.md")
end

Then %r{^I should see the deletion modal box$} do
  expect(page).to have_css("#delete-warning.modal")
end

Given %r{^I can see the document's deletion modal$} do
  steps %{
    Given I am on the document's show page
    When I click the "Delete" button
    Then I should see the deletion modal box
  }
end

When %r{^I (?:try|attempt) to delete the file(?: again)?$} do
  click_button "Delete"
  click_button "Confirm deletion"
end