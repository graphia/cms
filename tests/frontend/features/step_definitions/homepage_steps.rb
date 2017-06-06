Given %r{^I am on the homepage$} do
  path = "/cms/"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see a summary of recent changes$} do
  expect(page).to have_css("h4", text: "Recent Updates")
end

Then %r{^I should see a statistics section$} do
  expect(page).to have_css("h4", text: "Statistics")
end