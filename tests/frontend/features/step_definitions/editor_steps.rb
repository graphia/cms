SAMPLE_MARKDOWN = "this is a *minimalistic* markdown **document**"
SAMPLE_TEXT = "this is a minimalistic markdown document"

Given %r{^I am on the new document page$} do
  path = "/cms/documents/new"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on the new appendix page$} do
  path = "/cms/appendices/new"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on the edit document page for "(.*?)"$} do |document_filename|
  path = "/cms/documents/#{document_filename}/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on the edit appendix page for "(.*?)"$} do |appendix_filename|
  path = "/cms/appendices/#{appendix_filename}/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end

Given %r{^I am on the edit document page for a document$} do
  steps %q{Given I am on the edit document page for "document_1"}
end

Then %r{^I should see an editor with the following buttons:$} do |table|
  within(".editor") do
    expect(page).to have_css(".CodeMirror")

    table.transpose.raw.flatten.each do |button|
      expect(page).to have_css("a[title^='#{button}']")
    end
  end
end

Then %r{^I should see the following fields for document metadata:$} do |table|
  within(".metadata-fields") do
    table.transpose.raw.flatten.each do |field_name|
      input_name = page.find("label", text: field_name)[:for]
      # use the universal selector because some fields are
      # inputs and others are textarea
      expect(page).to have_css("*[name='#{input_name}']")
    end
  end
end

Then %r{^I should see a text area for the commit message$} do
  within(".metadata-fields") do
    expect(page).to have_css("label", text: 'Commit Message')
    expect(page).to have_css("textarea[name='commit-message']")
  end
end

When %r{^I enter some text into the editor$} do
  page.execute_script "$('.CodeMirror')[0].CodeMirror.setValue('#{SAMPLE_MARKDOWN}')"
end

When %r{^I fill in the document metadata$} do
  fill_in 'title', with: "Sample Document"
  fill_in 'commit-message', with: "Added sample document"
  fill_in 'version', with: "1.0.2"
  steps %{
    Then I add tags for Sales and Marketing
  }
end

Then %r{^I should see my correctly-formatted document$} do
  within(".content") do
    expect(page).to have_css("p", text: SAMPLE_TEXT)
    expect(page).to have_css("em", text: "minimalistic")
    expect(page).to have_css("strong", text: "document")
  end
  within(".document-metadata") do
    expect(page).to have_content("Sample Document")
    expect(page).to have_content("1.0.2")
  end
end

When %r{^I set (?:the)? "(.*)" to "(.*)"$} do |field, value|
  fill_in field.gsub(" ", "-"), with: value
end

When %r{^I have edited the document and commit message$} do
  steps %{
    When I enter some text into the editor
    And I set the "commit message" to "general updates"
  }
end

When %r{^I add content, tags and a commit message$} do
  steps %{
    Given I set the "title" to "sales and marketing"
    And I enter some text into the editor
    And I set the "commit message" to "general updates"
    And I add tags for Sales and Marketing
  }
end

When %r{^I have created a new document titled (.*?)$} do |title|
  steps %{
    Given I set the "title" to "#{title}"
    And I enter some text into the editor
    And I set the "commit message" to "new doc"
    When I submit the form
    Then I should see the document containing my recent changes
  }
end

Then %r{^I should see the document containing my recent changes$} do
  expect(page).to have_css("p", text: SAMPLE_TEXT)
end

Then %r{^the "([^"]*)" field should equal "([^"]*)"$} do |field_name, value|
  field = page.find("input[type=text][name='#{field_name}']")
  expect(field[:value]).to eql(value)
end

Then %r{^the "([^"]*)" field should be read only$} do |field_name|
  field = page.find("input[type=text][name='#{field_name}']")
  expect(field).to(be_readonly)
end

Then %r{^the "([^"]*)" field should not be read only$} do |field_name|
  field = page.find("input[type=text][name='#{field_name}']")
  expect(field).not_to(be_readonly)
end

Given %r{^I have entered my new document's details$} do
  steps %{
    When I set the "document" to "sample document 1"
    When I set the "commit message" to "added sample doc"
    And I enter some text into the editor
  }
end

When %r{^the "([^"]*)" is blank$} do |field_name|
  field = page.find("input[name='#{field_name}']")
  expect(field.value).to be_empty
end

Then %r{^the page heading should be "([^"]*)"$} do |text|
  expect(page).to have_css("h1", text: text)
end

