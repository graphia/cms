Given %r{^I visit the theme settings page$} do
  path = "/cms/settings/theme"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^the page should contain the following dynamic snippets:$} do |table|
  within(".theme-settings") do
    table.transpose.raw.flatten.each do |snippet|
      expect(page).to have_content(snippet)
    end
  end
end

Then %r{^the following subheadings:$} do |table|
  table.transpose.raw.flatten.each do |st|
    expect(page).to have_css("h2", text: st)
  end
end