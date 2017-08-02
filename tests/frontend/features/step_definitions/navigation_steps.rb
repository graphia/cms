When %r{^I click the "(.*?)" button$} do |button_text|
  # FIXME this works for a link styled as a button but not
  # for an actual `<button/>` or `<input type="submit"/>`
  button = page.find("a", text: button_text)
  scroll_into_view(button)
  page.click_link button_text
end

Then %r{^I should be redirected to the documents index$} do
  expect(page.current_path).to eql("/cms/documents")
end

Then %r{^I should (?:have been|be) redirected to "(.*)"$} do |path|
  expect(page.current_path).to eql(path)
end