When %r{^the "([^"]*)" is "([^"]*)"$} do |field_name, value|
  field = page.find("input[name='#{field_name}']")
  expect(field.value).to eql(value)
end

When %r{^I clear the "([^"]*)"$} do |field_name|
  field = page.find("input[name='#{field_name}']")
  field.send_keys(Array.new("document_1".length, :backspace))
end

Then %r{^I should not see the "([^"]*)" field$} do |field_name|
  within(".metadata-fields") do
    expect(page).not_to have_css("input[name='#{field_name}']")
  end
end

When %r{^I amend the text in the editor, modify the metadata and add a commit message$} do
  fill_in 'title', with: "Edited Document"
  fill_in 'version', with: "1.2.0"

  sample = "i have **modified** the *text*"
  steps %{
    When I set the editor text to "#{sample}"
    And I set the "commit message" to "general updates"
  }
end

Then %r{^I should see my updated document$} do
  within(".content") do
    expect(page).to have_css("p", text: "i have modified the text")
    expect(page).to have_css("em", text: "text")
    expect(page).to have_css("strong", text: "modified")
  end
  within(".document-metadata") do
    expect(page).to have_content("Edited Document")
    expect(page).to have_content("1.2.0")
  end
end

When %r{^I set the editor text to "(.*?)"$} do |text|
  page.execute_script "$('.CodeMirror')[0].CodeMirror.setValue('#{text}')"
end

Then %r{^the commit message validation feedback should be visible$} do
  within("div.commit-message") do
    expect(page).to have_css(".form-control-feedback.invalid-feedback")
  end
end

Then %r{^the commit message validation feedback should not be visible$} do
  within("div.commit-message") do
    expect(page).not_to have_css(".form-control-feedback.invalid-feedback")
  end
end

Given %r{^I haven't touched the '(.*)' field$} do |field|
  # do nothing
end

Given %r{^I haven't interacted with the form$} do
  # do nothing
end

Given %r{^I enter invalid information in the form$} do
  steps %{
    Given I enter 'a' in the 'Title' field
    And I enter 'abc' in the 'Commit Message' field
  }
end

Given %r{^I enter valid information in the form$} do
  steps %{
    Given I enter 'title' in the 'Title' field
    And I enter 'commit-message' in the 'Commit Message' field
  }
end

Then %r{^the submit button should be enabled$} do
  expect(page.find("input[type='submit']")).not_to be_disabled
end

Then %r{^the submit button should be disabled$} do
  expect(page.find("input[type='submit']")).to be_disabled
end

Then %r{^I should see a tags editing field$} do
  within(".metadata-fields") do
    expect(page).to have_css(".tags-input")
  end
end

Then %r{^the title validation feedback should be visible$} do
  within("div.document-title") do
    expect(page).to have_css(".form-control-feedback.invalid-feedback")
  end
end

Then %r{^the title validation feedback should not be visible$} do
  within("div.document-title") do
    expect(page).not_to have_css(".form-control-feedback.invalid-feedback")
  end
end

Then %r{^the date should be set to today$} do
  expect(page.find("#document-date").value).to eql(Date.today.iso8601)
end

When %r{^I select "([^"]*)" from the languages dropdown$} do |arg1|
  select "Finnish", from: "language"
end

Then %r{^the "(.*?)" file should be created and contain the correct information$} do |name|
  expect(page).to have_css("h1", text: "Sample Document")

  languages = {
    "Finnish" => "fi"
  }

  expect(File.exist?(File.join(REPO_PATH, "documents", "sample-document", "index.#{languages[name]}.md"))).to be true
  expect(page.current_path).to eql("/cms/documents/sample-document/fi")
end

Given %r{^I have filled in details for a new "([^"]*)" document$} do |lang|
  steps %{
    When I enter some text into the editor
		And I fill in the document metadata
		And I select "#{lang}" from the languages dropdown
  }
end

Given %r{^I have entered the custom document "([^"]*)"$} do |name|
  steps %{
    When I check the "custom-document-identifier" checkbox
  }
  fill_in "document", with: name
end

Then %r{^the "([^"]*)" custom file should be created and contain the correct information$} do |name|
  expect(page).to have_css("h1", text: "Sample Document")

  languages = {
    "Finnish" => "fi"
  }

  expect(File.exist?(File.join(REPO_PATH, "documents", "top-10-moomins", "index.#{languages[name]}.md"))).to be true
  expect(page.current_path).to eql("/cms/documents/top-10-moomins/fi")

end