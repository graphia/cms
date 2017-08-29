Given %r{^I try to navigate to the home page$} do
  visit("/cms")
end

Then %r{^I should be redirected to the initialize repository page$} do
  expect(page.current_path).to eql("/cms/setup/initialize_repo")
end

Given %r{^I am on the initialize repository page$} do
  path = "/cms/setup/initialize_repo"
  visit(path)
  expect(page.current_path).to eql(path)
end
