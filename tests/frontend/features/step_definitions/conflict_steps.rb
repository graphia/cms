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

Given %r{^my downloads directory is empty$} do
  FileUtils.rm_rf(Dir.glob("#{DOWNLOAD_DIR}/*"))
  expect(Dir.entries(DOWNLOAD_DIR)).to eql(%w{. ..})
end

Then %r{^I should see the conflict modal box$} do
  expect(page).to have_css("#conflict-warning.modal")
end

Given %r{^I have tried to save a file but the repository has been updated$} do
  steps %{
    Given I am on the new document page
    And I enter some text into the editor
    And I fill in the document metadata
    And I have added a new file
    And I submit the form
  }
end

When %r{^I click the "([^"]*)" button in the modal$} do |button_text|
  within("#conflict-warning.modal") do
    page.find(".btn", text: button_text).click
  end
end

Then %r{^I should download a copy of my file "(.*?)$} do |filename|
  # Wait for upto a second for this to complete
  1.upto(10) do
    begin
      expect(Dir.entries(DOWNLOAD_DIR)).to eql(%w{. .. sample-document.md})
    rescue RSpec::Expectations::ExpectationNotMetError
      sleep(0.1)
    end
  end
end

Then %r{^I should be on the documents directory's index page$} do
  expect(page.current_path).to eql("/cms/documents")
end