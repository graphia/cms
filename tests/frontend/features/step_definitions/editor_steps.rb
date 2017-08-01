Given %r{^I am on the new document page$} do
  path = "/cms/documents/new"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see an editor with the following buttons:$} do |table|
  within("#editor-container") do
    expect(page).to have_css(".CodeMirror")

    table.transpose.raw.flatten.each do |button|
      expect(page).to have_css("a[title^='#{button}']")
    end
  end
end

Then %r{^I should see the following fields for document metadata$} do |table|
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
  sample = "this is a *minimalistic* markdown **document**"
  page.execute_script "$('.CodeMirror')[0].CodeMirror.setValue('#{sample}')"
end

When %r{^I fill in the document metadata$} do
  fill_in 'title', with: "Sample Document"
  fill_in 'commit-message', with: "Added sample document"
end

Then %r{^I should see my correctly-formatted document$} do
  within(".content") do
    expect(page).to have_css("p", text: "this is a minimalistic markdown document")
    expect(page).to have_css("em", text: "minimalistic")
    expect(page).to have_css("strong", text: "document")
  end
end

When %r{^I set (?:the)? "(.*)" to "(.*)"$} do |field, value|
  fill_in field.gsub(" ", "-"), with: value
end

When %r{^I have edited the document and commit message$} do
  steps %{
    When I set the "commit message" to "added sample doc"
    And I enter some text into the editor
  }
end

Then %r{^I should be redirected to "(.*)"$} do |path|
  expect(page).to have_css("p", text: "this is a minimalistic markdown document")
  expect(page.current_path).to eql(path)
end

Then %r{^the "([^"]*)" should equal "([^"]*)"$} do |field_name, value|
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
    When I set the "filename" to "sample document 1"
    When I set the "commit message" to "added sample doc"
    And I enter some text into the editor
  }
end

When %r{^I check the "(.*)" checkbox$} do |checkbox|
  check(checkbox)
end