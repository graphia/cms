Given %r{^there are no users$} do
end

When %r{^I navigate to the login page$} do
  path = '/cms/login'
  visit(path)
end

Then %r{^I should be redirected to the initial setup page$} do
  expect(page.current_path).to eql('/cms/setup')
end