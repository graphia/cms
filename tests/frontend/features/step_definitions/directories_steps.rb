Given %r{^I am on the new directory page$} do
  path = "/cms/new"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should be redirected to the new directroy's index$} do
  expect(page).to have_content("Ice Creams")
  expect(page.current_path).to eql("/cms/ice-creams")
end

Then %r{^the directory should have been created with the correct information$} do
  file = File.read(File.join(REPO_PATH, "ice-creams", "_index.md"))
  expect(file).to include("title: Ice Creams")
  expect(file).to include("description: A description of ice cream related products")
  expect(file).to include("# Fabulous ices of all colours")
end

When %r{^I have filled in the required fields$} do
  fill_in "title", with: "Marge Gets a Job"
end

Then %r{^the "(.*?)" field's error text should contain "(.*?)"$} do |field, error|
  within(".directory-#{field.downcase}") do
    expect(page).to have_content(error)
  end
end

Then %r{^the "(.*?)" directory should have been deleted$} do |dir_name|
  within(".directories") do
    expect(page).not_to have_css("div[data-directory='#{dir_name}']")
  end
end

Then %r{^I should be able to see the directory's description$} do
  expect(page).to have_css(".directory-description", text: "Documents go here")
end

Then %r{^I should be able to see the directory's introduction$} do
  within(".directory-info-text") do
    expect(page).to have_content("These documents are amazing")
  end
end

Then %r{^there should be no "(.*?)" section$} do |css_class|
  expect(page).not_to have_css(css_class)
end

Given %r{^I am on the update directory page$} do
  path = "/cms/documents/edit"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^the directory index page should contain the newly\-updated information$} do
  expect(page).to have_css(".directory-description", text: "A description of ice cream related products")

  within(".directory-info-text") do
    expect(page).to have_css("h1", text: "Fabulous ices of all colours")
  end

end