Given %r{^I am on the page for non-existant directory "(.*?)"$} do |directory_name|
  path = "/cms/#{directory_name}"
  visit(path)
  expect(page.current_path).to eql(path)
end