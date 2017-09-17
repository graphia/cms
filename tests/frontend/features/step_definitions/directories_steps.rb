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