Given %r{^I am on the new document page$} do
  path = "/cms/documents/new"
  visit(path)
  expect(page.current_path).to eql(path)
end

Then %r{^I should see an editor with the following buttons:$} do |table|
  within("#editor-container") do
    expect(page).to have_css(".CodeMirror")

    table.transpose.raw.flatten.each do |button|
      expect(page).to have_css("a[title^='#{button}']")
    end
  end

end