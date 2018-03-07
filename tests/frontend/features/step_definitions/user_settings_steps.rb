Given %r{^I am on my settings page$} do
  path = "/cms/settings/my_profile"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see the heading "(.*?)"$} do |text|
  expect(page).to have_css("h1", text: text)
end

Then %r{^I should see subheadings:$} do |table|
  table.transpose.raw.each do |sh|
    expect(page).to have_css("h4", sh)
  end
end

When %r{^I change my name to "(.*?)"$} do |name|
  fill_in "name", with: name
end

Then %r{^my name should have changed to "(.*?)"$} do |name|
  within("nav.navbar") do
    expect(page).to have_css("li#user-dropdown", text: name)
  end
end