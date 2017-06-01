Given %r{I am on the login screen} do
  path = "/cms/login"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a '(.*)' field with type '(.*)'$} do |name, field_type|
  expect(page).to have_css("input[name='#{name.downcase}'][type='#{field_type}']")
end

Then %r{^the should be a submit button should be labelled 'Login'$} do
  pending # Write code here that turns the phrase above into concrete actions
end