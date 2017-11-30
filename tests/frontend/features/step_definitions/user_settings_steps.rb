Given %r{^I am on my settings page$} do
  path = "/cms/settings"
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