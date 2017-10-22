When %r{^I click the "(.*?)" button(?: again)?$} do |button_text|
  button = page.find("a,button", text: button_text)
  scroll_into_view(button)
  button.click
end

Then %r{^I should be redirected to the documents index$} do
  expect(page.current_path).to eql("/cms/documents")
end

Then %r{^I should (?:have been|be) redirected to "(.*)"$} do |path|
  expect(page.current_path).to eql(path)
end

Then %r{^there should be a '(.*?)' link$} do |text|
  expect(page).to have_css("a", text: text)
end